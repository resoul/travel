package cli

import (
	"fmt"
	"os"

	"github.com/resoul/travel/internal/domain"
	"github.com/resoul/travel/internal/usecase"
	"github.com/spf13/cobra"
)

func newDriiveMeCmd(
	searchUC *usecase.SearchFlightsUseCase,
	airportsUC *usecase.ListAirportsUseCase,
	presenter *Presenter,
) *cobra.Command {
	driivemeCmd := &cobra.Command{
		Use:   "driiveme",
		Short: "DriiveMe 1-euro car and van relocation search commands",
	}

	// driiveme login
	var (
		loginEmail    string
		loginPassword string
	)

	loginCmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate with DriiveMe using email and password",
		RunE: func(cmd *cobra.Command, args []string) error {
			if loginEmail == "" {
				loginEmail = os.Getenv("DRIIVEME_EMAIL")
			}
			if loginPassword == "" {
				loginPassword = os.Getenv("DRIIVEME_PASSWORD")
			}

			if loginEmail == "" || loginPassword == "" {
				return fmt.Errorf("email and password are required (via flags --email / --password or env DRIIVEME_EMAIL / DRIIVEME_PASSWORD)")
			}

			if err := searchUC.LoginDriiveMe(cmd.Context(), loginEmail, loginPassword); err != nil {
				return err
			}

			fmt.Fprintf(os.Stdout, "Successfully logged in to DriiveMe as %s\n", loginEmail)
			return nil
		},
	}

	loginCmd.Flags().StringVarP(&loginEmail, "email", "e", "", "DriiveMe account email")
	loginCmd.Flags().StringVarP(&loginPassword, "password", "p", "", "DriiveMe account password")

	// driiveme search
	var (
		from     string
		to       string
		date     string
		optEmail string
		optPass  string
	)

	searchCmd := &cobra.Command{
		Use:     "search",
		Aliases: []string{"offers"},
		Short:   "Search available 1-euro car relocation offers",
		RunE: func(cmd *cobra.Command, args []string) error {
			if optEmail == "" {
				optEmail = os.Getenv("DRIIVEME_EMAIL")
			}
			if optPass == "" {
				optPass = os.Getenv("DRIIVEME_PASSWORD")
			}

			if optEmail != "" && optPass != "" {
				_ = searchUC.LoginDriiveMe(cmd.Context(), optEmail, optPass)
			}

			criteria := domain.FlightSearchCriteria{
				Origin:        from,
				Destination:   to,
				DepartureDate: date,
			}

			results, err := searchUC.SearchDriiveMe(cmd.Context(), criteria)
			if err != nil {
				return err
			}

			presenter.PrintFlightOffers(results)
			return nil
		},
	}

	searchCmd.Flags().StringVar(&from, "from", "", "Pickup city name or ID (optional)")
	searchCmd.Flags().StringVar(&to, "to", "", "Dropoff city name or ID (optional)")
	searchCmd.Flags().StringVar(&date, "date", "", "Pickup min date (YYYY-MM-DD, optional)")
	searchCmd.Flags().StringVarP(&optEmail, "email", "e", "", "DriiveMe account email for enriched details (optional)")
	searchCmd.Flags().StringVarP(&optPass, "password", "p", "", "DriiveMe account password (optional)")

	// driiveme cities
	var cityQuery string

	citiesCmd := &cobra.Command{
		Use:     "cities",
		Aliases: []string{"locations"},
		Short:   "Search DriiveMe cities by query term",
		RunE: func(cmd *cobra.Command, args []string) error {
			if cityQuery == "" && len(args) > 0 {
				cityQuery = args[0]
			}
			if cityQuery == "" {
				cityQuery = "a"
			}

			cities, err := airportsUC.GetDriiveMeCities(cmd.Context(), cityQuery)
			if err != nil {
				return err
			}

			presenter.PrintAirports(cities)
			return nil
		},
	}

	citiesCmd.Flags().StringVarP(&cityQuery, "query", "q", "", "City name query (e.g. London, Barcelona, Paris)")

	// driiveme availabilities
	availCmd := &cobra.Command{
		Use:     "availabilities <transport_id>",
		Aliases: []string{"slots", "dates"},
		Short:   "Get available booking timeslots for a specific transport ID",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			transportID := args[0]
			slots, err := searchUC.GetDriiveMeAvailabilities(cmd.Context(), transportID)
			if err != nil {
				return err
			}

			fmt.Fprintf(os.Stdout, "Available booking slots for transport #%s (%d slots):\n", transportID, len(slots))
			for _, slot := range slots {
				fmt.Fprintf(os.Stdout, "  - %s\n", slot)
			}
			return nil
		},
	}

	driivemeCmd.AddCommand(loginCmd)
	driivemeCmd.AddCommand(searchCmd)
	driivemeCmd.AddCommand(citiesCmd)
	driivemeCmd.AddCommand(availCmd)

	return driivemeCmd
}
