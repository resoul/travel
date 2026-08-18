package wizzair

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/resoul/travel/internal/domain"
)

// GetMap fetches all connected cities/airports from the Wizzair map API.
func (c *Client) GetMap(ctx context.Context) ([]domain.City, error) {
	buildURL, err := c.FetchBuildURL(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get build URL for map: %w", err)
	}

	apiURL := fmt.Sprintf("%s/Api/asset/map?languageCode=en-gb&withConnections=false", buildURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create map request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch map: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wizzair map API error: status %d", resp.StatusCode)
	}

	var raw mapResponseDTO
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("failed to decode map response: %w", err)
	}

	cities := make([]domain.City, 0, len(raw.Cities))
	for _, dto := range raw.Cities {
		cities = append(cities, dto.toDomain())
	}

	return cities, nil
}
