package itabus

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/resoul/travel/internal/domain"
)

// GetStations fetches all ItaBus stations/cities.
func (c *Client) GetStations(ctx context.Context) ([]domain.Airport, error) {
	apiURL := baseURL + "/Api-Stations"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create stations request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stations: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ItaBus stations API error: status %d", resp.StatusCode)
	}

	var raw StationsResponseDTO
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("failed to decode stations response: %w", err)
	}

	airports := make([]domain.Airport, 0, len(raw.Data))
	for _, station := range raw.Data {
		airports = append(airports, domain.Airport{
			Code: station.Code,
			Name: station.Name,
			City: domain.City{
				Code: station.Code,
				Name: station.Name,
			},
			Country: domain.Country{
				Code: station.CountryCode,
			},
		})
	}

	return airports, nil
}
