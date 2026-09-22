package alsa

// StationDTO represents a station from ALSA API
type StationDTO struct {
	ID                  string       `json:"id"`
	Name                string       `json:"name"`
	ProvinceCountryName string       `json:"provincecountryName"`
	NameAutocomplete    string       `json:"nameAutocomplete"`
	Icon                string       `json:"icon"`
	IsMoveliaStop       bool         `json:"isMoveliaStop"`
	SimplifiedName      string       `json:"simplifiedName"`
	StopsList           []StationDTO `json:"stopsList"`
	Priority            string       `json:"priority"`
	Nomovelia           string       `json:"nomovelia"`
	AirportName         string       `json:"airportName"`
	IsFavorite          bool         `json:"isFavorite"`
	ShowFavoriteButton  bool         `json:"showFavoriteButton"`
}

// JourneyDTO represents a journey/trip result from search
type JourneyDTO struct {
	DepartureTime    string `json:"departureTime"`
	ArrivalTime      string `json:"arrivalTime"`
	DepartureStation string `json:"departureStation"`
	ArrivalStation   string `json:"arrivalStation"`
	Duration         string `json:"duration"`
	Price            string `json:"price"`
	Currency         string `json:"currency"`
	Stops            int    `json:"stops"`
	TransportNumber  string `json:"transportNumber"`
	OperatedBy       string `json:"operatedBy"`
	TravelClass      string `json:"travelClass"`
}

// SearchPayloadDTO is the request payload for journey search
type SearchPayloadDTO struct {
	OriginStationID      string         `json:"originStationId"`
	DestinationStationID string         `json:"destinationStationId"`
	DepartureDate        string         `json:"departureDate"`
	ReturnDate           string         `json:"returnDate"`
	Seats                int            `json:"seats"`
	PromoCode            string         `json:"promoCode"`
	YoungPromoCode       string         `json:"youngPromoCode"`
	TravelType           string         `json:"travelType"`
	Employee             bool           `json:"employee"`
	PassengerType        map[string]int `json:"passengerType"`
}
