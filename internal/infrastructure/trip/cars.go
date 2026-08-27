package trip

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/resoul/travel/internal/domain"
)

// SearchCars performs a live rental car search on Trip.com via Chromedp headless browser automation.
func (c *Client) SearchCars(ctx context.Context, criteria domain.CarHireCriteria) ([]domain.CarHireOffer, error) {
	countryID := criteria.CountryID
	if countryID == "" {
		countryID = "63" // Default Romania
	}

	pickupCityID := criteria.PickupCityID
	if pickupCityID == "" {
		pickupCityID = "39050" // Default Otopeni
	}

	pickupCityName := criteria.PickupCityName
	if pickupCityName == "" {
		pickupCityName = "Otopeni"
	}

	pickupCode := criteria.PickupCode
	if pickupCode == "" {
		pickupCode = "OTP"
	}

	pickupAddress := criteria.PickupAddress
	if pickupAddress == "" {
		pickupAddress = fmt.Sprintf("Bucharest Henri Coandă International Airport (%s)", pickupCode)
	}

	countryName := criteria.CountryName
	if countryName == "" {
		countryName = "Romania"
	}

	pickupDate := criteria.PickupDate
	if pickupDate == "" {
		pickupDate = time.Now().AddDate(0, 0, 14).Format("2006/01/02 10:00")
	} else if strings.Contains(pickupDate, "-") && !strings.Contains(pickupDate, "/") {
		pickupDate = strings.ReplaceAll(pickupDate, "-", "/")
		if !strings.Contains(pickupDate, ":") {
			pickupDate += " 10:00"
		}
	}

	returnDate := criteria.ReturnDate
	if returnDate == "" {
		pTime, err := time.Parse("2006/01/02 15:04", pickupDate)
		if err == nil {
			returnDate = pTime.AddDate(0, 0, 3).Format("2006/01/02 15:04")
		} else {
			returnDate = time.Now().AddDate(0, 0, 17).Format("2006/01/02 10:00")
		}
	} else if strings.Contains(returnDate, "-") && !strings.Contains(returnDate, "/") {
		returnDate = strings.ReplaceAll(returnDate, "-", "/")
		if !strings.Contains(returnDate, ":") {
			returnDate += " 10:00"
		}
	}

	returnCityID := criteria.ReturnCityID
	if returnCityID == "" {
		returnCityID = pickupCityID
	}

	returnCityName := criteria.ReturnCityName
	if returnCityName == "" {
		returnCityName = pickupCityName
	}

	returnCode := criteria.ReturnCode
	if returnCode == "" {
		returnCode = pickupCode
	}

	returnAddress := criteria.ReturnAddress
	if returnAddress == "" {
		returnAddress = pickupAddress
	}

	driverAge := criteria.DriverAge
	if driverAge == "" {
		driverAge = "30-60"
	}

	currency := criteria.Currency
	if currency == "" {
		currency = "USD"
	}

	params := url.Values{}
	params.Set("scountry", countryID)
	params.Set("locale", "en-XX")
	params.Set("curr", currency)
	params.Set("fromPage", "Home")
	params.Set("pcity", pickupCityID)
	params.Set("pcityname", pickupCityName)
	params.Set("ptime", pickupDate)
	params.Set("pcode", pickupCode)
	params.Set("paddress", pickupAddress)
	params.Set("ptype", "1")
	params.Set("pcountryname", countryName)
	params.Set("rcity", returnCityID)
	params.Set("rcityname", returnCityName)
	params.Set("rtime", returnDate)
	params.Set("rcode", returnCode)
	params.Set("raddress", returnAddress)
	params.Set("rtype", "1")
	params.Set("rcountryname", countryName)
	params.Set("age", driverAge)
	params.Set("channelid", "14409")

	searchURL := fmt.Sprintf("%s/carhire/online/list?%s", defaultBaseURL, params.Encode())

	// Chromedp headless allocator options
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

	browserCtx, cancelTimeout := context.WithTimeout(browserCtx, 35*time.Second)
	defer cancelTimeout()

	// JS evaluation script to extract live rental cars from Trip.com Car Hire DOM
	extractJS := `
	(() => {
		const cards = document.querySelectorAll("[class*='car-card'], [class*='vehicle-card'], [class*='vehicle-item'], [class*='car_card'], [class*='carCard'], [class*='card-item'], [class*='carItem'], .car-item");
		const results = [];

		cards.forEach(card => {
			const modelEl = card.querySelector("[class*='car-name'], [class*='carName'], [class*='vehicle-name'], h3, .name");
			let model = modelEl ? modelEl.innerText.replace(/\s+/g, ' ').trim() : "";
			if (!model) return;

			const typeEl = card.querySelector("[class*='car-type'], [class*='carType'], [class*='vehicle-type'], .type");
			let category = typeEl ? typeEl.innerText.replace(/\s+/g, ' ').trim() : "Standard";

			const supplierEl = card.querySelector("[class*='supplier-name'], [class*='vendor-name'], [class*='supplierName'], [class*='vendorName']");
			let supplier = supplierEl ? supplierEl.innerText.trim() : "";
			if (!supplier) {
				const logoImg = card.querySelector("img[alt*='logo'], img[class*='supplier'], img[class*='vendor']");
				if (logoImg) {
					supplier = logoImg.getAttribute("alt") || "";
					if (supplier.toLowerCase().includes("logo")) {
						supplier = supplier.replace(/logo/i, "").trim();
					}
				}
			}
			if (!supplier) supplier = "Rental Partner";

			const priceTotalEl = card.querySelector("[class*='total-price'], [class*='totalPrice'], [class*='price-num'], [class*='priceNum'], [class*='price_value']");
			let priceNum = 0;
			if (priceTotalEl) {
				const m = priceTotalEl.innerText.match(/(\d+(?:\.\d+)?)/);
				if (m) {
					priceNum = parseFloat(m[1]) || 0;
				}
			}

			const priceDayEl = card.querySelector("[class*='day-price'], [class*='price-day'], [class*='priceDay'], [class*='average-price']");
			let priceDayNum = 0;
			if (priceDayEl) {
				const m = priceDayEl.innerText.match(/(\d+(?:\.\d+)?)/);
				if (m) {
					priceDayNum = parseFloat(m[1]) || 0;
				}
			}

			const specsEl = card.querySelectorAll("[class*='spec'], [class*='feature'], [class*='tag'], [class*='attribute'], [class*='config-item']");
			const features = [];
			let transmission = "Manual";
			let seats = 5;
			let doors = 4;
			let bags = 2;

			specsEl.forEach(s => {
				const text = s.innerText.replace(/\s+/g, ' ').trim();
				if (!text || text.length > 50) return;
				if (!features.includes(text)) {
					features.push(text);
				}

				const lower = text.toLowerCase();
				if (lower.includes("auto")) transmission = "Automatic";
				if (lower.includes("manual")) transmission = "Manual";
				if (lower.includes("seat")) {
					const num = parseInt(text);
					if (num > 0) seats = num;
				}
				if (lower.includes("door")) {
					const num = parseInt(text);
					if (num > 0) doors = num;
				}
				if (lower.includes("bag")) {
					const num = parseInt(text);
					if (num > 0) bags = num;
				}
			});

			const imgEl = card.querySelector("img[class*='car'], img[class*='vehicle'], img");
			const imageUrl = imgEl ? (imgEl.getAttribute("src") || imgEl.getAttribute("data-src") || "") : "";

			const bookLinkEl = card.querySelector("a, button[class*='book'], [class*='btn']");
			const bookUrl = bookLinkEl ? (bookLinkEl.getAttribute("href") || "") : "";

			results.push({
				model: model,
				category: category,
				supplier: supplier,
				transmission: transmission,
				seats: seats,
				doors: doors,
				bags: bags,
				price_total: priceNum,
				price_day: priceDayNum,
				features: features,
				image_url: imageUrl,
				booking_url: bookUrl
			});
		});

		return results;
	})()
	`

	type rawCarDTO struct {
		Model        string   `json:"model"`
		Category     string   `json:"category"`
		Supplier     string   `json:"supplier"`
		Transmission string   `json:"transmission"`
		Seats        int      `json:"seats"`
		Doors        int      `json:"doors"`
		Bags         int      `json:"bags"`
		PriceTotal   float64  `json:"price_total"`
		PriceDay     float64  `json:"price_day"`
		Features     []string `json:"features"`
		ImageURL     string   `json:"image_url"`
		BookingURL   string   `json:"booking_url"`
	}

	var rawResults []rawCarDTO
	err := chromedp.Run(browserCtx,
		chromedp.Navigate(searchURL),
		chromedp.Sleep(4*time.Second),
		chromedp.Evaluate(`window.scrollBy(0, 500);`, nil),
		chromedp.Sleep(3*time.Second),
		chromedp.Evaluate(extractJS, &rawResults),
	)
	if err != nil {
		return nil, fmt.Errorf("Chromedp search failed for %s: %w", searchURL, err)
	}

	offers := make([]domain.CarHireOffer, 0, len(rawResults))
	for i, raw := range rawResults {
		if raw.Model == "" {
			continue
		}

		offer := domain.CarHireOffer{
			ID:           strconv.Itoa(i + 1),
			Model:        raw.Model,
			Category:     raw.Category,
			Transmission: raw.Transmission,
			Seats:        raw.Seats,
			Doors:        raw.Doors,
			Bags:         raw.Bags,
			Supplier:     raw.Supplier,
			PricePerDay: domain.Price{
				Amount:   raw.PriceDay,
				Currency: currency,
			},
			TotalPrice: domain.Price{
				Amount:   raw.PriceTotal,
				Currency: currency,
			},
			PickupDate:  pickupDate,
			ReturnDate:  returnDate,
			PickupPlace: pickupAddress,
			ReturnPlace: returnAddress,
			Features:    raw.Features,
			ImageURL:    raw.ImageURL,
			BookingURL:  raw.BookingURL,
		}

		offers = append(offers, offer)
		if criteria.Limit > 0 && len(offers) >= criteria.Limit {
			break
		}
	}

	return offers, nil
}
