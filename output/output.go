package output

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/vudoan2016/portfolio/input"
)

type page struct {
	Date      string
	t         *template.Template
	Pretaxes  portfolio
	Posttaxes portfolio
	Research  portfolio
}

type portfolio struct {
	Positions  []input.Position
	Sectors    []sector
	Value      float64 // market value of portfolio
	Gain       float64 // overall gain
	Percentage float64 // gain percentage
	Cash       float64 // cash available
	TodayGain  float64
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
	data.Date = time.Now().Format("Mon Jan 2 15:04:05 2006")
	data.Pretaxes = portfolio{Value: p.Pretaxes.Value,
		Gain:       p.Pretaxes.Gain,
		Percentage: p.Pretaxes.Gain / p.Pretaxes.Cost,
		Cash:       p.Pretaxes.Cash,
		TodayGain:  p.Pretaxes.TodayGain}
	data.Posttaxes = portfolio{Value: p.Posttaxes.Value,
		Gain:       p.Posttaxes.Gain,
		Percentage: p.Posttaxes.Gain / p.Posttaxes.Cost,
		Cash:       p.Posttaxes.Cash,
		TodayGain:  p.Posttaxes.TodayGain}

	for _, pos := range p.Positions {
		if pos.Type == "taxed" {
			data.Posttaxes.Positions = append(data.Posttaxes.Positions, pos)
		} else if pos.Type == "deferred" {
			data.Pretaxes.Positions = append(data.Pretaxes.Positions, pos)
		} else if pos.Type == "research" {
			data.Research.Positions = append(data.Research.Positions, pos)
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
