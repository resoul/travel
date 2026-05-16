package ryanair

import "encoding/json"

func (c *Client) FetchAirports() ([]Airport, error) {
	resp, err := c.HTTP.Get(
		baseURL + "/api/views/locate/5/airports/en/active",
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var airports []Airport
	if err := json.NewDecoder(resp.Body).Decode(&airports); err != nil {
		return nil, err
	}

	return airports, nil
}
