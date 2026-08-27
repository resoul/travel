package sata

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/resoul/travel/internal/domain"
)

const baseURL = "https://www.azoresairlines.pt/en"

var _ domain.SATAProvider = (*Client)(nil)

// Client handles communication with Azores Airlines / SATA API.
type Client struct {
	http *http.Client
}

// NewClient creates a new Azores Airlines client.
func NewClient(transport ...http.RoundTripper) *Client {
	var tr http.RoundTripper
	if len(transport) > 0 && transport[0] != nil {
		tr = transport[0]
	}

	return &Client{
		http: &http.Client{
			Transport: tr,
			Timeout:   20 * time.Second,
		},
	}
}

// GetAirports retrieves all airports served by Azores Airlines / SATA.
func (c *Client) GetAirports(ctx context.Context) ([]domain.Airport, error) {
	reqURL := fmt.Sprintf("%s/flight-search/airports", baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create airports request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch airports: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("airports request failed with status: %d", resp.StatusCode)
	}

	var dtos []AirportDTO
	if err := json.NewDecoder(resp.Body).Decode(&dtos); err != nil {
		return nil, fmt.Errorf("failed to decode airports response: %w", err)
	}

	airports := make([]domain.Airport, 0, len(dtos))
	for _, dto := range dtos {
		airports = append(airports, domain.Airport{
			Code: dto.Code,
			Name: dto.Name,
			City: domain.City{
				Code: dto.CityCode,
				Name: dto.CityName,
			},
			Country: domain.Country{
				Code: dto.CountryCode,
				Name: dto.CountryName,
			},
		})
	}

	sort.Slice(airports, func(i, j int) bool {
		return airports[i].Code < airports[j].Code
	})

	return airports, nil
}

// GetRoutesFromOrigin retrieves direct destination airport codes from the given origin.
func (c *Client) GetRoutesFromOrigin(ctx context.Context, origin string) ([]string, error) {
	reqURL := fmt.Sprintf("%s/flight-search/routes/%s", baseURL, strings.ToUpper(strings.TrimSpace(origin)))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create routes request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch routes: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("routes request failed with status: %d", resp.StatusCode)
	}

	// Routes endpoint returns a JSON with data.applicableDestinations[origin]
	var destinations []string
	var rawData map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&rawData); err == nil {
		if dataMap, ok := rawData["data"].(map[string]any); ok {
			if appDests, ok := dataMap["applicableDestinations"].(map[string]any); ok {
				for _, destList := range appDests {
					if list, ok := destList.([]any); ok {
						for _, item := range list {
							if m, ok := item.(map[string]any); ok {
								if code, ok := m["code"].(string); ok && code != "" {
									destinations = append(destinations, code)
								}
							} else if s, ok := item.(string); ok && s != "" {
								destinations = append(destinations, s)
							}
						}
					}
				}
			}
		}
	}

	sort.Strings(destinations)
	return destinations, nil
}

// GetFareCalendar retrieves Azores Airlines low-fare calendar across all available future dates.
func (c *Client) GetFareCalendar(ctx context.Context, origin, destination string) ([]domain.FlightOffer, error) {
	reqURL := fmt.Sprintf("%s/flight-search/outbound/%s/%s", baseURL,
		strings.ToUpper(strings.TrimSpace(origin)),
		strings.ToUpper(strings.TrimSpace(destination)),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create calendar request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch fare calendar: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fare calendar request failed with status: %d", resp.StatusCode)
	}

	var rawMap map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&rawMap); err != nil {
		return nil, fmt.Errorf("failed to decode fare calendar: %w", err)
	}

	itemsMap := make(map[string]CalendarItemDTO)

	// Check if Dates wrapper exists
	if datesRaw, ok := rawMap["Dates"].(map[string]any); ok {
		for dateStr, itemVal := range datesRaw {
			if itemBytes, err := json.Marshal(itemVal); err == nil {
				var dto CalendarItemDTO
				if err := json.Unmarshal(itemBytes, &dto); err == nil {
					itemsMap[dateStr] = dto
				}
			}
		}
	} else {
		// Try root map
		for dateStr, itemVal := range rawMap {
			if itemBytes, err := json.Marshal(itemVal); err == nil {
				var dto CalendarItemDTO
				if err := json.Unmarshal(itemBytes, &dto); err == nil && dto.Min > 0 {
					itemsMap[dateStr] = dto
				}
			}
		}
	}

	// Sort dates chronologically
	dateKeys := make([]string, 0, len(itemsMap))
	for d := range itemsMap {
		dateKeys = append(dateKeys, d)
	}
	sort.Strings(dateKeys)

	offers := make([]domain.FlightOffer, 0, len(dateKeys))
	for _, dateKey := range dateKeys {
		item := itemsMap[dateKey]
		var depT *time.Time
		if t, err := time.Parse("2006-01-02", dateKey); err == nil {
			depT = &t
		}

		currency := item.Cur
		if currency == "" {
			currency = "EUR"
		}

		offers = append(offers, domain.FlightOffer{
			TransportType:    domain.TransportTypeFlight,
			Airline:          "Azores Airlines / SATA",
			FlightNumber:     fmt.Sprintf("S4 %s->%s", item.From, item.To),
			DepartureStation: item.From,
			ArrivalStation:   item.To,
			DepartureTime:    depT,
			DepartureRaw:     dateKey,
			Price: domain.Price{
				Amount:   item.Min,
				Currency: currency,
			},
			IsAvailable: true,
			Status:      "AVAILABLE",
		})
	}

	return offers, nil
}
