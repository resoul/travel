package domain

import "time"

// Price represents monetary amount and currency.
type Price struct {
	Amount   float64
	Currency string
}

// TransportType represents the type of transportation mode.
type TransportType string

const (
	TransportTypeFlight TransportType = "Flight"
	TransportTypeBus    TransportType = "Bus"
	TransportTypeTrain  TransportType = "Train"
	TransportTypeCar    TransportType = "Car"
)

// FlightOffer represents a travel offer option found across airlines and ground operators.
type FlightOffer struct {
	TransportType    TransportType
	Airline          string
	FlightNumber     string
	DepartureStation string
	ArrivalStation   string
	DepartureTime    *time.Time
	ArrivalTime      *time.Time
	DepartureRaw     string
	ArrivalRaw       string
	Duration         string
	Price            Price
	SeatsLeft        int
	IsAvailable      bool
	Status           string
}

// FlightSearchCriteria defines search parameters for flights.
type FlightSearchCriteria struct {
	Origin         string
	Destination    string
	DepartureDate  string // Format: YYYY-MM-DD
	ReturnDate     string // Format: YYYY-MM-DD, optional
	Adults         int
	Teens          int
	Children       int
	Infants        int
	RoundTrip      bool
	FlexDaysBefore int
	FlexDaysAfter  int
}

// WizzairSearchCriteria defines search parameters specific to Wizzair searches.
type WizzairSearchCriteria struct {
	DepartureStation string
	ArrivalStation   string
	FromDate         string // Format: YYYY-MM-DD
	ToDate           string // Format: YYYY-MM-DD
	Adults           int
	Children         int
	Infants          int
	PriceType        string
}
