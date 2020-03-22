package main

import (
	"os"

	"github.com/vudoan2016/portfolio/analysis"
	"github.com/vudoan2016/portfolio/input"
	"github.com/vudoan2016/portfolio/output"
)

func main() {
	symbols := input.Get(os.Args[1])
	analysis.Analyze(&symbols)
	output.Render(symbols)
}
