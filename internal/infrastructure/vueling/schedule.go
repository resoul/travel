package vueling

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/resoul/travel/internal/domain"
)

// GetSchedule fetches all scheduled flights between origin and destination using browser session automation.
func (c *Client) GetSchedule(ctx context.Context, origin, destination string, year, month, monthsRange int) ([]domain.FlightOffer, error) {
	origin = strings.ToUpper(strings.TrimSpace(origin))
	destination = strings.ToUpper(strings.TrimSpace(destination))

	if origin == "" || destination == "" {
		return nil, fmt.Errorf("origin and destination are required")
	}

	now := time.Now()
	if year <= 0 {
		year = now.Year()
	}
	if month <= 0 {
		month = int(now.Month())
	}
	if monthsRange <= 0 {
		monthsRange = 12
	}

	url := "https://tickets.vueling.com"

	jsQuery := fmt.Sprintf(`
		(async () => {
			try {
				let token = "";
				const keys = Object.keys(sessionStorage).concat(Object.keys(localStorage));
				for (const k of keys) {
					const val = sessionStorage.getItem(k) || localStorage.getItem(k);
					if (val && val.includes("accessToken")) {
						try {
							const parsed = JSON.parse(val);
							if (parsed.accessToken) {
								token = parsed.accessToken;
								break;
							}
						} catch (e) {}
					}
				}

				const payload = {
					originCode: "%s",
					destinationCode: "%s",
					year: %d,
					month: %d,
					currencyCode: "EUR",
					monthsRange: %d,
					flightType: "OW"
				};

				const resp = await fetch('https://ams.vueling.com/avy/v3/AvailabilityServices/allFlights', {
					method: 'POST',
					headers: {
						'Content-Type': 'application/json',
						'Authorization': 'Bearer ' + token
					},
					body: JSON.stringify(payload)
				});

				if (!resp.ok) {
					return JSON.stringify({ error: "HTTP " + resp.status });
				}

				const data = await resp.json();
				return JSON.stringify(data);
			} catch (e) {
				return JSON.stringify({ error: e.message });
			}
		})()
	`, origin, destination, year, month, monthsRange)

	var jsOutput string
	err := c.executeInBrowser(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(4*time.Second),
		chromedp.Evaluate(jsQuery, &jsOutput, func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
			return p.WithAwaitPromise(true)
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to execute vueling schedule search in browser: %w", err)
	}

	if strings.Contains(jsOutput, `"error"`) {
		return nil, fmt.Errorf("vueling availability API error: %s", jsOutput)
	}

	var rawFlights []availabilityFlightDTO
	if err := json.Unmarshal([]byte(jsOutput), &rawFlights); err != nil {
		var objMap map[string][]availabilityFlightDTO
		if errMap := json.Unmarshal([]byte(jsOutput), &objMap); errMap == nil {
			rawFlights = objMap["flightOutboundList"]
		} else {
			return nil, fmt.Errorf("failed to decode vueling availability: %w", err)
		}
	}

	offers := make([]domain.FlightOffer, 0, len(rawFlights))
	for _, f := range rawFlights {
		if !f.IsAvailableDay && f.Price <= 0 {
			continue
		}
		flightCopy := f
		offers = append(offers, flightCopy.toDomain())
	}

	return offers, nil
}

// GetDates returns all unique scheduled departure dates between origin and destination.
func (c *Client) GetDates(ctx context.Context, origin, destination string, year, month, monthsRange int) ([]string, error) {
	offers, err := c.GetSchedule(ctx, origin, destination, year, month, monthsRange)
	if err != nil {
		return nil, err
	}

	seen := make(map[string]bool)
	var dates []string

	for _, o := range offers {
		if o.DepartureTime != nil {
			d := o.DepartureTime.Format("2006-01-02")
			if !seen[d] {
				seen[d] = true
				dates = append(dates, d)
			}
		} else if len(o.DepartureRaw) >= 10 {
			d := o.DepartureRaw[:10]
			if !seen[d] {
				seen[d] = true
				dates = append(dates, d)
			}
		}
	}

	sort.Strings(dates)
	return dates, nil
}
