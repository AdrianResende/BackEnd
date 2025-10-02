package routes

import (
	"net/http"

	"smartpicks-backend/internal/database"
	"smartpicks-backend/internal/handlers"

	"github.com/gorilla/mux"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// CORS simplificado para debug
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:9000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Origin")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Responder a requisições OPTIONS (preflight)
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func RegisterRoutes(r *mux.Router) {
	database.Connect()

	// Aplicar CORS globalmente
	r.Use(enableCORS)

	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/login", handlers.Login).Methods("POST", "OPTIONS")
	api.HandleFunc("/register", handlers.Register).Methods("POST", "OPTIONS")
	api.HandleFunc("/users", handlers.GetAllUsers).Methods("GET", "OPTIONS")
	api.HandleFunc("/users/check", handlers.CheckUsersTable).Methods("GET", "OPTIONS")
	api.HandleFunc("/users/permissions", handlers.CheckUserPermissions).Methods("GET", "OPTIONS")
	api.HandleFunc("/users/profile", handlers.GetUsersByProfile).Methods("GET", "OPTIONS")

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "API rodando", "version": "1.0.0"}`))
	}).Methods("GET")

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "healthy"}`))
	}).Methods("GET")

	r.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		http.ServeFile(w, r, "./docs/swagger.json")
	})

	r.HandleFunc("/swagger/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		http.ServeFile(w, r, "./swagger.html")
	})

	r.HandleFunc("/swagger/index.html", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		http.ServeFile(w, r, "./swagger.html")
	})

	// Rota para testar CORS
	r.HandleFunc("/test-cors", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		http.ServeFile(w, r, "./test-cors.html")
	})

	// Handler de teste simples
	r.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "test ok", "method": "` + r.Method + `"}`))
	})
}
