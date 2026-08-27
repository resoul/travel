package cli

import (
	"fmt"
	"strings"

	"github.com/resoul/travel/internal/domain"
	"github.com/resoul/travel/internal/usecase"
	"github.com/spf13/cobra"
)

func newTripCmd(searchUC *usecase.SearchFlightsUseCase) *cobra.Command {
	tripCmd := &cobra.Command{
		Use:   "trip",
		Short: "Trip.com hotel and accommodation search via fast Server-Side Rendered (SSR) HTTP requests",
	}

	tripCmd.AddCommand(newTripSearchCmd(searchUC))
	tripCmd.AddCommand(newTripDetailsCmd(searchUC))
	tripCmd.AddCommand(newTripCarsCmd(searchUC))

	return tripCmd
}

func newTripSearchCmd(searchUC *usecase.SearchFlightsUseCase) *cobra.Command {
	var (
		cityID   string
		cityName string
		checkIn  string
		checkOut string
		rooms    int
		adults   int
		children int
		currency string
		limit    int
	)

	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search live hotel and accommodation deals on Trip.com",
		Long: `Search hotels on Trip.com with live prices and ratings via fast SSR HTTP requests.

Examples:
  travel trip search --city-id 40795 --city Barcelona --check-in 2026-11-19 --check-out 2026-11-23 --adults 2 --currency USD
  travel trip search --city-id 35556 --city Thessaloniki --check-in 2026-09-10 --check-out 2026-09-12 --adults 2 --currency EUR
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			criteria := domain.HotelSearchCriteria{
				CityID:   cityID,
				CityName: cityName,
				CheckIn:  checkIn,
				CheckOut: checkOut,
				Rooms:    rooms,
				Adults:   adults,
				Children: children,
				Currency: currency,
				Limit:    limit,
			}

			fmt.Printf("🏨 Searching Trip.com accommodations for City ID: %s (%s) from %s to %s (%d adults, %s)...\n\n",
				cityID, cityName, checkIn, checkOut, adults, currency)

			offers, err := searchUC.SearchTripHotels(ctx, criteria)
			if err != nil {
				return fmt.Errorf("failed to search Trip.com hotels: %w", err)
			}

			if len(offers) == 0 {
				fmt.Println("No hotels found matching criteria.")
				return nil
			}

			fmt.Printf("Found %d accommodation deals on Trip.com:\n", len(offers))
			fmt.Println(strings.Repeat("-", 100))

			for i, h := range offers {
				ratingStr := "N/A"
				if h.Rating > 0 {
					ratingStr = fmt.Sprintf("%.1f/10 (%d reviews)", h.Rating, h.ReviewCount)
				}

				priceStr := "Price unavailable"
				if h.Price.Amount > 0 {
					priceStr = fmt.Sprintf("%.2f %s (%d nights, incl. taxes & fees)", h.Price.Amount, h.Price.Currency, h.Nights)
				}

				fmt.Printf("[%2d] %s (ID: %s)\n", i+1, h.Name, h.ID)
				if h.Address != "" {
					fmt.Printf("     📍 Location: %s\n", h.Address)
				}
				fmt.Printf("     ⭐ Score: %s | 💰 %s\n", ratingStr, priceStr)
				if h.RoomType != "" {
					fmt.Printf("     🛏️  Room: %s\n", h.RoomType)
				}
				if h.URL != "" {
					fmt.Printf("     🔗 Link: %s\n", h.URL)
				}
				fmt.Println()
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&cityID, "city-id", "40795", "Trip.com City ID (e.g., 40795 for Barcelona, 35556 for Thessaloniki)")
	cmd.Flags().StringVar(&cityName, "city", "Barcelona", "City name query")
	cmd.Flags().StringVar(&checkIn, "check-in", "2026-11-19", "Check-in date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&checkOut, "check-out", "2026-11-23", "Check-out date (YYYY-MM-DD)")
	cmd.Flags().IntVar(&rooms, "rooms", 1, "Number of rooms")
	cmd.Flags().IntVar(&adults, "adults", 2, "Number of adults")
	cmd.Flags().IntVar(&children, "children", 0, "Number of children")
	cmd.Flags().StringVar(&currency, "currency", "USD", "Currency code (USD, EUR, etc.)")
	cmd.Flags().IntVar(&limit, "limit", 15, "Maximum number of results to display")

	return cmd
}

func newTripDetailsCmd(searchUC *usecase.SearchFlightsUseCase) *cobra.Command {
	var (
		checkIn  string
		checkOut string
		adults   int
		rooms    int
		currency string
	)

	cmd := &cobra.Command{
		Use:   "details <hotel_id>",
		Short: "View full room breakdown, bed configurations, and amenities for a Trip.com hotel",
		Args:  cobra.ExactArgs(1),
		Example: `  travel trip details 134130947 --check-in 2026-11-19 --check-out 2026-11-23
  travel trip details 1603006 --check-in 2026-11-19 --check-out 2026-11-23`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			hotelID := args[0]

			fmt.Printf("🏨 Fetching room options and amenities for Hotel ID: %s (%s to %s)...\n\n",
				hotelID, checkIn, checkOut)

			roomsList, err := searchUC.GetTripHotelDetails(ctx, hotelID, checkIn, checkOut, adults, rooms, currency)
			if err != nil {
				return fmt.Errorf("failed to get hotel details: %w", err)
			}

			if len(roomsList) == 0 {
				fmt.Println("No detailed room configurations found.")
				return nil
			}

			fmt.Printf("Found %d room types:\n", len(roomsList))
			fmt.Println(strings.Repeat("-", 100))

			for i, r := range roomsList {
				fmt.Printf("[%2d] %s (Room ID: %s)\n", i+1, r.Name, r.ID)
				fmt.Printf("     👥 Guests: %d | 📐 Area: %s | 🛌 Beds: %s\n", r.Guests, r.Area, r.Beds)
				fmt.Printf("     🪟 Window: %s | 🚭 Smoke: %s\n", r.HasWindow, r.Smoking)
				if len(r.Amenities) > 0 {
					fmt.Printf("     ✨ Amenities: %s\n", strings.Join(r.Amenities, ", "))
				}
				if r.ImageURL != "" {
					fmt.Printf("     🖼️  Photo: %s\n", r.ImageURL)
				}
				fmt.Println()
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&checkIn, "check-in", "2026-11-19", "Check-in date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&checkOut, "check-out", "2026-11-23", "Check-out date (YYYY-MM-DD)")
	cmd.Flags().IntVar(&adults, "adults", 2, "Number of adults")
	cmd.Flags().IntVar(&rooms, "rooms", 1, "Number of rooms")
	cmd.Flags().StringVar(&currency, "currency", "USD", "Currency code (USD, EUR, etc.)")

	return cmd
}

func newTripCarsCmd(searchUC *usecase.SearchFlightsUseCase) *cobra.Command {
	var (
		countryID      string
		pickupCityID   string
		pickupCityName string
		pickupCode     string
		pickupAddress  string
		pickupDate     string
		returnDate     string
		driverAge      string
		currency       string
		limit          int
	)

	cmd := &cobra.Command{
		Use:   "cars",
		Short: "Search live rental cars and vehicle hire deals on Trip.com",
		Long: `Search rental cars across top rental agencies (Avis, Hertz, Europcar, Sixt, Surprice, Klass Wagen) via Trip.com.

Examples:
  travel trip cars --country-id 63 --city-id 39050 --city Otopeni --code OTP --pickup 2026-08-29 --return 2026-09-01 --currency USD
  travel trip cars --country-id 95 --city-id 40795 --city Barcelona --code BCN --pickup 2026-11-19 --return 2026-11-23 --currency EUR
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			criteria := domain.CarHireCriteria{
				CountryID:      countryID,
				PickupCityID:   pickupCityID,
				PickupCityName: pickupCityName,
				PickupCode:     pickupCode,
				PickupAddress:  pickupAddress,
				PickupDate:     pickupDate,
				ReturnDate:     returnDate,
				DriverAge:      driverAge,
				Currency:       currency,
				Limit:          limit,
			}

			fmt.Printf("🚗 Searching Trip.com rental cars at %s (%s) from %s to %s (Age: %s, %s)...\n\n",
				pickupCityName, pickupCode, pickupDate, returnDate, driverAge, currency)

			offers, err := searchUC.SearchTripCars(ctx, criteria)
			if err != nil {
				return fmt.Errorf("failed to search Trip.com rental cars: %w", err)
			}

			if len(offers) == 0 {
				fmt.Println("No rental cars found matching criteria.")
				return nil
			}

			fmt.Printf("Found %d rental car options on Trip.com:\n", len(offers))
			fmt.Println(strings.Repeat("-", 100))

			for i, car := range offers {
				priceDayStr := "N/A"
				if car.PricePerDay.Amount > 0 {
					priceDayStr = fmt.Sprintf("%.2f %s/day", car.PricePerDay.Amount, car.PricePerDay.Currency)
				}

				priceTotalStr := "N/A"
				if car.TotalPrice.Amount > 0 {
					priceTotalStr = fmt.Sprintf("%.2f %s total", car.TotalPrice.Amount, car.TotalPrice.Currency)
				}

				fmt.Printf("[%2d] %s (%s)\n", i+1, car.Model, car.Category)
				fmt.Printf("     🏢 Supplier: %s | 🕹️ Transmission: %s\n", car.Supplier, car.Transmission)
				fmt.Printf("     👥 Seats: %d | 🚪 Doors: %d | 🧳 Bags: %d\n", car.Seats, car.Doors, car.Bags)
				fmt.Printf("     💰 Price: %s | 💵 Total: %s\n", priceDayStr, priceTotalStr)
				if len(car.Features) > 0 {
					fmt.Printf("     ✨ Features: %s\n", strings.Join(car.Features, " • "))
				}
				if car.ImageURL != "" {
					fmt.Printf("     🖼️  Photo: %s\n", car.ImageURL)
				}
				fmt.Println()
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&countryID, "country-id", "63", "Country ID (e.g. 63 for Romania, 95 for Spain, 75 for France)")
	cmd.Flags().StringVar(&pickupCityID, "city-id", "39050", "City ID (e.g. 39050 for Otopeni, 40795 for Barcelona)")
	cmd.Flags().StringVar(&pickupCityName, "city", "Otopeni", "City/Airport name")
	cmd.Flags().StringVar(&pickupCode, "code", "OTP", "IATA airport code (e.g. OTP, BCN, ATH)")
	cmd.Flags().StringVar(&pickupAddress, "address", "Bucharest Henri Coandă International Airport (OTP)", "Full pickup address")
	cmd.Flags().StringVar(&pickupDate, "pickup", "2026-08-29 10:00", "Pickup date and time (YYYY-MM-DD HH:mm)")
	cmd.Flags().StringVar(&returnDate, "return", "2026-09-01 10:00", "Return date and time (YYYY-MM-DD HH:mm)")
	cmd.Flags().StringVar(&driverAge, "age", "30-60", "Driver age range (e.g. 30-60, 26-29, 21-25)")
	cmd.Flags().StringVar(&currency, "currency", "USD", "Currency code (USD, EUR, etc.)")
	cmd.Flags().IntVar(&limit, "limit", 15, "Maximum number of results to display")

	return cmd
}
