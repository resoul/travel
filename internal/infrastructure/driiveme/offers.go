package driiveme

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/resoul/travel/internal/domain"
)

var (
	poisTransportsRegex = regexp.MustCompile(`POIS_TRANSPORTS\s*=\s*(\[[\s\S]*?\]);`)
	cardBlockRegex      = regexp.MustCompile(`(?s)<div class="block-trajet filtr-item">([\s\S]*?)(?:<div class="block-trajet filtr-item">|</div>\s*</div>\s*<div class="pagination|$)`)
	detailURLRegex      = regexp.MustCompile(`href="([^"]*details-(\d+)\.html)"`)
	depCityRegex        = regexp.MustCompile(`name="filterDepartureValue"\s+value="([^"]*)"`)
	arrCityRegex        = regexp.MustCompile(`name="filterDestinationValue"\s+value="([^"]*)"`)
	startDateRegex      = regexp.MustCompile(`name="filterDateBeginAvailabilityValue"\s+value="([^"]*)"`)
	endDateRegex        = regexp.MustCompile(`name="filterDateEndAvailabilityValue"\s+value="([^"]*)"`)
	vehicleValRegex     = regexp.MustCompile(`name="filterVehiculeValue"\s+value="([^"]*)"`)
	seatsRegex          = regexp.MustCompile(`title="(\d+)\s+seats"`)
	hoursRegex          = regexp.MustCompile(`title="(\d+)\s+rental hours"`)
	modelDetailRegex    = regexp.MustCompile(`(?i)Model[\s\S]*?<span class="large">([^<]+)</span>`)
	transDetailRegex    = regexp.MustCompile(`(?i)Transmission[\s\S]*?<span class="large">([^<]+)</span>`)
	depositDetailRegex  = regexp.MustCompile(`(?i)pre-authorisation of ([^<]+) will be`)
	distIncludedRegex   = regexp.MustCompile(`(?i)<span class="[^"]*white">([^<]+(mi|km) included)</span>`)
	timeIncludedRegex   = regexp.MustCompile(`(?i)<span class="[^"]*white">([^<]+h included)</span>`)
	offeredByRegex      = regexp.MustCompile(`(?i)offered by ([^"]+)"`)
	insuranceRegex      = regexp.MustCompile(`(?i)insured by ([^"]+)"`)
)

// GetAvailabilities retrieves booking timeslots for a specific transport ID.
func (c *Client) GetAvailabilities(ctx context.Context, transportID string) ([]string, error) {
	if transportID == "" {
		return nil, fmt.Errorf("transport ID is required")
	}

	endpoint := fmt.Sprintf("/%s/transport/get-availabilities/%s?forBooking=true", c.locale, transportID)
	body, err := c.get(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch availabilities for transport %s: %w", transportID, err)
	}

	var resp availabilitiesResponseDTO
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to decode availabilities response: %w", err)
	}

	return resp.Availabilities, nil
}

// GetOffers searches and retrieves 1-euro relocation trips from DriiveMe.
func (c *Client) GetOffers(ctx context.Context, criteria domain.FlightSearchCriteria) ([]domain.FlightOffer, error) {
	depCityID := ""
	if criteria.Origin != "" {
		if _, err := strconv.Atoi(criteria.Origin); err == nil {
			depCityID = criteria.Origin
		} else {
			suggestions, err := c.SearchCities(ctx, criteria.Origin)
			if err == nil && len(suggestions) > 0 {
				depCityID = strconv.Itoa(suggestions[0].ID)
			}
		}
	}

	arrCityID := ""
	if criteria.Destination != "" {
		if _, err := strconv.Atoi(criteria.Destination); err == nil {
			arrCityID = criteria.Destination
		} else {
			suggestions, err := c.SearchCities(ctx, criteria.Destination)
			if err == nil && len(suggestions) > 0 {
				arrCityID = strconv.Itoa(suggestions[0].ID)
			}
		}
	}

	params := url.Values{}
	params.Set("departureCityId", depCityID)
	params.Set("arrivalCityId", arrCityID)
	params.Set("minDate", criteria.DepartureDate)
	params.Set("vehicleType", "0")
	params.Set("driverType", "")
	params.Set("page", "1")

	endpoint := fmt.Sprintf("/%s/component/transport-list?%s", c.locale, params.Encode())
	htmlBytes, err := c.get(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch DriiveMe transport list: %w", err)
	}

	cards := c.parseTransportCards(string(htmlBytes))
	if len(cards) == 0 {
		return nil, nil
	}

	offers := make([]domain.FlightOffer, 0, len(cards))
	for _, card := range cards {
		// If authenticated or detail page is available, enrich with detailed vehicle attributes
		if card.Slug != "" {
			c.enrichCardWithDetails(ctx, &card)
		}

		depTime, depFormatted := parseOfferDate(card.AvailabilityStart)
		arrTime, arrFormatted := parseOfferDate(card.AvailabilityEnd)

		vehicleDisplay := card.VehicleCategory
		if card.VehicleModel != "" {
			vehicleDisplay = card.VehicleModel
			if card.VehicleCategory != "" && !strings.EqualFold(card.VehicleCategory, card.VehicleModel) {
				vehicleDisplay = fmt.Sprintf("%s (%s)", card.VehicleModel, card.VehicleCategory)
			}
		}
		if card.Transmission != "" {
			vehicleDisplay = fmt.Sprintf("%s - %s", vehicleDisplay, card.Transmission)
		}
		if vehicleDisplay == "" {
			vehicleDisplay = "Car / Van"
		}

		durationDesc := fmt.Sprintf("%dh included", card.RentalHours)
		if card.IncludedMiles > 0 {
			durationDesc = fmt.Sprintf("%s, %d mi", durationDesc, card.IncludedMiles)
		}
		if card.Deposit != "" {
			durationDesc = fmt.Sprintf("%s (deposit: %s)", durationDesc, card.Deposit)
		}

		priceAmount := card.Price
		if priceAmount <= 0 {
			priceAmount = 1.0
		}
		currency := card.Currency
		if currency == "" {
			currency = "EUR"
		}

		offer := domain.FlightOffer{
			TransportType:    domain.TransportTypeCar,
			Airline:          "DriiveMe",
			FlightNumber:     vehicleDisplay,
			DepartureStation: card.DepartureCity,
			ArrivalStation:   card.ArrivalCity,
			DepartureTime:    depTime,
			ArrivalTime:      arrTime,
			DepartureRaw:     depFormatted,
			ArrivalRaw:       arrFormatted,
			Duration:         durationDesc,
			Price: domain.Price{
				Amount:   priceAmount,
				Currency: currency,
			},
			SeatsLeft:   card.Seats,
			IsAvailable: true,
			Status:      "available",
		}

		offers = append(offers, offer)
	}

	sort.Slice(offers, func(i, j int) bool {
		return offers[i].DepartureRaw < offers[j].DepartureRaw
	})

	return offers, nil
}

