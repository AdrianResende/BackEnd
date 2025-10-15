package main

import (
	"log"
	"net/http"
	"os"

	"smartpicks-backend/internal/routes"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	routes.RegisterRoutes(r)

	// Rotas de documentação Swagger (apenas em desenvolvimento)
	r.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		http.ServeFile(w, r, "./docs/swagger.json")
	}).Methods("GET")

	r.HandleFunc("/swagger/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		http.ServeFile(w, r, "./swagger.html")
	}).Methods("GET")

	port := getEnv("PORT", "8080")

	log.Printf("Servidor rodando na porta %s", port)
	log.Printf("Swagger disponível em: http://localhost:%s/swagger/", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
