package driiveme

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/resoul/travel/internal/domain"
)

// SearchCities queries DriiveMe city autocomplete endpoint POST /search/cities.
func (c *Client) SearchCities(ctx context.Context, term string) ([]CitySuggestionDTO, error) {
	if term == "" {
		term = "a"
	}

	data := url.Values{}
	data.Set("query", term)

	endpoint := fmt.Sprintf("/%s/search/cities", c.locale)
	body, err := c.post(ctx, endpoint, data)
	if err != nil {
		return nil, fmt.Errorf("failed to search DriiveMe cities: %w", err)
	}

	var resp citySearchResponseDTO
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode DriiveMe cities response: %w", err)
	}

	return resp.Suggestions, nil
}

// GetCities returns city suggestions as domain.Airport objects.
func (c *Client) GetCities(ctx context.Context, query string) ([]domain.Airport, error) {
	suggestions, err := c.SearchCities(ctx, query)
	if err != nil {
		return nil, err
	}

	airports := make([]domain.Airport, 0, len(suggestions))
	for _, s := range suggestions {
		airports = append(airports, domain.Airport{
			Code: strconv.Itoa(s.ID),
			Name: s.SubName,
			City: domain.City{
				Code: strconv.Itoa(s.ID),
				Name: s.Name,
			},
			Country: domain.Country{
				Name: s.SubName,
			},
		})
	}

	return airports, nil
}
