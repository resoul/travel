package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/resoul/travel/internal/domain"
	"github.com/resoul/travel/internal/usecase"
	"github.com/spf13/cobra"
)

func newPitchupCmd(searchUC *usecase.SearchFlightsUseCase, presenter *Presenter) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pitchup",
		Short: "Pitchup campsites, glamping, holiday parks, and caravan pitches search",
	}

	cmd.AddCommand(newPitchupSearchCmd(searchUC, presenter))

	return cmd
}

func newPitchupSearchCmd(searchUC *usecase.SearchFlightsUseCase, presenter *Presenter) *cobra.Command {
	var (
		country string
		region  string
		arrive  string
		depart  string
		adults  int
		limit   int
	)

	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search campsites, glamping, and holiday parks on Pitchup",
		Example: `  travel pitchup search --country france --arrive 2026-09-10 --depart 2026-09-12
  travel pitchup search --country england --arrive 2026-09-15 --depart 2026-09-18
  travel pitchup search --country italy --arrive 2026-10-01 --depart 2026-10-05 --adults 2`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if arrive == "" {
				arrive = time.Now().AddDate(0, 0, 14).Format("2006-01-02")
			}
			if depart == "" {
				if t, err := time.Parse("2006-01-02", arrive); err == nil {
					depart = t.AddDate(0, 0, 2).Format("2006-01-02")
				} else {
					depart = time.Now().AddDate(0, 0, 16).Format("2006-01-02")
				}
			}

			criteria := domain.PitchupSearchCriteria{
				Country:    country,
				Region:     region,
				ArriveDate: arrive,
				DepartDate: depart,
				Adults:     adults,
				Limit:      limit,
			}

			countryName := strings.Title(strings.ToLower(country))
			if countryName == "" {
				countryName = "France"
			}

			fmt.Printf("⛺ Searching Pitchup campsites in %s (%s -> %s for %d adults)...\n\n",
				countryName, arrive, depart, adults)

			offers, err := searchUC.SearchPitchupCampsites(ctx, criteria)
			if err != nil {
				return fmt.Errorf("failed to search Pitchup campsites: %w", err)
			}

			if len(offers) == 0 {
				fmt.Printf("No campsites found in %s for the selected dates.\n", countryName)
				return nil
			}

			fmt.Printf("Found %d campsites & glamping spots on Pitchup in %s:\n", len(offers), countryName)
			fmt.Println(strings.Repeat("-", 90))

			for i, offer := range offers {
				priceStr := "Price unavailable"
				if offer.Price.Amount > 0 {
					priceStr = fmt.Sprintf("%.2f %s", offer.Price.Amount, offer.Price.Currency)
				}

				fmt.Printf("[%2d] %s\n", i+1, offer.FlightNumber)
				if offer.DepartureStation != "" {
					fmt.Printf("     📍 Location: %s\n", offer.DepartureStation)
				}
				fmt.Printf("     📅 Dates: %s | 🏕️ Unit: %s\n", offer.DepartureRaw, offer.Status)
				fmt.Printf("     💰 Total Price: %s\n", priceStr)
				fmt.Println()
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&country, "country", "france", "Destination country (e.g. france, england, spain, italy, germany, scotland, wales, usa)")
	cmd.Flags().StringVar(&region, "region", "", "Optional region or area")
	cmd.Flags().StringVar(&arrive, "arrive", "", "Arrival / Check-in date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&depart, "depart", "", "Departure / Check-out date (YYYY-MM-DD)")
	cmd.Flags().IntVar(&adults, "adults", 2, "Number of adult guests")
	cmd.Flags().IntVar(&limit, "limit", 10, "Maximum number of results to display")

	return cmd
}
