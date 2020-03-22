package output

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"os"

	"github.com/vudoan2016/portfolio/input"
)

type page struct {
	PageTitle string
	Pretaxes  portfolio
	Posttaxes portfolio
}

type portfolio struct {
	Positions  []position
	Value      float64
	Gain       float64
	Percentage float64
}

type position struct {
	Name          string
	Price         float64
	Value         float64
	Weight        float64
	Gain          float64
	Percentage    float64
	Shares        float64
	PurchasePrice float64
	PurchaseDate  string
	SalePrice     float64
	SaleDate      string
}

// Render formats the data to a html file
func Render(p input.Portfolio) {
	output, err := os.Create("portfolio.html")
	if err != nil {
		log.Println(err)
	}
	data := page{
		PageTitle: "Portfolio",
		Pretaxes: portfolio{Value: math.Floor(p.Pretaxes.Value*100) / 100,
			Gain:       math.Floor((p.Pretaxes.Value-p.Pretaxes.Cost)*100) / 100,
			Percentage: math.Floor((100*(p.Pretaxes.Value-p.Pretaxes.Cost)/p.Pretaxes.Value)*100) / 100},
		Posttaxes: portfolio{Value: math.Floor(p.Posttaxes.Value*100) / 100,
			Gain:       math.Floor((p.Posttaxes.Value-p.Posttaxes.Cost)*100) / 100,
			Percentage: math.Floor((100*(p.Posttaxes.Value-p.Posttaxes.Cost)/p.Posttaxes.Value)*100) / 100},
	}
	for _, pos := range p.Positions {
		if pos.Taxed {
			data.Posttaxes.Positions = append(data.Posttaxes.Positions,
				position{pos.Name, pos.Price, math.Floor(pos.Value*100) / 100, math.Floor(pos.Weight*100) / 100,
					math.Floor(pos.Gain*100) / 100, math.Floor(pos.Percentage*100) / 100,
					pos.Shares, pos.BuyPrice, pos.BuyDate, pos.SalePrice, pos.SaleDate})

		} else {
			data.Pretaxes.Positions = append(data.Pretaxes.Positions,
				position{pos.Name, pos.Price, math.Floor(pos.Value*100) / 100, math.Floor(pos.Weight*100) / 100,
					math.Floor(pos.Gain*100) / 100, math.Floor(pos.Percentage*100) / 100,
					pos.Shares, pos.BuyPrice, pos.BuyDate, pos.SalePrice, pos.SaleDate})
		}
	}
	t, er := template.ParseFiles("output/layout.html")
	if er != nil {
		fmt.Println(er)
	} else {
		tmpl := template.Must(t, er)
		tmpl.Execute(output, data)
	}
	output.Close()
}
