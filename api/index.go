package handler

import (
	"encoding/json"
	"net/http"
	"sync"

	"smartpicks-backend/internal/database"
	"smartpicks-backend/internal/handlers"

	"github.com/gorilla/mux"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Version string `json:"version"`
}

var (
	router *mux.Router
	once   sync.Once
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		allowedOrigins := []string{
			"https://smartpicks-88709.web.app",
			"https://smartpicks-88709.firebaseapp.com",
			"http://localhost:9000",
		}

		// Verificar se a origem está permitida
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Origin, Accept")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "3600")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// Inicializar o router apenas uma vez
	once.Do(func() {
		router = mux.NewRouter()
		router.Use(enableCORS)

		// Conectar ao banco de dados
		database.Connect()

		// Rotas da API
		api := router.PathPrefix("/api").Subrouter()
		api.HandleFunc("/login", handlers.Login).Methods("POST", "OPTIONS")
		api.HandleFunc("/register", handlers.Register).Methods("POST", "OPTIONS")
		api.HandleFunc("/users", handlers.GetAllUsers).Methods("GET", "OPTIONS")
		api.HandleFunc("/users/permissions", handlers.CheckUserPermissions).Methods("GET", "OPTIONS")
		api.HandleFunc("/users/profile", handlers.GetUsersByProfile).Methods("GET", "OPTIONS")
		api.HandleFunc("/users/avatar", handlers.UpdateAvatar).Methods("POST", "PUT", "OPTIONS")
		api.HandleFunc("/users/avatar", handlers.DeleteAvatar).Methods("DELETE", "OPTIONS")

		// Rotas de status
		router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			response := Response{
				Status:  "online",
				Message: "Backend SmartPicks rodando na Vercel",
				Version: "1.0.0",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}).Methods("GET")

		router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"status": "healthy"}`))
		}).Methods("GET")
	})

	router.ServeHTTP(w, r)
}
