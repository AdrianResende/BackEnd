package routes

import (
	"net/http"

	"smartpicks-backend/internal/database"
	"smartpicks-backend/internal/handlers"

	"github.com/gorilla/mux"
)

// Middleware de CORS
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func RegisterRoutes(r *mux.Router) {
	// Conectar ao banco
	database.Connect()

	// Aplicar middleware de CORS
	r.Use(enableCORS)

	// Rotas da API
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/login", handlers.Login).Methods("POST")
	api.HandleFunc("/register", handlers.Register).Methods("POST")
	api.HandleFunc("/users", handlers.GetAllUsers).Methods("GET")
	api.HandleFunc("/users/check", handlers.CheckUsersTable).Methods("GET")
	api.HandleFunc("/users/permissions", handlers.CheckUserPermissions).Methods("GET")
	api.HandleFunc("/users/profile", handlers.GetUsersByProfile).Methods("GET")

	// Rota de health check
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "API rodando", "version": "1.0.0"}`))
	}).Methods("GET")

	// Rota de health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "healthy"}`))
	}).Methods("GET")
}
