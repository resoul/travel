package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/resoul/travel/internal/domain"
)

// Presenter handles formatting and printing domain entities to the terminal.
type Presenter struct {
	out io.Writer
}

// NewPresenter creates a new Presenter outputting to stdout by default.
func NewPresenter(w io.Writer) *Presenter {
	if w == nil {
		w = os.Stdout
	}
	return &Presenter{out: w}
}

// PrintFlightOffers prints a formatted list of flight and ground travel offers.
func (p *Presenter) PrintFlightOffers(offers []domain.FlightOffer) {
	if len(offers) == 0 {
		fmt.Fprintln(p.out, "No travel offers found.")
		return
	}

	for _, f := range offers {
		typeTag := ""
		if f.TransportType != "" {
			typeTag = fmt.Sprintf("[%s] ", f.TransportType)
		} else {
			typeTag = "[Flight] "
		}

		timeInfo := f.DepartureRaw
		if f.ArrivalRaw != "" && f.Duration != "" {
			timeInfo = fmt.Sprintf("%s — %s (%s)", f.DepartureRaw, f.ArrivalRaw, f.Duration)
		} else if f.ArrivalRaw != "" {
			timeInfo = fmt.Sprintf("%s — %s", f.DepartureRaw, f.ArrivalRaw)
		}

		availInfo := ""
		if f.SeatsLeft > 0 {
			availInfo = fmt.Sprintf(" | seats: %d", f.SeatsLeft)
		} else if f.Status != "" {
			availInfo = fmt.Sprintf(" | %s", f.Status)
		} else if f.IsAvailable {
			availInfo = " | available"
		}

		flightNumInfo := ""
		if f.FlightNumber != "" {
			flightNumInfo = fmt.Sprintf(" | %s", f.FlightNumber)
		}

		fmt.Fprintf(
			p.out,
			"%s%s: %s -> %s%s | %s | %.2f %s%s\n",
			typeTag,
			f.Airline,
			f.DepartureStation,
			f.ArrivalStation,
			flightNumInfo,
			timeInfo,
			f.Price.Amount,
			f.Price.Currency,
			availInfo,
		)
	}
}

// PrintAirports prints a formatted list of airports.
func (p *Presenter) PrintAirports(airports []domain.Airport) {
	if len(airports) == 0 {
		fmt.Fprintln(p.out, "No airports found.")
		return
	}

	for _, a := range airports {
		fmt.Fprintf(
			p.out,
			"[%s] %s — %s, %s\n",
			a.Code,
			a.Name,
			a.City.Name,
			a.Country.Name,
		)
	}
}

// PrintCities prints a formatted list of cities (e.g. from Wizzair map).
func (p *Presenter) PrintCities(cities []domain.City) {
	if len(cities) == 0 {
		fmt.Fprintln(p.out, "No cities found.")
		return
	}

	for _, c := range cities {
		fmt.Fprintf(p.out, "[%s] %s\n", c.Code, c.Name)
	}
}

// PrintDates prints scheduled flight dates.
func (p *Presenter) PrintDates(origin, destination string, dates []string) {
	if len(dates) == 0 {
		fmt.Fprintf(p.out, "No flights scheduled from %s to %s\n", origin, destination)
		return
	}

	fmt.Fprintf(p.out, "Flights from %s → %s (%d dates):\n", origin, destination, len(dates))
	for _, d := range dates {
		fmt.Fprintf(p.out, "  %s\n", d)
	}
}

// PrintCountries prints a formatted list of countries.
func (p *Presenter) PrintCountries(countries []domain.Country) {
	if len(countries) == 0 {
		fmt.Fprintln(p.out, "No countries found.")
		return
	}

	for _, c := range countries {
		if c.Currency != "" {
			fmt.Fprintf(p.out, "[%s] %s (%s)\n", c.Code, c.Name, c.Currency)
		} else {
			fmt.Fprintf(p.out, "[%s] %s\n", c.Code, c.Name)
		}
	}
}
