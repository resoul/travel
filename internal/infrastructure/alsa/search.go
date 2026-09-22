package alsa

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/resoul/travel/internal/domain"
)

// SearchTrips searches for bus journeys using Chromedp browser automation.
func (c *Client) SearchTrips(ctx context.Context, criteria domain.FlightSearchCriteria) ([]domain.FlightOffer, error) {
	originStationID := criteria.Origin
	destStationID := criteria.Destination
	departDate := criteria.DepartureDate

	if originStationID == "" || destStationID == "" || departDate == "" {
		return nil, fmt.Errorf("origin, destination, and departure date are required")
	}

	// Parse departure date (expects format: YYYY-MM-DD)
	depTime, err := time.Parse("2006-01-02", departDate)
	if err != nil {
		return nil, fmt.Errorf("invalid departure date format: %w", err)
	}

	// Format for ALSA URL (DD/MM/YYYY)
	formattedDate := depTime.Format("02/01/2006")

	// Build search URL
	searchURL := fmt.Sprintf(
		"%s/checkout?p_p_id=com_babel_alsa_espania_purchase_web_PurchasePortlet&"+
			"p_p_lifecycle=1&p_p_state=normal&p_p_mode=view&"+
			"_com_babel_alsa_espania_purchase_web_PurchasePortlet_javax.portlet.action=SearchJourneysAction&"+
			"alsaspam=&originStationId=%s&destinationStationId=%s&departureDate=%s&"+
			"seats=1&travelType=OUTWARD&employee=false&passengerType-1=1",
		baseURL, originStationID, destStationID, formattedDate,
	)

	// Set up Chromedp headless browser
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
	)

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(ctx, opts...)
	defer cancelAlloc()

	browserCtx, cancelBrowser := chromedp.NewContext(allocCtx)
	defer cancelBrowser()

	timeoutCtx, cancelTimeout := context.WithTimeout(browserCtx, 40*time.Second)
	defer cancelTimeout()

	// JavaScript to extract journey results from DOM
	const extractJS = `
	(() => {
		const results = [];

		// Try multiple selectors for journey cards
		const cards = document.querySelectorAll(
			"[class*='journey'], [class*='trip'], .trip-item, .journey-item, " +
			"[data-testid*='journey'], [data-testid*='trip'], " +
			"table tbody tr, .result-row"
		);

		cards.forEach(card => {
			const depTimeEl = card.querySelector("[class*='departure'], [class*='time'], .dep-time");
			const arrTimeEl = card.querySelector("[class*='arrival'], .arr-time");
			const priceEl = card.querySelector("[class*='price'], .fare, [class*='cost']");
			const durationEl = card.querySelector("[class*='duration'], [class*='length']");
			const stopsEl = card.querySelector("[class*='stops'], [class*='transfers']");

			const depTime = depTimeEl ? depTimeEl.innerText.trim() : "";
			const arrTime = arrTimeEl ? arrTimeEl.innerText.trim() : "";
			const priceStr = priceEl ? priceEl.innerText.trim() : "";
			const duration = durationEl ? durationEl.innerText.trim() : "";
			const stops = stopsEl ? stopsEl.innerText.trim() : "0";

			// Only add if we have at least departure time and price
			if (depTime && priceStr) {
				results.push({
					departureTime: depTime,
					arrivalTime: arrTime,
					price: priceStr,
					duration: duration,
					stops: stops
				});
			}
		});

		return results;
	})()
	`

	var rawResults []map[string]interface{}

	err = chromedp.Run(timeoutCtx,
		chromedp.Navigate(searchURL),
		chromedp.Sleep(4*time.Second), // Wait for page and JavaScript to load
		chromedp.Evaluate(extractJS, &rawResults),
	)

	if err != nil {
		return nil, fmt.Errorf("Chromedp search failed: %w", err)
	}

	// Convert raw results to domain FlightOffers
	offers := make([]domain.FlightOffer, 0, len(rawResults))
	priceRegex := regexp.MustCompile(`([\d.,]+)`)

	for _, raw := range rawResults {
		depTime := fmt.Sprintf("%v", raw["departureTime"])
		arrTime := fmt.Sprintf("%v", raw["arrivalTime"])
		priceStr := fmt.Sprintf("%v", raw["price"])

		// Extract numeric price
		priceMatch := priceRegex.FindString(priceStr)
		amount := 0.0
		if priceMatch != "" {
			cleanPrice := strings.ReplaceAll(priceMatch, ",", "")
			amount, _ = strconv.ParseFloat(cleanPrice, 64)
		}

		currency := "EUR" // ALSA uses EUR
		if strings.Contains(priceStr, "€") {
			currency = "EUR"
		} else if strings.Contains(priceStr, "£") {
			currency = "GBP"
		}

		offers = append(offers, domain.FlightOffer{
			TransportType:    domain.TransportTypeBus,
			Airline:          "ALSA",
			DepartureStation: originStationID,
			ArrivalStation:   destStationID,
			DepartureRaw:     depTime,
			ArrivalRaw:       arrTime,
			Duration:         fmt.Sprintf("%v", raw["duration"]),
			Price: domain.Price{
				Amount:   amount,
				Currency: currency,
			},
			IsAvailable: true,
			Status:      "available",
		})
	}

	return offers, nil
}
