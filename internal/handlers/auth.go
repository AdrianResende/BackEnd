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
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if loginData.Email == "" || loginData.Password == "" {
		http.Error(w, "Email e password são obrigatórios", http.StatusBadRequest)
		return
	}

	var dbUser models.User
	err := database.DB.QueryRow(`
		SELECT id, nome, email, password, cpf,
			   DATE_FORMAT(data_nascimento, '%Y-%m-%d') as data_nascimento,
			   perfil, created_at, updated_at 
		FROM users WHERE email=?`, loginData.Email).
		Scan(&dbUser.ID, &dbUser.Nome, &dbUser.Email, &dbUser.Password,
			&dbUser.CPF, &dbUser.DataNascimento, &dbUser.Perfil, &dbUser.CreatedAt, &dbUser.UpdatedAt)
	if err != nil {
		http.Error(w, "Usuário não encontrado", http.StatusUnauthorized)
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(loginData.Password)) != nil {
		http.Error(w, "Senha inválida", http.StatusUnauthorized)
		return
	}

	userResponse := dbUser.ToUserResponse()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userResponse)
}

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User

	log.Printf("Register: Recebida requisição de registro")

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("Register: Erro ao decodificar JSON: %v", err)
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	log.Printf("Register: Dados recebidos - Nome: %s, Email: %s, CPF: %s, Perfil: %s",
		user.Nome, user.Email, user.CPF, user.Perfil)

	if user.Nome == "" || user.Email == "" || user.Password == "" || user.CPF == "" || user.DataNascimento == "" {
		http.Error(w, "Todos os campos são obrigatórios (nome, email, password, cpf, data_nascimento)", http.StatusBadRequest)
		return
	}

	if user.Perfil == "" {
		user.Perfil = models.PERFIL_USER
	}

	if !models.IsValidPerfil(user.Perfil) {
		http.Error(w, "Perfil inválido. Use 'admin' ou 'user'", http.StatusBadRequest)
		return
	}

	var parsedDate time.Time
	var parseErr error

	parsedDate, parseErr = time.Parse("2006-01-02", user.DataNascimento)
	if parseErr != nil {
		if t2, err2 := time.Parse("02/01/2006", user.DataNascimento); err2 == nil {
			parsedDate = t2
		} else {
			http.Error(w, "Formato de data inválido. Use 'YYYY-MM-DD' ou 'DD/MM/YYYY'", http.StatusBadRequest)
			return
		}
	}
	normalizedBirth := parsedDate.Format("2006-01-02")

	var existingID int
	err := database.DB.QueryRow("SELECT id FROM users WHERE email=?", user.Email).Scan(&existingID)
	if err == nil {
		http.Error(w, "Email já cadastrado", http.StatusConflict)
		return
	}

	err = database.DB.QueryRow("SELECT id FROM users WHERE cpf=?", user.CPF).Scan(&existingID)
	if err == nil {
		http.Error(w, "CPF já cadastrado", http.StatusConflict)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	log.Printf("Register: Tentando inserir usuário no banco")
	result, err := database.DB.Exec(`
		INSERT INTO users (nome, email, password, cpf, data_nascimento, perfil) 
		VALUES (?, ?, ?, ?, ?, ?)`,
		user.Nome, user.Email, string(hashedPassword), user.CPF, normalizedBirth, user.Perfil)
	if err != nil {
		log.Printf("Register: Erro ao inserir no banco: %v", err)
		http.Error(w, "Erro ao cadastrar usuário: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Register: Usuário inserido com sucesso")

	userID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Erro ao obter ID do usuário", http.StatusInternalServerError)
		return
	}

	log.Printf("Register: Buscando dados do usuário criado com ID: %d", userID)
	var newUser models.User
	err = database.DB.QueryRow(`
		SELECT id, nome, email, cpf,
			   DATE_FORMAT(data_nascimento, '%Y-%m-%d') as data_nascimento,
			   perfil, created_at, updated_at 
		FROM users WHERE id=?`, userID).
		Scan(&newUser.ID, &newUser.Nome, &newUser.Email, &newUser.CPF,
			&newUser.DataNascimento, &newUser.Perfil, &newUser.CreatedAt, &newUser.UpdatedAt)
	if err != nil {
		log.Printf("Register: Erro ao buscar dados do usuário: %v", err)
		http.Error(w, "Erro ao buscar dados do usuário: "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Register: Dados do usuário recuperados com sucesso")

	userResponse := newUser.ToUserResponse()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userResponse)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query(`
		SELECT id, nome, email, cpf, DATE_FORMAT(data_nascimento, '%Y-%m-%d') as data_nascimento, 
		       perfil, DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s') as created_at,
		       DATE_FORMAT(updated_at, '%Y-%m-%d %H:%i:%s') as updated_at
		FROM users 
		ORDER BY created_at DESC`)
	if err != nil {
		http.Error(w, "Erro ao buscar usuários: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.UserResponse
	for rows.Next() {
		var user models.User
		var createdAtStr, updatedAtStr string

		err := rows.Scan(&user.ID, &user.Nome, &user.Email, &user.CPF,
			&user.DataNascimento, &user.Perfil, &createdAtStr, &updatedAtStr)
		if err != nil {
			http.Error(w, "Erro ao processar dados dos usuários: "+err.Error(), http.StatusInternalServerError)
			return
		}

		user.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
		user.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAtStr)

		users = append(users, user.ToUserResponse())
	}

	if err = rows.Err(); err != nil {
		http.Error(w, "Erro na consulta: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"users":   users,
		"total":   len(users),
		"message": "Usuários encontrados com sucesso",
	})
}

func CheckUsersTable(w http.ResponseWriter, r *http.Request) {

	var tableName string
	err := database.DB.QueryRow("SHOW TABLES LIKE 'users'").Scan(&tableName)
	if err != nil {
		http.Error(w, "Tabela users não encontrada: "+err.Error(), http.StatusNotFound)
		return
	}

	var count int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		http.Error(w, "Erro ao contar usuários: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rows, err := database.DB.Query("DESCRIBE users")
	if err != nil {
		http.Error(w, "Erro ao verificar estrutura da tabela: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var field, fieldType, null, key, defaultValue, extra string
		err := rows.Scan(&field, &fieldType, &null, &key, &defaultValue, &extra)
		if err != nil {
			continue
		}
		columns = append(columns, field)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"table_exists": true,
		"table_name":   tableName,
		"users_count":  count,
		"columns":      columns,
		"message":      "Tabela users está funcionando corretamente",
	})
}

func CheckUserPermissions(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email é obrigatório", http.StatusBadRequest)
		return
	}

	var user models.User
	err := database.DB.QueryRow(`
		SELECT id, nome, email, cpf,
			   DATE_FORMAT(data_nascimento, '%Y-%m-%d') as data_nascimento,
			   perfil, created_at, updated_at 
		FROM users WHERE email=?`, email).
		Scan(&user.ID, &user.Nome, &user.Email, &user.CPF,
			&user.DataNascimento, &user.Perfil, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		http.Error(w, "Usuário não encontrado", http.StatusNotFound)
		return
	}

	response := user.ToUserResponse()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user": response,
		"permissions": map[string]bool{
			"is_admin":       response.IsAdmin,
			"has_permission": response.HasPermission,
			"can_create":     response.HasPermission,
			"can_read":       response.HasPermission,
			"can_update":     response.IsAdmin,
			"can_delete":     response.IsAdmin,
			"can_manage_all": response.IsAdmin,
		},
		"message": "Permissões verificadas com sucesso",
	})
}

