package wizzair

import (
	"net/http"
	"time"

	"github.com/resoul/travel/internal/domain"
)

var _ domain.WizzairProvider = (*Client)(nil)

// Client handles communication with the Wizzair API.
type Client struct {
	http *http.Client
}

// NewClient creates a new Wizzair API client.
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
