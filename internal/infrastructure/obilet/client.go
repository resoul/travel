package obilet

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

var _ domain.OBiletProvider = (*Client)(nil)

// KnownCityIDs maps popular Turkish and regional cities to their official oBilet Location IDs.
var KnownCityIDs = map[string]int{
	"istanbul":         349,
	"istanbul avrupa":  349,
	"istanbul anadolu": 350,
	"ankara":           356,
	"izmir":            355,
	"antalya":          360,
	"bursa":            357,
	"adana":            358,
	"konya":            361,
	"trabzon":          370,
	"gaziantep":        362,
	"eskisehir":        364,
	"kayseri":          363,
	"cappadocia":       363,
	"kapadokya":        363,
	"denizli":          365,
	"pamukkale":        365,
	"mugla":            367,
	"bodrum":           475,
	"marmaris":         476,
	"fethiye":          477,
	"alanya":           480,
	"mersin":           368,
	"samsun":           369,
	"canakkale":        371,
	"edirne":           373,
	"rize":             374,
	"artvin":           375,
	"batumi":           481,
}

// Client handles communication with oBilet via Chromedp headless automation.
type Client struct {
	http *http.Client
}

// NewClient creates a new oBilet client.
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

// RawJourneyDTO represents a raw bus journey extracted from the DOM.
type RawJourneyDTO struct {
	Partner      string `json:"partner"`
	Time         string `json:"time"`
	Duration     string `json:"duration"`
	Price        string `json:"price"`
	Origin       string `json:"origin"`
	Dest         string `json:"dest"`
	SeatType     string `json:"seatType"`
	FeatureBadge string `json:"featureBadge"`
}

// SearchBuses searches available intercity and international bus journeys on oBilet.
func (c *Client) SearchBuses(ctx context.Context, criteria domain.OBiletSearchCriteria) ([]domain.FlightOffer, error) {
	originID := criteria.OriginID
	if originID == 0 {
		normalizedOrigin := strings.ToLower(strings.TrimSpace(criteria.OriginName))
		if id, found := KnownCityIDs[normalizedOrigin]; found {
			originID = id
		} else {
			// Default to Istanbul Avrupa
			originID = 349
		}
	}

	destID := criteria.DestinationID
	if destID == 0 {
		normalizedDest := strings.ToLower(strings.TrimSpace(criteria.DestinationName))
		if id, found := KnownCityIDs[normalizedDest]; found {
			destID = id
		} else {
			// Default to Ankara
			destID = 356
		}
	}

	depDate := criteria.DepartureDate
	if depDate == "" {
		depDate = time.Now().AddDate(0, 0, 7).Format("2006-01-02")
	}

	targetURL := fmt.Sprintf("https://www.obilet.com/seferler/%d-%d/%s", originID, destID, depDate)

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
		const cards = document.querySelectorAll(".journey, [class*='journey-item'], [class*='journey_item'], .journey-card, tr.journey");
		cards.forEach(c => {
			const img = c.querySelector("img.partner-logo, img.logo, .partner-logo img, [class*='partner'] img, img[alt]");
			const partner = img?.getAttribute("alt") || img?.getAttribute("title") ||
			                c.querySelector(".partner-name, [class*='partner-name'], [class*='partner']")?.innerText || "Bus";
			const time = c.querySelector(".time, .departure, [class*='departure-time'], [class*='time']")?.innerText || "";
			const duration = c.querySelector(".duration, [class*='duration']")?.innerText || "";
			const price = c.querySelector(".amount, .price, [class*='price-container'], [class*='money-unit']")?.innerText || "";
			const origin = c.querySelector(".origin, .from, [class*='origin']")?.innerText || "";
			const dest = c.querySelector(".destination, .to, [class*='destination']")?.innerText || "";
			const seatType = c.querySelector(".seat-type, [class*='seat'], [class*='feature']")?.innerText || "";

			if (time || price) {
				results.push({
					partner: partner.trim(),
					time: time.trim(),
					duration: duration.trim(),
					price: price.trim(),
					origin: origin.trim(),
					dest: dest.trim(),
					seatType: seatType.trim()
				});
			}
		});
		return results;
	})()
	`

	var rawJourneys []RawJourneyDTO

	err := chromedp.Run(timeoutCtx,
		chromedp.Navigate(targetURL),
		chromedp.Sleep(6*time.Second),
		chromedp.Evaluate(extractJS, &rawJourneys),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to extract oBilet journeys: %w", err)
	}

	limit := criteria.Limit
	if limit <= 0 {
		limit = 20
	}

	priceRegex := regexp.MustCompile(`([\d.,]+)`)
	timeRegex := regexp.MustCompile(`(\d{2}:\d{2})`)

	offers := make([]domain.FlightOffer, 0, len(rawJourneys))
	for _, j := range rawJourneys {
		if len(offers) >= limit {
			break
		}

		timeMatch := timeRegex.FindString(j.Time)
		if timeMatch == "" {
			timeMatch = "00:00"
		}

		priceStr := ""
		m := priceRegex.FindString(j.Price)
		if m != "" {
			priceStr = strings.ReplaceAll(m, ".", "")
			priceStr = strings.ReplaceAll(priceStr, ",", ".")
		}
		amount, _ := strconv.ParseFloat(priceStr, 64)

		fullDepStr := fmt.Sprintf("%s %s", depDate, timeMatch)
		var depT *time.Time
		if t, err := time.Parse("2006-01-02 15:04", fullDepStr); err == nil {
			depT = &t
		}

		partnerName := j.Partner
		if partnerName == "" || partnerName == "Bus" {
			partnerName = "Intercity Bus"
		}

		seatClass := ""
		if strings.Contains(j.SeatType, "2+1") {
			seatClass = "2+1"
		} else if strings.Contains(j.SeatType, "2+2") {
			seatClass = "2+2"
		}

		flightNum := partnerName
		if seatClass != "" {
			flightNum += fmt.Sprintf(" (%s)", seatClass)
		}

		dur := strings.Trim(j.Duration, "()* ")
		dur = strings.ReplaceAll(dur, "Saat", "h")
		dur = strings.ReplaceAll(dur, "Dakika", "m")
		dur = strings.ReplaceAll(dur, "dakika", "m")

		originName := j.Origin
		if originName == "" {
			originName = criteria.OriginName
			if originName == "" {
				originName = fmt.Sprintf("Location #%d", originID)
			}
		}

		destName := j.Dest
		if destName == "" {
			destName = criteria.DestinationName
			if destName == "" {
				destName = fmt.Sprintf("Location #%d", destID)
			}
		}

		offers = append(offers, domain.FlightOffer{
			TransportType:    domain.TransportTypeBus,
			Airline:          partnerName,
			FlightNumber:     flightNum,
			DepartureStation: originName,
			ArrivalStation:   destName,
			DepartureTime:    depT,
			DepartureRaw:     fullDepStr,
			Duration:         dur,
			Price: domain.Price{
				Amount:   amount,
				Currency: "TRY",
			},
			IsAvailable: true,
			Status:      "AVAILABLE",
		})
	}

	return offers, nil
}
