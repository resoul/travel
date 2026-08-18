package flixbus

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	"github.com/resoul/travel/internal/domain"
)

var uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// GetCities queries the autocomplete endpoint to find FlixBus cities matching a query.
func (c *Client) GetCities(ctx context.Context, query string) ([]domain.Airport, error) {
	if query == "" {
		query = "Berlin"
	}

	endpoint := fmt.Sprintf("%s/search/autocomplete/cities?q=%s&lang=en_US&country=us", baseURL, url.QueryEscape(query))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create autocomplete request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to query flixbus autocomplete: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("flixbus autocomplete error: status %d", resp.StatusCode)
	}

	var dtos []autocompleteCityDTO
	if err := json.NewDecoder(resp.Body).Decode(&dtos); err != nil {
		return nil, fmt.Errorf("failed to decode flixbus autocomplete response: %w", err)
	}

	airports := make([]domain.Airport, 0, len(dtos))
	for _, dto := range dtos {
		airports = append(airports, dto.toDomain())
	}

	return airports, nil
}

// ResolveCityID resolves a city name or UUID to a valid FlixBus city UUID.
func (c *Client) ResolveCityID(ctx context.Context, queryOrID string) (string, string, error) {
	if uuidRegex.MatchString(queryOrID) {
		return queryOrID, queryOrID, nil
	}

	cities, err := c.GetCities(ctx, queryOrID)
	if err != nil {
		return "", "", fmt.Errorf("failed to resolve city %q: %w", queryOrID, err)
	}

	if len(cities) == 0 {
		return "", "", fmt.Errorf("no flixbus city found for %q", queryOrID)
	}

	return cities[0].Code, cities[0].Name, nil
}
