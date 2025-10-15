package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"smartpicks-backend/internal/database"
	"smartpicks-backend/internal/models"

	"golang.org/x/crypto/bcrypt"
)

// Login @Summary Login de usuário
// @Description Autentica um usuário com email e senha
// @Tags Autenticação
// @Accept json
// @Produce json
// @Param loginData body models.UserLogin true "Dados de login"
// @Success 200 {object} models.UserResponse "Login realizado com sucesso"
// @Failure 400 {object} map[string]string "JSON inválido ou campos obrigatórios ausentes"
// @Failure 401 {object} map[string]string "Credenciais inválidas"
// @Router /login [post]
func Login(w http.ResponseWriter, r *http.Request) {
	var loginData models.UserLogin

	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		sendErrorResponse(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if loginData.Email == "" || loginData.Password == "" {
		sendErrorResponse(w, "Email e password são obrigatórios", http.StatusBadRequest)
		return
	}

	var user models.User
	err := database.DB.QueryRow(`
		SELECT id, nome, email, password, cpf,
			   TO_CHAR(data_nascimento, 'YYYY-MM-DD') as data_nascimento,
			   perfil, COALESCE(avatar, '') as avatar, created_at, updated_at 
		FROM users WHERE email = $1`, loginData.Email).
		Scan(&user.ID, &user.Nome, &user.Email, &user.Password,
			&user.CPF, &user.DataNascimento, &user.Perfil, &user.Avatar, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		sendErrorResponse(w, "Email ou senha incorretos", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		sendErrorResponse(w, "Email ou senha incorretos", http.StatusUnauthorized)
		return
	}

	sendSuccessResponse(w, user.ToResponse())
}

// Register @Summary Cadastro de usuário
// @Description Cadastra um novo usuário no sistema
// @Tags Autenticação
// @Accept json
// @Produce json
// @Param userData body models.User true "Dados do usuário (perfil é opcional, padrão: 'user')"
// @Success 201 {object} models.UserResponse "Usuário cadastrado com sucesso"
// @Failure 400 {object} map[string]string "JSON inválido, campos obrigatórios ausentes ou perfil inválido"
// @Failure 409 {object} map[string]string "Email ou CPF já cadastrado"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Router /register [post]
func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		sendErrorResponse(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if user.Nome == "" || user.Email == "" || user.Password == "" || user.CPF == "" || user.DataNascimento == "" {
		sendErrorResponse(w, "Nome, email, password, CPF e data de nascimento são obrigatórios", http.StatusBadRequest)
		return
	}

	if user.Perfil == "" {
		user.Perfil = "user"
	}

	if !models.IsValidPerfil(user.Perfil) {
		sendErrorResponse(w, "Perfil inválido. Use 'admin' ou 'user'", http.StatusBadRequest)
		return
	}

	// Validar formato de data
	parsedDate, err := time.Parse("2006-01-02", user.DataNascimento)
	if err != nil {
		if parsedDate, err = time.Parse("02/01/2006", user.DataNascimento); err != nil {
			sendErrorResponse(w, "Formato de data inválido. Use YYYY-MM-DD ou DD/MM/YYYY", http.StatusBadRequest)
			return
		}
		user.DataNascimento = parsedDate.Format("2006-01-02")
	}

	if userExists("email", user.Email) {
		sendErrorResponse(w, "Email já cadastrado", http.StatusConflict)
		return
	}

	if userExists("cpf", user.CPF) {
		sendErrorResponse(w, "CPF já cadastrado", http.StatusConflict)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		sendErrorResponse(w, "Erro ao processar password", http.StatusInternalServerError)
		return
	}

	var userID int
	err = database.DB.QueryRow(`
		INSERT INTO users (nome, email, password, cpf, data_nascimento, perfil, avatar)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`,
		user.Nome, user.Email, string(hashedPassword), user.CPF, user.DataNascimento, user.Perfil, user.Avatar).Scan(&userID)

	if err != nil {
		sendErrorResponse(w, "Erro ao cadastrar usuário", http.StatusInternalServerError)
		return
	}

	err = database.DB.QueryRow(`
		SELECT id, nome, email, cpf,
			   TO_CHAR(data_nascimento, 'YYYY-MM-DD') as data_nascimento,
			   perfil, COALESCE(avatar, '') as avatar, created_at, updated_at 
		FROM users WHERE id = $1`, userID).
		Scan(&user.ID, &user.Nome, &user.Email, &user.CPF,
			&user.DataNascimento, &user.Perfil, &user.Avatar, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		sendErrorResponse(w, "Erro ao buscar usuário cadastrado", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	sendSuccessResponse(w, user.ToResponse())
}
