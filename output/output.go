package output

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vudoan2016/portfolio/input"
)

type page struct {
	Date      string
	t         *template.Template
	Positions []input.Position
	Pretaxes  portfolio
	Posttaxes portfolio
	Research  portfolio
}

type portfolio struct {
	Sectors   []sector
	Value     float64 // market value of portfolio
	Gain      float64 // overall gain
	Cash      float64 // cash available
	TodayGain float64
}

type sector struct {
	Name  string
	Value float64
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
		Gain:      p.Pretaxes.Gain,
		Cash:      p.Pretaxes.Cash,
		TodayGain: p.Pretaxes.TodayGain}
	data.Posttaxes = portfolio{Value: p.Posttaxes.Value,
		Gain:      p.Posttaxes.Gain,
		Cash:      p.Posttaxes.Cash,
		TodayGain: p.Posttaxes.TodayGain}
	data.Positions = p.Positions

	for key, value := range p.Posttaxes.Sectors {
		data.Posttaxes.Sectors = append(data.Posttaxes.Sectors, sector{Name: key, Value: value})
	}

	for key, value := range p.Pretaxes.Sectors {
		data.Pretaxes.Sectors = append(data.Pretaxes.Sectors, sector{Name: key, Value: value})
	}
}

func Respond(ctx *gin.Context) {
	// Call the HTML method of the Context to render a template
	ctx.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"layout.html",
		// Pass the data that the page uses
		gin.H{
			"Date":      data.Date,
			"Positions": data.Positions,
			"Pretaxes":  data.Pretaxes,
			"Posttaxes": data.Posttaxes,
			"Research":  data.Research,
		},
	)
}
