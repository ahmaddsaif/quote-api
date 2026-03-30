package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"quote-api/cache"
	"quote-api/logger"
	"quote-api/service"
)

func GetQuote(w http.ResponseWriter, r *http.Request) {
	quote := service.RandomQuote()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"quote": quote,
	})
}
