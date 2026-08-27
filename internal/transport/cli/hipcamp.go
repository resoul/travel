package cli

import (
	"fmt"
	"strings"

	"github.com/resoul/travel/internal/domain"
	"github.com/resoul/travel/internal/usecase"
	"github.com/spf13/cobra"
)

func newHipcampCmd(searchUC *usecase.SearchFlightsUseCase, presenter *Presenter) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hipcamp",
		Short: "Hipcamp outdoor stays, private land camping, and glamping search",
	}

	cmd.AddCommand(newHipcampSearchCmd(searchUC, presenter))

	return cmd
}

func newHipcampSearchCmd(searchUC *usecase.SearchFlightsUseCase, presenter *Presenter) *cobra.Command {
	var (
		country string
		region  string
		limit   int
	)

	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search outdoor camping, glamping, and farm stays on Hipcamp",
		Example: `  travel hipcamp search --country united-states --region california
  travel hipcamp search --country united-kingdom --region england
  travel hipcamp search --country canada --region ontario
  travel hipcamp search --country australia --region new-south-wales`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			criteria := domain.HipcampSearchCriteria{
				Country: country,
				Region:  region,
				Limit:   limit,
			}

			countryName := strings.Title(strings.ReplaceAll(country, "-", " "))
			regionName := strings.Title(strings.ReplaceAll(region, "-", " "))
			if countryName == "" {
				countryName = "United States"
			}
			if regionName == "" {
				regionName = "California"
			}

			fmt.Printf("🌲 Searching Hipcamp outdoor spots in %s, %s...\n\n", regionName, countryName)

			offers, err := searchUC.SearchHipcampSpots(ctx, criteria)
			if err != nil {
				return fmt.Errorf("failed to search Hipcamp spots: %w", err)
			}

			if len(offers) == 0 {
				fmt.Printf("No spots found in %s, %s.\n", regionName, countryName)
				return nil
			}

			fmt.Printf("Found %d outdoor spots on Hipcamp in %s, %s:\n", len(offers), regionName, countryName)
			fmt.Println(strings.Repeat("-", 90))

			for i, offer := range offers {
				priceStr := "Price on site"
				if offer.Price.Amount > 0 {
					priceStr = fmt.Sprintf("%.2f %s/night", offer.Price.Amount, offer.Price.Currency)
				}

				fmt.Printf("[%2d] %s\n", i+1, offer.FlightNumber)
				fmt.Printf("     📍 Region: %s, %s\n", offer.DepartureStation, offer.ArrivalStation)
				fmt.Printf("     ⛺ Type: %s | 💰 Price: %s\n", offer.DepartureRaw, priceStr)
				fmt.Println()
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&country, "country", "united-states", "Country (e.g. united-states, united-kingdom, canada, australia, france)")
	cmd.Flags().StringVar(&region, "region", "california", "Region or State (e.g. california, england, ontario, new-south-wales)")
	cmd.Flags().IntVar(&limit, "limit", 10, "Maximum number of results to display")

	return cmd
}
