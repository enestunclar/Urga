package main

import (
	"Urga/controllers"
	"log"
	"net/http"
)

func main() {
	// Statik dosyaları servis etme
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/logout", controllers.Logout) // Login sayfası için bir handler ekleyin

	// Login sayfası için route
	http.HandleFunc("/login", controllers.LoginPage) // Login sayfası için bir handler ekleyin

	// Ana sayfayı ve diğer rotaları koruma altına alıyoruz
	http.HandleFunc("/", controllers.AuthMiddleware(controllers.HomePage))        // Auth middleware ile korunan rota
	http.HandleFunc("/add_call", controllers.AuthMiddleware(controllers.AddCall)) // Auth middleware ile korunan rota
	http.HandleFunc("/profile", controllers.AuthMiddleware(controllers.Profile))  // Auth middleware ile korunan rota

	log.Println("Sunucu başlatıldı: http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Sunucu başlatılamadı:", err)
	}
}
