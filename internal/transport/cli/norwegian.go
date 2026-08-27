package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/resoul/travel/internal/usecase"
	"github.com/spf13/cobra"
)

func newNorwegianCmd(datesUC *usecase.FlightDatesUseCase, presenter *Presenter) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "norwegian",
		Short: "Norwegian Air Shuttle low-fare calendars and Scandinavian route network",
	}

	cmd.AddCommand(newNorwegianCalendarCmd(datesUC))

	return cmd
}

func newNorwegianCalendarCmd(datesUC *usecase.FlightDatesUseCase) *cobra.Command {
	var (
		origin      string
		destination string
		year        int
		month       int
		currency    string
	)

	cmd := &cobra.Command{
		Use:   "calendar",
		Short: "Display Norwegian Air Shuttle daily lowest flight prices for a whole month",
		Example: `  travel norwegian calendar --origin OSL --destination BCN --year 2026 --month 9
  travel norwegian calendar --origin ARN --destination LGW --year 2026 --month 10
  travel norwegian calendar --origin CPH --destination FCO --year 2026 --month 11 --currency EUR`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if origin == "" || destination == "" {
				return fmt.Errorf("--origin and --destination flags are required (e.g. --origin OSL --destination BCN)")
			}

			if year <= 0 {
				year = time.Now().Year()
			}
			if month <= 0 {
				month = int(time.Now().Month())
			}
			if currency == "" {
				currency = "EUR"
			}

			fmt.Printf("📅 Fetching Norwegian Low-Fare Calendar (%s -> %s for %04d-%02d in %s)...\n\n",
				strings.ToUpper(origin), strings.ToUpper(destination), year, month, strings.ToUpper(currency))

			offers, err := datesUC.GetNorwegianFareCalendar(ctx, origin, destination, year, month, currency)
			if err != nil {
				return fmt.Errorf("failed to fetch Norwegian fare calendar: %w", err)
			}

			if len(offers) == 0 {
				fmt.Println("No flight fares found for the specified route and month.")
				return nil
			}

			fmt.Printf("Lowest Norwegian daily fares for %04d-%02d:\n", year, month)
			fmt.Println(strings.Repeat("-", 75))

			var lowestFare float64 = 999999
			var lowestDate string

			for _, offer := range offers {
				dateStr := offer.DepartureRaw
				if offer.DepartureTime != nil {
					dateStr = offer.DepartureTime.Format("2006-01-02")
				}

				if offer.Price.Amount < lowestFare {
					lowestFare = offer.Price.Amount
					lowestDate = dateStr
				}

				fmt.Printf("  %s | ✈️ NORWEGIAN | 💰 %7.2f %s | Status: %s\n",
					dateStr, offer.Price.Amount, offer.Price.Currency, offer.Status)
			}

			fmt.Println(strings.Repeat("-", 75))
			if lowestFare < 999999 {
				fmt.Printf("🌟 Best Price Deal: %.2f %s on %s\n", lowestFare, currency, lowestDate)
			}
			fmt.Printf("Total scheduled departure days found: %d\n", len(offers))

			return nil
		},
	}

	cmd.Flags().StringVar(&origin, "origin", "OSL", "Origin 3-letter IATA code (e.g. OSL, ARN, CPH, LGW, BCN)")
	cmd.Flags().StringVar(&destination, "destination", "BCN", "Destination 3-letter IATA code (e.g. BCN, FCO, CDG, OSL, ALC)")
	cmd.Flags().IntVar(&year, "year", 2026, "Year (e.g. 2026)")
	cmd.Flags().IntVar(&month, "month", 9, "Month (1-12)")
	cmd.Flags().StringVar(&currency, "currency", "EUR", "Currency code (e.g. EUR, USD, NOK, SEK)")

	return cmd
}
