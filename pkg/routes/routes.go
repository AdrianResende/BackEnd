package routes

import (
	internalroutes "smartpicks-backend/internal/routes"

	"github.com/gorilla/mux"
)

// RegisterRoutes re-exports internal routes for external consumers (e.g., serverless)
func RegisterRoutes(r *mux.Router) {
	internalroutes.RegisterRoutes(r)
}
