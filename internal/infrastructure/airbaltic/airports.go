package airbaltic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/resoul/travel/internal/domain"
)

// fetchOrigDest loads the full network map from /api/orig-dest/en.
func (c *Client) fetchOrigDest(ctx context.Context) (*origDestResponseDTO, error) {
	url := fmt.Sprintf("%s/api/orig-dest/en", baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create orig-dest request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch airBaltic network: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("airBaltic orig-dest API error: status %d", resp.StatusCode)
	}

	var data origDestResponseDTO
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode airBaltic network: %w", err)
	}

	return &data, nil
}

// GetAirports returns all primary origin airports in the airBaltic network.
func (c *Client) GetAirports(ctx context.Context) ([]domain.Airport, error) {
	data, err := c.fetchOrigDest(ctx)
	if err != nil {
		return nil, err
	}

	airports := make([]domain.Airport, 0, len(data.OrigData.BtOrigins))
	for _, dto := range data.OrigData.BtOrigins {
		airports = append(airports, dto.toDomain())
	}

	sort.Slice(airports, func(i, j int) bool {
		return airports[i].Code < airports[j].Code
	})

	return airports, nil
}

// GetRoutes returns all destinations reachable from a given origin IATA code in the airBaltic network.
func (c *Client) GetRoutes(ctx context.Context, originIATA string) ([]domain.Airport, error) {
	originIATA = strings.ToUpper(originIATA)

	data, err := c.fetchOrigDest(ctx)
	if err != nil {
		return nil, err
	}

	// In airBaltic destinData, keys are formatted as "{IATA}A" (e.g. "RIXA", "ALCA") or "{IATA}"
	key := originIATA + "A"
	destinItem, ok := data.DestinData[key]
	if !ok {
		destinItem, ok = data.DestinData[originIATA]
	}

	if !ok {
		return nil, fmt.Errorf("origin airport %q not found in airBaltic network", originIATA)
	}

	seen := make(map[string]bool)
	var destinations []domain.Airport

	// 1. Direct airBaltic destinations
	for _, dto := range destinItem.BtDest {
		code := strings.ToUpper(dto.Code)
		if code != "" && !seen[code] {
			seen[code] = true
			destinations = append(destinations, dto.toDomain())
		}
	}

	// 2. Partner / non-BT destinations
	for _, dto := range destinItem.NonBtDest {
		code := strings.ToUpper(dto.Code)
		if code != "" && !seen[code] {
			seen[code] = true
			destinations = append(destinations, dto.toDomain())
		}
	}

	sort.Slice(destinations, func(i, j int) bool {
		return destinations[i].Code < destinations[j].Code
	})

	return destinations, nil
}
