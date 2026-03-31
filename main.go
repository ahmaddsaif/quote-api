package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"quote-api/cache"
	"quote-api/config"
	"quote-api/handler"
	"quote-api/logger"
	"quote-api/service"
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

	// routes
	http.HandleFunc("/quote", handler.GetQuote)
	http.HandleFunc("/health", handler.HealthCheck)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: nil,
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
