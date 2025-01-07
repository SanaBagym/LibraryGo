package internal

import (
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRoutes(r *mux.Router) {
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		RenderBooks(w, r, "", nil)
	}).Methods("GET")

	r.HandleFunc("/books", HandleAddOrUpdate).Methods("POST")
	r.HandleFunc("/books/delete", HandleDelete).Methods("POST")
	r.HandleFunc("/books/get", HandleGetByID).Methods("POST") // Ensure this route exists
}
