package cmd

import (
	"fmt"

	"github.com/resoul/travel/internal/wizzair"

	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use: "search",
	RunE: func(cmd *cobra.Command, args []string) error {

		client := wizzair.New()

		buildURL, err := client.FetchBuildURL()
		if err != nil {
			return err
		}

		result, err := client.Search(
			buildURL,
			wizzair.FlightRequest{
				FlightList: []wizzair.Flight{
					{
						DepartureStation: "OTP",
						ArrivalStation:   "CRL",
						From:             "2026-06-01",
						To:               "2026-07-05",
					},
				},
				PriceType:          "regular",
				AdultCount:         1,
				ChildCount:         0,
				InfantCount:        0,
				MacStationsAllowed: false,
			},
		)
		if err != nil {
			return err
		}

		for _, flight := range result.OutboundFlights {
			fmt.Printf(
				"%s -> %s | %.2f %s\n",
				flight.DepartureStation,
				flight.ArrivalStation,
				flight.Price.Amount,
				flight.Price.CurrencyCode,
			)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
