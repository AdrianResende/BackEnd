// handlers/upload.go
package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// UploadImageHandler @Summary Fazer upload de imagem
// @Description Faz upload de uma imagem e retorna a URL
// @Tags Upload
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "Arquivo de imagem"
// @Success 200 {object} models.UploadResponse "Upload realizado com sucesso"
// @Failure 400 {object} map[string]string "Erro no upload"
// @Failure 500 {object} map[string]string "Erro interno do servidor"
// @Router /upload [post]
func UploadImageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendErrorResponse(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Limitar tamanho do arquivo (5MB)
	const maxUploadSize = 5 << 20 // 5MB
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		sendErrorResponse(w, "Arquivo muito grande. Máximo 5MB", http.StatusBadRequest)
		return
	}

	// Obter arquivo do form
	file, handler, err := r.FormFile("image")
	if err != nil {
		sendErrorResponse(w, "Erro ao receber arquivo: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validar tipo do arquivo
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		sendErrorResponse(w, "Erro ao ler arquivo", http.StatusBadRequest)
		return
	}
	file.Seek(0, 0) // Voltar ao início do arquivo

	contentType := http.DetectContentType(buffer)
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}

	if !allowedTypes[contentType] {
		sendErrorResponse(w, "Tipo de arquivo não suportado. Use JPEG, PNG, GIF ou WebP", http.StatusBadRequest)
		return
	}

	// Validar extensão do arquivo
	ext := strings.ToLower(filepath.Ext(handler.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}

	if !allowedExts[ext] {
		sendErrorResponse(w, "Extensão de arquivo não permitida", http.StatusBadRequest)
		return
	}

	// Gerar nome único para o arquivo
	timestamp := time.Now().UnixNano()
	newFileName := fmt.Sprintf("palpite_%d%s", timestamp, ext)

	// Diretório de upload (criar se não existir)
	uploadDir := "./uploads/palpites"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		sendErrorResponse(w, "Erro ao criar diretório: "+err.Error(), http.StatusInternalServerError)
		return
	}

	filePath := filepath.Join(uploadDir, newFileName)

	// Salvar arquivo
	dst, err := os.Create(filePath)
	if err != nil {
		sendErrorResponse(w, "Erro ao salvar arquivo: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copiar arquivo
	_, err = io.Copy(dst, file)
	if err != nil {
		sendErrorResponse(w, "Erro ao salvar arquivo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Retornar URL da imagem (ajuste conforme seu domínio)
	imageURL := fmt.Sprintf("/uploads/palpites/%s", newFileName)

	// Response de sucesso
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"image_url": imageURL,
		"message":   "Upload realizado com sucesso",
	})
}

// ServeUploadedFiles serve arquivos estáticos da pasta uploads
func ServeUploadedFiles(w http.ResponseWriter, r *http.Request) {
	// Remover o prefixo /uploads/ para obter o caminho do arquivo
	filePath := strings.TrimPrefix(r.URL.Path, "/uploads/")

	// Construir o caminho completo
	fullPath := filepath.Join("./uploads", filePath)

	// Servir arquivo
	http.ServeFile(w, r, fullPath)
}
