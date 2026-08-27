package cli

import (
	"fmt"
	"strings"

	"github.com/resoul/travel/internal/usecase"
	"github.com/spf13/cobra"
)

func newSATACmd(airportsUC *usecase.ListAirportsUseCase, datesUC *usecase.FlightDatesUseCase, presenter *Presenter) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "sata",
		Aliases: []string{"azores"},
		Short:   "Azores Airlines / SATA Air Açores route network and 365-day fare calendar",
	}

	cmd.AddCommand(newSATAAirportsCmd(airportsUC, presenter))
	cmd.AddCommand(newSATARoutesCmd(airportsUC))
	cmd.AddCommand(newSATACalendarCmd(datesUC, presenter))

	return cmd
}

func newSATAAirportsCmd(airportsUC *usecase.ListAirportsUseCase, presenter *Presenter) *cobra.Command {
	return &cobra.Command{
		Use:   "airports",
		Short: "List all airports in the Azores Airlines / SATA network",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			airports, err := airportsUC.GetSATAAirports(ctx)
			if err != nil {
				return fmt.Errorf("failed to get Azores Airlines airports: %w", err)
			}
			presenter.PrintAirports(airports)
			return nil
		},
	}
}

func newSATARoutesCmd(airportsUC *usecase.ListAirportsUseCase) *cobra.Command {
	var origin string

	cmd := &cobra.Command{
		Use:   "routes",
		Short: "List direct destination airport codes on Azores Airlines / SATA from an origin airport",
		Example: `  travel sata routes --origin PDL
  travel sata routes --origin LIS
  travel sata routes --origin TER
  travel sata routes --origin BOS`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			origin = strings.ToUpper(strings.TrimSpace(origin))
			if origin == "" {
				return fmt.Errorf("origin airport code is required (e.g. --origin PDL)")
			}

			destinations, err := airportsUC.GetSATARoutesFromOrigin(ctx, origin)
			if err != nil {
				return fmt.Errorf("failed to get Azores Airlines routes for %s: %w", origin, err)
			}

			if len(destinations) == 0 {
				fmt.Printf("No direct destinations found from %s.\n", origin)
				return nil
			}

			fmt.Printf("Direct destinations from %s on Azores Airlines / SATA (%d destinations):\n", origin, len(destinations))
			fmt.Println(strings.Join(destinations, ", "))
			return nil
		},
	}

	cmd.Flags().StringVar(&origin, "origin", "PDL", "Origin airport IATA code (e.g. PDL, LIS, TER, OPO, BOS, JFK)")

	return cmd
}

func newSATACalendarCmd(datesUC *usecase.FlightDatesUseCase, presenter *Presenter) *cobra.Command {
	var (
		origin      string
		destination string
		limit       int
	)

	cmd := &cobra.Command{
		Use:   "calendar",
		Short: "View Azores Airlines low-fare calendar across all available dates",
		Example: `  travel sata calendar --origin LIS --destination PDL
  travel sata calendar --origin PDL --destination LIS
  travel sata calendar --origin BOS --destination PDL
  travel sata calendar --origin OPO --destination TER`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			origin = strings.ToUpper(strings.TrimSpace(origin))
			destination = strings.ToUpper(strings.TrimSpace(destination))

			if origin == "" || destination == "" {
				return fmt.Errorf("origin and destination airport codes are required (e.g. --origin LIS --destination PDL)")
			}

			fmt.Printf("📅 Fetching Azores Airlines / SATA low-fare calendar for %s -> %s...\n\n", origin, destination)

			offers, err := datesUC.GetSATAFareCalendar(ctx, origin, destination)
			if err != nil {
				return fmt.Errorf("failed to get Azores Airlines fare calendar: %w", err)
			}

			if len(offers) == 0 {
				fmt.Printf("No scheduled fares found for %s -> %s.\n", origin, destination)
				return nil
			}

			if limit > 0 && len(offers) > limit {
				offers = offers[:limit]
			}

			presenter.PrintFlightOffers(offers)
			return nil
		},
	}

	cmd.Flags().StringVar(&origin, "origin", "", "Origin airport IATA code (e.g. LIS, PDL, OPO, BOS, JFK, FRA)")
	cmd.Flags().StringVar(&destination, "destination", "", "Destination airport IATA code (e.g. PDL, LIS, TER, HOR, BOS)")
	cmd.Flags().IntVar(&limit, "limit", 30, "Maximum number of daily fares to display (0 for all)")

	_ = cmd.MarkFlagRequired("origin")
	_ = cmd.MarkFlagRequired("destination")

	return cmd
}
