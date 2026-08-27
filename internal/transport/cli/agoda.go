package cli

import (
	"fmt"
	"strings"

	"github.com/resoul/travel/internal/domain"
	"github.com/resoul/travel/internal/usecase"
	"github.com/spf13/cobra"
)

func newAgodaCmd(searchUC *usecase.SearchFlightsUseCase, airportUC *usecase.ListAirportsUseCase) *cobra.Command {
	agodaCmd := &cobra.Command{
		Use:   "agoda",
		Short: "Agoda hotel and accommodation search via Chromedp headless browser and CDN",
	}

	agodaCmd.AddCommand(newAgodaSearchCmd(searchUC))
	agodaCmd.AddCommand(newAgodaCountriesCmd(airportUC))

	return agodaCmd
}

func newAgodaSearchCmd(searchUC *usecase.SearchFlightsUseCase) *cobra.Command {
	var (
		cityID   string
		cityName string
		checkIn  string
		checkOut string
		rooms    int
		adults   int
		children int
		currency string
		sort     string
		limit    int
	)

	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search live hotel and accommodation deals on Agoda",
		Long: `Search hotels on Agoda with live prices and ratings via Chromedp headless browser.

Examples:
  travel agoda search --city-id 19216 --city Mamaia --check-in 2026-09-06 --check-out 2026-09-08 --adults 2 --sort priceLowToHigh
  travel agoda search --city-id 17336 --city Thessaloniki --check-in 2026-08-28 --check-out 2026-08-29
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
				Sort:     sort,
				Limit:    limit,
			}

			fmt.Printf("🏨 Searching Agoda accommodations for City ID: %s (%s) from %s to %s (%d adults, %s)...\n\n",
				cityID, cityName, checkIn, checkOut, adults, currency)

			offers, err := searchUC.SearchAgodaHotels(ctx, criteria)
			if err != nil {
				return fmt.Errorf("failed to search Agoda hotels: %w", err)
			}

			if len(offers) == 0 {
				fmt.Println("No hotels found matching criteria.")
				return nil
			}

			fmt.Printf("Found %d accommodation deals on Agoda:\n", len(offers))
			fmt.Println(strings.Repeat("-", 100))

			for i, h := range offers {
				ratingStr := "N/A"
				if h.Rating > 0 {
					ratingStr = fmt.Sprintf("%.1f/10 (%d reviews)", h.Rating, h.ReviewCount)
				}

				priceStr := "Price unavailable"
				if h.Price.Amount > 0 {
					priceStr = fmt.Sprintf("%.2f %s (%d nights)", h.Price.Amount, h.Price.Currency, h.Nights)
				}

				fmt.Printf("[%2d] %s\n", i+1, h.Name)
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

	cmd.Flags().StringVar(&cityID, "city-id", "19216", "Agoda City ID (e.g., 19216 for Mamaia, 17336 for Thessaloniki)")
	cmd.Flags().StringVar(&cityName, "city", "Mamaia", "City name query")
	cmd.Flags().StringVar(&checkIn, "check-in", "2026-09-06", "Check-in date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&checkOut, "check-out", "2026-09-08", "Check-out date (YYYY-MM-DD)")
	cmd.Flags().IntVar(&rooms, "rooms", 1, "Number of rooms")
	cmd.Flags().IntVar(&adults, "adults", 2, "Number of adults")
	cmd.Flags().IntVar(&children, "children", 0, "Number of children")
	cmd.Flags().StringVar(&currency, "currency", "EUR", "Currency code (EUR, USD, etc.)")
	cmd.Flags().StringVar(&sort, "sort", "priceLowToHigh", "Sort mode (priceLowToHigh, rank, etc.)")
	cmd.Flags().IntVar(&limit, "limit", 20, "Maximum number of results to display")

	return cmd
}

func newAgodaCountriesCmd(airportUC *usecase.ListAirportsUseCase) *cobra.Command {
	var langID int

	cmd := &cobra.Command{
		Use:   "countries",
		Short: "List world countries from Agoda CDN directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			countries, err := airportUC.GetAgodaCountries(ctx, langID)
			if err != nil {
				return fmt.Errorf("failed to get countries: %w", err)
			}

			fmt.Printf("Retrieved %d countries from Agoda CDN:\n", len(countries))
			for i, c := range countries {
				if i < 25 || i >= len(countries)-5 {
					fmt.Printf("- [%s] %s (ID: %s)\n", c.ISO3Code, c.Name, c.Code)
				} else if i == 25 {
					fmt.Printf("... (%d more countries) ...\n", len(countries)-30)
				}
			}

			return nil
		},
	}

	cmd.Flags().IntVar(&langID, "lang", 11, "Language ID (11 for RU, 1 for EN)")

	return cmd
}
