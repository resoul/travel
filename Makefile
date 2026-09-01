.DEFAULT_GOAL := help

BINARY_NAME ?= travel
BIN_DIR     ?= bin
BIN_PATH    := $(BIN_DIR)/$(BINARY_NAME)

# Colors
CYAN  := \033[36m
GREEN := \033[32m
RESET := \033[0m

.PHONY: help
help: ## Display this help screen
	@echo ""
	@echo "$(CYAN)Travel CLI$(RESET) — Flight & Bus Search Utility"
	@echo ""
	@echo "Usage: make $(GREEN)<target>$(RESET) [OPTIONS]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-18s$(RESET) %s\n", $$1, $$2}'
	@echo ""

.PHONY: build
build: ## Compile the travel CLI binary into ./bin/travel
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_PATH) ./cmd
	@echo "Built binary: $(BIN_PATH)"

.PHONY: install
install: ## Build and install binary into Go bin directory
	@go install ./cmd
	@echo "Installed binary to GOPATH/bin"

.PHONY: clean
clean: ## Remove build artifacts and temporary binaries
	@rm -rf $(BIN_DIR)
	@echo "Cleaned build artifacts"

.PHONY: vet
vet: ## Run Go static analysis
	@go vet ./...

.PHONY: test
test: ## Run unit and integration tests
	@go test -v ./...

.PHONY: run
run: ## Run travel CLI with custom ARGS (e.g. make run ARGS="--help")
	@go run ./cmd $(ARGS)

# --- Cache Management ---

.PHONY: cache-info
cache-info: ## Show cache directory status and size
	@go run ./cmd cache info

.PHONY: cache-clear
cache-clear: ## Clear all cached provider responses
	@go run ./cmd cache clear

# --- Quick Operator Shortcuts ---

.PHONY: ryanair-search
ryanair-search: ## Search Ryanair (e.g. make ryanair-search ORIGIN=BBU DEST=GRO DATE=2026-08-22)
	@go run ./cmd ryanair search --origin $(or $(ORIGIN),BBU) --destination $(or $(DEST),GRO) --date $(or $(DATE),2026-08-22)

.PHONY: volotea-search
volotea-search: ## Search Volotea (e.g. make volotea-search ORIGIN=NTE DEST=VAR DATE=2026-08-25)
	@go run ./cmd volotea search --origin $(or $(ORIGIN),NTE) --destination $(or $(DEST),VAR) --date $(or $(DATE),2026-08-25)

.PHONY: vueling-search
vueling-search: ## Search Vueling (e.g. make vueling-search ORIGIN=BCN DEST=FCO DATE=2026-08-19)
	@go run ./cmd vueling search --origin $(or $(ORIGIN),BCN) --destination $(or $(DEST),FCO) --date $(or $(DATE),2026-08-19)

.PHONY: airbaltic-search
airbaltic-search: ## Search airBaltic (e.g. make airbaltic-search ORIGIN=ALC DEST=RIX DATE=2026-08-28)
	@go run ./cmd airbaltic search --origin $(or $(ORIGIN),ALC) --destination $(or $(DEST),RIX) --date $(or $(DATE),2026-08-28)

.PHONY: flyone-search
flyone-search: ## Search FlyOne (e.g. make flyone-search ORIGIN=OTP DEST=BRU DATE=2026-08-28)
	@go run ./cmd flyone search --origin $(or $(ORIGIN),OTP) --destination $(or $(DEST),BRU) --date $(or $(DATE),2026-08-28)

.PHONY: flixbus-search
flixbus-search: ## Search FlixBus (e.g. make flixbus-search FROM="Bucharest" TO="Brasov" DATE=2026-08-27)
	@go run ./cmd flixbus search --from "$(or $(FROM),Bucharest)" --to "$(or $(TO),Brasov)" --date $(or $(DATE),2026-08-27)

.PHONY: movacar-offers
movacar-offers: ## Search Movacar 1-euro car relocations (e.g. make movacar-offers FROM=Berlin TO=Reutlingen)
	@go run ./cmd movacar search $(if $(FROM),--from "$(FROM)",) $(if $(TO),--to "$(TO)",) $(if $(DATE),--date $(DATE),)

.PHONY: imoova-offers
imoova-offers: ## Search imoova 1-dollar campervan relocations (e.g. make imoova-offers FROM=Vancouver TO="San Francisco")
	@go run ./cmd imoova search $(if $(FROM),--from "$(FROM)",) $(if $(TO),--to "$(TO)",) $(if $(DATE),--date $(DATE),)

.PHONY: driiveme-offers
driiveme-offers: ## Search DriiveMe 1-euro car relocations (e.g. make driiveme-offers FROM=Barcelona TO=Torremolinos)
	@go run ./cmd driiveme search $(if $(FROM),--from "$(FROM)",) $(if $(TO),--to "$(TO)",) $(if $(DATE),--date $(DATE),) $(if $(EMAIL),--email "$(EMAIL)",) $(if $(PASSWORD),--password "$(PASSWORD)",)

.PHONY: driiveme-cities
driiveme-cities: ## Search DriiveMe cities (e.g. make driiveme-cities QUERY=London)
	@go run ./cmd driiveme cities $(if $(QUERY),--query "$(QUERY)",)

.PHONY: driiveme-login
driiveme-login: ## Login to DriiveMe (e.g. make driiveme-login EMAIL=user@example.com PASSWORD=pass)
	@go run ./cmd driiveme login $(if $(EMAIL),--email "$(EMAIL)",) $(if $(PASSWORD),--password "$(PASSWORD)",)

.PHONY: indigo-radar
indigo-radar: ## Get IndiGo lowest fare recommendations from origin (e.g. make indigo-radar ORIGIN=DEL)
	@go run ./cmd indigo radar $(or $(ORIGIN),DEL)

.PHONY: indigo-calendar
indigo-calendar: ## Search IndiGo date-by-date fare calendar (e.g. make indigo-calendar ORIGIN=DXB DEST=DEL CURRENCY=AED)
	@go run ./cmd indigo calendar --origin $(or $(ORIGIN),DXB) --destination $(or $(DEST),DEL) $(if $(CURRENCY),--currency $(CURRENCY),) $(if $(START_DATE),--start-date $(START_DATE),) $(if $(END_DATE),--end-date $(END_DATE),)

.PHONY: indigo-dates
indigo-dates: ## List scheduled IndiGo flight dates (e.g. make indigo-dates ORIGIN=DXB DEST=DEL)
	@go run ./cmd indigo dates $(or $(ORIGIN),DXB) $(or $(DEST),DEL)

.PHONY: tap-airports
tap-airports: ## List all TAP Air Portugal origin airports (e.g. make tap-airports)
	@go run ./cmd tap airports

.PHONY: tap-routes
tap-routes: ## List TAP destination airports from origin (e.g. make tap-routes ORIGIN=LIS)
	@go run ./cmd tap routes $(or $(ORIGIN),LIS)

.PHONY: tap-calendar
tap-calendar: ## Search TAP date-by-date fare calendar (e.g. make tap-calendar ORIGIN=LIS DEST=BCN MONTH=9 YEAR=2026)
	@go run ./cmd tap calendar --origin $(or $(ORIGIN),LIS) --destination $(or $(DEST),BCN) $(if $(MONTH),--month $(MONTH),) $(if $(YEAR),--year $(YEAR),) $(if $(MARKET),--market $(MARKET),)

.PHONY: tap-dates
tap-dates: ## List scheduled TAP flight dates (e.g. make tap-dates ORIGIN=LIS DEST=BCN)
	@go run ./cmd tap dates $(or $(ORIGIN),LIS) $(or $(DEST),BCN)

.PHONY: cruise-search
cruise-search: ## Search cruises (e.g. make cruise-search DESTINATION=54 MONTH=11 YEAR=2026 LIMIT=10)
	@go run ./cmd cruise search $(if $(DESTINATION),--destination $(DESTINATION),) $(if $(LINE),--cruise-line $(LINE),) $(if $(MONTH),--month $(MONTH),) $(if $(YEAR),--year $(YEAR),) $(if $(LIMIT),--limit $(LIMIT),)

.PHONY: cruise-lines
cruise-lines: ## List all cruise operators and matrix IDs (e.g. make cruise-lines)
	@go run ./cmd cruise lines

.PHONY: cruise-destinations
cruise-destinations: ## List all cruise destination regions and matrix IDs (e.g. make cruise-destinations)
	@go run ./cmd cruise destinations

.PHONY: agoda-search
agoda-search: ## Search Agoda hotel deals (e.g. make agoda-search CITY_ID=19216 CITY=Mamaia CHECKIN=2026-09-06 CHECKOUT=2026-09-08)
	@go run ./cmd agoda search $(if $(CITY_ID),--city-id $(CITY_ID),) $(if $(CITY),--city $(CITY),) $(if $(CHECKIN),--check-in $(CHECKIN),) $(if $(CHECKOUT),--check-out $(CHECKOUT),) $(if $(ADULTS),--adults $(ADULTS),) $(if $(CURRENCY),--currency $(CURRENCY),) $(if $(LIMIT),--limit $(LIMIT),)

.PHONY: agoda-countries
agoda-countries: ## List world countries directory from Agoda CDN (e.g. make agoda-countries LANG=11)
	@go run ./cmd agoda countries $(if $(LANG),--lang $(LANG),)

.PHONY: trip-search
trip-search: ## Search Trip.com hotel deals (e.g. make trip-search CITY_ID=40795 CITY=Barcelona CHECKIN=2026-11-19 CHECKOUT=2026-11-23)
	@go run ./cmd trip search $(if $(CITY_ID),--city-id $(CITY_ID),) $(if $(CITY),--city $(CITY),) $(if $(CHECKIN),--check-in $(CHECKIN),) $(if $(CHECKOUT),--check-out $(CHECKOUT),) $(if $(ADULTS),--adults $(ADULTS),) $(if $(CURRENCY),--currency $(CURRENCY),) $(if $(LIMIT),--limit $(LIMIT),)

.PHONY: trip-details
trip-details: ## View full room breakdown and amenities for a Trip.com hotel (e.g. make trip-details HOTEL=134130947 CHECKIN=2026-11-19 CHECKOUT=2026-11-23)
	@go run ./cmd trip details $(or $(HOTEL),134130947) $(if $(CHECKIN),--check-in $(CHECKIN),) $(if $(CHECKOUT),--check-out $(CHECKOUT),) $(if $(ADULTS),--adults $(ADULTS),) $(if $(CURRENCY),--currency $(CURRENCY),)

.PHONY: trip-cars
trip-cars: ## Search Trip.com rental cars (e.g. make trip-cars COUNTRY_ID=63 CITY_ID=39050 CITY=Otopeni CODE=OTP PICKUP=2026-08-29 RETURN=2026-09-01)
	@go run ./cmd trip cars $(if $(COUNTRY_ID),--country-id $(COUNTRY_ID),) $(if $(CITY_ID),--city-id $(CITY_ID),) $(if $(CITY),--city $(CITY),) $(if $(CODE),--code $(CODE),) $(if $(ADDRESS),--address "$(ADDRESS)",) $(if $(PICKUP),--pickup "$(PICKUP)",) $(if $(RETURN),--return "$(RETURN)",) $(if $(AGE),--age $(AGE),) $(if $(CURRENCY),--currency $(CURRENCY),) $(if $(LIMIT),--limit $(LIMIT),)

.PHONY: tictactrip-cities
tictactrip-cities: ## Search and autocomplete European cities and stations on Tictactrip (e.g. make tictactrip-cities QUERY=bar)
	@go run ./cmd tictactrip cities $(if $(QUERY),--query $(QUERY),)

.PHONY: tictactrip-popular
tictactrip-popular: ## List popular destinations on Tictactrip (e.g. make tictactrip-popular FROM=barcelone LIMIT=7)
	@go run ./cmd tictactrip popular $(if $(FROM),--from $(FROM),) $(if $(LIMIT),--limit $(LIMIT),)

.PHONY: tictactrip-calendar
tictactrip-calendar: ## View monthly lowest train and bus fare calendar on Tictactrip (e.g. make tictactrip-calendar ORIGIN_ID=76 DEST_ID=542 MONTH=2026-11)
	@go run ./cmd tictactrip calendar $(if $(ORIGIN_ID),--origin-id $(ORIGIN_ID),) $(if $(DEST_ID),--dest-id $(DEST_ID),) $(if $(MONTH),--month $(MONTH),)

.PHONY: trenitalia-search
trenitalia-search: ## Search Trenitalia & Le Frecce trains (e.g. make trenitalia-search ORIGIN="Roma Termini" DEST="Milano Centrale" DATE=2026-09-10)
	@go run ./cmd trenitalia search $(if $(ORIGIN),--origin "$(ORIGIN)",) $(if $(DEST),--destination "$(DEST)",) $(if $(ORIGIN_ID),--origin-id $(ORIGIN_ID),) $(if $(DEST_ID),--dest-id $(DEST_ID),) $(if $(DATE),--date $(DATE),) $(if $(TIME),--time $(TIME),) $(if $(ADULTS),--adults $(ADULTS),) $(if $(FRECCE_ONLY),--frecce-only,) $(if $(NO_CHANGES),--no-changes,) $(if $(LIMIT),--limit $(LIMIT),)

.PHONY: trenitalia-stations
trenitalia-stations: ## Search Trenitalia Italian stations catalog (e.g. make trenitalia-stations QUERY=roma)
	@go run ./cmd trenitalia stations $(if $(QUERY),--query $(QUERY),) $(if $(FRECCE),--frecce,)

.PHONY: norwegian-calendar
norwegian-calendar: ## View Norwegian Air Shuttle low fare calendar (e.g. make norwegian-calendar ORIGIN=OSL DEST=BCN YEAR=2026 MONTH=9 CURRENCY=EUR)
	@go run ./cmd norwegian calendar $(if $(ORIGIN),--origin $(ORIGIN),) $(if $(DEST),--destination $(DEST),) $(if $(YEAR),--year $(YEAR),) $(if $(MONTH),--month $(MONTH),) $(if $(CURRENCY),--currency $(CURRENCY),)

.PHONY: obilet-search
obilet-search: ## Search oBilet Turkish and regional bus tickets (e.g. make obilet-search ORIGIN=istanbul DEST=ankara DATE=2026-09-10)
	@go run ./cmd obilet search $(if $(ORIGIN),--origin $(ORIGIN),) $(if $(DEST),--destination $(DEST),) $(if $(ORIGIN_ID),--origin-id $(ORIGIN_ID),) $(if $(DEST_ID),--dest-id $(DEST_ID),) $(if $(DATE),--date $(DATE),) $(if $(LIMIT),--limit $(LIMIT),)

.PHONY: eurowings-airports
eurowings-airports: ## List Eurowings airports (e.g. make eurowings-airports COUNTRY=DE)
	@go run ./cmd eurowings airports $(if $(COUNTRY),--country $(COUNTRY),)

.PHONY: eurowings-routes
eurowings-routes: ## List direct destinations from origin on Eurowings (e.g. make eurowings-routes ORIGIN=OTP)
	@go run ./cmd eurowings routes $(if $(ORIGIN),--origin $(ORIGIN),)

.PHONY: eurowings-dates
eurowings-dates: ## List scheduled flight dates on Eurowings (e.g. make eurowings-dates ORIGIN=OTP DEST=DUS)
	@go run ./cmd eurowings dates $(if $(ORIGIN),--origin $(ORIGIN),) $(if $(DEST),--destination $(DEST),)

.PHONY: transavia-calendar
transavia-calendar: ## View Transavia low-fare calendar (e.g. make transavia-calendar ORIGIN=AMS DEST=BCN YEAR=2026 MONTH=9 ADULTS=1)
	@go run ./cmd transavia calendar $(if $(ORIGIN),--origin $(ORIGIN),) $(if $(DEST),--destination $(DEST),) $(if $(YEAR),--year $(YEAR),) $(if $(MONTH),--month $(MONTH),) $(if $(ADULTS),--adults $(ADULTS),)

.PHONY: pitchup-search
pitchup-search: ## Search Pitchup campsites and glamping (e.g. make pitchup-search COUNTRY=france ARRIVE=2026-09-10 DEPART=2026-09-12 ADULTS=2)
	@go run ./cmd pitchup search $(if $(COUNTRY),--country $(COUNTRY),) $(if $(REGION),--region $(REGION),) $(if $(ARRIVE),--arrive $(ARRIVE),) $(if $(DEPART),--depart $(DEPART),) $(if $(ADULTS),--adults $(ADULTS),) $(if $(LIMIT),--limit $(LIMIT),)

.PHONY: hipcamp-search
hipcamp-search: ## Search Hipcamp outdoor spots and glamping (e.g. make hipcamp-search COUNTRY=united-states REGION=california)
	@go run ./cmd hipcamp search $(if $(COUNTRY),--country $(COUNTRY),) $(if $(REGION),--region $(REGION),) $(if $(LIMIT),--limit $(LIMIT),)

.PHONY: campspace-search
campspace-search: ## Search Campspace micro-camping spots (e.g. make campspace-search CATEGORY=tent-pitches)
	@go run ./cmd campspace search $(if $(CATEGORY),--category $(CATEGORY),) $(if $(LIMIT),--limit $(LIMIT),)

.PHONY: sata-airports
sata-airports: ## List Azores Airlines / SATA network airports
	@go run ./cmd sata airports

.PHONY: sata-routes
sata-routes: ## List direct destination airport codes on Azores Airlines / SATA (e.g. make sata-routes ORIGIN=PDL)
	@go run ./cmd sata routes $(if $(ORIGIN),--origin $(ORIGIN),)

.PHONY: sata-calendar
sata-calendar: ## View Azores Airlines / SATA 365-day low-fare calendar (e.g. make sata-calendar ORIGIN=LIS DEST=PDL LIMIT=15)
	@go run ./cmd sata calendar $(if $(ORIGIN),--origin $(ORIGIN),) $(if $(DEST),--destination $(DEST),) $(if $(LIMIT),--limit $(LIMIT),)




