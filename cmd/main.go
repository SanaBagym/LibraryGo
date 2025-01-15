package main

import (
	"log"
	"net/http"

	i "github.com/SanaBagym/KitapSana/internal"

	"github.com/gorilla/mux"
)

func main() {
	i.ConnectDatabase()

	r := mux.NewRouter()
	i.SetupRoutes(r)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
