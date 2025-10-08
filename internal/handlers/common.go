package handlers

import (
	"encoding/json"
	"net/http"

	"smartpicks-backend/internal/database"
)

// sendErrorResponse envia uma resposta de erro padronizada
func sendErrorResponse(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

// sendSuccessResponse envia uma resposta de sucesso padronizada
func sendSuccessResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// userExists verifica se um usuário existe baseado em um campo específico
func userExists(field, value string) bool {
	var count int
	query := "SELECT COUNT(*) FROM users WHERE " + field + "=?"
	database.DB.QueryRow(query, value).Scan(&count)
	return count > 0
}
