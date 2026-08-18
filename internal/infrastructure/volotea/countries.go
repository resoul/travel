package volotea

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/resoul/travel/internal/domain"
)

// GetCountries fetches all countries from Volotea dist endpoint.
func (c *Client) GetCountries(ctx context.Context) ([]domain.Country, error) {
	url := fmt.Sprintf("%s/dist/countries/countries.json", baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create countries request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch volotea countries: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("volotea countries API error: status %d", resp.StatusCode)
	}

	var dtos []countryDTO
	if err := json.NewDecoder(resp.Body).Decode(&dtos); err != nil {
		return nil, fmt.Errorf("failed to decode volotea countries: %w", err)
	}

	countries := make([]domain.Country, 0, len(dtos))
	for _, dto := range dtos {
		countries = append(countries, dto.toDomain())
	}

	return countries, nil
}
