package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type CacheStore interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
}

type QuoteService interface {
	RandomQuote(ctx context.Context) (string, error)
}

type Logger interface {
	Info(msg string, kv ...any)
	Error(msg string, kv ...any)
	Sync() error
}

type QuoteHandler struct {
	service QuoteService
	cache   CacheStore
	logger  Logger
}

func NewQuoteHandler(s QuoteService, c CacheStore, l Logger) *QuoteHandler {
	return &QuoteHandler{
		service: s,
		cache:   c,
		logger:  l,
	}
}

func (h *QuoteHandler) GetQuote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json") // move here
	// try cache
	val, err := h.cache.Get(ctx, "quote")
	if err != nil {
		h.logger.Error("Redis GET failed", "error", err)
	} else {
		h.logger.Info("Redis HIT", "value", val)
	}
	if err == nil {
		json.NewEncoder(w).Encode(map[string]string{
			"quote":  val,
			"source": "cache",
		})
		return
	}

	quote, err := h.service.RandomQuote(ctx)
	if err != nil {
		http.Error(w, "failed to get quote", http.StatusInternalServerError)
		return
	}

	_ = h.cache.Set(ctx, "quote", quote, time.Minute)

	json.NewEncoder(w).Encode(map[string]string{
		"quote":  quote,
		"source": "fresh",
	})
}

func (h *QuoteHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
