package flytap

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/resoul/travel/internal/domain"
)

// GetAirports retrieves operating origin airports in the TAP Air Portugal network.
func (c *Client) GetAirports(ctx context.Context) ([]domain.Airport, error) {
	reqPayload := AirportSearchRequest{
		AirlineIds:   []string{"TP"},
		Language:     "en",
		Market:       "PT",
		PayWithMiles: false,
		TripType:     "O",
	}

	bodyBytes, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal origin search payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, originSearchURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create origin search request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Origin", originHeader)
	req.Header.Set("Referer", refererHeader)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("origin search request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("origin search API returned status %d", resp.StatusCode)
	}

	var dtos []AirportItemDTO
	if err := json.NewDecoder(resp.Body).Decode(&dtos); err != nil {
		return nil, fmt.Errorf("failed to decode origin search response: %w", err)
	}

	seen := make(map[string]bool)
	airports := make([]domain.Airport, 0, len(dtos))

	for _, dto := range dtos {
		if dto.Airport == "" || seen[dto.Airport] {
			continue
		}
		// Include TAP operated routes or major stations
		if dto.TapRoute || dto.Airport == "LIS" || dto.Airport == "OPO" || dto.Airport == "FNC" {
			seen[dto.Airport] = true
			airports = append(airports, dto.toDomain())
		}
	}

	// If filtered list is empty, fallback to all returned
	if len(airports) == 0 {
		for _, dto := range dtos {
			if dto.Airport == "" || seen[dto.Airport] {
				continue
			}
			seen[dto.Airport] = true
			airports = append(airports, dto.toDomain())
		}
	}

	sort.Slice(airports, func(i, j int) bool {
		return airports[i].Code < airports[j].Code
	})

	return airports, nil
}

// GetRoutes retrieves destination airports reachable from a given origin.
func (c *Client) GetRoutes(ctx context.Context, originIATA string) ([]domain.Airport, error) {
	originIATA = strings.ToUpper(strings.TrimSpace(originIATA))
	if originIATA == "" {
		return nil, fmt.Errorf("origin IATA code is required")
	}

	reqPayload := AirportSearchRequest{
		AirlineIds:   []string{"TP"},
		Language:     "en",
		Market:       "PT",
		Origin:       originIATA,
		PayWithMiles: false,
		TripType:     "O",
	}

	bodyBytes, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal destination search payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, destSearchURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create destination search request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Origin", originHeader)
	req.Header.Set("Referer", refererHeader)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("destination search request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("destination search API returned status %d", resp.StatusCode)
	}

	var dtos []AirportItemDTO
	if err := json.NewDecoder(resp.Body).Decode(&dtos); err != nil {
		return nil, fmt.Errorf("failed to decode destination search response: %w", err)
	}

	seen := make(map[string]bool)
	airports := make([]domain.Airport, 0, len(dtos))

	for _, dto := range dtos {
		if dto.Airport == "" || seen[dto.Airport] {
			continue
		}
		seen[dto.Airport] = true
		airports = append(airports, dto.toDomain())
	}

	sort.Slice(airports, func(i, j int) bool {
		return airports[i].Code < airports[j].Code
	})

	return airports, nil
}
