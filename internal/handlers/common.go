package handlers

import (
	"encoding/json"
	"fmt"
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

// userExists verifica se um usuário existe baseado em um campo específico (PostgreSQL)
func userExists(field, value string) bool {
	var count int
<<<<<<< HEAD
	query := fmt.Sprintf("SELECT COUNT(*) FROM users WHERE %s = $1", field)
	database.DB.QueryRow(query, value).Scan(&count)
=======
	// Usa $1 para Postgres e valida o nome do campo para evitar SQL Injection
	allowedFields := map[string]bool{"email": true, "cpf": true}
	if !allowedFields[field] {
		return false
	}
	query := "SELECT COUNT(*) FROM users WHERE " + field + " = $1"
	err := database.DB.QueryRow(query, value).Scan(&count)
	if err != nil {
		return false
	}
>>>>>>> development
	return count > 0
}
