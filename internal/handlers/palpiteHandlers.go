// handlers/palpite.go
package handlers

import (
	"encoding/json"
	"net/http"
	"smartpicks-backend/internal/database"
	"smartpicks-backend/internal/models"
)

// PostPalpite @Summary Criar um novo palpite
// @Description Cria um novo palpite no sistema
// @Tags Palpites
// @Accept json
// @Produce json
// @Param palpite body models.CreatePalpiteRequest true "Dados do palpite"
// @Success 201 {object} map[string]interface{} "Palpite criado com sucesso"
// @Failure 400 {object} map[string]string "Dados inválidos"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Router /palpites [post]
func PostPalpite(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendErrorResponse(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req models.CreatePalpiteRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		sendErrorResponse(w, "Requisição inválida: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validações
	if req.UserID == 0 {
		sendErrorResponse(w, "user_id é obrigatório", http.StatusBadRequest)
		return
	}

	if req.ImgURL == "" {
		sendErrorResponse(w, "img_url é obrigatório", http.StatusBadRequest)
		return
	}

	palpite := req.ToPalpite()

	// Inserir no banco
	err = database.DB.QueryRow(`
		INSERT INTO palpites (user_id, titulo, img_url, link, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`,
		palpite.UserID,
		palpite.Titulo,
		palpite.ImgURL,
		palpite.Link,
		palpite.CreatedAt,
		palpite.UpdatedAt,
	).Scan(&palpite.ID)

	if err != nil {
		sendErrorResponse(w, "Erro ao salvar palpite: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := palpite.ToResponse()
	sendSuccessResponse(w, map[string]interface{}{
		"palpite": resp,
		"message": "Palpite criado com sucesso",
	})
}
