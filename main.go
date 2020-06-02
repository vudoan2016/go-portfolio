package main

import (
	"fmt"
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
	// Initialize the logger
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
			fmt.Println(file, "not found")
			os.Exit(1)
		}
	} else {
		file = os.Args[1]
	}

	// Initialize the database
	db := models.ConnectDataBase()
	defer db.Close()

	// Initialize the router
	router := gin.Default()
	router.LoadHTMLGlob("output/*.html")

	router.Use(func(ctx *gin.Context) {
		ctx.Set("db", db)
		ctx.Next()
	})

	// Load portfolio data
	portfolio := input.Get(file)

	// Initialize the output module to render template
	output.Init()

	// Poll stock prices & perform simple analysis
	go analysis.Run(&portfolio, db)

	// Ready to serve
	router.GET("/", output.Respond)
	router.GET("/:id", output.RespondOne)
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
