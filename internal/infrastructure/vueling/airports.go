package vueling

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/resoul/travel/internal/domain"
)

const (
	stationsCDNURL  = "https://tickets.vueling.com/assets/1303296/stations/en-GB.json"
	countriesCDNURL = "https://tickets.vueling.com/assets/1303296/countries/en-GB.json"
)

// GetAirports retrieves all active airports in the Vueling network via public CDN.
func (c *Client) GetAirports(ctx context.Context) ([]domain.Airport, error) {
	// 1. Fetch countries
	reqCountries, err := http.NewRequestWithContext(ctx, http.MethodGet, countriesCDNURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create countries request: %w", err)
	}
	reqCountries.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	respCountries, err := c.http.Do(reqCountries)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch vueling countries: %w", err)
	}
	defer respCountries.Body.Close()

	countryMap := make(map[string]string)
	if respCountries.StatusCode == http.StatusOK {
		var countriesDTO []countryItemDTO
		if err := json.NewDecoder(respCountries.Body).Decode(&countriesDTO); err == nil {
			for _, c := range countriesDTO {
				countryMap[strings.ToUpper(c.CountryCode)] = c.Name
			}
		}
	}

	// 2. Fetch stations
	reqStations, err := http.NewRequestWithContext(ctx, http.MethodGet, stationsCDNURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create stations request: %w", err)
	}
	reqStations.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	respStations, err := c.http.Do(reqStations)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch vueling stations: %w", err)
	}
	defer respStations.Body.Close()

	if respStations.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("vueling stations API error: status %d", respStations.StatusCode)
	}

	var stationDTOs []stationItemDTO
	if err := json.NewDecoder(respStations.Body).Decode(&stationDTOs); err != nil {
		return nil, fmt.Errorf("failed to decode vueling stations: %w", err)
	}

	airports := make([]domain.Airport, 0, len(stationDTOs))
	for _, s := range stationDTOs {
		if s.InActive || s.StationCode == "" {
			continue
		}

		countryCode := strings.ToUpper(s.LocationDetails.CountryCode)
		countryName := countryMap[countryCode]
		if countryName == "" {
			countryName = countryCode
		}

		cityName := s.ShortName
		if cityName == "" {
			cityName = s.FullName
		}

		airports = append(airports, domain.Airport{
			Code: strings.ToUpper(s.StationCode),
			Name: s.FullName,
			City: domain.City{
				Name: cityName,
			},
			Country: domain.Country{
				Code: countryCode,
				Name: countryName,
			},
		})
	}

	sort.Slice(airports, func(i, j int) bool {
		return airports[i].Code < airports[j].Code
	})

	return airports, nil
}
