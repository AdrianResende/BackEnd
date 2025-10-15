package main

import (
	"net/http"

	"smartpicks-backend/internal/routes"

	"github.com/gorilla/mux"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	router := mux.NewRouter()
	routes.RegisterRoutes(router)
	router.ServeHTTP(w, r)
}
