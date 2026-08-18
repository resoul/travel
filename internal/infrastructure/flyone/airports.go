package flyone

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/resoul/travel/internal/domain"
)

// fetchRoutes retrieves the full route graph from /api/Routes/get-routes.
func (c *Client) fetchRoutes(ctx context.Context) (*getRoutesResponseDTO, error) {
	url := fmt.Sprintf("%s/api/Routes/get-routes", baseURL)
	body, err := c.postJSON(ctx, url, []byte("{}"))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch FlyOne routes: %w", err)
	}

	var data getRoutesResponseDTO
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to decode FlyOne routes: %w", err)
	}

	return &data, nil
}

// GetAirports returns all departure airports available in FlyOne network.
func (c *Client) GetAirports(ctx context.Context) ([]domain.Airport, error) {
	data, err := c.fetchRoutes(ctx)
	if err != nil {
		return nil, err
	}

	airports := make([]domain.Airport, 0, len(data.Routes))
	for _, r := range data.Routes {
		airports = append(airports, r.toDomain())
	}

	sort.Slice(airports, func(i, j int) bool {
		return airports[i].Code < airports[j].Code
	})

	return airports, nil
}

// GetRoutes returns all destinations reachable from originIATA in FlyOne network.
func (c *Client) GetRoutes(ctx context.Context, originIATA string) ([]domain.Airport, error) {
	originIATA = strings.ToUpper(originIATA)

	data, err := c.fetchRoutes(ctx)
	if err != nil {
		return nil, err
	}

	airportNames := make(map[string]domain.Airport)
	for _, r := range data.Routes {
		airportNames[strings.ToUpper(r.DepCode)] = r.toDomain()
	}

	var targetRoute *routeItemDTO
	for _, r := range data.Routes {
		if strings.ToUpper(r.DepCode) == originIATA {
			targetRoute = &r
			break
		}
	}

	if targetRoute == nil {
		return nil, fmt.Errorf("origin airport %q not found in FlyOne network", originIATA)
	}

	seen := make(map[string]bool)
	var destinations []domain.Airport

	for _, arrCode := range targetRoute.ArrCodes {
		arrCode = strings.ToUpper(arrCode)
		if arrCode == "" || seen[arrCode] {
			continue
		}
		seen[arrCode] = true

		if apt, ok := airportNames[arrCode]; ok {
			destinations = append(destinations, apt)
		} else {
			destinations = append(destinations, domain.Airport{
				Code: arrCode,
				Name: arrCode,
			})
		}
	}

	sort.Slice(destinations, func(i, j int) bool {
		return destinations[i].Code < destinations[j].Code
	})

	return destinations, nil
}
