package airbaltic

import (
	"net/http"
	"time"

	"github.com/resoul/travel/internal/domain"
)

const (
	baseURL = "https://www.airbaltic.com"
)

var _ domain.AirBalticProvider = (*Client)(nil)

// Client handles communication with airBaltic API endpoints.
type Client struct {
	http *http.Client
}

// NewClient creates a new airBaltic API client with transport support.
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
