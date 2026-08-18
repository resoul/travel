package volotea

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/resoul/travel/internal/domain"
)

// fetchStations is an internal helper that downloads and decodes the stations map.
func (c *Client) fetchStations(ctx context.Context) (map[string]stationDTO, error) {
	url := fmt.Sprintf("%s/dist/stations/stations.json", baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create stations request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch volotea stations: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("volotea stations API error: status %d", resp.StatusCode)
	}

	var stations map[string]stationDTO
	if err := json.NewDecoder(resp.Body).Decode(&stations); err != nil {
		return nil, fmt.Errorf("failed to decode volotea stations: %w", err)
	}

	return stations, nil
}

// GetAirports returns all airports serviced by Volotea.
func (c *Client) GetAirports(ctx context.Context) ([]domain.Airport, error) {
	stations, err := c.fetchStations(ctx)
	if err != nil {
		return nil, err
	}

	var airports []domain.Airport
	for code, s := range stations {
		stationCopy := s
		airports = append(airports, stationCopy.toDomain(code))
	}

	sort.Slice(airports, func(i, j int) bool {
		return airports[i].Code < airports[j].Code
	})

	return airports, nil
}

// GetRoutes returns all reachable airports from the specified origin IATA.
func (c *Client) GetRoutes(ctx context.Context, originIATA string) ([]domain.Airport, error) {
	stations, err := c.fetchStations(ctx)
	if err != nil {
		return nil, err
	}

	originIATA = strings.ToUpper(originIATA)
	originStation, exists := stations[originIATA]
	if !exists {
		return nil, fmt.Errorf("origin airport %q not found in Volotea network", originIATA)
	}

	var reachable []domain.Airport
	for destCode, market := range originStation.Markets {
		if !market.Enabled {
			continue
		}
		if destStation, found := stations[destCode]; found {
			destCopy := destStation
			reachable = append(reachable, destCopy.toDomain(destCode))
		} else {
			reachable = append(reachable, domain.Airport{
				Code: destCode,
				Name: destCode,
			})
		}
	}

	sort.Slice(reachable, func(i, j int) bool {
		return reachable[i].Code < reachable[j].Code
	})

	return reachable, nil
}
