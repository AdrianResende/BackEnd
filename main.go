package main

import (
	"log"
	"net/http"
	"os"

	"smartpicks-backend/internal/routes"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Carregar variáveis do arquivo .env (apenas para desenvolvimento local)
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  Arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}

	r := mux.NewRouter()
	routes.RegisterRoutes(r)

	port := getEnv("PORT", "8080")

	log.Printf("Servidor rodando na porta %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
