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

	// Validação dos campos obrigatórios
	if req.UserID == 0 || req.ImgURL == "" {
		sendErrorResponse(w, "Campos obrigatórios faltando: user_id e img_url são obrigatórios", http.StatusBadRequest)
		return
	}

	// Validar se a URL da imagem é válida (opcional)
	if !isValidImageURL(req.ImgURL) {
		sendErrorResponse(w, "URL da imagem inválida", http.StatusBadRequest)
		return
	}

	// Converter request para modelo Palpite
	palpite := req.ToPalpite()

	// Inserir no banco
	err = database.DB.QueryRow(`
		INSERT INTO palpites (user_id, titulo, img_url, link, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`,
		palpite.UserID,
		palpite.Titulo,
		palpite.ImgURL, // Agora usando ImgURL
		palpite.Link,
		palpite.CreatedAt,
		palpite.UpdatedAt,
	).Scan(&palpite.ID)

	if err != nil {
		sendErrorResponse(w, "Erro ao salvar palpite: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Converter para response
	resp := palpite.ToResponse()

	sendSuccessResponse(w, map[string]interface{}{
		"palpite": resp,
		"message": "Palpite criado com sucesso",
	})
}

// GetPalpites @Summary Listar todos os palpites
// @Description Retorna a lista de todos os palpites
// @Tags Palpites
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Lista de palpites"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Router /palpites [get]
func GetPalpites(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendErrorResponse(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	rows, err := database.DB.Query(`
		SELECT id, user_id, titulo, img_url, link, created_at, updated_at 
		FROM palpites 
		ORDER BY created_at DESC`)
	if err != nil {
		sendErrorResponse(w, "Erro ao buscar palpites: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var palpites []models.PalpiteResponse

	for rows.Next() {
		var p models.Palpite
		var titulo *int
		var link *string

		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&titulo,
			&p.ImgURL, // Agora scan para ImgURL
			&link,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			sendErrorResponse(w, "Erro ao ler palpite: "+err.Error(), http.StatusInternalServerError)
			return
		}

		p.Titulo = titulo
		p.Link = link

		resp := p.ToResponse()
		palpites = append(palpites, resp)
	}

	sendSuccessResponse(w, map[string]interface{}{
		"palpites": palpites,
		"total":    len(palpites),
		"message":  "Palpites listados com sucesso",
	})
}

// GetPalpiteByID @Summary Buscar palpite por ID
// @Description Retorna um palpite específico pelo seu ID
// @Tags Palpites
// @Accept json
// @Produce json
// @Param id path int true "ID do palpite"
// @Success 200 {object} map[string]interface{} "Palpite encontrado"
// @Failure 404 {object} map[string]string "Palpite não encontrado"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Router /palpites/{id} [get]
func GetPalpiteByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendErrorResponse(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Extrair ID da URL (depende do seu router)
	id := r.URL.Path[len("/palpites/"):]
	if id == "" {
		sendErrorResponse(w, "ID não fornecido", http.StatusBadRequest)
		return
	}

	var p models.Palpite
	var titulo *int
	var link *string

	err := database.DB.QueryRow(`
		SELECT id, user_id, titulo, img_url, link, created_at, updated_at 
		FROM palpites 
		WHERE id = $1`, id).Scan(
		&p.ID,
		&p.UserID,
		&titulo,
		&p.ImgURL, // Agora scan para ImgURL
		&link,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if err != nil {
		sendErrorResponse(w, "Palpite não encontrado", http.StatusNotFound)
		return
	}

	p.Titulo = titulo
	p.Link = link

	resp := p.ToResponse()

	sendSuccessResponse(w, map[string]interface{}{
		"palpite": resp,
		"message": "Palpite encontrado com sucesso",
	})
}

// Função auxiliar para validar URL de imagem (opcional)
func isValidImageURL(url string) bool {
	// Implementação básica - você pode expandir conforme necessário
	validExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	for _, ext := range validExtensions {
		if len(url) > len(ext) && url[len(url)-len(ext):] == ext {
			return true
		}
	}
	// Se for base64 ou outra URL, aceita (ou implemente validação mais robusta)
	return len(url) > 0
}
