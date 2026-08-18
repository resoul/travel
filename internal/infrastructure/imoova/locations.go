package imoova

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/resoul/travel/internal/domain"
)

const (
	relocationPointsQuery = `
query RelocationPoints($input: RelocationPointInput!) {
  relocationPoints(input: $input) {
    city_name
    city_slug
    region
    lat
    lng
    count
    relocation_ids
    secondary_points {
      city_name
      count
      relocation_ids
    }
  }
}`

	homepageQuery = `
query HomepageQuery {
  relocationsByRegion {
    count
    region
  }
}`
)

// GetLocations returns all active departure cities and routes in imoova network.
func (c *Client) GetLocations(ctx context.Context) ([]domain.Airport, error) {
	req := graphQLRequest{
		Query:         relocationPointsQuery,
		OperationName: "RelocationPoints",
		Variables: map[string]interface{}{
			"input": map[string]string{
				"type": "DEPART_FROM",
			},
		},
	}

	body, err := c.query(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch imoova locations: %w", err)
	}

	var resp relocationPointsResponseDTO
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode imoova locations: %w", err)
	}

	var airports []domain.Airport
	for _, p := range resp.Data.RelocationPoints {
		airports = append(airports, p.toDomain())
	}

	sort.Slice(airports, func(i, j int) bool {
		return airports[i].City.Name < airports[j].City.Name
	})

	return airports, nil
}

// GetRegions returns deal counts per region (US, CA, EU, AU, NZ, SA, etc.).
func (c *Client) GetRegions(ctx context.Context) ([]domain.Country, error) {
	req := graphQLRequest{
		Query:         homepageQuery,
		OperationName: "HomepageQuery",
	}

	body, err := c.query(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch imoova regions: %w", err)
	}

	var resp homepageQueryResponseDTO
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode imoova regions: %w", err)
	}

	regionNames := map[string]string{
		"US": "United States",
		"CA": "Canada",
		"EU": "Europe & UK",
		"AU": "Australia",
		"NZ": "New Zealand",
		"SA": "South America (Chile)",
		"JP": "Japan",
		"TW": "Taiwan",
		"KR": "South Korea",
		"MX": "Mexico",
	}

	var countries []domain.Country
	for _, r := range resp.Data.RelocationsByRegion {
		name := regionNames[r.Region]
		if name == "" {
			name = r.Region
		}
		countries = append(countries, domain.Country{
			Code:     r.Region,
			Name:     fmt.Sprintf("%s (%d deals)", name, r.Count),
			Currency: "",
		})
	}

	return countries, nil
}
