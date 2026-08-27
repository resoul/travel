package cli

import (
	"fmt"
	"strings"

	"github.com/resoul/travel/internal/domain"
	"github.com/resoul/travel/internal/usecase"
	"github.com/spf13/cobra"
)

func newCampspaceCmd(searchUC *usecase.SearchFlightsUseCase, presenter *Presenter) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "campspace",
		Short: "Campspace sustainable micro-camping, camper sites, and nature stays search",
	}

	cmd.AddCommand(newCampspaceSearchCmd(searchUC, presenter))

	return cmd
}

func newCampspaceSearchCmd(searchUC *usecase.SearchFlightsUseCase, presenter *Presenter) *cobra.Command {
	var (
		category string
		limit    int
	)

	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search micro-camping spots by category on Campspace",
		Example: `  travel campspace search --category tent-pitches
  travel campspace search --category camper-sites
  travel campspace search --category glamping
  travel campspace search --category treehouses
  travel campspace search --category yurts`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			criteria := domain.CampspaceSearchCriteria{
				Category: category,
				Limit:    limit,
			}

			catName := strings.Title(strings.ReplaceAll(category, "-", " "))
			if catName == "" {
				catName = "Tent Pitches"
			}

			fmt.Printf("🌿 Searching Campspace micro-camping spots for %s...\n\n", catName)

			offers, err := searchUC.SearchCampspaceSpots(ctx, criteria)
			if err != nil {
				return fmt.Errorf("failed to search Campspace spots: %w", err)
			}

			if len(offers) == 0 {
				fmt.Printf("No spots found for category %s.\n", catName)
				return nil
			}

			fmt.Printf("Found %d micro-camping spots on Campspace for %s:\n", len(offers), catName)
			fmt.Println(strings.Repeat("-", 90))

			for i, offer := range offers {
				priceStr := "Price on site"
				if offer.Price.Amount > 0 {
					priceStr = fmt.Sprintf("%.2f %s/night", offer.Price.Amount, offer.Price.Currency)
				}

				fmt.Printf("[%2d] %s\n", i+1, offer.FlightNumber)
				if offer.DepartureStation != "" {
					fmt.Printf("     📍 Location: %s\n", offer.DepartureStation)
				}
				fmt.Printf("     🏕️ Category: %s | 💰 Price: %s\n", offer.ArrivalStation, priceStr)
				fmt.Println()
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&category, "category", "tent-pitches", "Category (e.g. tent-pitches, camper-sites, glamping, treehouses, yurts)")
	cmd.Flags().IntVar(&limit, "limit", 10, "Maximum number of results to display")

	return cmd
}
