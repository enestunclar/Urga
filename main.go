package main

import (
	"Urga/controllers"
	"Urga/models"
	"log"
	"net/http"
)

func main() {
	// Veritabanına bağlan
	models.ConnectDatabase()

	// GORM ile veritabanı tablolarını otomatik oluşturma
	err := models.DB.AutoMigrate(&models.VulnerabilityAdd{}, &models.VulnerabilityMedia{})
	if err != nil {
		log.Fatal("Migration hatası:", err)
	}

	// Statik dosyaları servis etme
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Yüklenen dosyaların servis edilmesi
	uploads := http.FileServer(http.Dir("./uploads"))
	http.Handle("/uploads/", http.StripPrefix("/uploads/", uploads))

	// Yeni rotayı ekle
	http.HandleFunc("/upload_image", controllers.UploadImage)

	http.HandleFunc("/logout", controllers.Logout)
	http.HandleFunc("/login", controllers.LoginPage)

	// Ana sayfayı ve diğer rotaları koruma altına alıyoruz
	http.HandleFunc("/", controllers.AuthMiddleware(controllers.HomePage))
	http.HandleFunc("/add_call", controllers.AuthMiddleware(controllers.AddCall))
	http.HandleFunc("/profile", controllers.AuthMiddleware(controllers.Profile))

	log.Println("Sunucu başlatıldı: http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Sunucu başlatılamadı:", err)
	}
}
