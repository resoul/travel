package pitchup

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/resoul/travel/internal/domain"
)

var _ domain.PitchupProvider = (*Client)(nil)

// Client handles Pitchup campsite and glamping search via Chromedp headless automation.
type Client struct {
	http *http.Client
}

// NewClient creates a new Pitchup client.
func NewClient(transport ...http.RoundTripper) *Client {
	var tr http.RoundTripper
	if len(transport) > 0 && transport[0] != nil {
		tr = transport[0]
	}

	return &Client{
		http: &http.Client{
			Transport: tr,
			Timeout:   15 * time.Second,
		},
	}
}

type rawCampsite struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Rating   string `json:"rating"`
	Price    string `json:"price"`
	UnitType string `json:"unitType"`
}

// SearchCampsites searches campsites, glamping, and holiday parks on Pitchup.
func (c *Client) SearchCampsites(ctx context.Context, criteria domain.PitchupSearchCriteria) ([]domain.FlightOffer, error) {
	country := strings.TrimSpace(criteria.Country)
	if country == "" {
		country = "France"
	}
	countryTitle := strings.Title(strings.ToLower(country))

	arriveDate := criteria.ArriveDate
	if arriveDate == "" {
		arriveDate = time.Now().AddDate(0, 0, 14).Format("2006-01-02")
	}

	departDate := criteria.DepartDate
	if departDate == "" {
		if t, err := time.Parse("2006-01-02", arriveDate); err == nil {
			departDate = t.AddDate(0, 0, 2).Format("2006-01-02")
		} else {
			departDate = time.Now().AddDate(0, 0, 16).Format("2006-01-02")
		}
	}

	adults := criteria.Adults
	if adults <= 0 {
		adults = 2
	}

	targetURL := fmt.Sprintf("https://www.pitchup.com/campsites/%s/?arrive=%s&depart=%s&adults=%d",
		countryTitle, arriveDate, departDate, adults)

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
	)

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(ctx, opts...)
	defer cancelAlloc()

	chromCtx, cancelChrom := chromedp.NewContext(allocCtx)
	defer cancelChrom()

	timeoutCtx, cancelTimeout := context.WithTimeout(chromCtx, 35*time.Second)
	defer cancelTimeout()

	var rawList []rawCampsite

	const extractJS = `
	(() => {
		const results = [];
		const cards = document.querySelectorAll("article, .campsite-card, .search-result, [class*='campsite-item']");
		const seen = new Set();

		cards.forEach(c => {
			const nameEl = c.querySelector("h2, h3, [class*='name'], [class*='title']");
			const name = nameEl ? nameEl.innerText.trim() : "";
			
			if (!name || name.startsWith("Location:") || seen.has(name)) {
				return;
			}

			const locEl = c.querySelector("[class*='location'], [class*='place'], [class*='address']");
			const location = locEl ? locEl.innerText.trim() : "";

			const ratingEl = c.querySelector("[class*='rating'], [class*='score']");
			let rating = ratingEl ? ratingEl.innerText.trim() : "";
			const ratingMatch = rating.match(/(\d+(\.\d+)?)/);
			if (ratingMatch) {
				rating = ratingMatch[1];
			}

			const priceEl = c.querySelector("[class*='price'], [class*='amount'], [class*='total']");
			const price = priceEl ? priceEl.innerText.replace(/\n+/g, ' ').trim() : "";

			const unitEl = c.querySelector("[class*='pitch-type'], [class*='unit-type'], [class*='feature']");
			const unitType = unitEl ? unitEl.innerText.trim() : "Pitch";

			if (price || rating) {
				seen.add(name);
				results.push({
					name: name,
					location: location,
					rating: rating,
					price: price,
					unitType: unitType
				});
			}
		});

		return results;
	})()
	`

	err := chromedp.Run(timeoutCtx,
		chromedp.Navigate(targetURL),
		chromedp.Sleep(6*time.Second),
		chromedp.Evaluate(extractJS, &rawList),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to extract campsites from Pitchup: %w", err)
	}

	limit := criteria.Limit
	if limit <= 0 {
		limit = 15
	}

	priceRegex := regexp.MustCompile(`([\d.,]+)`)

	offers := make([]domain.FlightOffer, 0, len(rawList))
	for _, item := range rawList {
		if len(offers) >= limit {
			break
		}

		currency := "EUR"
		if strings.Contains(item.Price, "£") {
			currency = "GBP"
		} else if strings.Contains(item.Price, "$") {
			currency = "USD"
		}

		m := priceRegex.FindString(item.Price)
		amount := 0.0
		if m != "" {
			cleanStr := strings.ReplaceAll(m, ",", "")
			amount, _ = strconv.ParseFloat(cleanStr, 64)
		}

		ratingInfo := ""
		if item.Rating != "" {
			ratingInfo = fmt.Sprintf(" ⭐ %s/10", item.Rating)
		}

		status := "AVAILABLE"
		if item.UnitType != "" {
			status = item.UnitType
		}

		var depT *time.Time
		if t, err := time.Parse("2006-01-02", arriveDate); err == nil {
			depT = &t
		}

		offers = append(offers, domain.FlightOffer{
			TransportType:    domain.TransportTypeHotel,
			Airline:          "Pitchup",
			FlightNumber:     fmt.Sprintf("⛺ %s%s", item.Name, ratingInfo),
			DepartureStation: item.Location,
			ArrivalStation:   countryTitle,
			DepartureTime:    depT,
			DepartureRaw:     fmt.Sprintf("%s -> %s", arriveDate, departDate),
			Price: domain.Price{
				Amount:   amount,
				Currency: currency,
			},
			IsAvailable: true,
			Status:      status,
		})
	}

	return offers, nil
}
