package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var (
	db   *sql.DB
	once sync.Once
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Version string `json:"version"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type UserResponse struct {
	ID              int       `json:"id"`
	Nome            string    `json:"nome"`
	Email           string    `json:"email"`
	CPF             string    `json:"cpf,omitempty"`
	DataNascimento  string    `json:"data_nascimento,omitempty"`
	Perfil          string    `json:"perfil"`
	Avatar          string    `json:"avatar,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func connectDB() {
	once.Do(func() {
		databaseURL := os.Getenv("DATABASE_URL")
		if databaseURL == "" {
			fmt.Println("DATABASE_URL não configurada")
			return
		}

		var err error
		db, err = sql.Open("postgres", databaseURL)
		if err != nil {
			fmt.Printf("Erro ao conectar ao banco: %v\n", err)
			return
		}

		if err = db.Ping(); err != nil {
			fmt.Printf("Erro ao verificar conexão: %v\n", err)
			db = nil
			return
		}

		fmt.Println("Banco de dados conectado com sucesso!")
	})
}

func enableCORS(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	allowedOrigins := []string{
		"https://smartpicks-88709.web.app",
		"https://smartpicks-88709.firebaseapp.com",
		"http://localhost:9000",
	}

	for _, allowedOrigin := range allowedOrigins {
		if origin == allowedOrigin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			break
		}
	}

	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Origin, Accept")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Max-Age", "3600")
}

func sendError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:   "error",
		Message: message,
	})
}

