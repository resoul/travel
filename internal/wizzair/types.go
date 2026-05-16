package wizzair

type Connection struct {
	IATA string `json:"iata"`
}

type City struct {
	IATA        string       `json:"iata"`
	ShortName   string       `json:"shortName"`
	CountryName string       `json:"countryName"`
	Connections []Connection `json:"connections"`
}

type MapResponse struct {
	Cities []City `json:"cities"`
}

type Flight struct {
	DepartureStation string `json:"departureStation"`
	ArrivalStation   string `json:"arrivalStation"`
	From             string `json:"from"`
	To               string `json:"to"`
}

type FlightRequest struct {
	FlightList []Flight `json:"flightList"`

	PriceType          string `json:"priceType"`
	AdultCount         int    `json:"adultCount"`
	ChildCount         int    `json:"childCount"`
	InfantCount        int    `json:"infantCount"`
	MacStationsAllowed bool   `json:"macStationsAllowed"`
}

type Price struct {
	Amount       float64 `json:"amount"`
	CurrencyCode string  `json:"currencyCode"`
}

type FlightResult struct {
	ArrivalStation   string `json:"arrivalStation"`
	DepartureStation string `json:"departureStation"`
	DepartureDate    string `json:"departureDate"`

	Price Price `json:"price"`
}

type TimetableResponse struct {
	OutboundFlights []FlightResult `json:"outboundFlights"`
	ReturnFlights   []FlightResult `json:"returnFlights"`
}
