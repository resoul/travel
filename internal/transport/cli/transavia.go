package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/resoul/travel/internal/usecase"
	"github.com/spf13/cobra"
)

func newTransaviaCmd(datesUC *usecase.FlightDatesUseCase, presenter *Presenter) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transavia",
		Short: "Transavia (Air France-KLM Group) low-fare calendar search",
	}

	cmd.AddCommand(newTransaviaCalendarCmd(datesUC, presenter))

	return cmd
}

func newTransaviaCalendarCmd(datesUC *usecase.FlightDatesUseCase, presenter *Presenter) *cobra.Command {
	var (
		origin      string
		destination string
		year        int
		month       int
		adults      int
	)

	cmd := &cobra.Command{
		Use:   "calendar",
		Short: "View full month low-fare calendar with daily prices on Transavia",
		Example: `  travel transavia calendar --origin AMS --destination BCN --year 2026 --month 9
  travel transavia calendar --origin ORY --destination MAD --year 2026 --month 10
  travel transavia calendar --origin RTM --destination ALC --year 2026 --month 11`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if origin == "" || destination == "" {
				return fmt.Errorf("--origin and --destination flags are required (e.g. --origin AMS --destination BCN)")
			}

			if year <= 0 {
				year = time.Now().Year()
			}
			if month <= 0 || month > 12 {
				month = int(time.Now().Month())
			}

			originUpper := strings.ToUpper(origin)
			destUpper := strings.ToUpper(destination)

			fmt.Printf("📅 Fetching Transavia low-fare calendar for %s -> %s (%04d-%02d)...\n\n",
				originUpper, destUpper, year, month)

			offers, err := datesUC.GetTransaviaFareCalendar(ctx, originUpper, destUpper, year, month, adults)
			if err != nil {
				return fmt.Errorf("failed to get Transavia low-fare calendar: %w", err)
			}

			if len(offers) == 0 {
				fmt.Printf("No fares found for %s -> %s in %04d-%02d.\n", originUpper, destUpper, year, month)
				return nil
			}

			fmt.Printf("Found %d daily fare options on Transavia for %04d-%02d:\n", len(offers), year, month)
			fmt.Println(strings.Repeat("-", 80))

			var lowestPrice float64
			var lowestDate string

			for _, offer := range offers {
				priceStr := fmt.Sprintf("%.2f %s", offer.Price.Amount, offer.Price.Currency)
				if lowestPrice == 0 || (offer.Price.Amount > 0 && offer.Price.Amount < lowestPrice) {
					lowestPrice = offer.Price.Amount
					lowestDate = offer.DepartureRaw
				}

				tag := "  "
				if strings.Contains(offer.Status, "LOW") {
					tag = "🔥"
				}

				fmt.Printf("%s 📅 %s  |  💰 %10s  |  %s\n",
					tag, offer.DepartureRaw, priceStr, offer.Status)
			}

			fmt.Println(strings.Repeat("-", 80))
			if lowestPrice > 0 {
				fmt.Printf("🏆 Best Fare: %.2f EUR on %s\n", lowestPrice, lowestDate)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&origin, "origin", "AMS", "Origin 3-letter IATA airport code")
	cmd.Flags().StringVar(&destination, "destination", "BCN", "Destination 3-letter IATA airport code")
	cmd.Flags().IntVar(&year, "year", time.Now().Year(), "Year (e.g. 2026)")
	cmd.Flags().IntVar(&month, "month", int(time.Now().Month()), "Month (1-12)")
	cmd.Flags().IntVar(&adults, "adults", 1, "Number of adult passengers")

	return cmd
}
