package transavia

// CalendarFaresResponseDTO represents the JSON returned by /start/api/calendar-fares.
type CalendarFaresResponseDTO struct {
	Data []FareItemDTO `json:"data"`
}

// FareItemDTO represents daily price and tariff type on Transavia.
type FareItemDTO struct {
	Date  string  `json:"date"`
	Price float64 `json:"price"`
	Type  string  `json:"type"` // e.g. "lowFare", "regularFare"
}
