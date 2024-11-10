// controllers/controllers.go
package controllers

import (
	"html/template"
	"log"
	"net/http"
)

// Şablonları cache'lemek için global bir değişken
var templates = template.Must(template.ParseGlob("templates/**/*.html"))

// Kullanıcının oturum açıp açmadığını kontrol eden bir işlev
func isLoggedIn(r *http.Request) bool {
	sessionCookie, err := r.Cookie("session_id") // Örnek çerez adı
	return err == nil && sessionCookie.Value != ""
}

// Logout İşlemi: Kullanıcı oturumunu kapatır
func Logout(w http.ResponseWriter, r *http.Request) {
	// Oturum çerezini iptal etmek için aynı adla ve süresi geçmiş bir çerez gönderiyoruz
	http.SetCookie(w, &http.Cookie{
		Name:   "session_id",
		Value:  "",
		Path:   "/",
		MaxAge: -1, // Süresi geçmiş olarak işaretlemek için negatif bir değer kullanıyoruz
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther) // Oturum kapatıldıktan sonra login sayfasına yönlendir
}

// Middleware: Kullanıcı oturumunu kontrol eder
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isLoggedIn(r) {
			http.Redirect(w, r, "/login", http.StatusSeeOther) // Login sayfasına yönlendir
			return
		}
		next.ServeHTTP(w, r)
	})
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index.html", map[string]interface{}{
		"Title": "Home Page",
	})
}

func Profile(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "profile.html", map[string]interface{}{
		"Title": "Profile",
	})
}

func AddCall(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "add_call.html", map[string]interface{}{
		"Title": "Add Call Page",
	})
}

// LoginPage Kullanıcı giriş sayfasını render eder
func LoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Giriş formundan gelen verileri işleyin
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Burada kullanıcı doğrulama işlemlerini gerçekleştirin
		if username == "admin@admin.com" && password == "123456" { // Örnek kontrol
			http.SetCookie(w, &http.Cookie{
				Name:  "session_id",
				Value: "some_session_value", // Gerçek oturum değeri burada olmalı
				Path:  "/",
			})
			http.Redirect(w, r, "/index", http.StatusSeeOther) // Giriş başarılıysa add_call rotasına yönlendir
			return
		}
		// Giriş başarısızsa hata mesajı göster
		http.Error(w, "Geçersiz kullanıcı adı veya şifre", http.StatusUnauthorized)
		return
	}

	renderTemplate(w, "login.html", map[string]interface{}{
		"Title": "Login Page",
	})
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := templates.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		log.Println("Şablon render hatası:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		log.Println("Şablon başarıyla render edildi:", tmpl)
	}
}
