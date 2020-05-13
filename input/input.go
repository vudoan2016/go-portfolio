package input

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/piquette/finance-go"
)

type Portfolio struct {
	Positions []Position `json:"companies"`
	Pretaxes  Summary
	Posttaxes Summary
}

type Summary struct {
	Value      float64
	Cost       float64
	Gain       float64
	Cash       float64
	TodayGain  float64
	Industries map[string]float64
	Sectors    map[string]float64
}

type Position struct {
	Ticker                      string  `json:"symbol"`
	Shares                      float64 `json:"shares"`
	BuyDate                     string  `json:"buydate"`
	BuyPrice                    float64 `json:"purchaseprice"`
	SaleDate                    string  `json:"saledate"`
	SalePrice                   float64 `json:"saleprice"`
	Type                        string  `json:"type"`
	Name                        string
	Price                       float64
	Value                       float64
	Weight                      float64
	Cost                        float64
	Gain                        float64
	Percentage                  float64
	ForwardPE                   float64
	ForwardEPS                  float64
	TrailingAnnualDividendYield float64
	FiftyDayAverage             float64
	TwoHundredDayAverage        float64
	RegularMarketChangePercent  float64
	EarningsTimestamp           string
	MarketState                 finance.MarketState
	QuoteType                   finance.QuoteType
}

// Get portfolio from a json file
func Get(fileName string) Portfolio {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		log.Println(fileName, err)
	}
	byteStream, _ := ioutil.ReadAll(jsonFile)

	var p Portfolio
	err = json.Unmarshal(byteStream, &p)
	if err != nil {
		log.Println(err)
	}

	jsonFile.Close()
	return p
}
