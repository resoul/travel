package ryanair

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/resoul/travel/internal/domain"
)

// GetRoutes returns all airports reachable from the given IATA code.
func (c *Client) GetRoutes(ctx context.Context, originIATA string) ([]domain.Airport, error) {
	url := fmt.Sprintf("%s/api/views/locate/searchWidget/routes/en/airport/%s", baseURL, originIATA)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create routes request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch routes: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ryanair routes API error: status %d", resp.StatusCode)
	}

	var dtos []routeItemDTO
	if err := json.NewDecoder(resp.Body).Decode(&dtos); err != nil {
		return nil, fmt.Errorf("failed to decode routes response: %w", err)
	}

	airports := make([]domain.Airport, 0, len(dtos))
	for _, dto := range dtos {
		airports = append(airports, dto.ArrivalAirport.toDomain())
	}

	return airports, nil
}
