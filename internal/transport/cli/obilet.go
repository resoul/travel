package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/resoul/travel/internal/domain"
	"github.com/resoul/travel/internal/usecase"
	"github.com/spf13/cobra"
)

func newOBiletCmd(searchUC *usecase.SearchFlightsUseCase, presenter *Presenter) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "obilet",
		Short: "oBilet Turkish and regional intercity bus journey search",
	}

	cmd.AddCommand(newOBiletSearchCmd(searchUC, presenter))

	return cmd
}

func newOBiletSearchCmd(searchUC *usecase.SearchFlightsUseCase, presenter *Presenter) *cobra.Command {
	var (
		origin        string
		destination   string
		originID      int
		destinationID int
		date          string
		limit         int
	)

	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search intercity and regional bus journeys on oBilet",
		Example: `  travel obilet search --origin istanbul --destination ankara --date 2026-09-10
  travel obilet search --origin izmir --destination antalya --date 2026-09-15
  travel obilet search --origin antalya --destination cappadocia --date 2026-09-20`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if date == "" {
				date = time.Now().AddDate(0, 0, 7).Format("2006-01-02")
			}

			criteria := domain.OBiletSearchCriteria{
				OriginID:        originID,
				DestinationID:   destinationID,
				OriginName:      origin,
				DestinationName: destination,
				DepartureDate:   date,
				Limit:           limit,
			}

			fmt.Printf("🚌 Searching oBilet buses from %s to %s on %s...\n\n",
				strings.Title(origin), strings.Title(destination), date)

			offers, err := searchUC.SearchOBiletBuses(ctx, criteria)
			if err != nil {
				return fmt.Errorf("failed to search oBilet buses: %w", err)
			}

			if len(offers) == 0 {
				fmt.Println("No bus journeys found for this route and date.")
				return nil
			}

			fmt.Printf("Found %d bus journeys on oBilet:\n", len(offers))
			fmt.Println(strings.Repeat("-", 100))

			for i, offer := range offers {
				depStr := offer.DepartureRaw
				if offer.DepartureTime != nil {
					depStr = offer.DepartureTime.Format("15:04")
				}

				priceStr := "Price unavailable"
				if offer.Price.Amount > 0 {
					priceStr = fmt.Sprintf("%.2f %s", offer.Price.Amount, offer.Price.Currency)
				}

				fmt.Printf("[%2d] 🚌 %s\n", i+1, offer.FlightNumber)
				fmt.Printf("     📍 Route: %s -> %s\n", offer.DepartureStation, offer.ArrivalStation)
				fmt.Printf("     🕒 Departure: %s | ⏱️ Duration: %s\n", depStr, offer.Duration)
				fmt.Printf("     💰 Price: %s | Status: %s\n", priceStr, offer.Status)
				fmt.Println()
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&origin, "origin", "istanbul", "Origin city or station (e.g. istanbul, ankara, izmir, antalya, bursa, bodrum)")
	cmd.Flags().StringVar(&destination, "destination", "ankara", "Destination city or station (e.g. ankara, izmir, antalya, cappadocia)")
	cmd.Flags().IntVar(&originID, "origin-id", 0, "Optional numeric Origin Location ID (e.g. 349 for Istanbul, 356 for Ankara)")
	cmd.Flags().IntVar(&destinationID, "dest-id", 0, "Optional numeric Destination Location ID")
	cmd.Flags().StringVar(&date, "date", "", "Departure date (YYYY-MM-DD)")
	cmd.Flags().IntVar(&limit, "limit", 15, "Maximum number of results to display")

	return cmd
}
