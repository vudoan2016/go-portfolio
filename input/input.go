package input

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/piquette/finance-go"
)

const (
	investment int = 0
	deferred   int = 1
	research   int = 2
	reportSize int = 3
)

type Portfolio struct {
	Positions map[PositionKey][]Position
	Reports   [reportSize]Report
}

type PositionKey struct {
	Ticker string
	Type   string
	Active bool
}

type Report struct {
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
	TaxType           int
	Active            bool
	Value             float64
	Weight            float64
	Cost              float64
	Gain              float64
	EarningsTimestamp string
	RegularMarketTime time.Time
}

type portfolio struct {
	Positions []Position `json:"companies"`
}

// Get portfolio from a json file
func Get(fileName string) Portfolio {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		log.Println(fileName, err)
	}
	defer jsonFile.Close()

	var p portfolio
	byteStream, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteStream, &p)
	if err != nil {
		log.Println(err)
	}

	var portfolio Portfolio
	portfolio.Positions = make(map[PositionKey][]Position)
	for _, pos := range p.Positions {
		switch pos.Type {
		case "taxed":
			pos.TaxType = investment
		case "deferred":
			pos.TaxType = deferred
		case "research":
			pos.TaxType = research
		}
		switch pos.SaleDate {
		case "":
			pos.Active = true
		default:
			pos.Active = false
		}
		key := PositionKey{Ticker: pos.Ticker, Type: pos.Type, Active: pos.Active}
		portfolio.Positions[key] = append(portfolio.Positions[key], pos)
	}
	for i := range portfolio.Reports {
		portfolio.Reports[i].Sectors = make(map[string]float64)
	}

	return portfolio
}
