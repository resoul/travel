package cruise

import "github.com/resoul/travel/internal/domain"

// TokenResponse represents the response from /api/auth/get-access-token.
type TokenResponse struct {
	Token string `json:"token"`
}

// SearchMatrixResponse represents the response from /api/search/get-search-matrix.
type SearchMatrixResponse struct {
	CruiseLines  []MatrixItemDTO `json:"cruiseLines"`
	Destinations []MatrixItemDTO `json:"destinations"`
}

// MatrixItemDTO represents an item in the search matrix (cruise line or destination).
type MatrixItemDTO struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	IsAvailable bool   `json:"isAvailable"`
}

func (dto MatrixItemDTO) toCruiseLine() domain.CruiseLine {
	return domain.CruiseLine{
		ID:   dto.ID,
		Name: dto.Name,
	}
}

func (dto MatrixItemDTO) toCruiseDestination() domain.CruiseDestination {
	return domain.CruiseDestination{
		ID:   dto.ID,
		Name: dto.Name,
	}
}

// SearchResultsResponse represents the response from /api/search/get-search-results.
type SearchResultsResponse struct {
	SearchResults struct {
		CruiseSpecialList []CruiseSpecialDTO `json:"cruiseSpecialList"`
	} `json:"searchResults"`
}

// CruiseSpecialDTO represents a single cruise item in search results.
type CruiseSpecialDTO struct {
	Recno             int        `json:"recno"`
	FromPrice         *float64   `json:"fromPrice"`
	FromPriceGovtFees *float64   `json:"fromPriceGovtFees"`
	OurPrice          *float64   `json:"ourPrice"`
	Sailing           SailingDTO `json:"sailing"`
}

// SailingDTO represents sailing details for a cruise.
type SailingDTO struct {
	Recno         int          `json:"recno"`
	Duration      int          `json:"duration"`
	DaysOrNights  string       `json:"daysOrNights"`
	SailingDate   string       `json:"sailingDate"`
	CruiseLine    NamedItemDTO `json:"cruiseLine"`
	Ship          NamedItemDTO `json:"ship"`
	Itinerary     ItineraryDTO `json:"itinerary"`
	DeparturePort NamedItemDTO `json:"departurePort"`
	ArrivalPort   NamedItemDTO `json:"arrivalPort"`
}

// NamedItemDTO represents a named entity with recno and name.
type NamedItemDTO struct {
	Recno int    `json:"recno"`
	Name  string `json:"name"`
}

// ItineraryDTO represents an itinerary inside sailing.
type ItineraryDTO struct {
	Name  string    `json:"name"`
	Ports []PortDTO `json:"ports"`
}

// PortDTO represents a port of call.
type PortDTO struct {
	DayNumber string `json:"dayNumber"`
	Name      string `json:"name"`
	PortType  string `json:"portType"`
}
