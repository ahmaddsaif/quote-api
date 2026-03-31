package logger

import (
	"go.uber.org/zap"
)

type ZapLogger struct {
	log *zap.SugaredLogger
}

func NewZapLogger() (*ZapLogger, error) {
	l, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return &ZapLogger{
		log: l.Sugar(),
	}, nil
}

func (z *ZapLogger) Info(msg string, kv ...any) {
	z.log.Infow(msg, kv...)
}

func (z *ZapLogger) Error(msg string, kv ...any) {
	z.log.Errorw(msg, kv...)
}

func (z *ZapLogger) Sync() error {
	return z.log.Sync()
}
