// controllers/controllers.go
package controllers

import (
	"Urga/models"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

func NewCall(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "new_call.html", map[string]interface{}{
		"Title": "Profile",
	})
}

func AddCall(w http.ResponseWriter, r *http.Request) {
	// Eğer kullanıcı oturum açmamışsa
	if !isLoggedIn(r) {
		http.Error(w, "Oturum açılmadı. Lütfen oturum açın.", http.StatusUnauthorized)
		http.Redirect(w, r, "/login", http.StatusSeeOther)

		return
	}

	if r.Method == http.MethodPost {
		// Kullanıcı bilgisini oturumdan alın (örneğin session id'den)
		sessionCookie, err := r.Cookie("session_id")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		addedBy := "Anonymous" // Eğer kullanıcı yoksa varsayılan değer
		if sessionCookie != nil {
			addedBy = "admin@admin.com" // Örnek kullanıcı adı
		}

		// Form verilerini al
		title := r.FormValue("vulnerabilityTitle")
		url := r.FormValue("vulnerabilityUrl")
		description := r.FormValue("vulnerabilityDescription")
		applicationName := r.FormValue("applicationName")
		vendor := r.FormValue("vendor")
		applicationCategory := r.FormValue("applicationCategory")
		vulnerabilityCategory := r.FormValue("vulnerabilityCategory")
		severityLevel := r.FormValue("severityLevel")
		tags := r.FormValue("vulnerabilityTags")
		cvssScore, _ := strconv.ParseFloat(r.FormValue("cvssScore"), 64)

		// Yeni zafiyet oluştur
		vulnerability := models.VulnerabilityAdd{
			Title:                 title,
			Url:                   url,
			Description:           description,
			ApplicationName:       applicationName,
			Vendor:                vendor,
			ApplicationCategory:   applicationCategory,
			VulnerabilityCategory: vulnerabilityCategory,
			SeverityLevel:         severityLevel,
			CvssScore:             cvssScore,
			Tags:                  tags,
			AddedBy:               addedBy, // Ekleyen kişi bilgisi
		}

		// Veritabanına kaydet
		result := models.DB.Create(&vulnerability)
		if result.Error != nil {
			http.Error(w, "Zafiyet eklenirken bir hata oluştu.", http.StatusInternalServerError)
			return
		}

		// Başarılı ekleme sonrası yönlendirme
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Eğer GET isteği ise formu render et
	renderTemplate(w, "add_call.html", map[string]interface{}{
		"Title": "Add Vulnerability",
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

func UploadImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Resim dosyasını al
	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Resim dosyası alınamadı.", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Dosya adı üzerinde XSS koruması sağlama
	imageName := strings.ReplaceAll(handler.Filename, "<", "")
	imageName = strings.ReplaceAll(imageName, ">", "")
	imageName = strings.ReplaceAll(imageName, "&", "")
	imageName = strings.ReplaceAll(imageName, "\"", "")
	imageName = strings.ReplaceAll(imageName, "'", "")
	imageName = strings.ReplaceAll(imageName, "/", "")

	// Dosya dizinini ve adını belirleyin
	mediaDirectory := "./uploads/images"
	os.MkdirAll(mediaDirectory, os.ModePerm) // Eğer dizin yoksa oluştur
	filePath := filepath.Join(mediaDirectory, imageName)

	// Dosyayı belirli bir dizine kaydet
	f, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Dosya oluşturulamadı.", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// Dosya verilerini yaz
	_, err = io.Copy(f, file)
	if err != nil {
		http.Error(w, "Dosya kopyalanırken bir hata oluştu.", http.StatusInternalServerError)
		return
	}

	// Başarıyla kaydedilen resmin adını geri döndür
	response := map[string]string{
		"success":   "true",
		"imageName": imageName,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
