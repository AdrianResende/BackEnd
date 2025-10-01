package main

import (
	"log"
	"net/http"

	_ "smartpicks-backend/docs" // Import para registrar documentação Swagger
	"smartpicks-backend/internal/routes"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	routes.RegisterRoutes(r)

	log.Println("Servidor rodando na porta 8080")
	log.Println("Documentação Swagger disponível em: http://localhost:8080/swagger/")
	log.Fatal(http.ListenAndServe(":8080", r))
}
