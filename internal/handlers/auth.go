package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"smartpicks-backend/internal/database"
	"smartpicks-backend/internal/models"

	"golang.org/x/crypto/bcrypt"
)

// Helper function to send standardized error response
func sendErrorResponse(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

// Helper function to check if user exists
func userExists(field, value string) bool {
	var count int
	query := "SELECT COUNT(*) FROM users WHERE " + field + "=?"
	database.DB.QueryRow(query, value).Scan(&count)
	return count > 0
}

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
			   DATE_FORMAT(data_nascimento, '%Y-%m-%d') as data_nascimento,
			   perfil, avatar, created_at, updated_at 
		FROM users WHERE email=?`, loginData.Email).
		Scan(&user.ID, &user.Nome, &user.Email, &user.Password,
			&user.CPF, &user.DataNascimento, &user.Perfil, &user.Avatar, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		sendErrorResponse(w, "Email ou senha incorretos", http.StatusUnauthorized)
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)) != nil {
		sendErrorResponse(w, "Email ou senha incorretos", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user.ToResponse())
}

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		sendErrorResponse(w, "Dados inválidos fornecidos", http.StatusBadRequest)
		return
	}

	// Validações básicas
	if user.Nome == "" || user.Email == "" || user.Password == "" || user.CPF == "" || user.DataNascimento == "" {
		sendErrorResponse(w, "Todos os campos são obrigatórios", http.StatusBadRequest)
		return
	}

	if user.Perfil == "" {
		user.Perfil = models.PERFIL_USER
	}

	if !models.IsValidPerfil(user.Perfil) {
		sendErrorResponse(w, "Perfil inválido. Use 'admin' ou 'user'", http.StatusBadRequest)
		return
	}

	// Validar e normalizar data
	parsedDate, err := time.Parse("2006-01-02", user.DataNascimento)
	if err != nil {
		if parsedDate, err = time.Parse("02/01/2006", user.DataNascimento); err != nil {
			sendErrorResponse(w, "Formato de data inválido. Use 'YYYY-MM-DD' ou 'DD/MM/YYYY'", http.StatusBadRequest)
			return
		}
	}
	normalizedBirth := parsedDate.Format("2006-01-02")

	// Verificar se email ou CPF já existem
	if userExists("email", user.Email) {
		sendErrorResponse(w, "Email já cadastrado", http.StatusConflict)
		return
	}

	if userExists("cpf", user.CPF) {
		sendErrorResponse(w, "CPF já cadastrado", http.StatusConflict)
		return
	}

	// Hash da senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		sendErrorResponse(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	// Inserir usuário
	result, err := database.DB.Exec(`
		INSERT INTO users (nome, email, password, cpf, data_nascimento, perfil, avatar) 
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		user.Nome, user.Email, string(hashedPassword), user.CPF, normalizedBirth, user.Perfil, user.Avatar)
	if err != nil {
		sendErrorResponse(w, "Erro ao cadastrar usuário", http.StatusInternalServerError)
		return
	}

	userID, err := result.LastInsertId()
	if err != nil {
		sendErrorResponse(w, "Erro ao criar usuário", http.StatusInternalServerError)
		return
	}

	// Buscar usuário criado
	var newUser models.User
	err = database.DB.QueryRow(`
		SELECT id, nome, email, cpf,
			   DATE_FORMAT(data_nascimento, '%Y-%m-%d') as data_nascimento,
			   perfil, avatar, created_at, updated_at 
		FROM users WHERE id=?`, userID).
		Scan(&newUser.ID, &newUser.Nome, &newUser.Email, &newUser.CPF,
			&newUser.DataNascimento, &newUser.Perfil, &newUser.Avatar, &newUser.CreatedAt, &newUser.UpdatedAt)
	if err != nil {
		sendErrorResponse(w, "Erro ao buscar dados do usuário", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser.ToResponse())
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query(`
		SELECT id, nome, email, cpf, DATE_FORMAT(data_nascimento, '%Y-%m-%d') as data_nascimento, 
		       perfil, avatar, DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s') as created_at,
		       DATE_FORMAT(updated_at, '%Y-%m-%d %H:%i:%s') as updated_at
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
		var createdAtStr, updatedAtStr string

		err := rows.Scan(&user.ID, &user.Nome, &user.Email, &user.CPF,
			&user.DataNascimento, &user.Perfil, &user.Avatar, &createdAtStr, &updatedAtStr)
		if err != nil {
			http.Error(w, "Erro ao processar dados dos usuários: "+err.Error(), http.StatusInternalServerError)
			return
		}

		user.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
		user.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAtStr)

		users = append(users, user.ToResponse())
	}

	if err = rows.Err(); err != nil {
		sendErrorResponse(w, "Erro na consulta de usuários", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"users":   users,
		"total":   len(users),
		"message": "Usuários encontrados com sucesso",
	})
}

func CheckUserPermissions(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		sendErrorResponse(w, "Email é obrigatório", http.StatusBadRequest)
		return
	}

	var user models.User
	err := database.DB.QueryRow(`
		SELECT id, nome, email, cpf,
			   DATE_FORMAT(data_nascimento, '%Y-%m-%d') as data_nascimento,
			   perfil, avatar, created_at, updated_at 
		FROM users WHERE email=?`, email).
		Scan(&user.ID, &user.Nome, &user.Email, &user.CPF,
			&user.DataNascimento, &user.Perfil, &user.Avatar, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		sendErrorResponse(w, "Usuário não encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user.ToResponse())
}

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
		SELECT id, nome, email, cpf, DATE_FORMAT(data_nascimento, '%Y-%m-%d') as data_nascimento, 
		       perfil, avatar, DATE_FORMAT(created_at, '%Y-%m-%d %H:%i:%s') as created_at,
		       DATE_FORMAT(updated_at, '%Y-%m-%d %H:%i:%s') as updated_at
		FROM users 
		WHERE perfil = ?
		ORDER BY created_at DESC`, profile)
	if err != nil {
		sendErrorResponse(w, "Erro ao buscar usuários por perfil", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.UserResponse
	for rows.Next() {
		var user models.User
		var createdAtStr, updatedAtStr string

		err := rows.Scan(&user.ID, &user.Nome, &user.Email, &user.CPF,
			&user.DataNascimento, &user.Perfil, &user.Avatar, &createdAtStr, &updatedAtStr)
		if err != nil {
			sendErrorResponse(w, "Erro ao processar dados dos usuários", http.StatusInternalServerError)
			return
		}

		// Converter strings para time.Time
		user.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
		user.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAtStr)

		users = append(users, user.ToResponse())
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"users":   users,
		"total":   len(users),
		"profile": profile,
		"message": "Usuários encontrados com sucesso",
	})
}

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

	// Validar tamanho do avatar (aceitar até ~5MB em base64)
	if requestData.Avatar != "" {
		// Base64 aumenta ~33% o tamanho. 5MB binário ~ 6.7MB base64. Usamos 7MB como margem.
		const maxBase64Len = 7 * 1024 * 1024 // ~7MB
		if len(requestData.Avatar) > maxBase64Len {
			sendErrorResponse(w, "Avatar muito grande. Máximo 5MB", http.StatusBadRequest)
			return
		}
	}

	// Atualizar avatar no banco
	var avatarPtr *string
	if requestData.Avatar != "" {
		avatarPtr = &requestData.Avatar
	}

	result, err := database.DB.Exec("UPDATE users SET avatar = ? WHERE id = ?", avatarPtr, requestData.UserID)
	if err != nil {
		sendErrorResponse(w, "Erro na query SQL: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Verificar se alguma linha foi afetada
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		sendErrorResponse(w, "Erro ao verificar linhas afetadas: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		sendErrorResponse(w, "Usuário não encontrado ou não foi possível atualizar", http.StatusNotFound)
		return
	}

	// Buscar usuário atualizado
	var user models.User
	err = database.DB.QueryRow(`
		SELECT id, nome, email, cpf,
			   DATE_FORMAT(data_nascimento, '%Y-%m-%d') as data_nascimento,
			   perfil, avatar, created_at, updated_at 
		FROM users WHERE id=?`, requestData.UserID).
		Scan(&user.ID, &user.Nome, &user.Email, &user.CPF,
			&user.DataNascimento, &user.Perfil, &user.Avatar, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		sendErrorResponse(w, "Usuário não encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user":    user.ToResponse(),
		"message": "Avatar atualizado com sucesso",
	})
}
