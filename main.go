package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"quote-api/cache"
	"quote-api/handler"
	"quote-api/logger"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	logger.Init()
	cache.Init()
	defer logger.Log.Sync()

	mux := http.NewServeMux()
	mux.HandleFunc("/quote", handler.GetQuote)
	mux.HandleFunc("/health", handler.HealthCheck)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		logger.Log.Info("Server starting on: 8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	logger.Log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	server.Shutdown(ctx)
	logger.Log.Info("Server exited")
}
