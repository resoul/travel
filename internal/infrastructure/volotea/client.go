package volotea

import (
	"net/http"
	"time"

	"github.com/resoul/travel/internal/domain"
)

const baseURL = "https://json.volotea.com"

var _ domain.VoloteaProvider = (*Client)(nil)

// Client handles communication with the Volotea JSON API endpoints.
type Client struct {
	http *http.Client
}

// NewClient creates a new Volotea API client.
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
