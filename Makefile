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




