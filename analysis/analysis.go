package analysis

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/piquette/finance-go"
	"github.com/piquette/finance-go/equity"
	"github.com/vudoan2016/portfolio/finhub"
	"github.com/vudoan2016/portfolio/input"
	"github.com/vudoan2016/portfolio/models"
	"github.com/vudoan2016/portfolio/output"
)

const (
	mutualFundETF       = "Mutual fund/ETF"
	cash                = "Cash"
	openingBellHour int = 8
	refreshInterval int = 5 // 5-minute
)

// Run polls stock prices from the slice and performs basic analysis
func Run(portfolio input.Portfolio, db *gorm.DB) {
	// Retrieve profile for each company from the database.
	// If the database doesn't have the profile then use the FinHub API to get it.
	profiles := getProfiles(portfolio.Positions, db)

	// consolidate() modifies positions hence passing a copy to analyze() instead
	tmpPortfolio := portfolio
	analyze(&tmpPortfolio, profiles)
	output.Render(tmpPortfolio)

	ticker := time.NewTicker(time.Duration(refreshInterval) * time.Minute)
	for {
		select {
		case <-ticker.C:
			tmpPortfolio := portfolio
			analyze(&tmpPortfolio, profiles)
			output.Render(tmpPortfolio)
		}
	}
}

// Analyze calculates the portfolio's performance
func analyze(portfolio *input.Portfolio, profiles map[string]models.Company) {
	for _, position := range portfolio.Positions {
		getFinancial(position)
		for i := range position {
			populatePosition(&position[i])
		}
	}

	// Update portfolios
	for key, position := range portfolio.Positions {
		for _, holding := range position {
			// Update portfolios
			analyzePortfolio(portfolio, holding)
			// Reflect on sector distribution
			applySectorDistribution(holding, profiles[key.Ticker].P.Sector, portfolio)
		}
	}

	// Run report for sectors
	for i := range portfolio.Reports {
		for sector := range portfolio.Reports[i].Sectors {
			portfolio.Reports[i].Sectors[sector] = 100 * portfolio.Reports[i].Sectors[sector] / portfolio.Reports[i].Value
		}
	}
}

// populatePosition populate value, cost & gain for each position
func populatePosition(pos *input.Position) {
	// Active holding
	if pos.SaleDate == "" {
		pos.Value = pos.RegularMarketPrice * pos.Shares
		pos.Cost = pos.BuyPrice * pos.Shares
		pos.Gain = (pos.RegularMarketPrice - pos.BuyPrice) * pos.Shares
	} else { // Past holding
		pos.Cost = pos.BuyPrice * pos.Shares
		pos.Gain = (pos.SalePrice - pos.BuyPrice) * pos.Shares
	}
}

func applySectorDistribution(pos input.Position, sectorName string, portfolio *input.Portfolio) {
	sector := portfolio.Reports[input.ConvertTypeToVal(pos.Type)].Sectors
	if pos.SaleDate == "" {
		if len(sectorName) > 0 {
			sector[sectorName] += pos.Value
		} else {
			if pos.Ticker != "fidelity" && pos.Ticker != "vanguard" &&
				pos.Ticker != "etrade" && pos.Ticker != "merrill" {
				sector[mutualFundETF] += pos.Value
			} else {
				sector[cash] += pos.Value
			}
		}
	}
}

func analyzePortfolio(portfolio *input.Portfolio, pos input.Position) {
	report := &portfolio.Reports[input.ConvertTypeToVal(pos.Type)]

	report.Cost += pos.Cost
	report.Gain += pos.Gain
	report.Value += pos.Value
	report.TodayGain += calcTodayGain(pos)
	if pos.Ticker == "etrade" || pos.Ticker == "merrill" || pos.Ticker == "vanguard" || pos.Ticker == "fidelity" {
		report.Cash += pos.Value
	}
}

