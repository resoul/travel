package flixbus

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/resoul/travel/internal/domain"
)

// GetReachable returns destinations reachable from a city in the FlixBus network.
func (c *Client) GetReachable(ctx context.Context, cityQueryOrID string, limit int) ([]domain.Airport, error) {
	cityUUID, _, err := c.ResolveCityID(ctx, cityQueryOrID)
	if err != nil {
		return nil, err
	}

	if limit <= 0 {
		limit = 50
	}

	endpoint := fmt.Sprintf("%s/cms/cities/%s/reachable?language=en-us&country=US&limit=%d", baseURL, cityUUID, limit)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create reachable request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch reachable destinations: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("flixbus reachable error: status %d", resp.StatusCode)
	}

	var response reachableResponseDTO
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode reachable response: %w", err)
	}

	results := make([]domain.Airport, 0, len(response.Result))
	for _, item := range response.Result {
		results = append(results, item.toDomain())
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})

	return results, nil
}
