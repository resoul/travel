package trenitalia

// StationRawDTO represents parsed item from en.cruscotto-stations.json.
type StationRawDTO struct {
	Value string `json:"value"`
	Text  string `json:"text"`
	IsF   int    `json:"isF"`
	FB    int    `json:"FB"`
	FA    int    `json:"FA"`
	IsE   int    `json:"isE"`
}

// SolutionsRequestDTO represents payload sent to /website/ticket/solutions.
type SolutionsRequestDTO struct {
	DepartureLocationID int               `json:"departureLocationId"`
	ArrivalLocationID   int               `json:"arrivalLocationId"`
	DepartureTime       string            `json:"departureTime"`
	Adults              int               `json:"adults"`
	Children            int               `json:"children"`
	Criteria            SearchCriteriaDTO `json:"criteria"`
}

// SearchCriteriaDTO represents query filters.
type SearchCriteriaDTO struct {
	FrecceOnly   bool   `json:"frecceOnly"`
	RegionalOnly bool   `json:"regionalOnly"`
	NoChanges    bool   `json:"noChanges"`
	Order        string `json:"order"`
	Limit        int    `json:"limit"`
	Offset       int    `json:"offset"`
}

// SolutionsResponseDTO represents top-level response from ticket/solutions.
type SolutionsResponseDTO struct {
	SearchID  string               `json:"searchId"`
	Solutions []SolutionWrapperDTO `json:"solutions"`
}

// SolutionWrapperDTO wraps individual train solution.
type SolutionWrapperDTO struct {
	Solution SolutionDTO `json:"solution"`
}

// SolutionDTO represents train journey option.
type SolutionDTO struct {
	ID            string     `json:"id"`
	Origin        string     `json:"origin"`
	Destination   string     `json:"destination"`
	DepartureTime string     `json:"departureTime"`
	ArrivalTime   string     `json:"arrivalTime"`
	Duration      string     `json:"duration"`
	Status        string     `json:"status"`
	Trains        []TrainDTO `json:"trains"`
	Price         *PriceDTO  `json:"price"`
	Nodes         []NodeDTO  `json:"nodes"`
	Grids         []GridDTO  `json:"grids"`
}

// TrainDTO represents train specifications.
type TrainDTO struct {
	TrainCategory string `json:"trainCategory"`
	TrainAcronym  string `json:"trainAcronym"`
	Description   string `json:"description"`
	Denomination  string `json:"denomination"`
}

// PriceDTO represents monetary price and currency.
type PriceDTO struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

// NodeDTO represents individual trip segment.
type NodeDTO struct {
	DepartureLocation string   `json:"departureLocation"`
	ArrivalLocation   string   `json:"arrivalLocation"`
	DepartureTime     string   `json:"departureTime"`
	ArrivalTime       string   `json:"arrivalTime"`
	Duration          string   `json:"duration"`
	Train             TrainDTO `json:"train"`
}

// GridDTO represents service tiers.
type GridDTO struct {
	Services []ServiceDTO `json:"services"`
}

// ServiceDTO represents travel class (Standard, Business, etc.).
type ServiceDTO struct {
	Name   string     `json:"name"`
	Offers []OfferDTO `json:"offers"`
}

// OfferDTO represents tariff (Base, Economy, Super Economy).
type OfferDTO struct {
	Name  string    `json:"name"`
	Price *PriceDTO `json:"price"`
}
