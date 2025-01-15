package internal

import (
	"log"
	"net/http"
	"strconv"
	"sync"
	"text/template"
	"time"
)

var deleteRateLimit = make(map[string]time.Time)
var getRateLimit = make(map[string]time.Time)
var rateLimitMutex sync.Mutex

// Helper to render templates
func renderTemplate(w http.ResponseWriter, filename string, data any) {
	tmpl, err := template.ParseFiles("templates/" + filename)
	if err != nil {
		log.Printf("[ERROR] Failed to render template %s: %v", filename, err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

// Home Page
func RenderHome(w http.ResponseWriter, r *http.Request) {
	var books []Book
	DB.Find(&books)
	renderTemplate(w, "all_books.html", books)
}

// Add Book Page
func RenderAddPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "add.html", nil)
}

func HandleAdd(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	author := r.FormValue("author")
	year, _ := strconv.Atoi(r.FormValue("year"))

	if title == "" || author == "" || year == 0 {
		log.Printf("[ERROR] Invalid form input: title='%s', author='%s', year='%d'", title, author, year)
		http.Redirect(w, r, "/add?error=All fields are required", http.StatusSeeOther)
		return
	}

	book := Book{Title: title, Author: author, Year: year}
	DB.Create(&book)
	log.Printf("[INFO] Book added: %+v", book)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Update Book Page
func RenderUpdatePage(w http.ResponseWriter, r *http.Request) {
	errorMessage := r.URL.Query().Get("error")
	successMessage := r.URL.Query().Get("success")

	renderTemplate(w, "update.html", map[string]interface{}{
		"Error":   errorMessage,
		"Success": successMessage,
	})
}

func HandleUpdate(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	title := r.FormValue("title")
	author := r.FormValue("author")
	year, _ := strconv.Atoi(r.FormValue("year"))

	var book Book
	if err := DB.First(&book, id).Error; err != nil {
		log.Printf("[ERROR] Book not found with ID %s: %v", id, err)
		http.Redirect(w, r, "/update?error=Book not found with ID: "+id, http.StatusSeeOther)
		return
	}

	book.Title = title
	book.Author = author
	book.Year = year
	DB.Save(&book)
	log.Printf("[INFO] Book updated: %+v", book)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Delete Book Page
func RenderDeletePage(w http.ResponseWriter, r *http.Request) {
	errorMessage := r.URL.Query().Get("error")
	successMessage := r.URL.Query().Get("success")

	renderTemplate(w, "delete.html", map[string]interface{}{
		"Error":   errorMessage,
		"Success": successMessage,
	})
}

func HandleDelete(w http.ResponseWriter, r *http.Request) {
	userIP := r.Header.Get("X-Forwarded-For")
	if userIP == "" {
		userIP = r.RemoteAddr
	}

	id := r.FormValue("id")

	rateLimitMutex.Lock()
	defer rateLimitMutex.Unlock()

	lastActionTime, exists := deleteRateLimit[userIP]
	if exists && time.Since(lastActionTime) < 1*time.Second {
		log.Printf("[ERROR] Rate limit exceeded for IP %s", userIP)
		http.Redirect(w, r, "/delete?error=You can delete only one book every second.", http.StatusSeeOther)
		return
	}

	deleteRateLimit[userIP] = time.Now()

	var book Book
	if err := DB.First(&book, id).Error; err != nil {
		log.Printf("[ERROR] Book not found with ID %s: %v", id, err)
		http.Redirect(w, r, "/delete?error=Book not found with ID: "+id, http.StatusSeeOther)
		return
	}

	DB.Delete(&book)
	log.Printf("[INFO] Book deleted: %+v", book)
	http.Redirect(w, r, "/delete?success=Book with ID "+id+" successfully deleted", http.StatusSeeOther)
}

// Get Book Page
func RenderGetPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "get.html", nil)
}

func HandleGet(w http.ResponseWriter, r *http.Request) {
	userIP := r.RemoteAddr
	id := r.FormValue("id")

	rateLimitMutex.Lock()
	lastActionTime, exists := getRateLimit[userIP]
	if exists && time.Since(lastActionTime) < 5*time.Second {
		rateLimitMutex.Unlock()
		log.Printf("[ERROR] Rate limit exceeded for IP %s", userIP)
		renderTemplate(w, "get.html", map[string]interface{}{
			"Error": "You can get only one book every 5 seconds.",
		})
		return
	}
	getRateLimit[userIP] = time.Now()
	rateLimitMutex.Unlock()

	var book Book
	if err := DB.First(&book, id).Error; err != nil {
		log.Printf("[ERROR] Book not found with ID %s: %v", id, err)
		renderTemplate(w, "get.html", map[string]interface{}{
			"Error": "Book not found with ID: " + id,
		})
		return
	}

	log.Printf("[INFO] Book fetched: %+v", book)
	renderTemplate(w, "get.html", map[string]interface{}{
		"Book": book,
	})
}

// Get All Books Page
func RenderAllBooksPage(w http.ResponseWriter, r *http.Request) {
	sortParam := r.URL.Query().Get("sort")
	var books []Book

	switch sortParam {
	case "id":
		DB.Order("id").Find(&books)
	case "title":
		DB.Order("title").Find(&books)
	case "author":
		DB.Order("author").Find(&books)
	case "year":
		DB.Order("year").Find(&books)
	default:
		DB.Find(&books)
	}

	log.Printf("[INFO] All books fetched with sort param: %s", sortParam)
	renderTemplate(w, "all_books.html", books)
}
