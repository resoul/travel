package cmd

import (
	"fmt"

	"github.com/resoul/travel/internal/ryanair"

	"github.com/spf13/cobra"
)

var ryanairCmd = &cobra.Command{
	Use:   "ryanair",
	Short: "Search Ryanair flights",
}

var ryanairSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search available Ryanair fares",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ryanair.New()

		results, err := client.Search(ryanair.FlightRequest{
			Origin:      "BBU",
			Destination: "GRO",
			DateOut:     "2026-06-24",
			Adults:      1,
			RoundTrip:   false,
		})
		if err != nil {
			return err
		}

		if len(results) == 0 {
			fmt.Println("No flights found.")
			return nil
		}

		for _, f := range results {
			fmt.Printf(
				"%s -> %s | %s | %s — %s (%s) | %.2f %s | seats: %d\n",
				f.DepartureStation,
				f.ArrivalStation,
				f.FlightNumber,
				f.DepartureLocal,
				f.ArrivalLocal,
				f.Duration,
				f.Price,
				f.Currency,
				f.FaresLeft,
			)
		}

		return nil
	},
}

var ryanairAirportsCmd = &cobra.Command{
	Use:   "airports",
	Short: "List all active Ryanair airports",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ryanair.New()

		airports, err := client.FetchAirports()
		if err != nil {
			return err
		}

		for _, a := range airports {
			fmt.Printf("[%s] %s — %s, %s\n",
				a.Code, a.Name, a.City.Name, a.Country.Name,
			)
		}

		return nil
	},
}

var ryanairRoutesCmd = &cobra.Command{
	Use:   "routes [IATA]",
	Short: "List airports reachable from a given airport",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ryanair.New()

		airports, err := client.FetchRoutes(args[0])
		if err != nil {
			return err
		}

		if len(airports) == 0 {
			fmt.Printf("No routes found from %s\n", args[0])
			return nil
		}

		for _, a := range airports {
			fmt.Printf("[%s] %s — %s, %s\n",
				a.Code, a.Name, a.City.Name, a.Country.Name,
			)
		}

		return nil
	},
}

var ryanairAvailabilitiesCmd = &cobra.Command{
	Use:   "dates [ORIGIN] [DESTINATION]",
	Short: "List scheduled flight dates between two airports",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := ryanair.New()

		dates, err := client.FetchAvailabilities(args[0], args[1])
		if err != nil {
			return err
		}

		if len(dates) == 0 {
			fmt.Printf("No flights scheduled from %s to %s\n", args[0], args[1])
			return nil
		}

		fmt.Printf("Flights from %s → %s (%d dates):\n", args[0], args[1], len(dates))
		for _, d := range dates {
			fmt.Println(" ", d)
		}

		return nil
	},
}

func init() {
	ryanairCmd.AddCommand(ryanairSearchCmd)
	ryanairCmd.AddCommand(ryanairAirportsCmd)
	ryanairCmd.AddCommand(ryanairRoutesCmd)
	ryanairCmd.AddCommand(ryanairAvailabilitiesCmd)
	rootCmd.AddCommand(ryanairCmd)
}
