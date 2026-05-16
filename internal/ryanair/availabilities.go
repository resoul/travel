package ryanair

import (
	"encoding/json"
	"fmt"
)

// FetchAvailabilities returns all dates with scheduled flights
// between origin and destination (one-way).
func (c *Client) FetchAvailabilities(origin, destination string) ([]string, error) {
	url := fmt.Sprintf(
		"%s/api/farfnd/3/oneWayFares/%s/%s/availabilities",
		baseURL,
		origin,
		destination,
	)

	resp, err := c.HTTP.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("ryanair availabilities API: status %d", resp.StatusCode)
	}

	var dates []string
	if err := json.NewDecoder(resp.Body).Decode(&dates); err != nil {
		return nil, err
	}

	return dates, nil
}
