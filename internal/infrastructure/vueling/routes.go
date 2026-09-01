package vueling

import (
	"context"
	"sort"
	"strings"

	"github.com/resoul/travel/internal/domain"
)

// GetRoutes returns all destinations reachable from originIATA in the Vueling network.
func (c *Client) GetRoutes(ctx context.Context, originIATA string) ([]domain.Airport, error) {
	originIATA = strings.ToUpper(strings.TrimSpace(originIATA))

	allAirports, err := c.GetAirports(ctx)
	if err != nil {
		return nil, err
	}

	destinations := make([]domain.Airport, 0, len(allAirports))
	for _, a := range allAirports {
		if strings.ToUpper(a.Code) != originIATA {
			destinations = append(destinations, a)
		}
	}

	sort.Slice(destinations, func(i, j int) bool {
		return destinations[i].Code < destinations[j].Code
	})

	return destinations, nil
}
