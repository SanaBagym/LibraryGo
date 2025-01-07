package internal

import (
	"net/http"
	"strconv"
	"text/template"
)

type TemplateData struct {
	Books  []Book
	Error  string
	FormID string
	Form   Book
}

// RenderBooks renders the main page with the list of books
func RenderBooks(w http.ResponseWriter, r *http.Request, errorMsg string, formBook *Book) {
	var books []Book
	DB.Find(&books)

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}

	data := TemplateData{
		Books: books,
		Error: errorMsg,
	}

	if formBook != nil {
		data.Form = *formBook
		data.FormID = strconv.Itoa(int(formBook.ID))
	}

	tmpl.Execute(w, data)
}

// HandleAddOrUpdate processes adding or updating books
func HandleAddOrUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		id := r.FormValue("id")
		title := r.FormValue("title")
		author := r.FormValue("author")
		year, err := strconv.Atoi(r.FormValue("year"))

		if err != nil || title == "" || author == "" || year == 0 {
			RenderBooks(w, r, "All fields are required", nil)
			return
		}

		if id == "" {
			// Add a new book
			book := Book{Title: title, Author: author, Year: year}
			DB.Create(&book)
		} else {
			// Update an existing book
			var book Book
			if err := DB.First(&book, id).Error; err != nil {
				RenderBooks(w, r, "Book not found", nil)
				return
			}
			book.Title = title
			book.Author = author
			book.Year = year
			DB.Save(&book)
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// HandleDelete processes deleting books by ID
func HandleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		id := r.FormValue("delete-id")
		var book Book
		if err := DB.First(&book, id).Error; err != nil {
			RenderBooks(w, r, "Book not found", nil)
			return
		}
		DB.Delete(&book)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// HandleGetByID processes retrieving a book by ID
func HandleGetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		id := r.FormValue("get-id")
		bookID, err := strconv.Atoi(id)
		if err != nil {
			RenderBooks(w, r, "Invalid Book ID", nil)
			return
		}

		var book Book
		if err := DB.First(&book, bookID).Error; err != nil {
			RenderBooks(w, r, "Book not found", nil)
			return
		}

		// Render page with the retrieved book's data in the form without redirect
		RenderBooks(w, r, "", &book)
	}
}
