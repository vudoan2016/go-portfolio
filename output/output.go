package output

import (
	"html/template"
	"net/http"
	"sort"
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
		consolidated.Shares += p.Shares
		consolidated.Gain += p.Gain
		consolidated.Cost += p.Cost
		consolidated.Value += p.Value
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

type positions []input.Position

func (c positions) Len() int {
	return len(c)
}

func (c positions) Less(i, j int) bool {
	return c[i].RegularMarketChangePercent < c[j].RegularMarketChangePercent
}

func (c positions) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func sortByRegularMarketChangePercent(pos positions) {
	sort.Sort(sort.Reverse(positions(pos)))
}

const (
	none   int = 1
	active int = 2
)

type fn func(pos positions)

// Returns a slice of consolidated (combined lots of the same asset) positions.
func filterPositions(pos map[input.PositionKey][]input.Position, filter int, sortFn fn) positions {
	var positions []input.Position

	switch filter {
	case none:
		for _, p := range pos {
			positions = append(positions, consolidate(p))
		}
	case active:
		for key, p := range pos {
			if key.Active == true {
				positions = append(positions, consolidate(p))
				sortFn(positions)
			}
		}
	}
	return positions
}

// Respond processes '/' route
func Respond(ctx *gin.Context) {
	var positions []input.Position

	switch ctx.Request.Header.Get("Accept") {
	case "application/json":
		positions = filterPositions(data.Positions, active, sortByRegularMarketChangePercent)
		ctx.JSON(http.StatusOK, positions)
	default:
		positions = filterPositions(data.Positions, none, nil)

		// Call the HTML method of the Context to render a template
		ctx.HTML(
			// Set the HTTP status to 200 (OK)
			http.StatusOK,
			// Use the layout.html template
			"layout.html",
			// Pass the data that layout.html uses
			gin.H{
				"Date":      data.Date,
				"Positions": positions,
				"Pretaxes":  data.Reports[deferred],
				"Posttaxes": data.Reports[investment],
				"Research":  data.Reports[research],
			},
		)
	}
}

func RespondEquity(ctx *gin.Context) {
	switch ctx.Request.Header.Get("Accept") {
	case "application/json":
		// Respond with JSON format
		pos := data.Positions[input.PositionKey{Ticker: ctx.Param("id"),
			Type: input.ConvertTypeToVal(ctx.Param("type")), Active: true}]
		ctx.JSON(http.StatusOK, pos)
	default:
		// Respond with HTML format
		ctx.HTML(
			// Set the HTTP status to 200 (OK)
			http.StatusOK,
			// Use the equity.html template
			"equity.html",
			// Pass the data that equity.html uses
			gin.H{
				"Date": data.Date,
				"Equity": data.Positions[input.PositionKey{Ticker: ctx.Param("id"),
					Type: input.ConvertTypeToVal(ctx.Param("type")), Active: true}],
			},
		)
	}
}
