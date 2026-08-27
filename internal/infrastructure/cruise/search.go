package cruise

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/resoul/travel/internal/domain"
)

// SearchCruises searches available cruises matching given criteria.
func (c *Client) SearchCruises(ctx context.Context, criteria domain.CruiseSearchCriteria) ([]domain.FlightOffer, error) {
	token, err := c.getToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to obtain cruise access token: %w", err)
	}

	limit := criteria.Limit
	if limit <= 0 {
		limit = 25
	}

	now := time.Now()
	month := criteria.Month
	if month <= 0 {
		month = int(now.Month())
	}
	year := criteria.Year
	if year <= 0 {
		year = now.Year()
	}

	dateParam := fmt.Sprintf("%02dX%04d", month, year)

	queryParams := url.Values{}
	queryParams.Set("cobrandId", cobrandID)
	queryParams.Set("pin", pin)
	queryParams.Set("languageId", "1")
	queryParams.Set("fromApplication", "3")
	queryParams.Set("isGuest", "true")
	queryParams.Set("isDri", "false")
	queryParams.Set("displayType", "rewards")
	queryParams.Set("primarySort", "bestvalue")
	queryParams.Set("numberOfRecords", strconv.Itoa(limit))
	queryParams.Set("offset", "0")
	queryParams.Set("startDate", dateParam)
	queryParams.Set("date", dateParam)

	if criteria.DestinationID != "" {
		queryParams.Set("destinationList", criteria.DestinationID)
	}
	if criteria.CruiseLineID != "" {
		queryParams.Set("cruiseLineList", criteria.CruiseLineID)
	}
	if criteria.DurationMin > 0 {
		queryParams.Set("durationStart", strconv.Itoa(criteria.DurationMin))
	}
	if criteria.DurationMax > 0 {
		queryParams.Set("durationEnd", strconv.Itoa(criteria.DurationMax))
	}

	fullURL := searchURL + "?" + queryParams.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create search request: %w", err)
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Origin", originHeader)
	req.Header.Set("Referer", refererHeader)
	req.Header.Set("CruiseWebApiKey", apiKey)
	req.Header.Set("AppDomain", appDomain)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cruise search request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cruise search API returned status %d", resp.StatusCode)
	}

	var searchResp SearchResultsResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to decode search results: %w", err)
	}

	specials := searchResp.SearchResults.CruiseSpecialList
	offers := make([]domain.FlightOffer, 0, len(specials))

	for _, sp := range specials {
		s := sp.Sailing

		depPort := s.DeparturePort.Name
		if depPort == "" {
			depPort = "Port of Departure"
		}
		arrPort := s.ArrivalPort.Name
		if arrPort == "" {
			arrPort = "Port of Arrival"
		}

		var depTime *time.Time
		if parsed, err := time.Parse("01/02/2006", s.SailingDate); err == nil {
			depTime = &parsed
		} else if parsed, err := time.Parse("2006-01-02", s.SailingDate); err == nil {
			depTime = &parsed
		}

		var price float64
		if sp.FromPrice != nil {
			price = *sp.FromPrice
		}

		durationStr := fmt.Sprintf("%d %s", s.Duration, s.DaysOrNights)
		shipInfo := s.Ship.Name
		if shipInfo != "" {
			shipInfo = fmt.Sprintf("%s (%s)", shipInfo, durationStr)
		} else {
			shipInfo = durationStr
		}

		offers = append(offers, domain.FlightOffer{
			TransportType:    domain.TransportTypeCruise,
			Airline:          s.CruiseLine.Name,
			FlightNumber:     shipInfo,
			DepartureStation: depPort,
			ArrivalStation:   arrPort,
			DepartureTime:    depTime,
			DepartureRaw:     s.SailingDate,
			Duration:         durationStr,
			Price: domain.Price{
				Amount:   price,
				Currency: "USD",
			},
			IsAvailable: price > 0,
			Status:      s.Itinerary.Name,
		})
	}

	return offers, nil
}
