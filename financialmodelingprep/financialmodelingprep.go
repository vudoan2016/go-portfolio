package financialmodelingprep

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/vudoan2016/portfolio/models"
)

const (
	profileURL = "https://financialmodelingprep.com/api/v3/company/profile/"
	ratingURL  = "https://financialmodelingprep.com/api/v3/models.Company/rating/AAPL"
)

// GetProfile returns key data of a models.Company
func GetProfile(symbol string) models.Company {
	var c models.Company

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
