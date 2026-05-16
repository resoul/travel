package wizzair

import "encoding/json"

func (c *Client) FetchMap(buildURL string) (*MapResponse, error) {
	url := buildURL +
		"/Api/asset/map?languageCode=en-gb&withConnections=false"

	resp, err := c.HTTP.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result MapResponse

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
