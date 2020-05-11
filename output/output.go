package output

import (
	"html/template"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/vudoan2016/portfolio/input"
)

type page struct {
	Date      string
	t         *template.Template
	Pretaxes  portfolio
	Posttaxes portfolio
}

type portfolio struct {
	Positions  []input.Position
	Sectors    []sector
	Value      float64 // market value of portfolio
	Gain       float64 // overall gain
	Percentage float64 // gain percentage
	Cash       float64 // cash available
}

type sector struct {
	Name   string
	Weight float64
}

var data page

func Init() {
	t, err := template.ParseFiles("output/layout.html")
	if err != nil {
		log.Println("Failed to parse file", err)
	} else {
		data.t = template.Must(t, err)
		if err != nil {
			log.Println(err)
		}
	}
}

// Render formats the data & writes it to a html file
func Render(p input.Portfolio) {
	data.Date = time.Now().Format("Mon Jan _2 15:04:05 2006")
	data.Pretaxes = portfolio{Value: math.Floor(p.Pretaxes.Value*100) / 100,
		Gain:       math.Floor(p.Pretaxes.Gain*100) / 100,
		Percentage: math.Floor((100*(p.Pretaxes.Gain)/p.Pretaxes.Cost)*100) / 100,
		Cash:       math.Floor(100*p.Pretaxes.Cash) / 100}
	data.Posttaxes = portfolio{Value: math.Floor(p.Posttaxes.Value*100) / 100,
		Gain:       math.Floor((p.Posttaxes.Gain)*100) / 100,
		Percentage: math.Floor((100*(p.Posttaxes.Gain)/p.Posttaxes.Cost)*100) / 100,
		Cash:       math.Floor(100*p.Posttaxes.Cash) / 100}

	for _, pos := range p.Positions {
		if pos.Taxed {
			data.Posttaxes.Positions = append(data.Posttaxes.Positions, pos)
		} else {
			data.Pretaxes.Positions = append(data.Pretaxes.Positions, pos)
		}
	}
	for key, value := range p.Posttaxes.Sectors {
		data.Posttaxes.Sectors = append(data.Posttaxes.Sectors, sector{Name: key, Weight: value})
	}

	for key, value := range p.Pretaxes.Sectors {
		data.Pretaxes.Sectors = append(data.Pretaxes.Sectors, sector{Name: key, Weight: value})
	}
}

func Respond(w http.ResponseWriter, r *http.Request) {
	if data.t != nil {
		err := data.t.Execute(w, data)
		if err != nil {
			log.Println("Executed template with error", err)
		}
	}
}
