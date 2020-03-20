package main

import (
	"os"

	symbol "github.com/vudoan2016/portfolio/input"
	html "github.com/vudoan2016/portfolio/output"
)

func main() {
	symbols := symbol.Get(os.Args[1])
	html.Render(symbols)
}
