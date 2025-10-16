package models

import "time"

// Struct principal que representa a tabela no banco
type Palpite struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Titulo    *int      `json:"titulo,omitempty"` // Optional
	ImgURL    string    `json:"img_url"`          // URL da imagem
	Link      *string   `json:"link,omitempty"`   // Optional
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Request para criação
type CreatePalpiteRequest struct {
	UserID int     `json:"user_id" binding:"required"`
	Titulo *int    `json:"titulo,omitempty"`
	ImgURL string  `json:"img_url" binding:"required"` // Agora chamando de ImgURL para consistência
	Link   *string `json:"link,omitempty"`
}

// Response
type PalpiteResponse struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Titulo    *int      `json:"titulo,omitempty"`
	ImgURL    string    `json:"img_url"` // Mantém o mesmo nome
	Link      *string   `json:"link,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Método para converter Palpite para Response
func (p *Palpite) ToResponse() PalpiteResponse {
	return PalpiteResponse{
		ID:        p.ID,
		UserID:    p.UserID,
		Titulo:    p.Titulo,
		ImgURL:    p.ImgURL,
		Link:      p.Link,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

// Método para converter Request em Palpite
func (req *CreatePalpiteRequest) ToPalpite() Palpite {
	now := time.Now()
	return Palpite{
		UserID:    req.UserID,
		Titulo:    req.Titulo,
		ImgURL:    req.ImgURL, // Agora usando ImgURL
		Link:      req.Link,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
