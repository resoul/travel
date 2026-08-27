package trenitalia

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/resoul/travel/internal/domain"
)

const (
	defaultStationsURL = "https://www.trenitalia.com/content/trenitalia/en.cruscotto-stations.json"
	defaultBFFURL      = "https://www.lefrecce.it/Channels.Website.BFF.WEB/website/ticket/solutions"
	defaultUserAgent   = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
)

var _ domain.TrenitaliaProvider = (*Client)(nil)

// KnownStationIDs maps common Italian station names and aliases to their official Trenitalia Location IDs.
var KnownStationIDs = map[string]int{
	"roma termini":                830008409,
	"rome termini":                830008409,
	"roma tiburtina":              830008217,
	"rome tiburtina":              830008217,
	"milano centrale":             830001700,
	"milan centrale":              830001700,
	"milano porta garibaldi":      830001645,
	"milano rogoredo":             830001820,
	"firenze santa maria novella": 830006421,
	"firenze smn":                 830006421,
	"florence":                    830006421,
	"venezia santa lucia":         830002593,
	"venice santa lucia":          830002593,
	"venezia mestre":              830002589,
	"venice mestre":               830002589,
	"napoli centrale":             830009218,
	"naples centrale":             830009218,
	"napoli afragola":             830009956,
	"torino porta nuova":          830000219,
	"turin porta nuova":           830000219,
	"torino porta susa":           830000222,
	"bologna centrale":            830005043,
	"verona porta nuova":          830003073,
	"genova piazza principe":      830000803,
	"genoa piazza principe":       830000803,
	"genova brignole":             830000820,
	"bari centrale":               830011118,
	"pisa centrale":               830006500,
	"palermo centrale":            830012284,
	"salerno":                     830009818,
	"padova":                      830002166,
	"padua":                       830002166,
	"trieste centrale":            830004127,
	"rimini":                      830005315,
	"bergamo":                     830001007,
	"brescia":                     830001046,
	"reggio emilia av":            830005254,
	"ancona":                      830007140,
	"pescara centrale":            830007626,
	"foggia":                      830011030,
	"lecce":                       830011384,
	"taranto":                     830011500,
	"messina centrale":            830012000,
	"catania centrale":            830012111,
	"bolzano":                     830003200,
	"trento":                      830003150,
}

// Client handles communication with the Trenitalia and Le Frecce APIs.
type Client struct {
	http *http.Client
}

// NewClient creates a new Trenitalia client.
func NewClient(transport ...http.RoundTripper) *Client {
	var tr http.RoundTripper
	if len(transport) > 0 && transport[0] != nil {
		tr = transport[0]
	}

	return &Client{
		http: &http.Client{
			Transport: tr,
			Timeout:   15 * time.Second,
		},
	}
}

// GetStations downloads the global stations list from Trenitalia CDN.
func (c *Client) GetStations(ctx context.Context) ([]domain.TrenitaliaStation, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, defaultStationsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", defaultUserAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request to %s failed: %w", defaultStationsURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Trenitalia stations error: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var dtos []StationRawDTO
	if err := json.Unmarshal(body, &dtos); err != nil {
		return nil, fmt.Errorf("failed to parse stations JSON: %w", err)
	}

	stations := make([]domain.TrenitaliaStation, 0, len(dtos))
	for _, d := range dtos {
		valLower := strings.ToLower(strings.TrimSpace(d.Value))
		id := KnownStationIDs[valLower]

		stations = append(stations, domain.TrenitaliaStation{
			ID:         id,
			Name:       d.Text,
			Value:      d.Value,
			IsFrecce:   d.IsF == 1,
			IsEurocity: d.IsE == 1,
		})
	}

	return stations, nil
}

