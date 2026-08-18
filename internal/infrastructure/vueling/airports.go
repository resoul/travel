package vueling

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/resoul/travel/internal/domain"
)

// GetAirports retrieves all active airports in the Vueling network via AirTRFX.
func (c *Client) GetAirports(ctx context.Context) ([]domain.Airport, error) {
	url := fmt.Sprintf("%s/hangar-service/v2/vy/airports/search", airtrfxBaseURL)

	payload := map[string]any{
		"outputFields": []string{"locationLabel", "name", "cityName", "country", "iataCode"},
		"setting": map[string]string{
			"airportSource": "TRFX",
			"routeSource":   "TRFX",
		},
		"sortingDetails": []map[string]string{
			{"field": "cityName", "order": "ASC"},
		},
		"from":     0,
		"size":     6000,
		"language": "en",
		"routeOption": map[string]string{
			"airportType": "ORIGIN",
		},
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal airports payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create airports request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("em-api-key", c.apiKey)
	req.Header.Set("Origin", "https://www.vueling.com")
	req.Header.Set("Referer", "https://www.vueling.com/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch vueling airports: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("vueling airports API error: status %d", resp.StatusCode)
	}

	var dtos []airtrfxAirportDTO
	if err := json.NewDecoder(resp.Body).Decode(&dtos); err != nil {
		return nil, fmt.Errorf("failed to decode vueling airports: %w", err)
	}

	airports := make([]domain.Airport, 0, len(dtos))
	for _, dto := range dtos {
		airports = append(airports, dto.toDomain())
	}

	sort.Slice(airports, func(i, j int) bool {
		return airports[i].Code < airports[j].Code
	})

	return airports, nil
}
