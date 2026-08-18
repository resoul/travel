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

// GetRoutes returns all destinations reachable from originIATA in the Vueling network.
func (c *Client) GetRoutes(ctx context.Context, originIATA string) ([]domain.Airport, error) {
	originIATA = strings.ToUpper(originIATA)

	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get vueling auth token for routes: %w", err)
	}

	url := fmt.Sprintf("%s/res/v1/Markets/ByOrigin/%s", amsBaseURL, originIATA)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create routes request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch vueling routes: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("vueling routes API error: status %d", resp.StatusCode)
	}

	var markets []marketDTO
	if err := json.NewDecoder(resp.Body).Decode(&markets); err != nil {
		return nil, fmt.Errorf("failed to decode vueling routes: %w", err)
	}

	// Fetch airports to enrich names, city, country (benefiting from 1-hr file cache)
	airportsMap := make(map[string]domain.Airport)
	if allAirports, err := c.GetAirports(ctx); err == nil {
		for _, a := range allAirports {
			airportsMap[strings.ToUpper(a.Code)] = a
		}
	}

	seen := make(map[string]bool)
	var destinations []domain.Airport

	for _, m := range markets {
		dest := strings.ToUpper(m.ToCode)
		if dest != "" && !seen[dest] {
			seen[dest] = true
			if airportInfo, found := airportsMap[dest]; found {
				destinations = append(destinations, airportInfo)
			} else {
				destinations = append(destinations, domain.Airport{
					Code: dest,
					Name: dest,
				})
			}
		}
	}

	sort.Slice(destinations, func(i, j int) bool {
		return destinations[i].Code < destinations[j].Code
	})

	return destinations, nil
}
