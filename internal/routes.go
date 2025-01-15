package internal

import (
	"github.com/gorilla/mux"
)

func SetupRoutes(r *mux.Router) {
	r.HandleFunc("/", RenderHome).Methods("GET")

	r.HandleFunc("/add", RenderAddPage).Methods("GET")
	r.HandleFunc("/add", HandleAdd).Methods("POST")

	r.HandleFunc("/update", RenderUpdatePage).Methods("GET")
	r.HandleFunc("/update", HandleUpdate).Methods("POST")

	r.HandleFunc("/delete", RenderDeletePage).Methods("GET")
	r.HandleFunc("/delete", HandleDelete).Methods("POST")

	r.HandleFunc("/get", RenderGetPage).Methods("GET")
	r.HandleFunc("/get", HandleGet).Methods("POST")

	r.HandleFunc("/books", RenderAllBooksPage).Methods("GET")
}
