package ryanair

import (
	"encoding/json"
	"fmt"
)

// FetchRoutes returns all airports reachable from the given IATA code.
func (c *Client) FetchRoutes(origin string) ([]Airport, error) {
	url := fmt.Sprintf(
		"%s/api/views/locate/searchWidget/routes/en/airport/%s",
		baseURL,
		origin,
	)

	resp, err := c.HTTP.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("ryanair routes API: status %d", resp.StatusCode)
	}

	var airports []Airport
	if err := json.NewDecoder(resp.Body).Decode(&airports); err != nil {
		return nil, err
	}

	return airports, nil
}
