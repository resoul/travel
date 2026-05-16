package ryanair

// --- Airports ---

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type City struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type Region struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type Country struct {
	Code               string `json:"code"`
	ISO3Code           string `json:"iso3code"`
	Name               string `json:"name"`
	Currency           string `json:"currency"`
	DefaultAirportCode string `json:"defaultAirportCode"`
	Schengen           bool   `json:"schengen"`
}

type Airport struct {
	Code        string      `json:"code"`
	Name        string      `json:"name"`
	SeoName     string      `json:"seoName"`
	Aliases     []string    `json:"aliases"`
	Base        bool        `json:"base"`
	City        City        `json:"city"`
	Region      Region      `json:"region"`
	Country     Country     `json:"country"`
	Coordinates Coordinates `json:"coordinates"`
	TimeZone    string      `json:"timeZone"`
}

// --- Availability API ---

type FlightRequest struct {
	Origin      string
	Destination string
	DateOut     string // YYYY-MM-DD
	DateIn      string // YYYY-MM-DD, empty for one-way
	Adults      int
	Teens       int
	Children    int
	Infants     int
	RoundTrip   bool
	// flex window around DateOut/DateIn (0–6 days)
	FlexDaysBefore int
	FlexDaysAfter  int
}

type Fare struct {
	Type             string  `json:"type"`
	Amount           float64 `json:"amount"`
	Count            int     `json:"count"`
	HasDiscount      bool    `json:"hasDiscount"`
	PublishedFare    float64 `json:"publishedFare"`
	DiscountAmount   float64 `json:"discountAmount"`
	DiscountInPct    int     `json:"discountInPercent"`
	HasPromoDiscount bool    `json:"hasPromoDiscount"`
	IsPrime          bool    `json:"isPrime"`
	HasBogof         bool    `json:"hasBogof"`
}

type RegularFare struct {
	FareKey string `json:"fareKey"`
	Fares   []Fare `json:"fares"`
}

type Segment struct {
	SegmentNr    int    `json:"segmentNr"`
	Origin       string `json:"origin"`
	Destination  string `json:"destination"`
	FlightNumber string `json:"flightNumber"`
}

type RawFlight struct {
	FaresLeft    int         `json:"faresLeft"`
	FlightKey    string      `json:"flightKey"`
	FlightNumber string      `json:"flightNumber"`
	InfantsLeft  int         `json:"infantsLeft"`
	Duration     string      `json:"duration"`
	IsSSIMLoad   bool        `json:"isSSIMLoad"`
	OperatedBy   string      `json:"operatedBy"`
	RegularFare  RegularFare `json:"regularFare"`
	Segments     []Segment   `json:"segments"`
	Time         []string    `json:"time"`    // local: ["2026-06-22T10:55:00.000", ...]
	TimeUTC      []string    `json:"timeUTC"` // utc:   ["2026-06-22T07:55:00.000Z", ...]
}

type RawDate struct {
	DateOut string      `json:"dateOut"`
	Flights []RawFlight `json:"flights"`
}

type RawTrip struct {
	Origin          string    `json:"origin"`
	OriginName      string    `json:"originName"`
	Destination     string    `json:"destination"`
	DestinationName string    `json:"destinationName"`
	Dates           []RawDate `json:"dates"`
}

type AvailabilityResponse struct {
	Currency string    `json:"currency"`
	Trips    []RawTrip `json:"trips"`
}

// FlightResult is the flat, easy-to-use result returned by Search.
type FlightResult struct {
	DepartureStation string
	ArrivalStation   string
	FlightNumber     string
	DepartureLocal   string // e.g. "2026-06-22T10:55:00.000"
	ArrivalLocal     string
	Duration         string
	Price            float64 // ADT regular fare
	Currency         string
	FaresLeft        int
}
