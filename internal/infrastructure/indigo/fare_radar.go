package indigo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/resoul/travel/internal/domain"
)

// GetFareRadar retrieves lowest fares and destination recommendations from a given origin.
func (c *Client) GetFareRadar(ctx context.Context, originIATA string) ([]domain.IndiGoRadarFare, error) {
	originIATA = strings.ToUpper(strings.TrimSpace(originIATA))
	if originIATA == "" {
		return nil, fmt.Errorf("origin airport IATA code is required")
	}

	endpoint := fmt.Sprintf("%s?origin=%s", fareRadarBaseURL, url.QueryEscape(originIATA))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create fare radar request: %w", err)
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json, text/plain, */*")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fare radar request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fare radar returned status %d for origin %s", resp.StatusCode, originIATA)
	}

	var radarResp FareRadarResponse
	if err := json.NewDecoder(resp.Body).Decode(&radarResp); err != nil {
		return nil, fmt.Errorf("failed to decode fare radar response: %w", err)
	}

	results := make([]domain.IndiGoRadarFare, 0, len(radarResp.Fares))
	currency := radarResp.Currency
	if currency == "" {
		currency = "INR"
	}

	for _, item := range radarResp.Fares {
		results = append(results, domain.IndiGoRadarFare{
			Origin:      radarResp.Origin,
			OriginCity:  radarResp.OriginCity,
			Destination: item.IATA,
			DestCity:    item.City,
			TravelDate:  radarResp.TravelDate,
			FlightTime:  item.Time,
			Price: domain.Price{
				Amount:   item.Fare,
				Currency: currency,
			},
		})
	}

	return results, nil
}
