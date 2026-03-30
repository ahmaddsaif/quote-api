package cache

import (
	"context"
	"os"
	"quote-api/logger"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var Ctx = context.Background()
var Rdb *redis.Client

func Init() {
	addr := os.Getenv("REDIS_ADDR")
	logger.Log.Info("Redis addr", zap.String("addr", addr))
	if addr == "" {
		addr = "localhost:6379"
	}
	Rdb = redis.NewClient(&redis.Options{
		Addr: addr,
	})
}
