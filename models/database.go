package models

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "user=postgres password=2231 dbname=urga host=localhost sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Veritabanı bağlantı hatası:", err)
	}

	log.Println("Veritabanına başarıyla bağlanıldı")
	DB = db
}
