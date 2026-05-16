package ryanair

import (
	"net/http"
	"net/http/cookiejar"
	"time"

	"golang.org/x/net/publicsuffix"
)

const baseURL = "https://www.ryanair.com"

type Client struct {
	HTTP *http.Client
}

func New() *Client {
	jar, _ := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})

	return &Client{
		HTTP: &http.Client{
			Timeout: 30 * time.Second,
			Jar:     jar,
		},
	}
}
