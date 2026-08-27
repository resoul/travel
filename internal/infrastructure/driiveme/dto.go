package driiveme

// citySearchResponseDTO represents response from POST /search/cities.
type citySearchResponseDTO struct {
	Errors      []any               `json:"errors"`
	Suggestions []CitySuggestionDTO `json:"suggestions"`
}

// CitySuggestionDTO represents an individual city item returned by search/cities.
type CitySuggestionDTO struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	SubName string `json:"subName"`
	Icon    string `json:"icon"`
}

// availabilitiesResponseDTO represents response from /transport/get-availabilities/{id}.
type availabilitiesResponseDTO struct {
	Valid          bool     `json:"valid"`
	Availabilities []string `json:"availabilities"`
	IsDepartureRdv bool     `json:"isDepartureRdv"`
	IsArrivalRdv   bool     `json:"isArrivalRdv"`
}

// loginResponseDTO represents response from POST /login.html.
type loginResponseDTO struct {
	User      string `json:"user"`
	ReturnURL string `json:"returnUrl"`
	Message   string `json:"message"`
	Error     string `json:"error"`
}

// POITransportDTO represents point of interest transport parsed from POIS_TRANSPORTS JS object.
type POITransportDTO struct {
	A POILocationDTO `json:"a"`
	B POILocationDTO `json:"b"`
}

// POILocationDTO represents departure or arrival coordinate and tooltip link.
type POILocationDTO struct {
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
	Tooltip string  `json:"tooltip"`
}

// TransportCardDTO represents a parsed trip card from HTML.
type TransportCardDTO struct {
	ID                  string
	Slug                string
	DepartureCity       string
	ArrivalCity         string
	AvailabilityStart   string
	AvailabilityEnd     string
	VehicleCategory     string
	VehicleModel        string
	Transmission        string
	Seats               int
	RentalHours         int
	Price               float64
	Currency            string
	DistanceMiles       int
	IncludedMiles       int
	Deposit             string
	Insurance           string
	OfferedBy           string
	AvailabilitiesCount int
}
