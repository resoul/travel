package cache

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type httpCachePayload struct {
	StatusCode int                 `json:"status_code"`
	Header     map[string][]string `json:"header"`
	Body       []byte              `json:"body"`
}

// CachedTransport wraps an http.RoundTripper with disk-based caching for GET requests.
type CachedTransport struct {
	base http.RoundTripper
	c    Cache
	ttl  time.Duration
}

// NewCachedTransport creates an http.RoundTripper that caches GET responses.
func NewCachedTransport(base http.RoundTripper, c Cache, ttl time.Duration) http.RoundTripper {
	if base == nil {
		base = http.DefaultTransport
	}
	if ttl <= 0 {
		ttl = 1 * time.Hour
	}
	return &CachedTransport{
		base: base,
		c:    c,
		ttl:  ttl,
	}
}

// RoundTrip executes the HTTP request or returns cached response if available.
func (t *CachedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Only cache GET and HEAD requests
	if req.Method != http.MethodGet && req.Method != http.MethodHead {
		return t.base.RoundTrip(req)
	}

	cacheKey := "http:" + req.URL.String()

	if cachedData, found, err := t.c.Get(cacheKey); err == nil && found {
		var payload httpCachePayload
		if err := json.Unmarshal(cachedData, &payload); err == nil {
			header := http.Header{}
			for k, vv := range payload.Header {
				for _, v := range vv {
					header.Add(k, v)
				}
			}
			header.Set("X-Cache", "HIT")

			resp := &http.Response{
				StatusCode: payload.StatusCode,
				Status:     http.StatusText(payload.StatusCode),
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header:     header,
				Body:       io.NopCloser(bytes.NewReader(payload.Body)),
				Request:    req,
			}
			return resp, nil
		}
	}

	resp, err := t.base.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	// Only cache 200 OK responses
	if resp.StatusCode == http.StatusOK && resp.Body != nil {
		bodyBytes, readErr := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if readErr != nil {
			return nil, readErr
		}

		payload := httpCachePayload{
			StatusCode: resp.StatusCode,
			Header:     resp.Header,
			Body:       bodyBytes,
		}

		if payloadBytes, err := json.Marshal(payload); err == nil {
			_ = t.c.Set(cacheKey, payloadBytes, t.ttl)
		}

		resp.Header.Set("X-Cache", "MISS")
		resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	}

	return resp, nil
}
