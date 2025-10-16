package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"smartpicks-backend/internal/database"
	"smartpicks-backend/internal/models"

	"golang.org/x/crypto/bcrypt"
)

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
		log.Printf("Erro ao buscar usuário: %v", err)
		sendErrorResponse(w, "Email ou senha incorretos", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		sendErrorResponse(w, "Email ou senha incorretos", http.StatusUnauthorized)
		return
	}

	sendSuccessResponse(w, user.ToResponse())
}

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

	// Validar formato da data
	parsedDate, err := time.Parse("2006-01-02", user.DataNascimento)
	if err != nil {
		parsedDate, err = time.Parse("02/01/2006", user.DataNascimento)
		if err != nil {
			sendErrorResponse(w, "Formato de data inválido. Use YYYY-MM-DD ou DD/MM/YYYY", http.StatusBadRequest)
			return
		}
	}
	user.DataNascimento = parsedDate.Format("2006-01-02")

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
		log.Printf("Erro ao cadastrar usuário: %v", err)
		sendErrorResponse(w, "Erro ao cadastrar usuário: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Buscar usuário recém-criado
	err = database.DB.QueryRow(`
		SELECT id, nome, email, password, cpf,
			   TO_CHAR(data_nascimento, 'YYYY-MM-DD') as data_nascimento,
			   perfil, avatar, created_at, updated_at
		FROM users WHERE id=$1
	`, user.ID).Scan(
		&user.ID, &user.Nome, &user.Email, &user.Password,
		&user.CPF, &user.DataNascimento, &user.Perfil, &user.Avatar,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		log.Printf("Erro ao buscar usuário cadastrado: %v", err)
		sendErrorResponse(w, "Erro ao buscar usuário cadastrado: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	sendSuccessResponse(w, user.ToResponse())
}
