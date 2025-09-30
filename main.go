package main

import (
	"log"
	"net/http"

	"smartpicks-backend/internal/routes"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	routes.RegisterRoutes(r)

	log.Println("Servidor rodando na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
