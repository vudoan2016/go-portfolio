package output

import (
	"html/template"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vudoan2016/portfolio/input"
)

const (
	investment int = 0
	deferred   int = 1
	research   int = 2
	reportSize int = 3
)

type page struct {
	Date      string
	t         *template.Template
	Reports   [reportSize]report
	Positions map[input.PositionKey][]input.Position
}

type report struct {
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

// Combine lots of the same holding
func consolidate(pos []input.Position) input.Position {
	consolidated := pos[0]

	for _, p := range pos[1:] {
		if (p.SaleDate == "" && consolidated.SaleDate == "") || (p.SaleDate != "" && consolidated.SaleDate != "") {
			consolidated.Shares += p.Shares
			consolidated.Gain += p.Gain
			consolidated.Cost += p.Cost
			consolidated.Value += p.Value
		}
	}
	return consolidated
}

// Render formats the data & writes it to a html file
func Render(p input.Portfolio) {
	data.Date = time.Now().Format("Mon Jan 2 15:04:05 2006")
	for i := 0; i < reportSize-1; i++ {
		data.Reports[i] = report{Value: p.Reports[i].Value,
			Gain:      p.Reports[i].Gain,
			Cash:      p.Reports[i].Cash,
			TodayGain: p.Reports[i].TodayGain}
		for key, value := range p.Reports[i].Sectors {
			data.Reports[i].Sectors = append(data.Reports[i].Sectors, sector{Name: key, Value: value})
		}
	}
	data.Positions = p.Positions
}

func Respond(ctx *gin.Context) {
	var positions []input.Position
	for _, position := range data.Positions {
		positions = append(positions, consolidate(position))
	}
	// Call the HTML method of the Context to render a template
	ctx.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the layout.html template
		"layout.html",
		// Pass the data that the page uses
		gin.H{
			"Date":      data.Date,
			"Positions": positions,
			"Pretaxes":  data.Reports[deferred],
			"Posttaxes": data.Reports[investment],
			"Research":  data.Reports[research],
		},
	)
}

func RespondEquity(ctx *gin.Context) {
	ctx.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the layout.html template
		"equity.html",
		// Pass the data that the page uses
		gin.H{
			"Date": data.Date,
			"Equity": data.Positions[input.PositionKey{Ticker: ctx.Param("id"),
				Type: input.ConvertTypeToVal(ctx.Param("type")), Active: true}],
		},
	)
}
