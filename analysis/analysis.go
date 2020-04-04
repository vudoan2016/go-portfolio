package analysis

import (
	"log"
	"math"
	"sort"

	"github.com/piquette/finance-go/equity"
	"github.com/vudoan2016/portfolio/input"
)

func getFinancial(pos *input.Position) {
	if pos.Ticker == "etrade" || pos.Ticker == "merrill" || pos.Ticker == "vanguard" || pos.Ticker == "fidelity" {
		pos.Name = pos.Ticker
		pos.Price = pos.BuyPrice
	} else {
		equity, err := equity.Get(pos.Ticker)

		if err != nil {
			log.Println(pos.Ticker, err)
		} else {
			pos.Name = equity.ShortName
			pos.Price = equity.RegularMarketPrice
			pos.ForwardPE = equity.ForwardPE
			pos.ForwardEPS = equity.EpsForward
			pos.TrailingAnnualDividendRate = equity.TrailingAnnualDividendRate
		}
	}
}

// Analyze calculates the portfolio's performance
func Analyze(portfolio *input.Portfolio) {
	for i, pos := range portfolio.Positions {
		getFinancial(&portfolio.Positions[i])

		if pos.SaleDate == "" {
			portfolio.Positions[i].Value = portfolio.Positions[i].Price * portfolio.Positions[i].Shares
			portfolio.Positions[i].Cost = portfolio.Positions[i].BuyPrice * portfolio.Positions[i].Shares
			portfolio.Positions[i].Gain = (portfolio.Positions[i].Price - portfolio.Positions[i].BuyPrice) * portfolio.Positions[i].Shares
			if pos.Taxed {
				portfolio.Posttaxes.Value += portfolio.Positions[i].Value
			} else {
				portfolio.Pretaxes.Value += portfolio.Positions[i].Value
			}
		} else {
			portfolio.Positions[i].Cost = portfolio.Positions[i].BuyPrice * portfolio.Positions[i].Shares
			portfolio.Positions[i].Gain = (portfolio.Positions[i].SalePrice - portfolio.Positions[i].BuyPrice) * portfolio.Positions[i].Shares
		}
		if portfolio.Positions[i].Taxed {
			portfolio.Posttaxes.Cost += portfolio.Positions[i].Cost
			portfolio.Posttaxes.Gain += portfolio.Positions[i].Gain
			if pos.Ticker == "etrade" || pos.Ticker == "merrill" || pos.Ticker == "vanguard" || pos.Ticker == "fidelity" {
				portfolio.Posttaxes.Cash += portfolio.Positions[i].Value
			}
		} else {
			portfolio.Pretaxes.Cost += portfolio.Positions[i].Cost
			portfolio.Pretaxes.Gain += portfolio.Positions[i].Gain
			if pos.Ticker == "etrade" || pos.Ticker == "merrill" || pos.Ticker == "vanguard" || pos.Ticker == "fidelity" {
				portfolio.Pretaxes.Cash += portfolio.Positions[i].Value
			}
		}
	}

	//portfolio.Positions = consolidate(portfolio.Positions)

	for i, pos := range portfolio.Positions {
		portfolio.Positions[i].Percentage = math.Floor(100*pos.Gain/pos.Cost*100) / 100
		// 2 trailing digits
		portfolio.Positions[i].ForwardPE = math.Floor(100*portfolio.Positions[i].ForwardPE) / 100
		portfolio.Positions[i].Value = math.Floor(100*portfolio.Positions[i].Value) / 100
		portfolio.Positions[i].Gain = math.Floor(100*portfolio.Positions[i].Gain) / 100
		portfolio.Positions[i].Shares = math.Floor(100*portfolio.Positions[i].Shares) / 100

		if pos.Taxed {
			portfolio.Positions[i].Weight = math.Floor(100*pos.Value/portfolio.Posttaxes.Value*100) / 100
		} else {
			portfolio.Positions[i].Weight = math.Floor(100*pos.Value/portfolio.Pretaxes.Value*100) / 100
		}
	}
	sortPositionsByWeight(portfolio.Positions)
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
			if p.Ticker == c.Ticker && p.Taxed == c.Taxed &&
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