func sendSuccess(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func userExists(field, value string) bool {
	if db == nil {
		return false
	}
	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM users WHERE %s = $1", field)
	db.QueryRow(query, value).Scan(&count)
	return count > 0
}

func Handler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w, r)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	connectDB()

	router := mux.NewRouter()

	// Rota principal
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		response := Response{
			Status:  "online",
			Message: "Backend SmartPicks rodando na Vercel",
			Version: "1.0.0",
		}
		sendSuccess(w, response, http.StatusOK)
	}).Methods("GET")

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		status := map[string]interface{}{
			"status":   "healthy",
			"database": db != nil,
		}
		sendSuccess(w, status, http.StatusOK)
	}).Methods("GET")

	// API Routes
	api := router.PathPrefix("/api").Subrouter()

	// Register
	api.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		var userData map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
			sendError(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		// Validar campos obrigatórios
		required := []string{"nome", "email", "senha", "cpf", "data_nascimento"}
		for _, field := range required {
			if _, ok := userData[field]; !ok {
				sendError(w, fmt.Sprintf("Campo %s é obrigatório", field), http.StatusBadRequest)
				return
			}
		}

		if db == nil {
			sendError(w, "Banco de dados não conectado", http.StatusServiceUnavailable)
			return
		}

		email := userData["email"].(string)
		cpf := userData["cpf"].(string)

		if userExists("email", email) {
			sendError(w, "Email já cadastrado", http.StatusConflict)
			return
		}

		if userExists("cpf", cpf) {
			sendError(w, "CPF já cadastrado", http.StatusConflict)
			return
		}

		perfil := "user"
		if p, ok := userData["perfil"].(string); ok && p != "" {
			perfil = p
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData["senha"].(string)), bcrypt.DefaultCost)
		if err != nil {
			sendError(w, "Erro ao processar senha", http.StatusInternalServerError)
			return
		}

		var userID int
		err = db.QueryRow(`
			INSERT INTO users (nome, email, password, cpf, data_nascimento, perfil)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id`,
			userData["nome"], email, string(hashedPassword), cpf,
			userData["data_nascimento"], perfil).Scan(&userID)

		if err != nil {
			sendError(w, fmt.Sprintf("Erro ao cadastrar usuário: %v", err), http.StatusInternalServerError)
			return
		}

		var user UserResponse
		err = db.QueryRow(`
			SELECT id, nome, email, cpf, 
				   TO_CHAR(data_nascimento, 'YYYY-MM-DD') as data_nascimento,
				   perfil, COALESCE(avatar, '') as avatar, created_at, updated_at
			FROM users WHERE id = $1`, userID).
			Scan(&user.ID, &user.Nome, &user.Email, &user.CPF,
				&user.DataNascimento, &user.Perfil, &user.Avatar, &user.CreatedAt, &user.UpdatedAt)

		if err != nil {
			sendError(w, "Erro ao buscar usuário cadastrado", http.StatusInternalServerError)
			return
		}

		sendSuccess(w, user, http.StatusCreated)
	}).Methods("POST")

	// Login
	api.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		var loginData map[string]string
		if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
			sendError(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		email, senha := loginData["email"], loginData["senha"]
		if email == "" || senha == "" {
			sendError(w, "Email e senha são obrigatórios", http.StatusBadRequest)
			return
		}

		if db == nil {
			sendError(w, "Banco de dados não conectado", http.StatusServiceUnavailable)
			return
		}

		var user UserResponse
		var hashedPassword string
		err := db.QueryRow(`
			SELECT id, nome, email, password, cpf,
				   TO_CHAR(data_nascimento, 'YYYY-MM-DD') as data_nascimento,
				   perfil, COALESCE(avatar, '') as avatar, created_at, updated_at
			FROM users WHERE email = $1`, email).
			Scan(&user.ID, &user.Nome, &user.Email, &hashedPassword, &user.CPF,
				&user.DataNascimento, &user.Perfil, &user.Avatar, &user.CreatedAt, &user.UpdatedAt)

		if err != nil {
			sendError(w, "Email ou senha incorretos", http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(senha)); err != nil {
			sendError(w, "Email ou senha incorretos", http.StatusUnauthorized)
			return
		}

		sendSuccess(w, user, http.StatusOK)
	}).Methods("POST")

	// Get All Users
	api.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if db == nil {
			sendError(w, "Banco de dados não conectado", http.StatusServiceUnavailable)
			return
		}

		rows, err := db.Query(`
			SELECT id, nome, email, cpf,
				   TO_CHAR(data_nascimento, 'YYYY-MM-DD') as data_nascimento,
				   perfil, COALESCE(avatar, '') as avatar, created_at, updated_at
			FROM users ORDER BY created_at DESC`)

		if err != nil {
			sendError(w, "Erro ao buscar usuários", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var users []UserResponse
		for rows.Next() {
			var user UserResponse
			err := rows.Scan(&user.ID, &user.Nome, &user.Email, &user.CPF,
				&user.DataNascimento, &user.Perfil, &user.Avatar, &user.CreatedAt, &user.UpdatedAt)
			if err != nil {
				continue
			}
			users = append(users, user)
		}

		sendSuccess(w, map[string]interface{}{
			"total":   len(users),
			"users":   users,
			"message": "Usuários encontrados com sucesso",
		}, http.StatusOK)
	}).Methods("GET")

	// Check User Permissions
	api.HandleFunc("/users/permissions", func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("email")
		if email == "" {
			sendError(w, "Email é obrigatório", http.StatusBadRequest)
			return
		}

		if db == nil {
			sendError(w, "Banco de dados não conectado", http.StatusServiceUnavailable)
			return
		}

		var perfil string
		err := db.QueryRow("SELECT perfil FROM users WHERE email = $1", email).Scan(&perfil)
		if err != nil {
			sendError(w, "Usuário não encontrado", http.StatusNotFound)
			return
		}

		sendSuccess(w, map[string]interface{}{
			"email":   email,
			"perfil":  perfil,
			"isAdmin": perfil == "admin",
		}, http.StatusOK)
	}).Methods("GET")

	// Get Users by Profile
	api.HandleFunc("/users/profile", func(w http.ResponseWriter, r *http.Request) {
		perfil := r.URL.Query().Get("perfil")
		if perfil == "" {
			sendError(w, "Perfil é obrigatório", http.StatusBadRequest)
			return
		}

		if db == nil {
			sendError(w, "Banco de dados não conectado", http.StatusServiceUnavailable)
			return
		}

		rows, err := db.Query(`
			SELECT id, nome, email, cpf,
				   TO_CHAR(data_nascimento, 'YYYY-MM-DD') as data_nascimento,
				   perfil, COALESCE(avatar, '') as avatar, created_at, updated_at
			FROM users WHERE perfil = $1 ORDER BY created_at DESC`, perfil)

		if err != nil {
			sendError(w, "Erro ao buscar usuários", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var users []UserResponse
		for rows.Next() {
			var user UserResponse
			err := rows.Scan(&user.ID, &user.Nome, &user.Email, &user.CPF,
				&user.DataNascimento, &user.Perfil, &user.Avatar, &user.CreatedAt, &user.UpdatedAt)
			if err != nil {
				continue
			}
			users = append(users, user)
		}

		sendSuccess(w, map[string]interface{}{
			"total":   len(users),
			"perfil":  perfil,
			"users":   users,
			"message": fmt.Sprintf("Usuários com perfil '%s' encontrados", perfil),
		}, http.StatusOK)
	}).Methods("GET")

	// Update Avatar (simplificado - aceita URL do avatar)
	api.HandleFunc("/users/avatar", func(w http.ResponseWriter, r *http.Request) {
		var data map[string]string
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			sendError(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		email, avatar := data["email"], data["avatar"]
		if email == "" {
			sendError(w, "Email é obrigatório", http.StatusBadRequest)
			return
		}

		if db == nil {
			sendError(w, "Banco de dados não conectado", http.StatusServiceUnavailable)
			return
		}

		_, err := db.Exec("UPDATE users SET avatar = $1, updated_at = NOW() WHERE email = $2", avatar, email)
		if err != nil {
			sendError(w, "Erro ao atualizar avatar", http.StatusInternalServerError)
			return
		}

		sendSuccess(w, map[string]interface{}{
			"message": "Avatar atualizado com sucesso",
			"avatar":  avatar,
		}, http.StatusOK)
	}).Methods("POST", "PUT")

	// Delete Avatar
	api.HandleFunc("/users/avatar", func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("email")
		if email == "" {
			sendError(w, "Email é obrigatório", http.StatusBadRequest)
			return
		}

		if db == nil {
			sendError(w, "Banco de dados não conectado", http.StatusServiceUnavailable)
			return
		}

		_, err := db.Exec("UPDATE users SET avatar = NULL, updated_at = NOW() WHERE email = $1", email)
		if err != nil {
			sendError(w, "Erro ao remover avatar", http.StatusInternalServerError)
			return
		}

		sendSuccess(w, map[string]string{
			"message": "Avatar removido com sucesso",
		}, http.StatusOK)
	}).Methods("DELETE")

	router.ServeHTTP(w, r)
}
