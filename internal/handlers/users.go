package handlers

import (
	"net/http"

	"smartpicks-backend/internal/database"
	"smartpicks-backend/internal/models"
)

// GetAllUsers @Summary Listar todos os usuários
// @Description Retorna a lista de todos os usuários cadastrados com informações de perfil e permissões
// @Tags Usuários
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Lista de usuários com total e mensagem"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Router /users [get]
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query(`
		SELECT id, nome, email, cpf,
			   TO_CHAR(data_nascimento, 'YYYY-MM-DD') as data_nascimento,
			   perfil, COALESCE(avatar, '') as avatar,
			   created_at,
			   updated_at
		FROM users
		ORDER BY created_at DESC`)
	if err != nil {
		sendErrorResponse(w, "Erro ao buscar usuários", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.UserResponse
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Nome, &user.Email, &user.CPF,
			&user.DataNascimento, &user.Perfil, &user.Avatar,
			&user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			sendErrorResponse(w, "Erro ao processar dados dos usuários", http.StatusInternalServerError)
			return
		}

		// timestamps already scanned into time.Time for Postgres

		users = append(users, user.ToResponse())
	}

	sendSuccessResponse(w, map[string]interface{}{
		"users":   users,
		"total":   len(users),
		"message": "Usuários listados com sucesso",
	})
}

// CheckUserPermissions @Summary Verificar permissões do usuário
// @Description Verifica as permissões de um usuário específico baseado no email
// @Tags Usuários
// @Accept json
// @Produce json
// @Param email query string true "Email do usuário"
// @Success 200 {object} map[string]interface{} "Informações do usuário e suas permissões"
// @Failure 400 {object} map[string]string "Email é obrigatório"
// @Failure 404 {object} map[string]string "Usuário não encontrado"
// @Router /users/permissions [get]
func CheckUserPermissions(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		sendErrorResponse(w, "Email é obrigatório", http.StatusBadRequest)
		return
	}

	var user models.User
	err := database.DB.QueryRow(`
		SELECT id, nome, email, cpf,
			   TO_CHAR(data_nascimento, 'YYYY-MM-DD') as data_nascimento,
			   perfil, COALESCE(avatar, '') as avatar, created_at, updated_at 
		FROM users WHERE email = $1`, email).
		Scan(&user.ID, &user.Nome, &user.Email, &user.CPF,
			&user.DataNascimento, &user.Perfil, &user.Avatar, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		sendErrorResponse(w, "Usuário não encontrado", http.StatusNotFound)
		return
	}

	sendSuccessResponse(w, user.ToResponse())
}

// GetUsersByProfile @Summary Listar usuários por perfil
// @Description Retorna usuários filtrados por perfil (admin ou user)
// @Tags Usuários
// @Accept json
// @Produce json
// @Param profile query string true "Perfil do usuário" Enums(admin, user)
// @Success 200 {object} map[string]interface{} "Lista de usuários do perfil especificado"
// @Failure 400 {object} map[string]string "Parâmetro profile obrigatório ou perfil inválido"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Router /users/profile [get]
func GetUsersByProfile(w http.ResponseWriter, r *http.Request) {
	profile := r.URL.Query().Get("profile")
	if profile == "" {
		sendErrorResponse(w, "Parâmetro 'profile' é obrigatório", http.StatusBadRequest)
		return
	}

	if !models.IsValidPerfil(profile) {
		sendErrorResponse(w, "Perfil inválido. Use 'admin' ou 'user'", http.StatusBadRequest)
		return
	}

	rows, err := database.DB.Query(`
		SELECT id, nome, email, cpf,
			   TO_CHAR(data_nascimento, 'YYYY-MM-DD') as data_nascimento,
			   perfil, COALESCE(avatar, '') as avatar, created_at, updated_at
		FROM users 
		WHERE perfil = $1
		ORDER BY created_at DESC`, profile)
	if err != nil {
		sendErrorResponse(w, "Erro ao buscar usuários por perfil", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.UserResponse
	for rows.Next() {
		var user models.User

		err := rows.Scan(&user.ID, &user.Nome, &user.Email, &user.CPF,
			&user.DataNascimento, &user.Perfil, &user.Avatar, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			sendErrorResponse(w, "Erro ao processar dados dos usuários", http.StatusInternalServerError)
			return
		}

		users = append(users, user.ToResponse())
	}

	sendSuccessResponse(w, map[string]interface{}{
		"users":   users,
		"total":   len(users),
		"profile": profile,
		"message": "Usuários encontrados com sucesso",
	})
}
