package analysis

import (
	"log"
	"sort"
	"time"

	"github.com/piquette/finance-go"
	"github.com/piquette/finance-go/equity"
	"github.com/vudoan2016/portfolio/financialmodelingprep"
	"github.com/vudoan2016/portfolio/input"
)

const (
	mutualFundETF       = "Mutual fund/ETF"
	cash                = "Cash"
	openingBellHour int = 8
)

// Analyze calculates the portfolio's performance
func Analyze(portfolio *input.Portfolio) {
	// Combine into one?
	portfolio.Posttaxes.Sectors = make(map[string]float64)
	portfolio.Pretaxes.Sectors = make(map[string]float64)

	// Todo: cache
	profiles := getProfiles(portfolio.Positions)
	getFinancial(portfolio.Positions)

	for i, pos := range portfolio.Positions {
		// Populate holding's financial data
		populatePosition(&portfolio.Positions[i])

		// Reflect on sector distribution
		reflectSectorDist(portfolio.Positions[i], profiles[pos.Ticker].P.Sector, portfolio)

		// Update portfolios
		updateSubPortfolio(portfolio, portfolio.Positions[i])
	}

	// Summarize posttaxed sectors
	for sector := range portfolio.Posttaxes.Sectors {
		portfolio.Posttaxes.Sectors[sector] = 100 * portfolio.Posttaxes.Sectors[sector] / portfolio.Posttaxes.Value
	}

	// Summarize pretaxed sectors
	for sector := range portfolio.Pretaxes.Sectors {
		portfolio.Pretaxes.Sectors[sector] = 100 * portfolio.Pretaxes.Sectors[sector] / portfolio.Pretaxes.Value
	}

	// Consolidate equity's holdings into one
	portfolio.Positions = consolidate(portfolio.Positions)

	// The sub-portfolios are populated. Weigh the equities.
	for i := range portfolio.Positions {
		weighEquity(portfolio, &portfolio.Positions[i])
	}

	sortPositionsByWeight(portfolio.Positions)
}

// populatePosition populate value, cost & gain for each position
func populatePosition(pos *input.Position) {
	// Active holding
	if pos.SaleDate == "" {
		pos.Value = pos.Price * pos.Shares
		pos.Cost = pos.BuyPrice * pos.Shares
		pos.Gain = (pos.Price - pos.BuyPrice) * pos.Shares
	} else { // Past holding
		pos.Cost = pos.BuyPrice * pos.Shares
		pos.Gain = (pos.SalePrice - pos.BuyPrice) * pos.Shares
	}
}

func reflectSectorDist(pos input.Position, sectorName string, portfolio *input.Portfolio) {
	var sector map[string]float64

	if pos.Type == "taxed" {
		sector = portfolio.Posttaxes.Sectors
	} else if pos.Type == "deferred" {
		sector = portfolio.Pretaxes.Sectors
	} else { // research type
		return
	}

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

func updateSubPortfolio(portfolio *input.Portfolio, pos input.Position) {
	var sub *input.Summary

	if pos.Type == "taxed" {
		sub = &portfolio.Posttaxes
	} else if pos.Type == "deferred" {
		sub = &portfolio.Pretaxes
	} else { // research type
		return
	}

	sub.Cost += pos.Cost
	sub.Gain += pos.Gain
	sub.Value += pos.Value
	sub.TodayGain += calcTodayGain(pos)
	if pos.Ticker == "etrade" || pos.Ticker == "merrill" || pos.Ticker == "vanguard" || pos.Ticker == "fidelity" {
		sub.Cash += pos.Value
	}
}

func calcTodayGain(pos input.Position) float64 {
	var gain float64

	if pos.SaleDate == "" && pos.Type != "research" &&
		(pos.QuoteType == finance.QuoteTypeEquity || pos.QuoteType == finance.QuoteTypeETF ||
			// Mutual funds are not updated until around 15:00 PDT on trading days.
			// Todo: weekend & holidays?
			(pos.QuoteType == finance.QuoteTypeMutualFund &&
				pos.MarketState != finance.MarketStateRegular && pos.RegularMarketTime.Hour() > openingBellHour)) {
		gain = pos.RegularMarketChangePercent * pos.RegularMarketPreviousClose * pos.Shares / 100
	}
	return gain
}

func weighEquity(portfolio *input.Portfolio, pos *input.Position) {
	if pos.Type == "deferred" || pos.Type == "taxed" {
		pos.Percentage = pos.Gain / pos.Cost * 100
		// Average buy price
		pos.BuyPrice = pos.Cost / pos.Shares
		if pos.Type == "taxed" {
			pos.Weight = pos.Value / portfolio.Posttaxes.Value * 100
		} else {
			pos.Weight = pos.Value / portfolio.Pretaxes.Value * 100
		}
	}
}

func getFinancial(positions []input.Position) {
	equities := make(map[string]*finance.Equity)
	var e *finance.Equity = nil

	start := time.Now()
	for index, pos := range positions {
		if pos.Ticker == "etrade" || pos.Ticker == "merrill" || pos.Ticker == "vanguard" || pos.Ticker == "fidelity" || pos.Ticker == "payflex" {
			positions[index].Name = pos.Ticker
			positions[index].Price = pos.BuyPrice
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
				positions[index].Price = e.RegularMarketPrice
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
			}
		}
	}
	log.Println("getFinancial() takes", time.Since(start))
}

func getProfiles(positions []input.Position) map[string]financialmodelingprep.Company {
	profiles := make(map[string]financialmodelingprep.Company)

	start := time.Now()
	for _, pos := range positions {
		if pos.Ticker != "fidelity" && pos.Ticker != "vanguard" &&
			pos.Ticker != "etrade" && pos.Ticker != "merrill" && pos.Ticker != "payflex" &&
			pos.Ticker != "vinix" && pos.Ticker != "sdscx" &&
			pos.Ticker != "seegx" && pos.Ticker != "sflnx" &&
			pos.SaleDate == "" {
			if _, exist := profiles[pos.Ticker]; !exist {
				profiles[pos.Ticker] = financialmodelingprep.GetProfile(pos.Ticker)
			}
		}
	}
	log.Println("getProfiles() takes", time.Since(start))
	return profiles
}

type positions []input.Position

func (c positions) Len() int {
	return len(c)
}

func (c positions) Less(i, j int) bool {
	return c[i].Weight < c[j].Weight
}

func (c positions) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func sortPositionsByWeight(pos []input.Position) {
	sort.Sort(sort.Reverse(positions(pos)))
}

// Combine lots of the same holding
func consolidate(pos positions) positions {
	var consolidated positions

	for _, p := range pos {
		found := false
		for i, c := range consolidated {
			if p.Ticker == c.Ticker && p.Type == c.Type &&
				((p.SaleDate == "" && c.SaleDate == "") || (p.SaleDate != "" && c.SaleDate != "")) {
				consolidated[i].Shares += p.Shares
				consolidated[i].Gain += p.Gain
				consolidated[i].Cost += p.Cost
				consolidated[i].Value += p.Value
				found = true
			}
		}
		if !found {
			consolidated = append(consolidated, p)
		}
	}
	return consolidated
}
