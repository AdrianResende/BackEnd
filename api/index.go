package main

import (
	"net/http"

<<<<<<< HEAD
	pkgroutes "smartpicks-backend/pkg/routes"
=======
	"smartpicks-backend/internal/routes"
>>>>>>> 371c546afd72c591e6eaf155f375028286126e15

	"github.com/gorilla/mux"
)

<<<<<<< HEAD
// Vercel serverless entrypoint delegating to internal routes
func Handler(w http.ResponseWriter, r *http.Request) {
	// Build router with the same wiring as local main.go
	router := mux.NewRouter()
	pkgroutes.RegisterRoutes(router)
=======
func Handler(w http.ResponseWriter, r *http.Request) {
	router := mux.NewRouter()
	routes.RegisterRoutes(router)
>>>>>>> 371c546afd72c591e6eaf155f375028286126e15
	router.ServeHTTP(w, r)
}
