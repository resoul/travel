package trip

import (
	"context"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/resoul/travel/internal/domain"
)

const (
	defaultBaseURL   = "https://www.trip.com"
	defaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
)

var _ domain.TripProvider = (*Client)(nil)

// Client handles communication with Trip.com via direct HTTP SSR requests.
type Client struct {
	http *http.Client
}

// NewClient creates a new Trip.com client.
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

// SearchHotels searches hotel accommodations on Trip.com via direct Server-Side Rendered (SSR) HTTP requests.
func (c *Client) SearchHotels(ctx context.Context, criteria domain.HotelSearchCriteria) ([]domain.HotelOffer, error) {
	cityID := criteria.CityID
	if cityID == "" {
		cityID = "40795" // Default Barcelona
	}

	cityName := criteria.CityName
	if cityName == "" {
		cityName = "Barcelona"
	}

	checkIn := criteria.CheckIn
	if checkIn == "" {
		checkIn = time.Now().AddDate(0, 1, 0).Format("2006-01-02")
	}

	checkOut := criteria.CheckOut
	if checkOut == "" {
		inDate, err := time.Parse("2006-01-02", checkIn)
		if err == nil {
			checkOut = inDate.AddDate(0, 0, 3).Format("2006-01-02")
		} else {
			checkOut = time.Now().AddDate(0, 1, 3).Format("2006-01-02")
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
		currency = "USD"
	}

	params := url.Values{}
	params.Set("cityId", cityID)
	params.Set("cityName", cityName)
	params.Set("destName", cityName)
	params.Set("searchWord", cityName)
	params.Set("searchType", "CT")
	params.Set("checkin", checkIn)
	params.Set("checkout", checkOut)
	params.Set("crn", strconv.Itoa(rooms))
	params.Set("adult", strconv.Itoa(adults))
	params.Set("children", strconv.Itoa(criteria.Children))
	params.Set("curr", currency)
	params.Set("locale", "en-XX")
	params.Set("old", "1")

	searchURL := fmt.Sprintf("%s/hotels/list?%s", defaultBaseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", defaultUserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request to %s failed: %w", searchURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Trip.com error: status %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	htmlContent := string(bodyBytes)

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

	// Split by hotel item blocks
	itemSplitter := regexp.MustCompile(`data-offline-hotelId="(\d+)"`)
	segments := itemSplitter.Split(htmlContent, -1)
	matches := itemSplitter.FindAllStringSubmatch(htmlContent, -1)

	var offers []domain.HotelOffer

	nameRe := regexp.MustCompile(`<span class="hotelName">([^<]+)</span>`)
	scoreRe := regexp.MustCompile(`<span class="score"[^>]*>([^<]+)</span>`)
	reviewRe := regexp.MustCompile(`<span class="comment-num">([^<]+)</span>`)
	posRe := regexp.MustCompile(`<span class="position-desc">([^<]+)</span>`)
	roomRe := regexp.MustCompile(`<div class="room-name">([^<]+)</div>`)
	totalPriceRe := regexp.MustCompile(`<span class="price-highlight">([^<]+)</span>`)
	nightPriceRe := regexp.MustCompile(`<span class="sale"[^>]*>([^<]+)</span>`)

	for i := 0; i < len(matches) && i+1 < len(segments); i++ {
		hotelID := matches[i][1]
		block := segments[i+1]

		nameMatch := nameRe.FindStringSubmatch(block)
		if len(nameMatch) < 2 {
			continue
		}
		hotelName := html.UnescapeString(strings.TrimSpace(nameMatch[1]))

		var rating float64
		if sm := scoreRe.FindStringSubmatch(block); len(sm) >= 2 {
			rating, _ = strconv.ParseFloat(sm[1], 64)
		}

		var reviews int
		if rm := reviewRe.FindStringSubmatch(block); len(rm) >= 2 {
			digits := regexp.MustCompile(`\d+`).FindString(rm[1])
			reviews, _ = strconv.Atoi(digits)
		}

		var location string
		if pm := posRe.FindStringSubmatch(block); len(pm) >= 2 {
			location = html.UnescapeString(strings.TrimSpace(pm[1]))
		}

		var roomType string
		if rm := roomRe.FindStringSubmatch(block); len(rm) >= 2 {
			roomType = html.UnescapeString(strings.TrimSpace(rm[1]))
		}

		var priceAmount float64
		if tpm := totalPriceRe.FindStringSubmatch(block); len(tpm) >= 2 {
			priceAmount = parseMonetaryValue(tpm[1])
		} else if npm := nightPriceRe.FindStringSubmatch(block); len(npm) >= 2 {
			nightVal := parseMonetaryValue(npm[1])
			priceAmount = nightVal * float64(nights)
		}

		hotelURL := fmt.Sprintf("%s/hotels/detail/?hotelId=%s&checkIn=%s&checkOut=%s&adult=%d&curr=%s&locale=en-XX",
			defaultBaseURL, hotelID, checkIn, checkOut, adults, currency)

		offer := domain.HotelOffer{
			ID:          hotelID,
			Name:        hotelName,
			City:        cityName,
			Address:     location,
			Rating:      rating,
			ReviewCount: reviews,
			Price: domain.Price{
				Amount:   priceAmount,
				Currency: currency,
			},
			Nights:   nights,
			RoomType: roomType,
			URL:      hotelURL,
		}

		offers = append(offers, offer)
		if criteria.Limit > 0 && len(offers) >= criteria.Limit {
			break
		}
	}

	return offers, nil
}

// GetHotelDetails retrieves all room options, room sizes, bed configurations, and facilities for a hotel.
func (c *Client) GetHotelDetails(ctx context.Context, hotelID, checkIn, checkOut string, adults, rooms int, currency string) ([]domain.TripRoomOffer, error) {
	if checkIn == "" {
		checkIn = time.Now().AddDate(0, 1, 0).Format("2006-01-02")
	}
	if checkOut == "" {
		inDate, err := time.Parse("2006-01-02", checkIn)
		if err == nil {
			checkOut = inDate.AddDate(0, 0, 3).Format("2006-01-02")
		} else {
			checkOut = time.Now().AddDate(0, 1, 3).Format("2006-01-02")
		}
	}
	if adults <= 0 {
		adults = 2
	}
	if rooms <= 0 {
		rooms = 1
	}
	if currency == "" {
		currency = "USD"
	}

	detailURL := fmt.Sprintf("%s/hotels/detail/?hotelId=%s&checkIn=%s&checkOut=%s&adult=%d&children=0&crn=%d&curr=%s&locale=en-XX",
		defaultBaseURL, hotelID, checkIn, checkOut, adults, rooms, currency)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, detailURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", defaultUserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request to %s failed: %w", detailURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Trip.com hotel detail error: status %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	raw := string(bodyBytes)
	// Unescape JSON quotes in next_f payload
	unescaped := strings.ReplaceAll(raw, `\"`, `"`)
	unescaped = strings.ReplaceAll(unescaped, `\\`, `\`)

	roomRe := regexp.MustCompile(`"id":\s*(\d+),\s*"name":\s*"([^"]+)",\s*"person":\s*(\d+)`)
	areaRe := regexp.MustCompile(`"areaInfo":\{"icon":"[^"]*","title":"([^"]+)"\}`)
	bedRe := regexp.MustCompile(`"bedInfo":\{[^}]*"title":"([^"]+)"`)
	winRe := regexp.MustCompile(`"windowInfo":\{[^}]*"title":"([^"]+)"`)
	smokeRe := regexp.MustCompile(`"smokeInfo":\{[^}]*"title":"([^"]+)"`)
	picRe := regexp.MustCompile(`"pictureInfo":\[\{"url":"([^"]+)"`)
	facRe := regexp.MustCompile(`"baseFacilityInfo":\[(.*?)\]`)
	facTitleRe := regexp.MustCompile(`"title":"([^"]+)"`)

	matches := roomRe.FindAllStringSubmatchIndex(unescaped, -1)
	var roomOffers []domain.TripRoomOffer

	seenIDs := make(map[string]bool)

	for _, loc := range matches {
		sub := unescaped[loc[0]:]
		if len(sub) > 1500 {
			sub = sub[:1500]
		}

		parts := roomRe.FindStringSubmatch(unescaped[loc[0]:loc[1]])
		if len(parts) < 4 {
			continue
		}

		rID := parts[1]
		if seenIDs[rID] {
			continue
		}
		seenIDs[rID] = true

		rName := html.UnescapeString(parts[2])
		guests, _ := strconv.Atoi(parts[3])

		area := "N/A"
		if am := areaRe.FindStringSubmatch(sub); len(am) >= 2 {
			area = am[1]
		}

		bed := "N/A"
		if bm := bedRe.FindStringSubmatch(sub); len(bm) >= 2 {
			bed = bm[1]
		}

		hasWindow := "N/A"
		if wm := winRe.FindStringSubmatch(sub); len(wm) >= 2 {
			hasWindow = wm[1]
		}

		smoking := "N/A"
		if sm := smokeRe.FindStringSubmatch(sub); len(sm) >= 2 {
			smoking = sm[1]
		}

		imageURL := ""
		if pm := picRe.FindStringSubmatch(sub); len(pm) >= 2 {
			imageURL = pm[1]
		}

		var amenities []string
		if fm := facRe.FindStringSubmatch(sub); len(fm) >= 2 {
			for _, m := range facTitleRe.FindAllStringSubmatch(fm[1], -1) {
				if len(m) >= 2 {
					amenities = append(amenities, html.UnescapeString(m[1]))
				}
			}
		}

		roomOffers = append(roomOffers, domain.TripRoomOffer{
			ID:        rID,
			Name:      rName,
			Area:      area,
			Beds:      bed,
			Guests:    guests,
			HasWindow: hasWindow,
			Smoking:   smoking,
			Amenities: amenities,
			ImageURL:  imageURL,
		})
	}

	return roomOffers, nil
}

func parseMonetaryValue(str string) float64 {
	re := regexp.MustCompile(`[0-9]+(?:\.[0-9]+)?`)
	m := re.FindString(str)
	if m != "" {
		val, _ := strconv.ParseFloat(m, 64)
		return val
	}
	return 0
}
