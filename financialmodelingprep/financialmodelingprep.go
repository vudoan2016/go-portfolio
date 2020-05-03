package financialmodelingprep

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	profileURL = "https://financialmodelingprep.com/api/v3/company/profile/"
	ratingURL  = "https://financialmodelingprep.com/api/v3/company/rating/AAPL"
)

type Company struct {
	P profile `json:"profile"`
}

type profile struct {
	Industry string `json:"industry"`
	Sector   string `json:"sector"`
}

// GetProfile returns key data of a company
func GetProfile(symbol string) Company {
	var c Company

	response, err := http.Get(profileURL + symbol)
	if err != nil {
		log.Println(err)
		return c
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return c
	}

	err = json.Unmarshal(responseData, &c)
	if err != nil {
		log.Println("Failed to unmarshal response")
	}

	return c
}
