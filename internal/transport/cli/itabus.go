package cli

import (
	"fmt"
	"time"

	"github.com/resoul/travel/internal/domain"
	"github.com/resoul/travel/internal/usecase"
	"github.com/spf13/cobra"
)

func newItaBusCmd(
	searchUC *usecase.SearchFlightsUseCase,
	airportsUC *usecase.ListAirportsUseCase,
	presenter *Presenter,
) *cobra.Command {
	itabusCmd := &cobra.Command{
		Use:   "itabus",
		Short: "ItaBus Italian bus network trip search, stations, and fare calendar",
	}

	// itabus search
	var (
		fromCity string
		toCity   string
		date     string
	)

	searchCmd := &cobra.Command{
		Use:   "search",
		Short: "Search available ItaBus trips on a specific date",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromCity == "" || toCity == "" {
				return fmt.Errorf("both --from and --to city codes are required (e.g. TOR_T, ROM_T)")
			}

			criteria := domain.FlightSearchCriteria{
				Origin:        fromCity,
				Destination:   toCity,
				DepartureDate: date,
			}

			results, err := searchUC.SearchItaBus(cmd.Context(), criteria)
			if err != nil {
				return err
			}

			presenter.PrintFlightOffers(results)
			return nil
		},
	}

	searchCmd.Flags().StringVar(&fromCity, "from", "", "Departure city code (e.g. TOR_T for Turin, ROM_T for Rome)")
	searchCmd.Flags().StringVar(&toCity, "to", "", "Arrival city code (e.g. ROM_T for Rome, MIL_T for Milan)")
	searchCmd.Flags().StringVar(&date, "date", time.Now().Format("2006-01-02"), "Departure date (YYYY-MM-DD)")
	_ = searchCmd.MarkFlagRequired("from")
	_ = searchCmd.MarkFlagRequired("to")

	// itabus stations
	stationsCmd := &cobra.Command{
		Use:   "stations",
		Short: "List all ItaBus cities and stations",
		RunE: func(cmd *cobra.Command, args []string) error {
			stations, err := airportsUC.GetItaBusStations(cmd.Context())
			if err != nil {
				return err
			}

			presenter.PrintAirports(stations)
			return nil
		},
	}

	// itabus calendar [FROM] [TO]
	var startDate, endDate string

	calendarCmd := &cobra.Command{
		Use:   "calendar [FROM] [TO]",
		Short: "Show fare calendar for a route (minimum prices per day)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			from := args[0]
			to := args[1]

			if startDate == "" {
				startDate = time.Now().Format("2006-01-02")
			}
			if endDate == "" {
				endDate = time.Now().AddDate(0, 1, 0).Format("2006-01-02")
			}

			offers, err := airportsUC.GetItaBusFareCalendar(cmd.Context(), from, to, startDate, endDate)
			if err != nil {
				return err
			}

			presenter.PrintFlightOffers(offers)
			return nil
		},
	}

	calendarCmd.Flags().StringVar(&startDate, "start", "", "Start date (YYYY-MM-DD, default: today)")
	calendarCmd.Flags().StringVar(&endDate, "end", "", "End date (YYYY-MM-DD, default: 30 days from today)")

	itabusCmd.AddCommand(searchCmd)
	itabusCmd.AddCommand(stationsCmd)
	itabusCmd.AddCommand(calendarCmd)

	return itabusCmd
}
