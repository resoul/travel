package ryanair

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetAvailabilities returns all dates with scheduled flights between origin and destination.
func (c *Client) GetAvailabilities(ctx context.Context, origin, destination string) ([]string, error) {
	url := fmt.Sprintf(
		"%s/api/farfnd/3/oneWayFares/%s/%s/availabilities",
		baseURL,
		origin,
		destination,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create availabilities request: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch availabilities: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ryanair availabilities API error: status %d", resp.StatusCode)
	}

	var dates []string
	if err := json.NewDecoder(resp.Body).Decode(&dates); err != nil {
		return nil, fmt.Errorf("failed to decode availabilities response: %w", err)
	}

	return dates, nil
}
