package handler

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Version string `json:"version"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// Configurar CORS para o frontend específico
	origin := r.Header.Get("Origin")
	allowedOrigins := []string{
		"https://smartpicks-88709.web.app",
		"https://smartpicks-88709.firebaseapp.com",
		"http://localhost:9000", // Para desenvolvimento local
	}
	
	// Verificar se a origem está permitida
	isAllowed := false
	for _, allowedOrigin := range allowedOrigins {
		if origin == allowedOrigin {
			isAllowed = true
			w.Header().Set("Access-Control-Allow-Origin", origin)
			break
		}
	}
	
	// Se não estiver na lista, permitir qualquer origem (remova isso se quiser mais segurança)
	if !isAllowed && origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Origin, Accept")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Max-Age", "3600")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	response := Response{
		Status:  "online",
		Message: "Backend SmartPicks rodando na Vercel",
		Version: "1.0.0",
	}

	json.NewEncoder(w).Encode(response)
}
