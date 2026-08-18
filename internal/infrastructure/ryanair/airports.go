package ryanair

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/resoul/travel/internal/domain"
)

// GetAirports retrieves all active Ryanair airports.
func (c *Client) GetAirports(ctx context.Context) ([]domain.Airport, error) {
	url := fmt.Sprintf("%s/api/views/locate/5/airports/en/active", baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create airports request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch airports: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ryanair airports API error: status %d", resp.StatusCode)
	}

	var dtos []airportDTO
	if err := json.NewDecoder(resp.Body).Decode(&dtos); err != nil {
		return nil, fmt.Errorf("failed to decode airports response: %w", err)
	}

	airports := make([]domain.Airport, 0, len(dtos))
	for _, dto := range dtos {
		airports = append(airports, dto.toDomain())
	}

	return airports, nil
}
