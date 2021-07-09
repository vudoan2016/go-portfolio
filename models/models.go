package models

import (
	"log"

	"github.com/jinzhu/gorm"
)

type Company struct {
	P Profile `json:"profile"`
}

type Profile struct {
	Symbol   string
	Industry string `json:"industry"`
	Sector   string `json:"sector"`
}

func ConnectDataBase() *gorm.DB {
	db, err := gorm.Open("sqlite3", "./models.db")

	if err != nil {
		log.Println("Failed to connect to database!")
	}

	db.AutoMigrate(&Profile{})

	// FinHub is unable to pull these profiles
	AddProfile(db, Profile{Symbol: "kl", Sector: "Basic Materials"})
	AddProfile(db, Profile{Symbol: "brk-b", Sector: "Financial Services"})
	AddProfile(db, Profile{Symbol: "tm", Sector: "Consumer Cyclical"})

	return db
}

func AddProfile(db *gorm.DB, prof Profile) {
	db.Create(&prof)
}

func FindProfile(db *gorm.DB, symbol string) (Profile, error) {
	var prof Profile
	err := db.Where("Symbol = ?", symbol).First(&prof).Error
	return prof, err
}
