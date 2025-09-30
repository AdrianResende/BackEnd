package models

import "time"

// Constantes para perfis de usuário
const (
	PERFIL_ADMIN = "admin"
	PERFIL_USER  = "user"
)

// ValidPerfis lista os perfis válidos
var ValidPerfis = []string{PERFIL_ADMIN, PERFIL_USER}

type User struct {
	ID             int       `json:"id"`
	Nome           string    `json:"nome"`
	Email          string    `json:"email"`
	Password       string    `json:"password,omitempty"` // omitempty para não retornar senha nos responses
	CPF            string    `json:"cpf"`
	DataNascimento string    `json:"data_nascimento"` // Formato: YYYY-MM-DD
	Perfil         string    `json:"perfil"`          // ENUM: 'admin', 'user'
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// UserLogin representa os dados necessários para login
type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserResponse representa os dados do usuário para retorno (sem senha)
type UserResponse struct {
	ID             int       `json:"id"`
	Nome           string    `json:"nome"`
	Email          string    `json:"email"`
	CPF            string    `json:"cpf"`
	DataNascimento string    `json:"data_nascimento"`
	Perfil         string    `json:"perfil"`
	IsAdmin        bool      `json:"is_admin"`       // Verificação de permissão
	HasPermission  bool      `json:"has_permission"` // Permissão para operações específicas
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// IsValidPerfil verifica se o perfil é válido
func IsValidPerfil(perfil string) bool {
	for _, validPerfil := range ValidPerfis {
		if perfil == validPerfil {
			return true
		}
	}
	return false
}

// HasAdminPermission verifica se o usuário tem permissão de admin
func (u *User) HasAdminPermission() bool {
	return u.Perfil == PERFIL_ADMIN
}

// HasUserPermission verifica se o usuário tem permissão de user ou superior
func (u *User) HasUserPermission() bool {
	return u.Perfil == PERFIL_USER || u.Perfil == PERFIL_ADMIN
}

// ToUserResponse converte User para UserResponse com verificações de permissão
func (u *User) ToUserResponse() UserResponse {
	return UserResponse{
		ID:             u.ID,
		Nome:           u.Nome,
		Email:          u.Email,
		CPF:            u.CPF,
		DataNascimento: u.DataNascimento,
		Perfil:         u.Perfil,
		IsAdmin:        u.HasAdminPermission(),
		HasPermission:  u.HasUserPermission(),
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
	}
}
