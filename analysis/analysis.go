package analysis

import (
	"log"
	"sort"

	"github.com/piquette/finance-go/equity"
	"github.com/vudoan2016/portfolio/input"
)

func getFinancial(pos *input.Position) {
	if pos.Ticker == "cash" {
		pos.Name = pos.Ticker
		pos.Price = pos.BuyPrice
	} else {
		equity, err := equity.Get(pos.Ticker)

		if err != nil {
			log.Println(pos.Ticker, err)
		} else {
			pos.Name = equity.ShortName
			pos.Price = equity.RegularMarketPrice
		}
	}
}

// Analyze calculates the portfolio's performance
func Analyze(portfolio *input.Portfolio) {
	for i, pos := range portfolio.Positions {
		getFinancial(&pos)
		portfolio.Positions[i].Name = pos.Name
		portfolio.Positions[i].Price = pos.Price
		if pos.SaleDate == "" {
			if pos.Taxed {
				portfolio.Posttaxes.Value += pos.Price * pos.Shares
				portfolio.Posttaxes.Cost += pos.BuyPrice * pos.Shares
			} else {
				portfolio.Pretaxes.Value += pos.Price * pos.Shares
				portfolio.Pretaxes.Cost += pos.BuyPrice * pos.Shares
			}
			portfolio.Positions[i].Gain = (pos.Price - pos.BuyPrice) * pos.Shares
			portfolio.Positions[i].Value = pos.Price * pos.Shares
		} else {
			portfolio.Positions[i].Gain = (pos.SalePrice - pos.BuyPrice) * pos.Shares
		}
	}
	for i, pos := range portfolio.Positions {
		portfolio.Positions[i].Percentage = 100 * pos.Gain / (pos.Shares * pos.BuyPrice)
		if pos.Taxed {
			portfolio.Positions[i].Weight = 100 * pos.Value / portfolio.Posttaxes.Value
		} else {
			portfolio.Positions[i].Weight = 100 * pos.Value / portfolio.Pretaxes.Value
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
