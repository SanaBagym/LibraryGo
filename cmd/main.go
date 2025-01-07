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

	log.Println("Server running on http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