func calcTodayGain(pos input.Position) float64 {
	var gain float64

	if pos.SaleDate == "" && pos.Type != "research" &&
		(pos.QuoteType == finance.QuoteTypeEquity || pos.QuoteType == finance.QuoteTypeETF ||
			// Mutual funds are not updated until around 15:00 PDT on trading days.
			// Todo: weekend & holidays?
			(pos.QuoteType == finance.QuoteTypeMutualFund &&
				pos.MarketState != finance.MarketStateRegular &&
				pos.RegularMarketTime.Hour() > openingBellHour)) {
		gain = pos.RegularMarketChangePercent * pos.RegularMarketPreviousClose * pos.Shares / 100
	}
	return gain
}

func weighEquity(portfolio *input.Portfolio, pos *input.Position) {
	if pos.Type == "deferred" || pos.Type == "taxed" {
		// Average buy price
		pos.BuyPrice = pos.Cost / pos.Shares
		pos.Weight = pos.Value / portfolio.Reports[input.ConvertTypeToVal(pos.Type)].Value * 100
	}
}

func getFinancial(positions []input.Position) {
	equities := make(map[string]*finance.Equity)
	var e *finance.Equity = nil

	start := time.Now()
	for index, pos := range positions {
		if pos.Ticker == "etrade" || pos.Ticker == "merrill" || pos.Ticker == "vanguard" ||
			pos.Ticker == "fidelity" || pos.Ticker == "payflex" {
			positions[index].Name = pos.Ticker
			positions[index].RegularMarketPrice = pos.BuyPrice
		} else {
			var exist bool
			var err error

			// Haven't looked up yet
			if e, exist = equities[pos.Ticker]; !exist {
				e, err = equity.Get(pos.Ticker)
				if err != nil {
					log.Println(pos.Ticker, err)
				} else {
					equities[pos.Ticker] = e
				}
			}
			if exist || err == nil {
				positions[index].Name = e.ShortName
				positions[index].RegularMarketPrice = e.RegularMarketPrice
				positions[index].ForwardPE = e.ForwardPE
				positions[index].ForwardEPS = e.EpsForward
				positions[index].TrailingAnnualDividendYield = e.TrailingAnnualDividendYield
				positions[index].FiftyDayAverage = e.FiftyDayAverage
				positions[index].TwoHundredDayAverage = e.TwoHundredDayAverage
				positions[index].RegularMarketChangePercent = e.RegularMarketChangePercent
				positions[index].QuoteType = e.QuoteType
				positions[index].MarketState = e.MarketState
				positions[index].RegularMarketPreviousClose = e.RegularMarketPreviousClose
				positions[index].EarningsTimestamp = time.Unix(int64(e.EarningsTimestamp), 0).Format("2006/01/02")
				positions[index].RegularMarketTime = time.Unix(int64(e.RegularMarketTime), 0)
				positions[index].FiftyTwoWeekLowChangePercent = e.FiftyTwoWeekLowChangePercent
				positions[index].FiftyTwoWeekHighChangePercent = e.FiftyTwoWeekHighChangePercent
				positions[index].RegularMarketVolume = e.RegularMarketVolume
				positions[index].AverageDailyVolume10Day = e.AverageDailyVolume10Day
				positions[index].AverageDailyVolume3Month = e.AverageDailyVolume3Month
			}
		}
	}
	log.Println("getFinancial() takes", time.Since(start))
}

func getProfiles(positions map[input.PositionKey][]input.Position, db *gorm.DB) map[string]models.Company {
	profiles := make(map[string]models.Company)

	for key := range positions {
		if key.Ticker != "fidelity" && key.Ticker != "vanguard" &&
			key.Ticker != "etrade" && key.Ticker != "merrill" && key.Ticker != "payflex" &&
			key.Ticker != "vinix" && key.Ticker != "sdscx" && key.Ticker != "vig" &&
			key.Ticker != "seegx" && key.Ticker != "sflnx" {
			if _, exist := profiles[key.Ticker]; !exist {
				profiles[key.Ticker] = finhub.GetProfile(key.Ticker, db)
			}
		}
	}
	return profiles
}
