package itabus

import (
	"net/http"
	"time"
)

const baseURL = "https://www.itabus.it/on/demandware.store/Sites-ITABUS-Site/en"

// Client handles communication with ItaBus REST API.
type Client struct {
	http *http.Client
}

// NewClient creates a new ItaBus API client.
func NewClient(transport ...http.RoundTripper) *Client {
	var tr http.RoundTripper
	if len(transport) > 0 && transport[0] != nil {
		tr = transport[0]
	}

	return &Client{
		http: &http.Client{
			Transport: tr,
			Timeout:   30 * time.Second,
		},
	}
}
