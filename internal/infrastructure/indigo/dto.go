package indigo

// TokenCreateRequest represents the payload for creating a session token.
type TokenCreateRequest struct {
	StrToken        string `json:"strToken"`
	SubscriptionKey string `json:"subscriptionKey"`
}

// TokenCreateResponse represents the response from token creation.
type TokenCreateResponse struct {
	Data struct {
		Token struct {
			Token                string `json:"token"`
			IdleTimeoutInMinutes int    `json:"idleTimeoutInMinutes"`
		} `json:"token"`
		RoleName string `json:"roleName"`
	} `json:"data"`
}

// FareRadarResponse represents the response from the 6ewai fare-radar endpoint.
type FareRadarResponse struct {
	Origin     string          `json:"origin"`
	OriginCity string          `json:"originCity"`
	Currency   string          `json:"currency"`
	TravelDate string          `json:"travelDate"`
	DateLabel  string          `json:"dateLabel"`
	Fares      []FareRadarItem `json:"fares"`
}

// FareRadarItem represents a destination in fare radar.
type FareRadarItem struct {
	IATA string  `json:"iata"`
	City string  `json:"city"`
	Fare float64 `json:"fare"`
	Time string  `json:"time"`
}

// FareCalendarRequest represents the request payload for getfarecalendar.
type FareCalendarRequest struct {
	StartDate    string `json:"startDate"`
	EndDate      string `json:"endDate"`
	Origin       string `json:"origin"`
	Destination  string `json:"destination"`
	CurrencyCode string `json:"currencyCode"`
	PromoCode    string `json:"promoCode"`
	LowestIn     string `json:"lowestIn"`
}

// FareCalendarResponse represents the response from getfarecalendar.
type FareCalendarResponse struct {
	Data struct {
		LowFares []FareCalendarItem `json:"lowFares"`
	} `json:"data"`
}

// FareCalendarItem represents a single day fare in getfarecalendar.
type FareCalendarItem struct {
	Date     string  `json:"date"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
}