func GetUsersByProfile(w http.ResponseWriter, r *http.Request) {
	profile := r.URL.Query().Get("profile")
	if profile == "" {
		http.Error(w, "Parâmetro 'profile' é obrigatório", http.StatusBadRequest)
		return
	}

	if !models.IsValidPerfil(profile) {
		http.Error(w, "Perfil inválido. Use 'admin' ou 'user'", http.StatusBadRequest)
		return
	}

	rows, err := database.DB.Query(`
		SELECT id, nome, email, cpf, DATE_FORMAT(data_nascimento, '%Y-%m-%d') as data_nascimento, 
		       perfil, DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s') as created_at,
		       DATE_FORMAT(updated_at, '%Y-%m-%d %H:%i:%s') as updated_at
		FROM users 
		WHERE perfil = ?
		ORDER BY created_at DESC`, profile)
	if err != nil {
		http.Error(w, "Erro ao buscar usuários: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.UserResponse
	for rows.Next() {
		var user models.User
		var createdAtStr, updatedAtStr string

		err := rows.Scan(&user.ID, &user.Nome, &user.Email, &user.CPF,
			&user.DataNascimento, &user.Perfil, &createdAtStr, &updatedAtStr)
		if err != nil {
			http.Error(w, "Erro ao processar dados dos usuários: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Converter strings para time.Time
		user.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
		user.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAtStr)

		users = append(users, user.ToUserResponse())
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"users":   users,
		"total":   len(users),
		"profile": profile,
		"message": "Usuários encontrados com sucesso",
	})
}
