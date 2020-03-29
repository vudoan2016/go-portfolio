package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/vudoan2016/portfolio/analysis"
	"github.com/vudoan2016/portfolio/input"
	"github.com/vudoan2016/portfolio/output"
)

func find(name string) string {
	dirs, err := ioutil.ReadDir(".")
	if err == nil {
		for _, dir := range dirs {
			pattern := dir.Name() + "/" + name
			// Windows use '\' hence returns pattern
			matches, err := filepath.Glob(pattern)
			if err == nil && len(matches) > 0 {
				return "./" + pattern
			}
		}
	}
	return ""
}

func main() {
	var file = "portfolio.json"
	if len(os.Args) < 2 {
		file = find(file)
		if len(file) == 0 {
			log.Fatalln("portfolio not found")
		}
	} else {
		file = os.Args[1]
	}
	symbols := input.Get(file)
	analysis.Analyze(&symbols)
	output.Render(symbols)
}
