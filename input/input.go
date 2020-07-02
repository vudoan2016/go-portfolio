package input

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/piquette/finance-go"
)

type Portfolio struct {
	Positions []Position `json:"companies"`
	Equities  map[string][]Position
	Pretaxes  Summary
	Posttaxes Summary
}

type Summary struct {
	Value         float64
	PreviousValue float64
	Cost          float64
	Gain          float64
	Cash          float64
	TodayGain     float64
	Industries    map[string]float64
	Sectors       map[string]float64
}

type Position struct {
	Ticker    string  `json:"symbol"`
	Shares    float64 `json:"shares"`
	BuyDate   string  `json:"buydate"`
	BuyPrice  float64 `json:"purchaseprice"`
	SaleDate  string  `json:"saledate"`
	SalePrice float64 `json:"saleprice"`
	Type      string  `json:"type"`

	// Populate using piquette library
	RegularMarketPrice            float64
	ForwardPE                     float64
	ForwardEPS                    float64
	TrailingAnnualDividendYield   float64
	FiftyDayAverage               float64
	TwoHundredDayAverage          float64
	RegularMarketChangePercent    float64
	MarketState                   finance.MarketState
	QuoteType                     finance.QuoteType
	RegularMarketPreviousClose    float64
	FiftyTwoWeekLowChangePercent  float64
	FiftyTwoWeekHighChangePercent float64
	RegularMarketVolume           int
	AverageDailyVolume10Day       int
	AverageDailyVolume3Month      int

	// Analysis fields
	Name              string
	Value             float64
	Weight            float64
	Cost              float64
	Gain              float64
	EarningsTimestamp string
	RegularMarketTime time.Time
}

// Get portfolio from a json file
func Get(fileName string) Portfolio {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		log.Println(fileName, err)
	}
	byteStream, _ := ioutil.ReadAll(jsonFile)

	var portfolio Portfolio
	err = json.Unmarshal(byteStream, &portfolio)
	if err != nil {
		log.Println(err)
	}
	jsonFile.Close()
	for _, pos := range portfolio.Positions {
		portfolio.Equities[pos.Ticker] = append(portfolio.Equities[pos.Ticker], pos)
	}
	return portfolio
}
