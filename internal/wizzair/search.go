package wizzair

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func (c *Client) Search(
	buildURL string,
	payload FlightRequest,
) (*TimetableResponse, error) {

	url := buildURL + "/Api/search/timetableV2"

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result TimetableResponse

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
