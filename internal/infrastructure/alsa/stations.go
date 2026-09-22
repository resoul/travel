package alsa

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/resoul/travel/internal/domain"
)

// GetOrigins fetches all origin stations via HTTP API.
func (c *Client) GetOrigins(ctx context.Context) ([]domain.Airport, error) {
	params := url.Values{}
	params.Set("p_p_id", "com_babel_alsa_espania_journeysearch_web_JourneySearchPortlet_INSTANCE_21651890")
	params.Set("p_p_lifecycle", "2")
	params.Set("p_p_state", "normal")
	params.Set("p_p_mode", "view")
	params.Set("p_p_resource_id", "GetOriginsResourceCommand")
	params.Set("p_p_cacheability", "cacheLevelPage")
	params.Set("usualStations", "true")

	apiURL := fmt.Sprintf("%s/home?%s", baseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create origins request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch origins: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ALSA origins API error: status %d", resp.StatusCode)
	}

	var raw []StationDTO
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("failed to decode origins response: %w", err)
	}

	return stationsToAirports(raw), nil
}

// GetDestinations fetches available destinations for a given origin via HTTP API.
func (c *Client) GetDestinations(ctx context.Context, originStationID string) ([]domain.Airport, error) {
	params := url.Values{}
	params.Set("p_p_id", "com_babel_alsa_espania_journeysearch_web_JourneySearchPortlet_INSTANCE_21651890")
	params.Set("p_p_lifecycle", "2")
	params.Set("p_p_state", "normal")
	params.Set("p_p_mode", "view")
	params.Set("p_p_resource_id", "GetDestinationsResourceCommand")
	params.Set("p_p_cacheability", "cacheLevelPage")
	params.Set("usualStations", "true")
	params.Set("originStationId", originStationID)

	apiURL := fmt.Sprintf("%s/home?%s", baseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create destinations request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch destinations: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ALSA destinations API error: status %d", resp.StatusCode)
	}

	var raw []StationDTO
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("failed to decode destinations response: %w", err)
	}

	return stationsToAirports(raw), nil
}

func stationsToAirports(stations []StationDTO) []domain.Airport {
	var airports []domain.Airport
	seen := make(map[string]bool)

	for _, station := range stations {
		if !seen[station.ID] {
			airports = append(airports, domain.Airport{
				Code: station.ID,
				Name: station.SimplifiedName,
				City: domain.City{
					Code: station.ID,
					Name: station.ProvinceCountryName,
				},
			})
			seen[station.ID] = true
		}

		// Add individual stops as well
		for _, stop := range station.StopsList {
			if !seen[stop.ID] {
				airports = append(airports, domain.Airport{
					Code: stop.ID,
					Name: stop.Name,
					City: domain.City{
						Code: stop.ID,
						Name: stop.ProvinceCountryName,
					},
				})
				seen[stop.ID] = true
			}
		}
	}

	return airports
}
