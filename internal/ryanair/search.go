package ryanair

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// Search fetches available fares and returns a flat list of FlightResult.
// It uses the flex-day window to catch nearby dates (default ±2 days).
func (c *Client) Search(req FlightRequest) ([]FlightResult, error) {
	flexBefore := req.FlexDaysBefore
	flexAfter := req.FlexDaysAfter
	if flexBefore == 0 && flexAfter == 0 {
		flexBefore, flexAfter = 2, 2
	}

	params := url.Values{}
	params.Set("ADT", strconv.Itoa(req.Adults))
	params.Set("TEEN", strconv.Itoa(req.Teens))
	params.Set("CHD", strconv.Itoa(req.Children))
	params.Set("INF", strconv.Itoa(req.Infants))
	params.Set("Origin", req.Origin)
	params.Set("Destination", req.Destination)
	params.Set("DateOut", req.DateOut)
	params.Set("FlexDaysBeforeOut", strconv.Itoa(flexBefore))
	params.Set("FlexDaysOut", strconv.Itoa(flexAfter))
	params.Set("RoundTrip", strconv.FormatBool(req.RoundTrip))
	params.Set("IncludeConnectingFlights", "false")
	params.Set("IncludePrimeFares", "false")
	params.Set("promoCode", "")
	params.Set("ToUs", "AGREED")

	if req.RoundTrip && req.DateIn != "" {
		params.Set("DateIn", req.DateIn)
		params.Set("FlexDaysBeforeIn", strconv.Itoa(flexBefore))
		params.Set("FlexDaysIn", strconv.Itoa(flexAfter))
	}

	apiURL := fmt.Sprintf(
		"%s/api/booking/v4/en-us/availability?%s",
		baseURL,
		params.Encode(),
	)

	httpReq, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := c.HTTP.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ryanair availability API: status %d", resp.StatusCode)
	}

	var raw AvailabilityResponse
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	var results []FlightResult
	for _, trip := range raw.Trips {
		for _, date := range trip.Dates {
			for _, flight := range date.Flights {
				price := adtPrice(flight.RegularFare.Fares)
				if price <= 0 {
					continue
				}

				dep, arr := "", ""
				if len(flight.Time) >= 2 {
					dep, arr = flight.Time[0], flight.Time[1]
				}

				results = append(results, FlightResult{
					DepartureStation: trip.Origin,
					ArrivalStation:   trip.Destination,
					FlightNumber:     flight.FlightNumber,
					DepartureLocal:   dep,
					ArrivalLocal:     arr,
					Duration:         flight.Duration,
					Price:            price,
					Currency:         raw.Currency,
					FaresLeft:        flight.FaresLeft,
				})
			}
		}
	}

	return results, nil
}

// adtPrice extracts the ADT (adult) fare amount from the fares list.
func adtPrice(fares []Fare) float64 {
	for _, f := range fares {
		if f.Type == "ADT" {
			return f.Amount
		}
	}
	// fallback: first fare
	if len(fares) > 0 {
		return fares[0].Amount
	}
	return 0
}
