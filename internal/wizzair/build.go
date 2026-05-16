package wizzair

import (
	"fmt"
	"io"
	"strings"
)

func (c *Client) FetchBuildURL() (string, error) {
	resp, err := c.HTTP.Get(
		"https://www.wizzair.com/buildnumber",
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	text := strings.TrimSpace(string(body))

	parts := strings.Split(text, " ")

	if len(parts) < 2 {
		return "", fmt.Errorf("invalid response")
	}

	return parts[1], nil
}