// SearchTrains searches available trains on Trenitalia matching given criteria.
func (c *Client) SearchTrains(ctx context.Context, criteria domain.TrenitaliaSearchCriteria) ([]domain.FlightOffer, error) {
	originID := criteria.OriginID
	if originID == 0 {
		normalizedOrigin := strings.ToLower(strings.TrimSpace(criteria.OriginName))
		if id, found := KnownStationIDs[normalizedOrigin]; found {
			originID = id
		} else {
			// Default to Roma Termini
			originID = 830008409
		}
	}

	destinationID := criteria.DestinationID
	if destinationID == 0 {
		normalizedDest := strings.ToLower(strings.TrimSpace(criteria.DestinationName))
		if id, found := KnownStationIDs[normalizedDest]; found {
			destinationID = id
		} else {
			// Default to Milano Centrale
			destinationID = 830001700
		}
	}

	depDate := criteria.DepartureDate
	if depDate == "" {
		depDate = time.Now().AddDate(0, 0, 7).Format("2006-01-02")
	}

	depTime := criteria.DepartureTime
	if depTime == "" {
		depTime = "08:00"
	}
	if !strings.Contains(depTime, ":") {
		depTime = depTime + ":00"
	}

	// Format ISO timestamp: YYYY-MM-DDTHH:mm:00.000+02:00
	isoDeparture := fmt.Sprintf("%sT%s:00.000+02:00", depDate, depTime)

	adults := criteria.Adults
	if adults <= 0 {
		adults = 1
	}

	limit := criteria.Limit
	if limit <= 0 {
		limit = 10
	}

	payload := SolutionsRequestDTO{
		DepartureLocationID: originID,
		ArrivalLocationID:   destinationID,
		DepartureTime:       isoDeparture,
		Adults:              adults,
		Children:            criteria.Children,
		Criteria: SearchCriteriaDTO{
			FrecceOnly:   criteria.FrecceOnly,
			RegionalOnly: criteria.RegionalOnly,
			NoChanges:    criteria.NoChanges,
			Order:        "DEPARTURE_DATE",
			Limit:        limit,
			Offset:       0,
		},
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, defaultBFFURL, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", defaultUserAgent)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Referer", "https://www.lefrecce.it/Channels.Website.WEB/")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request to %s failed: %w", defaultBFFURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Trenitalia BFF error: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var solResp SolutionsResponseDTO
	if err := json.Unmarshal(body, &solResp); err != nil {
		return nil, fmt.Errorf("failed to parse solutions response: %w", err)
	}

	offers := make([]domain.FlightOffer, 0, len(solResp.Solutions))
	for _, w := range solResp.Solutions {
		s := w.Solution

		var trainName string
		if len(s.Trains) > 0 {
			t := s.Trains[0]
			parts := []string{}
			if t.TrainCategory != "" {
				parts = append(parts, t.TrainCategory)
			}
			if t.TrainAcronym != "" {
				parts = append(parts, t.TrainAcronym)
			}
			if t.Description != "" {
				parts = append(parts, t.Description)
			}
			trainName = strings.Join(parts, " ")
		}
		if trainName == "" {
			trainName = "Trenitalia Train"
		}

		var depT, arrT *time.Time
		if t, err := time.Parse("2006-01-02T15:04:05.000-07:00", s.DepartureTime); err == nil {
			depT = &t
		}
		if t, err := time.Parse("2006-01-02T15:04:05.000-07:00", s.ArrivalTime); err == nil {
			arrT = &t
		}

		var priceVal float64
		var currency = "EUR"
		if s.Price != nil {
			priceVal = s.Price.Amount
			if s.Price.Currency != "" {
				currency = s.Price.Currency
			}
		}

		offer := domain.FlightOffer{
			TransportType:    domain.TransportTypeTrain,
			Airline:          "Trenitalia",
			FlightNumber:     trainName,
			DepartureStation: s.Origin,
			ArrivalStation:   s.Destination,
			DepartureTime:    depT,
			ArrivalTime:      arrT,
			DepartureRaw:     s.DepartureTime,
			ArrivalRaw:       s.ArrivalTime,
			Duration:         s.Duration,
			Price: domain.Price{
				Amount:   priceVal,
				Currency: currency,
			},
			IsAvailable: s.Status == "SALEABLE",
			Status:      s.Status,
		}

		offers = append(offers, offer)
	}

	return offers, nil
}
