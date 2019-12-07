package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func NewDB(mirDB string) (db *gorm.DB) {
	db, err := gorm.Open("sqlite3", mirDB)
	if err != nil {
		panic("failed to connect database")
	}
	return db
}
