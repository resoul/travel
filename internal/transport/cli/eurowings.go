package cli

import (
	"fmt"
	"strings"

	"github.com/resoul/travel/internal/usecase"
	"github.com/spf13/cobra"
)

func newEurowingsCmd(airportsUC *usecase.ListAirportsUseCase, datesUC *usecase.FlightDatesUseCase, presenter *Presenter) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "eurowings",
		Short: "Eurowings (Lufthansa Group) airports, route networks, and flight schedules",
	}

	cmd.AddCommand(newEurowingsAirportsCmd(airportsUC, presenter))
	cmd.AddCommand(newEurowingsRoutesCmd(airportsUC))
	cmd.AddCommand(newEurowingsDatesCmd(datesUC, presenter))

	return cmd
}

func newEurowingsAirportsCmd(airportsUC *usecase.ListAirportsUseCase, presenter *Presenter) *cobra.Command {
	var country string

	cmd := &cobra.Command{
		Use:   "airports",
		Short: "List active airports in the Eurowings network",
		Example: `  travel eurowings airports
  travel eurowings airports --country DE
  travel eurowings airports --country ES`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			fmt.Println("🔍 Fetching Eurowings airports catalog...")

			airports, err := airportsUC.GetEurowingsAirports(ctx)
			if err != nil {
				return fmt.Errorf("failed to list Eurowings airports: %w", err)
			}

			if country != "" {
				cUpper := strings.ToUpper(country)
				filtered := airports[:0]
				for _, a := range airports {
					if strings.ToUpper(a.Country.Code) == cUpper || strings.EqualFold(a.Country.Name, country) {
						filtered = append(filtered, a)
					}
				}
				airports = filtered
			}

			if len(airports) == 0 {
				fmt.Println("No airports found.")
				return nil
			}

			presenter.PrintAirports(airports)
			return nil
		},
	}

	cmd.Flags().StringVar(&country, "country", "", "Filter airports by 2-letter ISO country code (e.g. DE, ES, IT, GR)")

	return cmd
}

func newEurowingsRoutesCmd(airportsUC *usecase.ListAirportsUseCase) *cobra.Command {
	var origin string

	cmd := &cobra.Command{
		Use:   "routes",
		Short: "List direct destination airport codes available from an origin on Eurowings",
		Example: `  travel eurowings routes --origin OTP
  travel eurowings routes --origin DUS
  travel eurowings routes --origin BER`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if origin == "" {
				return fmt.Errorf("--origin flag is required (e.g. --origin OTP)")
			}

			originUpper := strings.ToUpper(origin)
			fmt.Printf("🔍 Fetching direct destination routes from %s on Eurowings...\n\n", originUpper)

			destinations, err := airportsUC.GetEurowingsRoutesFromOrigin(ctx, originUpper)
			if err != nil {
				return fmt.Errorf("failed to get routes from %s: %w", originUpper, err)
			}

			if len(destinations) == 0 {
				fmt.Printf("No direct routes found from %s.\n", originUpper)
				return nil
			}

			fmt.Printf("Found %d direct destination airports from %s:\n", len(destinations), originUpper)
			fmt.Println(strings.Repeat("-", 75))

			for i, d := range destinations {
				fmt.Printf("[%2d] %s -> %s\n", i+1, originUpper, d)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&origin, "origin", "OTP", "Origin 3-letter IATA airport code")

	return cmd
}

func newEurowingsDatesCmd(datesUC *usecase.FlightDatesUseCase, presenter *Presenter) *cobra.Command {
	var (
		origin      string
		destination string
	)

	cmd := &cobra.Command{
		Use:   "dates",
		Short: "List all scheduled flight dates for a route on Eurowings",
		Example: `  travel eurowings dates --origin OTP --destination DUS
  travel eurowings dates --origin DUS --destination BCN
  travel eurowings dates --origin BER --destination PMI`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			if origin == "" || destination == "" {
				return fmt.Errorf("--origin and --destination flags are required")
			}

			originUpper := strings.ToUpper(origin)
			destUpper := strings.ToUpper(destination)

			fmt.Printf("📅 Fetching scheduled flight dates for %s -> %s on Eurowings...\n\n",
				originUpper, destUpper)

			dates, err := datesUC.GetEurowingsDates(ctx, originUpper, destUpper)
			if err != nil {
				return fmt.Errorf("failed to get flight dates: %w", err)
			}

			if len(dates) == 0 {
				fmt.Printf("No scheduled flights found for %s -> %s.\n", originUpper, destUpper)
				return nil
			}

			presenter.PrintDates(originUpper, destUpper, dates)
			return nil
		},
	}

	cmd.Flags().StringVar(&origin, "origin", "OTP", "Origin 3-letter IATA airport code")
	cmd.Flags().StringVar(&destination, "destination", "DUS", "Destination 3-letter IATA airport code")

	return cmd
}
