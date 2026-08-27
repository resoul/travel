package hipcamp

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

var _ domain.HipcampProvider = (*Client)(nil)

// Client handles Hipcamp glamping and outdoor spots search via Chromedp headless automation.
type Client struct {
	http *http.Client
}

// NewClient creates a new Hipcamp client.
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

type rawSpot struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Rating   string `json:"rating"`
	Price    string `json:"price"`
	Type     string `json:"type"`
}

// SearchSpots searches outdoor campsites, glamping spots, and private land stays on Hipcamp.
func (c *Client) SearchSpots(ctx context.Context, criteria domain.HipcampSearchCriteria) ([]domain.FlightOffer, error) {
	country := strings.ToLower(strings.TrimSpace(criteria.Country))
	if country == "" {
		country = "united-states"
	}
	country = strings.ReplaceAll(country, " ", "-")

	region := strings.ToLower(strings.TrimSpace(criteria.Region))
	if region == "" {
		region = "california"
	}
	region = strings.ReplaceAll(region, " ", "-")

	targetURL := fmt.Sprintf("https://www.hipcamp.com/en-US/d/%s/%s/camping/all", country, region)

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

	var rawList []rawSpot

	const extractJS = `
	(() => {
		const results = [];
		const anchors = document.querySelectorAll("a[href*='/land/'], a[href*='/discover/']");
		const seen = new Set();

		anchors.forEach(a => {
			const text = a.innerText.trim();
			const lines = text.split("\n").map(l => l.trim()).filter(l => l.length > 0);
			if (lines.length === 0) return;

			const name = lines[0];
			if (seen.has(name) || name.length < 3) return;
			seen.add(name);

			let location = "";
			let rating = "";
			let price = "";
			let stayType = "Outdoor Stay";

			for (const l of lines) {
				if (l.includes("%") || l.includes("★") || l.includes("rating")) {
					rating = l;
				} else if (l.includes("$") || l.includes("£") || l.includes("€") || l.includes("/night")) {
					price = l;
				} else if (l.includes("Camp") || l.includes("Glamping") || l.includes("RV") || l.includes("Cabin")) {
					stayType = l;
				}
			}

			results.push({
				name: name,
				location: location,
				rating: rating,
				price: price,
				type: stayType
			});
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
		return nil, fmt.Errorf("failed to extract spots from Hipcamp: %w", err)
	}

	limit := criteria.Limit
	if limit <= 0 {
		limit = 10
	}

	priceRegex := regexp.MustCompile(`([\d.,]+)`)

	offers := make([]domain.FlightOffer, 0, len(rawList))
	for _, item := range rawList {
		if len(offers) >= limit {
			break
		}

		currency := "USD"
		if strings.Contains(item.Price, "£") {
			currency = "GBP"
		} else if strings.Contains(item.Price, "€") {
			currency = "EUR"
		}

		m := priceRegex.FindString(item.Price)
		amount := 0.0
		if m != "" {
			cleanStr := strings.ReplaceAll(m, ",", "")
			amount, _ = strconv.ParseFloat(cleanStr, 64)
		}

		ratingInfo := ""
		if item.Rating != "" {
			ratingInfo = fmt.Sprintf(" (%s)", item.Rating)
		}

		now := time.Now()

		offers = append(offers, domain.FlightOffer{
			TransportType:    domain.TransportTypeHotel,
			Airline:          "Hipcamp",
			FlightNumber:     fmt.Sprintf("🌲 %s%s", item.Name, ratingInfo),
			DepartureStation: strings.Title(strings.ReplaceAll(region, "-", " ")),
			ArrivalStation:   strings.Title(strings.ReplaceAll(country, "-", " ")),
			DepartureTime:    &now,
			DepartureRaw:     item.Type,
			Price: domain.Price{
				Amount:   amount,
				Currency: currency,
			},
			IsAvailable: true,
			Status:      "AVAILABLE",
		})
	}

	return offers, nil
}
