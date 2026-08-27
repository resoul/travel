package cruise

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"

	"github.com/resoul/travel/internal/domain"
)

// fetchMatrix loads the search matrix containing cruise lines and destinations.
func (c *Client) fetchMatrix(ctx context.Context) (*SearchMatrixResponse, error) {
	token, err := c.getToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to obtain cruise access token: %w", err)
	}

	queryParams := url.Values{}
	queryParams.Set("cobrandId", cobrandID)
	queryParams.Set("partnerId", partnerID)
	queryParams.Set("pin", pin)
	queryParams.Set("applicationId", "3")
	queryParams.Set("store", "")
	queryParams.Set("cert", "")
	queryParams.Set("auth", "")

	fullURL := matrixURL + "?" + queryParams.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create matrix request: %w", err)
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Origin", originHeader)
	req.Header.Set("Referer", refererHeader)
	req.Header.Set("CruiseWebApiKey", apiKey)
	req.Header.Set("AppDomain", appDomain)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("matrix request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("matrix API returned status %d", resp.StatusCode)
	}

	var matrixResp SearchMatrixResponse
	if err := json.NewDecoder(resp.Body).Decode(&matrixResp); err != nil {
		return nil, fmt.Errorf("failed to decode matrix response: %w", err)
	}

	return &matrixResp, nil
}

// GetCruiseLines retrieves all cruise operators from the matrix.
func (c *Client) GetCruiseLines(ctx context.Context) ([]domain.CruiseLine, error) {
	matrix, err := c.fetchMatrix(ctx)
	if err != nil {
		return nil, err
	}

	lines := make([]domain.CruiseLine, 0, len(matrix.CruiseLines))
	for _, l := range matrix.CruiseLines {
		if l.Name != "" {
			lines = append(lines, l.toCruiseLine())
		}
	}

	sort.Slice(lines, func(i, j int) bool {
		return lines[i].Name < lines[j].Name
	})

	return lines, nil
}

// GetCruiseDestinations retrieves all cruising destination regions from the matrix.
func (c *Client) GetCruiseDestinations(ctx context.Context) ([]domain.CruiseDestination, error) {
	matrix, err := c.fetchMatrix(ctx)
	if err != nil {
		return nil, err
	}

	destinations := make([]domain.CruiseDestination, 0, len(matrix.Destinations))
	for _, d := range matrix.Destinations {
		if d.Name != "" {
			destinations = append(destinations, d.toCruiseDestination())
		}
	}

	sort.Slice(destinations, func(i, j int) bool {
		return destinations[i].Name < destinations[j].Name
	})

	return destinations, nil
}
