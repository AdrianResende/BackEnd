package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"smartpicks-backend/internal/services"
)

func UploadImageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendErrorResponse(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	const maxUploadSize = 5 << 20 // 5MB
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		sendErrorResponse(w, "Arquivo muito grande. Máximo 5MB", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		sendErrorResponse(w, "Erro ao receber arquivo: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validar tipo e extensão
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		sendErrorResponse(w, "Erro ao ler arquivo", http.StatusBadRequest)
		return
	}
	file.Seek(0, 0)

	contentType := http.DetectContentType(buffer)
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}
	if !allowedTypes[contentType] {
		sendErrorResponse(w, "Tipo de arquivo não suportado", http.StatusBadRequest)
		return
	}

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

	// Nome único do arquivo
	timestamp := time.Now().UnixNano()
	newFileName := fmt.Sprintf("palpite_%d%s", timestamp, ext)

	// Upload para S3
	s3Service, err := services.NewS3Service()
	if err != nil {
		sendErrorResponse(w, "Erro ao configurar AWS S3: "+err.Error(), http.StatusInternalServerError)
		return
	}

	imageURL, err := s3Service.UploadFile(file, newFileName, contentType)
	if err != nil {
		sendErrorResponse(w, "Erro ao enviar para S3: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Resposta com a URL pública do S3
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"image_url": imageURL,
		"message":   "Upload realizado com sucesso para o S3",
	})
}
