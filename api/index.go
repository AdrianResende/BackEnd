package handler

import (
	"net/http"

	pkgroutes "smartpicks-backend/pkg/routes"

	"github.com/gorilla/mux"
)

// Vercel serverless entrypoint delegating to internal routes
func Handler(w http.ResponseWriter, r *http.Request) {
	// Build router with the same wiring as local main.go
	router := mux.NewRouter()
	pkgroutes.RegisterRoutes(router)
	router.ServeHTTP(w, r)
}
