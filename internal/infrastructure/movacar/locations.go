package movacar

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/resoul/travel/internal/domain"
)

// GetLocations returns all active cities and stations in Movacar network with their offer counts.
func (c *Client) GetLocations(ctx context.Context) ([]domain.Airport, error) {
	url := fmt.Sprintf("%s/v1/locations/offers?locale=en", baseURL)
	body, err := c.get(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Movacar locations: %w", err)
	}

	var resp locationsOffersResponseDTO
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode Movacar locations: %w", err)
	}

	seen := make(map[string]bool)
	var locations []domain.Airport

	for _, loc := range resp.Included {
		name := strings.TrimSpace(loc.Attributes.Name)
		if name == "" || seen[name] {
			continue
		}
		seen[name] = true
		locations = append(locations, loc.toDomain())
	}

	sort.Slice(locations, func(i, j int) bool {
		return locations[i].City.Name < locations[j].City.Name
	})

	return locations, nil
}
