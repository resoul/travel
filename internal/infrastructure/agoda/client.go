package agoda

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/resoul/travel/internal/domain"
)

const (
	defaultBaseURL   = "https://www.agoda.com"
	defaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
)

var _ domain.AgodaProvider = (*Client)(nil)

// Client handles communication with Agoda via Chromedp headless browser and static CDN endpoints.
type Client struct {
	http *http.Client
}

// NewClient creates a new Agoda client.
func NewClient(transport ...http.RoundTripper) *Client {
	var tr http.RoundTripper
	if len(transport) > 0 && transport[0] != nil {
		tr = transport[0]
	}

	return &Client{
		http: &http.Client{
			Transport: tr,
			Timeout:   30 * time.Second,
		},
	}
}

// GetCountries fetches the static country database from Agoda CDN.
func (c *Client) GetCountries(ctx context.Context, languageID int) ([]domain.Country, error) {
	if languageID <= 0 {
		languageID = 11 // Default Russian language ID
	}

	cdnURL := fmt.Sprintf("https://cdn6.agoda.net/js/static/v2/countries_list_%d_v4.json", languageID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cdnURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create countries request: %w", err)
	}

	req.Header.Set("User-Agent", defaultUserAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request to %s failed: %w", cdnURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Agoda CDN error at %s: status %d", cdnURL, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var items []CountryItemDTO
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, fmt.Errorf("failed to decode countries JSON: %w", err)
	}

	countries := make([]domain.Country, 0, len(items))
	for _, item := range items {
		countries = append(countries, domain.Country{
			Code:     item.ID,
			ISO3Code: item.ISO2,
			Name:     item.Name,
		})
	}

	return countries, nil
}

