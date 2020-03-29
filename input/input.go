package input

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Portfolio struct {
	Positions []Position `json:"companies"`
	Pretaxes  summary
	Posttaxes summary
}

type summary struct {
	Value float64
	Cost  float64
	Gain  float64
}

type Position struct {
	Ticker                     string  `json:"symbol"`
	Shares                     float64 `json:"shares"`
	BuyDate                    string  `json:"buy date"`
	BuyPrice                   float64 `json:"purchase price"`
	SaleDate                   string  `json:"sale date"`
	SalePrice                  float64 `json:"sale price"`
	Taxed                      bool    `json:taxed`
	Name                       string
	Price                      float64
	Value                      float64
	Weight                     float64
	Cost                       float64
	Gain                       float64
	Percentage                 float64
	ForwardPE                  float64
	ForwardEPS                 float64
	TrailingAnnualDividendRate float64
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
