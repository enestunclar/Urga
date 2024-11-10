package models

import (
	"database/sql"
	_ "github.com/lib/pq" // PostgreSQL sürücüsü
	"log"
)

var DB *sql.DB

// Veritabanı bağlantısını kur
func ConnectDatabase() {
	connStr := "user=diastek password=Diastek_0 dbname=urga host=localhost sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Veritabanı bağlantı hatası:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Veritabanına erişilemiyor:", err)
	}

	log.Println("Veritabanına başarıyla bağlanıldı")
	DB = db
}