// SearchHotels performs a live accommodation search on Agoda via Chromedp headless browser.
func (c *Client) SearchHotels(ctx context.Context, criteria domain.HotelSearchCriteria) ([]domain.HotelOffer, error) {
	cityID := criteria.CityID
	if cityID == "" {
		cityID = "19216" // Default Mamaia
	}

	checkIn := criteria.CheckIn
	if checkIn == "" {
		checkIn = time.Now().AddDate(0, 0, 7).Format("2006-01-02")
	}

	checkOut := criteria.CheckOut
	if checkOut == "" {
		inDate, err := time.Parse("2006-01-02", checkIn)
		if err == nil {
			checkOut = inDate.AddDate(0, 0, 2).Format("2006-01-02")
		} else {
			checkOut = time.Now().AddDate(0, 0, 9).Format("2006-01-02")
		}
	}

	rooms := criteria.Rooms
	if rooms <= 0 {
		rooms = 1
	}

	adults := criteria.Adults
	if adults <= 0 {
		adults = 2
	}

	currency := criteria.Currency
	if currency == "" {
		currency = "EUR"
	}

	sortField := criteria.Sort
	if sortField == "" {
		sortField = "priceLowToHigh"
	}

	params := url.Values{}
	params.Set("city", cityID)
	params.Set("checkIn", checkIn)
	params.Set("checkOut", checkOut)
	params.Set("rooms", strconv.Itoa(rooms))
	params.Set("adults", strconv.Itoa(adults))
	params.Set("children", strconv.Itoa(criteria.Children))
	params.Set("priceCur", currency)
	params.Set("currency", currency)
	params.Set("sort", sortField)
	if criteria.CityName != "" {
		params.Set("textToSearch", criteria.CityName)
	}

	searchURL := fmt.Sprintf("%s/ru-ru/search?%s", defaultBaseURL, params.Encode())

	// Set up Chromedp headless options
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.UserAgent(defaultUserAgent),
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("window-size", "1920,1080"),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("mute-audio", true),
	)

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(ctx, opts...)
	defer cancelAlloc()

	browserCtx, cancelBrowser := chromedp.NewContext(allocCtx)
	defer cancelBrowser()

	browserCtx, cancelTimeout := context.WithTimeout(browserCtx, 30*time.Second)
	defer cancelTimeout()

	// JS evaluation script to extract rendered hotel cards from Agoda DOM
	extractJS := `
	(() => {
		const cards = document.querySelectorAll("[data-selenium='hotel-item'], [data-selenium='property-card-container'], [data-element-name='property-card'], ol li[data-hotelid], .PropertyCard");
		const results = [];
		
		cards.forEach(card => {
			const nameEl = card.querySelector("[data-selenium='hotel-name'], [data-selenium='property-card-title'], h3, .PropertyCardTitle");
			const name = nameEl ? nameEl.innerText.trim() : "";
			if (!name) return;

			const priceEl = card.querySelector("[data-selenium='display-price'], [data-selenium='property-price'], [data-element-name='final-price'], .PropertyCardPrice__Value, .price-box, [data-selenium='price-box']");
			let priceStr = priceEl ? priceEl.innerText.trim() : "0";
			let priceNum = parseFloat(priceStr.replace(/[^0-9.]/g, '')) || 0;

			const ratingEl = card.querySelector("[data-selenium='review-score'], .ReviewScore, .PropertyCardRating");
			let rating = ratingEl ? parseFloat(ratingEl.innerText.trim()) || 0 : 0;

			const reviewsEl = card.querySelector("[data-selenium='review-count'], .ReviewCount");
			let reviews = 0;
			if (reviewsEl) {
				const numMatch = reviewsEl.innerText.match(/\d+/);
				if (numMatch) reviews = parseInt(numMatch[0]);
			}

			const addressEl = card.querySelector("[data-selenium='property-card-address'], [data-selenium='hotel-address'], .PropertyCardAddress");
			const address = addressEl ? addressEl.innerText.trim() : "";

			const linkEl = card.querySelector("a[data-selenium='hotel-name'], a[href*='/hotel/'], a[data-element-name='property-card-content']");
			let hotelUrl = linkEl ? linkEl.getAttribute("href") || "" : "";
			if (hotelUrl && hotelUrl.startsWith("/")) {
				hotelUrl = "https://www.agoda.com" + hotelUrl;
			}

			const imgEl = card.querySelector("img[data-selenium='property-card-image'], img");
			const imageUrl = imgEl ? (imgEl.getAttribute("src") || imgEl.getAttribute("data-src") || "") : "";

			const roomEl = card.querySelector("[data-selenium='room-type'], .RoomType");
			const roomType = roomEl ? roomEl.innerText.trim() : "";

			results.push({
				name: name,
				address: address,
				price: priceNum,
				rating: rating,
				reviews: reviews,
				url: hotelUrl,
				image_url: imageUrl,
				room_type: roomType
			});
		});

		return results;
	})()
	`

	var rawResults []HotelCardDTO
	err := chromedp.Run(browserCtx,
		chromedp.Navigate(searchURL),
		chromedp.Sleep(3*time.Second),
		chromedp.Evaluate(`window.scrollBy(0, 600);`, nil),
		chromedp.Sleep(3*time.Second),
		chromedp.Evaluate(extractJS, &rawResults),
	)
	if err != nil {
		return nil, fmt.Errorf("Chromedp search failed for %s: %w", searchURL, err)
	}

	// Calculate nights
	nights := 1
	if inT, err := time.Parse("2006-01-02", checkIn); err == nil {
		if outT, err := time.Parse("2006-01-02", checkOut); err == nil {
			days := int(outT.Sub(inT).Hours() / 24)
			if days > 0 {
				nights = days
			}
		}
	}

	offers := make([]domain.HotelOffer, 0, len(rawResults))
	for _, raw := range rawResults {
		if raw.Name == "" {
			continue
		}

		offer := domain.HotelOffer{
			Name:        raw.Name,
			City:        criteria.CityName,
			Address:     raw.Address,
			Rating:      raw.Rating,
			ReviewCount: raw.Reviews,
			Price: domain.Price{
				Amount:   raw.Price,
				Currency: currency,
			},
			Nights:   nights,
			RoomType: raw.RoomType,
			URL:      raw.URL,
			ImageURL: raw.ImageURL,
		}

		offers = append(offers, offer)
		if criteria.Limit > 0 && len(offers) >= criteria.Limit {
			break
		}
	}

	return offers, nil
}

func extractNumber(s string) float64 {
	re := regexp.MustCompile(`[0-9]+(?:\.[0-9]+)?`)
	m := re.FindString(s)
	if m != "" {
		val, _ := strconv.ParseFloat(m, 64)
		return val
	}
	return 0
}
