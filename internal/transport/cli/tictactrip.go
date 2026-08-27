package cli

import (
	"fmt"
	"strings"

	"github.com/resoul/travel/internal/usecase"
	"github.com/spf13/cobra"
)

func newTictactripCmd(airportsUC *usecase.ListAirportsUseCase, datesUC *usecase.FlightDatesUseCase) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tictactrip",
		Short: "Tictactrip (ComparaBUS) European multimodal trains, buses, and low-fare calendars",
	}

	cmd.AddCommand(newTictactripCitiesCmd(airportsUC))
	cmd.AddCommand(newTictactripPopularCmd(airportsUC))
	cmd.AddCommand(newTictactripCalendarCmd(datesUC))

	return cmd
}

func newTictactripCitiesCmd(airportsUC *usecase.ListAirportsUseCase) *cobra.Command {
	var query string

	cmd := &cobra.Command{
		Use:   "cities",
		Short: "Search and autocomplete European cities and railway/bus stations",
		Example: `  travel tictactrip cities --query bar
  travel tictactrip cities --query paris
  travel tictactrip cities --query montpellier`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if query == "" {
				return fmt.Errorf("--query flag is required")
			}

			fmt.Printf("🔍 Searching Tictactrip cities and stations for '%s'...\n\n", query)

			cities, err := airportsUC.AutocompleteTictactripCities(ctx, query)
			if err != nil {
				return fmt.Errorf("failed to autocomplete cities: %w", err)
			}

			if len(cities) == 0 {
				fmt.Println("No cities or stations found.")
				return nil
			}

			fmt.Printf("Found %d locations:\n", len(cities))
			fmt.Println(strings.Repeat("-", 90))

			for i, c := range cities {
				fmt.Printf("[%2d] %s (City ID: %d | Code: %s)\n", i+1, c.LocalName, c.ID, c.UniqueName)
				fmt.Printf("     📍 Coordinates: %.4f, %.4f | Score: %d\n", c.Latitude, c.Longitude, c.Score)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&query, "query", "q", "bar", "City or station query name")

	return cmd
}

func newTictactripPopularCmd(airportsUC *usecase.ListAirportsUseCase) *cobra.Command {
	var (
		fromCity string
		limit    int
	)

	cmd := &cobra.Command{
		Use:   "popular",
		Short: "List popular train and bus travel destinations",
		Example: `  travel tictactrip popular --from barcelone --limit 7
  travel tictactrip popular --from paris --limit 10
  travel tictactrip popular --limit 10`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if fromCity != "" {
				fmt.Printf("🌟 Fetching popular destinations from '%s'...\n\n", fromCity)
			} else {
				fmt.Printf("🌟 Fetching overall popular European destinations...\n\n")
			}

			cities, err := airportsUC.GetTictactripPopularDestinations(ctx, fromCity, limit)
			if err != nil {
				return fmt.Errorf("failed to get popular destinations: %w", err)
			}

			if len(cities) == 0 {
				fmt.Println("No popular destinations found.")
				return nil
			}

			fmt.Printf("Top %d popular destinations:\n", len(cities))
			fmt.Println(strings.Repeat("-", 90))

			for i, c := range cities {
				searchCount := c.NbSearch
				if searchCount == "" {
					searchCount = "N/A"
				}
				fmt.Printf("[%2d] %s (City ID: %d | Code: %s)\n", i+1, c.LocalName, c.ID, c.UniqueName)
				fmt.Printf("     🔍 Searches: %s | Coordinates: %.4f, %.4f\n", searchCount, c.Latitude, c.Longitude)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&fromCity, "from", "", "Origin city unique name (e.g. barcelone, paris, lyon)")
	cmd.Flags().IntVar(&limit, "limit", 7, "Maximum number of popular destinations")

	return cmd
}

func newTictactripCalendarCmd(datesUC *usecase.FlightDatesUseCase) *cobra.Command {
	var (
		originID int
		destID   int
		month    string
	)

	cmd := &cobra.Command{
		Use:   "calendar",
		Short: "Display daily lowest train and bus prices for a full month",
		Example: `  travel tictactrip calendar --origin-id 76 --dest-id 542 --month 2026-11
  travel tictactrip calendar --origin-id 628 --dest-id 485 --month 2026-09`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			fmt.Printf("📅 Fetching Tictactrip monthly lowest fare calendar (Origin ID: %d -> Dest ID: %d, Month: %s)...\n\n",
				originID, destID, month)

			days, err := datesUC.GetTictactripMonthlyPriceCalendar(ctx, originID, destID, month)
			if err != nil {
				return fmt.Errorf("failed to get price calendar: %w", err)
			}

			if len(days) == 0 {
				fmt.Println("No fare calendar data available for this month.")
				return nil
			}

			fmt.Printf("Daily lowest fares for %s:\n", month)
			fmt.Println(strings.Repeat("-", 90))

			var availableCount int
			for _, d := range days {
				if !d.HasTrip {
					fmt.Printf("  %s | ❌ No direct trips recorded\n", d.Date)
					continue
				}

				availableCount++
				companiesStr := strings.Join(d.Companies, ", ")
				if companiesStr == "" {
					companiesStr = "Standard"
				}

				typeIcon := "🚌"
				if strings.ToLower(d.TransportType) == "train" {
					typeIcon = "🚄"
				}

				hours := d.DurationMinutes / 60
				mins := d.DurationMinutes % 60
				durStr := fmt.Sprintf("%dh %02dm", hours, mins)

				timeStr := ""
				if !d.DepartureTime.IsZero() {
					timeStr = fmt.Sprintf("Dep: %s", d.DepartureTime.Format("15:04"))
				}

				fmt.Printf("  %s | %s %-6s | 💰 %6.2f %s | ⏱️ %-7s | %-10s | %s\n",
					d.Date, typeIcon, strings.ToUpper(d.TransportType), d.Price.Amount, d.Price.Currency, durStr, timeStr, companiesStr)
			}

			fmt.Printf("\nTotal days with available trips: %d / %d\n", availableCount, len(days))
			return nil
		},
	}

	cmd.Flags().IntVar(&originID, "origin-id", 76, "Origin City ID (e.g. 76 for Barcelona, 628 for Paris)")
	cmd.Flags().IntVar(&destID, "dest-id", 542, "Destination City ID (e.g. 542 for Montpellier, 485 for Lyon)")
	cmd.Flags().StringVar(&month, "month", "2026-11", "Target month (YYYY-MM)")

	return cmd
}
