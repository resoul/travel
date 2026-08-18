package flixbus

import (
	"net/http"
	"time"

	"github.com/resoul/travel/internal/domain"
)

const (
	baseURL = "https://global.api.flixbus.com"
)

var _ domain.FlixBusProvider = (*Client)(nil)

// Client handles communication with FlixBus API endpoints.
type Client struct {
	http *http.Client
}

// NewClient creates a new FlixBus API client with transport support.
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
