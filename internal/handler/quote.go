package handler

import (
	"context"
	"net/http"
	"quote-api/internal/logger"
	"time"
)

type CacheStore interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
}

type QuoteService interface {
	RandomQuote(ctx context.Context) (string, error)
}

type QuoteHandler struct {
	service QuoteService
	cache   CacheStore
	logger  logger.Logger
}

func NewQuoteHandler(s QuoteService, c CacheStore, l logger.Logger) *QuoteHandler {
	return &QuoteHandler{
		service: s,
		cache:   c,
		logger:  l,
	}
}

func (h *QuoteHandler) GetQuote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// try cache
	val, err := h.cache.Get(ctx, "quote")
	if err == nil && val != "" {
		writeJSON(w, http.StatusOK, map[string]string{
			"quote":  val,
			"source": "cache",
		}, h.logger)
		return
	} else if err != nil {
		h.logger.Info("cache get failed", "error", err)
	}

	quote, err := h.service.RandomQuote(ctx)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError,
			map[string]string{"error": "failed to get quote"},
			h.logger)
		return
	}

	if err := h.cache.Set(ctx, "quote", quote, time.Minute); err != nil {
		h.logger.Error("cache set failed", "error", err)
	}

	writeJSON(w, http.StatusOK,
		map[string]string{"quote": quote,
			"source": "fresh"},
		h.logger)
}

func (h *QuoteHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	}, h.logger)
}
