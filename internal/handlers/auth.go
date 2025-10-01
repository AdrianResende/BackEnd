package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"smartpicks-backend/internal/database"
	"smartpicks-backend/internal/models"

	"golang.org/x/crypto/bcrypt"
)

// Login godoc
// @Summary Login de usuário
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

	// Validação do JSON de entrada
	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Validação de campos obrigatórios
	if loginData.Email == "" || loginData.Password == "" {
		http.Error(w, "Email e password são obrigatórios", http.StatusBadRequest)
		return
	}

	var dbUser models.User
	err := database.DB.QueryRow(`
		SELECT id, nome, email, password, cpf, data_nascimento, perfil, created_at, updated_at 
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

	// Criar response sem senha com verificações de permissão
	userResponse := dbUser.ToUserResponse()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userResponse)
}

// Register godoc
// @Summary Cadastro de usuário
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

	// Validação do JSON de entrada
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Validação de campos obrigatórios
	if user.Nome == "" || user.Email == "" || user.Password == "" || user.CPF == "" || user.DataNascimento == "" {
		http.Error(w, "Todos os campos são obrigatórios (nome, email, password, cpf, data_nascimento)", http.StatusBadRequest)
		return
	}

	// Se perfil não foi fornecido, usar 'user' como padrão
	if user.Perfil == "" {
		user.Perfil = models.PERFIL_USER
	}

	// Validação do perfil
	if !models.IsValidPerfil(user.Perfil) {
		http.Error(w, "Perfil inválido. Use 'admin' ou 'user'", http.StatusBadRequest)
		return
	}

	// Verificar se email já existe
	var existingID int
	err := database.DB.QueryRow("SELECT id FROM users WHERE email=?", user.Email).Scan(&existingID)
	if err == nil {
		http.Error(w, "Email já cadastrado", http.StatusConflict)
		return
	}

	// Verificar se CPF já existe
	err = database.DB.QueryRow("SELECT id FROM users WHERE cpf=?", user.CPF).Scan(&existingID)
	if err == nil {
		http.Error(w, "CPF já cadastrado", http.StatusConflict)
		return
	}

	// Hash da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	// Inserir usuário no banco
	result, err := database.DB.Exec(`
		INSERT INTO users (nome, email, password, cpf, data_nascimento, perfil) 
		VALUES (?, ?, ?, ?, ?, ?)`,
		user.Nome, user.Email, string(hashedPassword), user.CPF, user.DataNascimento, user.Perfil)
	if err != nil {
		http.Error(w, "Erro ao cadastrar usuário", http.StatusInternalServerError)
		return
	}

	// Obter ID do usuário criado
	userID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Erro ao obter ID do usuário", http.StatusInternalServerError)
		return
	}

	// Buscar dados completos do usuário criado
	var newUser models.User
	err = database.DB.QueryRow(`
		SELECT id, nome, email, cpf, data_nascimento, perfil, created_at, updated_at 
		FROM users WHERE id=?`, userID).
		Scan(&newUser.ID, &newUser.Nome, &newUser.Email, &newUser.CPF,
			&newUser.DataNascimento, &newUser.Perfil, &newUser.CreatedAt, &newUser.UpdatedAt)
	if err != nil {
		http.Error(w, "Erro ao buscar dados do usuário", http.StatusInternalServerError)
		return
	}

	// Criar response sem senha com verificações de permissão
	userResponse := newUser.ToUserResponse()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userResponse)
}

// GetAllUsers godoc
// @Summary Listar todos os usuários
// @Description Retorna a lista de todos os usuários cadastrados com informações de perfil e permissões
// @Tags Usuários
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Lista de usuários com total e mensagem"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Router /users [get]
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

		// Converter strings para time.Time
		user.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
		user.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAtStr)

		// Usar método ToUserResponse para incluir verificações de permissão
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

// CheckUsersTable godoc
// @Summary Verificar status da tabela users
// @Description Endpoint de debug para verificar se a tabela users existe e está funcionando
// @Tags Debug
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Status da tabela users"
// @Failure 404 {object} map[string]string "Tabela users não encontrada"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Router /users/check [get]
func CheckUsersTable(w http.ResponseWriter, r *http.Request) {
	// Verificar se a tabela existe
	var tableName string
	err := database.DB.QueryRow("SHOW TABLES LIKE 'users'").Scan(&tableName)
	if err != nil {
		http.Error(w, "Tabela users não encontrada: "+err.Error(), http.StatusNotFound)
		return
	}

	// Contar quantos usuários existem
	var count int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		http.Error(w, "Erro ao contar usuários: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Verificar estrutura da tabela
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

// CheckUserPermissions godoc
// @Summary Verificar permissões do usuário
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
	// Obter email do usuário dos parâmetros da query
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email é obrigatório", http.StatusBadRequest)
		return
	}

	var user models.User
	err := database.DB.QueryRow(`
		SELECT id, nome, email, cpf, data_nascimento, perfil, created_at, updated_at 
		FROM users WHERE email=?`, email).
		Scan(&user.ID, &user.Nome, &user.Email, &user.CPF,
			&user.DataNascimento, &user.Perfil, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		http.Error(w, "Usuário não encontrado", http.StatusNotFound)
		return
	}

	// Criar response com informações de permissão
	response := user.ToUserResponse()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user": response,
		"permissions": map[string]bool{
			"is_admin":       response.IsAdmin,
			"has_permission": response.HasPermission,
			"can_create":     response.HasPermission,
			"can_read":       response.HasPermission,
			"can_update":     response.IsAdmin, // Apenas admin pode atualizar outros usuários
			"can_delete":     response.IsAdmin, // Apenas admin pode deletar usuários
			"can_manage_all": response.IsAdmin,
		},
		"message": "Permissões verificadas com sucesso",
	})
}

// GetUsersByProfile godoc
// @Summary Listar usuários por perfil
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
		http.Error(w, "Parâmetro 'profile' é obrigatório", http.StatusBadRequest)
		return
	}

	// Validar perfil
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
