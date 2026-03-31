package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"quote-api/internal/cache"
	"quote-api/internal/config"
	"quote-api/internal/handler"
	"quote-api/internal/logger"
	"quote-api/internal/middleware"
	"quote-api/internal/service"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.Load()

	// redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	// dependencies
	cache := cache.NewRedisCache(rdb)
	service := service.NewQuoteService()
	zapLogger, err := logger.NewZapLogger()
	if err != nil {
		log.Fatal(err)
	}
	handler := handler.NewQuoteHandler(service, cache, zapLogger)

	mux := http.NewServeMux()

	// routes
	mux.HandleFunc("/quote", handler.GetQuote)
	mux.HandleFunc("/health", handler.HealthCheck)

	// wrap with middleware
	wrapped := middleware.Recovery(zapLogger)(
		middleware.RequestID()(
			middleware.Recovery(zapLogger)(mux),
		),
	)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: wrapped,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("server shutdown failed: %v", err)
	}

	if err := rdb.Close(); err != nil {
		log.Printf("redis close error: %v", err)
	}

	_ = zapLogger.Sync()
}
