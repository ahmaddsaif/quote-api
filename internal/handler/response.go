package handler

import (
	"encoding/json"
	"net/http"
	"quote-api/internal/logger"
)

func writeJSON(w http.ResponseWriter, status int, data any, logger logger.Logger) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Error("failed to write response", "error", err)
	}
}
