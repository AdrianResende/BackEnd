package handlers

import (
	"encoding/json"
	"net/http"
	"smartpicks-backend/internal/database"
	"smartpicks-backend/internal/models"
)

func UpdateAvatar(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		UserID int    `json:"user_id"`
		Avatar string `json:"avatar"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		sendErrorResponse(w, "Dados inválidos fornecidos", http.StatusBadRequest)
		return
	}

	if requestData.UserID <= 0 {
		sendErrorResponse(w, "ID do usuário é obrigatório", http.StatusBadRequest)
		return
	}

	if requestData.Avatar != "" {
		const maxBase64Len = 7 * 1024 * 1024 // ~5MB em Base64
		if len(requestData.Avatar) > maxBase64Len {
			sendErrorResponse(w, "Avatar muito grande. Máximo 5MB", http.StatusBadRequest)
			return
		}
	}

	var avatarPtr *string
	if requestData.Avatar != "" {
		avatarPtr = &requestData.Avatar
	}

<<<<<<< HEAD
	result, err := database.DB.Exec("UPDATE users SET avatar = $1 WHERE id = $2", avatarPtr, requestData.UserID)
=======
	// Atualizar avatar usando PostgreSQL
	result, err := database.DB.Exec(
		"UPDATE users SET avatar = $1 WHERE id = $2",
		avatarPtr, requestData.UserID,
	)
>>>>>>> development
	if err != nil {
		sendErrorResponse(w, "Erro na query SQL: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		sendErrorResponse(w, "Erro ao verificar linhas afetadas: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		sendErrorResponse(w, "Usuário não encontrado ou não foi possível atualizar", http.StatusNotFound)
		return
	}

	var user models.User
	err = database.DB.QueryRow(`
<<<<<<< HEAD
		SELECT id, nome, email, cpf,
			   TO_CHAR(data_nascimento, 'YYYY-MM-DD') as data_nascimento,
			   perfil, COALESCE(avatar, '') as avatar, created_at, updated_at 
		FROM users WHERE id = $1`, requestData.UserID).
		Scan(&user.ID, &user.Nome, &user.Email, &user.CPF,
			&user.DataNascimento, &user.Perfil, &user.Avatar, &user.CreatedAt, &user.UpdatedAt)
=======
		SELECT id, nome, email, password, cpf,
		       to_char(data_nascimento, 'YYYY-MM-DD') as data_nascimento,
		       perfil, avatar, created_at, updated_at
		FROM users WHERE id=$1`, requestData.UserID).
		Scan(&user.ID, &user.Nome, &user.Email, &user.Password,
			&user.CPF, &user.DataNascimento, &user.Perfil, &user.Avatar, &user.CreatedAt, &user.UpdatedAt)
>>>>>>> development
	if err != nil {
		sendErrorResponse(w, "Usuário não encontrado", http.StatusNotFound)
		return
	}

	sendSuccessResponse(w, map[string]interface{}{
		"user":    user.ToResponse(),
		"message": "Avatar atualizado com sucesso",
	})
}

func DeleteAvatar(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		UserID int `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		sendErrorResponse(w, "Dados inválidos fornecidos", http.StatusBadRequest)
		return
	}

	if requestData.UserID <= 0 {
		sendErrorResponse(w, "ID do usuário é obrigatório", http.StatusBadRequest)
		return
	}

<<<<<<< HEAD
	var userExists bool
	err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", requestData.UserID).Scan(&userExists)
=======
	// Verificar se o usuário existe (PostgreSQL)
	var exists bool
	err := database.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)",
		requestData.UserID,
	).Scan(&exists)
>>>>>>> development
	if err != nil {
		sendErrorResponse(w, "Erro ao verificar usuário: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !exists {
		sendErrorResponse(w, "Usuário não encontrado", http.StatusNotFound)
		return
	}

<<<<<<< HEAD
	result, err := database.DB.Exec("UPDATE users SET avatar = NULL WHERE id = $1", requestData.UserID)
=======
	// Remover o avatar (definir como NULL)
	result, err := database.DB.Exec(
		"UPDATE users SET avatar = NULL WHERE id = $1",
		requestData.UserID,
	)
>>>>>>> development
	if err != nil {
		sendErrorResponse(w, "Erro na query SQL: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		sendErrorResponse(w, "Erro ao verificar linhas afetadas: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		sendErrorResponse(w, "Não foi possível remover o avatar", http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(w, map[string]string{
		"message": "Avatar removido com sucesso",
	})
}
