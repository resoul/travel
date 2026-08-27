package trip

// HotelItemDTO represents parsed hotel search card.
type HotelItemDTO struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	City       string  `json:"city"`
	Address    string  `json:"address"`
	Rating     float64 `json:"rating"`
	Reviews    int     `json:"reviews"`
	PriceNight float64 `json:"price_night"`
	TotalPrice float64 `json:"total_price"`
	Currency   string  `json:"currency"`
	RoomName   string  `json:"room_name"`
	URL        string  `json:"url"`
	ImageURL   string  `json:"image_url"`
	Nights     int     `json:"nights"`
}

// RoomDetailDTO represents detailed room attributes.
type RoomDetailDTO struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Area      string   `json:"area"`
	Beds      string   `json:"beds"`
	Guests    int      `json:"guests"`
	HasWindow string   `json:"has_window"`
	Smoking   string   `json:"smoking"`
	Amenities []string `json:"amenities"`
	ImageURL  string   `json:"image_url"`
}
