package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vudoan2016/portfolio/analysis"
	"github.com/vudoan2016/portfolio/input"
	"github.com/vudoan2016/portfolio/models"
	"github.com/vudoan2016/portfolio/output"
)

func main() {
	logger, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Close()
	log.SetOutput(logger)

	var file = "portfolio.json"
	if len(os.Args) < 2 {
		file = find(file)
		if len(file) == 0 {
			log.Fatalln("portfolio not found")
		}
	} else {
		file = os.Args[1]
	}

	db := models.ConnectDataBase()
	defer db.Close()

	router := gin.Default()
	router.LoadHTMLGlob("output/layout.html")

	symbols := input.Get(file)
	analysis.Analyze(&symbols, db)
	output.Init()
	output.Render(symbols)

	router.GET("/", output.Respond)
	router.Run()
}

// Find file in current directory and level-1 subdirectories
func find(name string) string {
	files, err := ioutil.ReadDir(".")
	if err == nil {
		for _, f := range files {
			if f.IsDir() {
				pattern := f.Name() + "/" + name + "*"
				matches, err := filepath.Glob(pattern)
				if err == nil && len(matches) > 0 {
					return "./" + strings.Replace(matches[0], "\\", "/", 1)
				}
			} else if strings.Contains(f.Name(), name) {
				return f.Name()
			}
		}
	}
	return ""
}
