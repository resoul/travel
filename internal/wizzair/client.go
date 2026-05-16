package wizzair

import (
	"net/http"
	"time"
)

type Client struct {
	HTTP *http.Client
}

func New() *Client {
	return &Client{
		HTTP: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}
