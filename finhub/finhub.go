package finhub

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/vudoan2016/portfolio/models"
)

const (
	finHubProfileURL = "https://finnhub.io/api/v1/stock/profile2?symbol="
	tokenID          = "br30ufvrh5re69qjuu80"
)

func GetProfile(symbol string, db *gorm.DB) models.Company {
	var c models.Company
	var err error

	c.P, err = models.FindProfile(db, symbol)
	if err != nil {
		client := http.Client{}
		request, err := http.NewRequest("GET", finHubProfileURL+symbol+"&token="+tokenID, nil)
		if err != nil {
			fmt.Println(err)
		}

		resp, err := client.Do(request)
		if err != nil {
			fmt.Println(err)
		}

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		if len(result) > 0 {
			c.P = models.Profile{Symbol: symbol, Sector: result["finnhubIndustry"].(string)}
			models.AddProfile(db, c.P)
		}
	}
	return c
}