// parseTransportCards extracts trip card data from transport-list HTML response.
func (c *Client) parseTransportCards(html string) []TransportCardDTO {
	var cards []TransportCardDTO

	matches := cardBlockRegex.FindAllStringSubmatch(html, -1)
	for _, m := range matches {
		block := m[1]

		var card TransportCardDTO
		if dm := detailURLRegex.FindStringSubmatch(block); len(dm) > 2 {
			card.Slug = dm[1]
			card.ID = dm[2]
		}
		if dep := depCityRegex.FindStringSubmatch(block); len(dep) > 1 {
			card.DepartureCity = dep[1]
		}
		if arr := arrCityRegex.FindStringSubmatch(block); len(arr) > 1 {
			card.ArrivalCity = arr[1]
		}
		if sDate := startDateRegex.FindStringSubmatch(block); len(sDate) > 1 {
			card.AvailabilityStart = sDate[1]
		}
		if eDate := endDateRegex.FindStringSubmatch(block); len(eDate) > 1 {
			card.AvailabilityEnd = eDate[1]
		}
		if veh := vehicleValRegex.FindStringSubmatch(block); len(veh) > 1 {
			card.VehicleCategory = veh[1]
		}
		if st := seatsRegex.FindStringSubmatch(block); len(st) > 1 {
			card.Seats, _ = strconv.Atoi(st[1])
		}
		if hr := hoursRegex.FindStringSubmatch(block); len(hr) > 1 {
			card.RentalHours, _ = strconv.Atoi(hr[1])
		}
		if card.RentalHours == 0 {
			card.RentalHours = 24
		}
		card.Price = 1.0
		card.Currency = "EUR"

		if card.DepartureCity != "" || card.ArrivalCity != "" || card.ID != "" {
			cards = append(cards, card)
		}
	}

	// Fallback to POIS_TRANSPORTS if HTML card block regex didn't catch everything
	if len(cards) == 0 {
		if pm := poisTransportsRegex.FindStringSubmatch(html); len(pm) > 1 {
			var pois []POITransportDTO
			if err := json.Unmarshal([]byte(pm[1]), &pois); err == nil {
				for _, poi := range pois {
					dm := detailURLRegex.FindStringSubmatch(poi.A.Tooltip)
					var card TransportCardDTO
					if len(dm) > 2 {
						card.Slug = dm[1]
						card.ID = dm[2]
					}
					card.Price = 1.0
					card.Currency = "EUR"
					card.RentalHours = 24
					if card.ID != "" {
						cards = append(cards, card)
					}
				}
			}
		}
	}

	return cards
}

// enrichCardWithDetails fetches the transport details page and updates vehicle model, deposit, mileage.
func (c *Client) enrichCardWithDetails(ctx context.Context, card *TransportCardDTO) {
	if card.Slug == "" {
		return
	}

	detailsHTML, err := c.get(ctx, card.Slug)
	if err != nil {
		return
	}
	content := string(detailsHTML)

	if m := modelDetailRegex.FindStringSubmatch(content); len(m) > 1 {
		card.VehicleModel = strings.TrimSpace(m[1])
	}
	if t := transDetailRegex.FindStringSubmatch(content); len(t) > 1 {
		card.Transmission = strings.TrimSpace(t[1])
	}
	if dep := depositDetailRegex.FindStringSubmatch(content); len(dep) > 1 {
		card.Deposit = strings.TrimSpace(dep[1])
	}
	if off := offeredByRegex.FindStringSubmatch(content); len(off) > 1 {
		card.OfferedBy = strings.TrimSpace(off[1])
	}
	if ins := insuranceRegex.FindStringSubmatch(content); len(ins) > 1 {
		card.Insurance = strings.TrimSpace(ins[1])
	}
	if dist := distIncludedRegex.FindStringSubmatch(content); len(dist) > 1 {
		numStr := regexp.MustCompile(`\d+`).FindString(dist[1])
		card.IncludedMiles, _ = strconv.Atoi(numStr)
	}
}

func parseOfferDate(raw string) (*time.Time, string) {
	if raw == "" {
		return nil, ""
	}
	// Try parsing YYYY-MM-DD
	if t, err := time.Parse("2006-01-02", raw); err == nil {
		return &t, t.Format("2006-01-02")
	}
	// Try parsing DD/MM/YYYY
	if t, err := time.Parse("02/01/2006", raw); err == nil {
		return &t, t.Format("2006-01-02")
	}
	return nil, raw
}
