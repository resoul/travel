package agoda

// HotelCardDTO represents extracted hotel properties from Agoda DOM or API.
type HotelCardDTO struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	City     string  `json:"city"`
	Address  string  `json:"address"`
	Stars    float64 `json:"stars"`
	Rating   float64 `json:"rating"`
	Reviews  int     `json:"reviews"`
	Price    float64 `json:"price"`
	Currency string  `json:"currency"`
	RoomType string  `json:"room_type"`
	URL      string  `json:"url"`
	ImageURL string  `json:"image_url"`
	Nights   int     `json:"nights"`
}

// CountryItemDTO represents a country entry from the static CDN list.
type CountryItemDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Longitude   string `json:"longitude"`
	Latitude    string `json:"latitude"`
	ISO2        string `json:"iso2"`
	CallingCode string `json:"callingCode"`
	LanguageID  int    `json:"languageId"`
}
