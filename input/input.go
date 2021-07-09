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
	Type   int
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
	Ticker    string  `json:"Symbol"`
	Shares    float64 `json:"shares"`
	BuyDate   string  `json:"buydate"`
	BuyPrice  float64 `json:"purchaseprice"`
	SaleDate  string  `json:"saledate"`
	SalePrice float64 `json:"saleprice"`
	Type      string  `json:"type"`

	// Populate using piquette library
	RegularMarketPrice            float64 `json:"RegularMarketPrice"`
	ForwardPE                     float64
	ForwardEPS                    float64
	TrailingAnnualDividendYield   float64
	FiftyDayAverage               float64
	TwoHundredDayAverage          float64
	RegularMarketChangePercent    float64 `json:"RegularMarketChangePercent"`
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
	Value             float64 `json:"Value"`
	Weight            float64
	Cost              float64 `json:"Cost"`
	Gain              float64 `json:"Gain"`
	EarningsTimestamp string
	RegularMarketTime time.Time
}

type portfolio struct {
	Positions []Position `json:"companies"`
}

func ConvertTypeToVal(t string) int {
	var val int
	switch t {
	case "taxed":
		val = investment
	case "deferred":
		val = deferred
	case "research":
		val = research
	}
	return val
}

func isActive(saleDate string) bool {
	var active bool
	switch saleDate {
	case "":
		active = true
	default:
		active = false
	}
	return active
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
		key := PositionKey{Ticker: pos.Ticker, Type: ConvertTypeToVal(pos.Type), Active: isActive(pos.SaleDate)}
		portfolio.Positions[key] = append(portfolio.Positions[key], pos)

		for i := range portfolio.Positions[key] {
			if pos.BuyDate < portfolio.Positions[key][i].BuyDate {
				copy(portfolio.Positions[key][i+1:], portfolio.Positions[key][i:])
				portfolio.Positions[key][i] = pos
				break
			}
		}
	}
	for i := range portfolio.Reports {
		portfolio.Reports[i].Sectors = make(map[string]float64)
	}

	return portfolio
}
