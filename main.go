package main

import (
	"log"
	"net/http"
	"os"

	_ "smartpicks-backend/docs"
	"smartpicks-backend/internal/routes"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	routes.RegisterRoutes(r)

	port := getEnv("PORT", "8080")

	log.Printf("Servidor rodando na porta %s", port)
	log.Printf("Documentação Swagger disponível em: http://localhost:%s/swagger/", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
