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
	Positions  []input.Position
	Value      float64
	Gain       float64
	Percentage float64
}

// Render formats the data & writes it to a html file
func Render(p input.Portfolio) {
	output, err := os.Create("portfolio.html")
	if err != nil {
		log.Println(err)
		return
	}
	data := page{
		PageTitle: "Portfolio",
		Pretaxes: portfolio{Value: math.Floor(p.Pretaxes.Value*100) / 100,
			Gain:       math.Floor(p.Pretaxes.Gain*100) / 100,
			Percentage: math.Floor((100*(p.Pretaxes.Gain)/p.Pretaxes.Cost)*100) / 100},

		Posttaxes: portfolio{Value: math.Floor(p.Posttaxes.Value*100) / 100,
			Gain:       math.Floor((p.Posttaxes.Gain)*100) / 100,
			Percentage: math.Floor((100*(p.Posttaxes.Gain)/p.Posttaxes.Cost)*100) / 100},
	}
	for _, pos := range p.Positions {
		if pos.Taxed {
			data.Posttaxes.Positions = append(data.Posttaxes.Positions, pos)
		} else {
			data.Pretaxes.Positions = append(data.Pretaxes.Positions, pos)
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
