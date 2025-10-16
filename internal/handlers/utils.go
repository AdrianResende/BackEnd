package handlers

import (
	"encoding/base64"
	"strings"
)

// isValidBase64 verifica se uma string é um Base64 válido
func isValidBase64(str string) bool {
	// Remover prefixo data URL se existir
	base64Data := str
	if strings.Contains(str, "base64,") {
		parts := strings.Split(str, "base64,")
		if len(parts) > 1 {
			base64Data = parts[1]
		}
	}

	// Tentar decodificar
	_, err := base64.StdEncoding.DecodeString(base64Data)
	return err == nil
}

// extractBase64Data extrai apenas a parte base64 de uma data URL
func extractBase64Data(dataURL string) string {
	if strings.Contains(dataURL, "base64,") {
		parts := strings.Split(dataURL, "base64,")
		if len(parts) > 1 {
			return parts[1]
		}
	}
	return dataURL
}
