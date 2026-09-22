package itabus

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/resoul/travel/internal/domain"
)

// GetFareCalendar fetches available dates and minimum fares between two cities.
func (c *Client) GetFareCalendar(ctx context.Context, origin, destination string, startDate, endDate string) ([]domain.FlightOffer, error) {
	if origin == "" || destination == "" {
		return nil, fmt.Errorf("origin and destination are required")
	}

	// Validate date formats
	if _, err := time.Parse("2006-01-02", startDate); err != nil {
		return nil, fmt.Errorf("invalid start date format: %w", err)
	}

	if _, err := time.Parse("2006-01-02", endDate); err != nil {
		return nil, fmt.Errorf("invalid end date format: %w", err)
	}

	params := url.Values{}
	params.Set("origin", origin)
	params.Set("destination", destination)
	params.Set("datestart", startDate)
	params.Set("dateend", endDate)
	params.Set("direction", "outbound")
	params.Set("adults", "1")
	params.Set("children", "0")

	apiURL := fmt.Sprintf("%s/Api-Dates?%s", baseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create dates request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch fare calendar: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ItaBus dates API error: status %d", resp.StatusCode)
	}

	var raw DatesResponseDTO
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("failed to decode dates response: %w", err)
	}

	offers := make([]domain.FlightOffer, 0, len(raw.Data))
	for _, date := range raw.Data {
		dateObj, err := time.Parse("2006-01-02", date.Date)
		if err != nil {
			continue
		}

		offers = append(offers, domain.FlightOffer{
			TransportType:    domain.TransportTypeBus,
			Airline:          "ItaBus",
			DepartureStation: origin,
			ArrivalStation:   destination,
			DepartureRaw:     date.Date,
			DepartureTime:    &dateObj,
			Price: domain.Price{
				Amount:   date.Amount,
				Currency: "EUR",
			},
			IsAvailable: true,
			Status:      "available",
		})
	}

	return offers, nil
}
