package finhub

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/vudoan2016/portfolio/models"
)

const (
	finHubProfileURL    = "https://finnhub.io/api/v1/stock/profile2?symbol="
	finHubTokenID       = "br30ufvrh5re69qjuu80"
	alphavantageURL     = "https://www.alphavantage.co/query?function=OVERVIEW&symbol="
	alphavantageTokenID = "MWFQJJUPRYPAQ5WX"
)

func GetProfile(symbol string, db *gorm.DB) (models.Company, error) {
	var c models.Company
	var err error

	c.P, err = models.FindProfile(db, symbol)
	if err != nil {
		client := http.Client{}
		request, err := http.NewRequest("GET", finHubProfileURL+symbol+"&token="+finHubTokenID, nil)
		//request, err := http.NewRequest("GET", alphavantageURL+symbol+"&apikey="+alphavantageTokenID, nil) // 5 APIs per minute
		if err != nil {
			fmt.Println(err)
		}

		resp, err := client.Do(request)
		if err != nil {
			fmt.Println(err)
		}

		if resp.StatusCode == 429 {
			log.Println(resp)
			return c, errors.New("API exceeds limit")
		}

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		if result["finnhubIndustry"] != nil {
			c.P = models.Profile{Symbol: symbol, Sector: result["finnhubIndustry"].(string)}
			models.AddProfile(db, c.P)
		} else {
			log.Println(symbol, result)
		}
	}
	return c, nil
}
