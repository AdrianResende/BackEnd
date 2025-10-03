package models

import "time"

const (
	PERFIL_ADMIN = "admin"
	PERFIL_USER  = "user"
)

var ValidPerfis = []string{PERFIL_ADMIN, PERFIL_USER}

type User struct {
	ID             int       `json:"id"`
	Nome           string    `json:"nome"`
	Email          string    `json:"email"`
	Password       string    `json:"password,omitempty"`
	CPF            string    `json:"cpf"`
	DataNascimento string    `json:"data_nascimento"`
	Perfil         string    `json:"perfil"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID             int       `json:"id"`
	Nome           string    `json:"nome"`
	Email          string    `json:"email"`
	CPF            string    `json:"cpf"`
	DataNascimento string    `json:"data_nascimento"`
	Perfil         string    `json:"perfil"`
	IsAdmin        bool      `json:"is_admin"`
	HasPermission  bool      `json:"has_permission"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func IsValidPerfil(perfil string) bool {
	for _, validPerfil := range ValidPerfis {
		if perfil == validPerfil {
			return true
		}
	}
	return false
}

func (u *User) HasAdminPermission() bool {
	return u.Perfil == PERFIL_ADMIN
}

func (u *User) HasUserPermission() bool {
	return u.Perfil == PERFIL_USER || u.Perfil == PERFIL_ADMIN
}

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
