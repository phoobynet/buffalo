package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var db *gorm.DB

func openDB() {
	d, err := gorm.Open(sqlite.Open("buffalo.db"), &gorm.Config{})

	if err != nil {
		log.Fatal("failed to connect database")
	}

	err = d.AutoMigrate(&AppConfiguration{})

	if err != nil {
		log.Fatal(err)
	}

	db = d
}

type AppConfiguration struct {
	gorm.Model
	Key    string
	X      int
	Y      int
	Width  int
	Height int
}
