package models

import "time"

type Palpite struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Titulo    *string   `json:"titulo,omitempty"`
	ImgURL    string    `json:"img_url"`
	Link      *string   `json:"link,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreatePalpiteRequest struct {
	UserID int     `json:"user_id" binding:"required"`
	Titulo *string `json:"titulo,omitempty"`
	ImgURL string  `json:"img_url" binding:"required"`
	Link   *string `json:"link,omitempty"`
}

type PalpiteResponse struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Titulo    *string   `json:"titulo,omitempty"`
	ImgURL    string    `json:"img_url"`
	Link      *string   `json:"link,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UploadResponse struct {
	ImageURL string `json:"image_url"`
	Message  string `json:"message"`
}

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

func (req *CreatePalpiteRequest) ToPalpite() Palpite {
	now := time.Now()
	return Palpite{
		UserID:    req.UserID,
		Titulo:    req.Titulo,
		ImgURL:    req.ImgURL,
		Link:      req.Link,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
