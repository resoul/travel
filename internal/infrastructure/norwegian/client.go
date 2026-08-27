package norwegian

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/resoul/travel/internal/domain"
)

var _ domain.NorwegianProvider = (*Client)(nil)

// Client handles communication with Norwegian Air Shuttle via Chromedp headless automation.
type Client struct {
	http *http.Client
}

// NewClient creates a new Norwegian client.
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

// CalendarItemDTO represents a single day's extracted fare.
type CalendarItemDTO struct {
	Day      int     `json:"day"`
	Price    float64 `json:"price"`
	RawPrice string  `json:"rawPrice"`
}

// GetFareCalendar retrieves Norwegian Air Shuttle daily lowest fares for a specific route, year, and month.
func (c *Client) GetFareCalendar(ctx context.Context, origin, destination string, year, month int, currency string) ([]domain.FlightOffer, error) {
	if origin == "" || destination == "" {
		return nil, fmt.Errorf("origin and destination are required")
	}

	if year <= 0 {
		year = time.Now().Year()
	}
	if month <= 0 || month > 12 {
		month = int(time.Now().Month())
	}
	if currency == "" {
		currency = "EUR"
	}

	monthStr := fmt.Sprintf("%04d%02d", year, month)

	// Format low fare calendar URL
	baseURL := "https://www.norwegian.com/en/low-fare-calendar/"
	params := url.Values{}
	params.Set("D_City", strings.ToUpper(origin))
	params.Set("A_City", strings.ToUpper(destination))
	params.Set("TripType", "1")
	params.Set("D_Day", "01")
	params.Set("D_Month", monthStr)
	params.Set("R_Day", "01")
	params.Set("R_Month", monthStr)
	params.Set("AgreementCodeInv", "0")
	params.Set("currencyCode", strings.ToUpper(currency))

	targetURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

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

	const extractJS = `
	(() => {
		const results = [];
		const seenDays = new Set();
		const cells = document.querySelectorAll("td[data-date], button[data-date], tr td, div.calendar-day");
		cells.forEach(td => {
			const priceEl = td.querySelector(".fare, .price, [class*='price'], .calendar-fare");
			if (!priceEl) return;
			const priceText = priceEl.innerText.trim();
			const rawText = td.innerText.trim();
			const dayMatch = rawText.match(/^(\d{1,2})\.?/);
			if (dayMatch && priceText) {
				const day = parseInt(dayMatch[1], 10);
				if (day >= 1 && day <= 31 && !seenDays.has(day)) {
					const cleanPrice = priceText.replace(/[^\d.,]/g, '').replace(',', '.');
					const amount = parseFloat(cleanPrice);
					if (amount > 0) {
						seenDays.add(day);
						results.push({ day, price: amount, rawPrice: priceText });
					}
				}
			}
		});
		return results;
	})()
	`

	var items []CalendarItemDTO

	err := chromedp.Run(timeoutCtx,
		chromedp.Navigate(targetURL),
		chromedp.Sleep(6*time.Second),
		chromedp.Evaluate(extractJS, &items),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to extract Norwegian fare calendar: %w", err)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Day < items[j].Day
	})

	offers := make([]domain.FlightOffer, 0, len(items))
	for _, item := range items {
		dateStr := fmt.Sprintf("%04d-%02d-%02d", year, month, item.Day)
		t, _ := time.Parse("2006-01-02", dateStr)

		offers = append(offers, domain.FlightOffer{
			TransportType:    domain.TransportTypeFlight,
			Airline:          "Norwegian",
			FlightNumber:     "DY Low-Fare",
			DepartureStation: strings.ToUpper(origin),
			ArrivalStation:   strings.ToUpper(destination),
			DepartureTime:    &t,
			DepartureRaw:     dateStr,
			Price: domain.Price{
				Amount:   item.Price,
				Currency: strings.ToUpper(currency),
			},
			IsAvailable: true,
			Status:      "AVAILABLE",
		})
	}

	return offers, nil
}
