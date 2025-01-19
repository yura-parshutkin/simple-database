package main

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os/signal"
	"simpledatabase/pkg/pumbkin/initialization"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logEntry := newStandardLogger(true)
	zap.ReplaceGlobals(logEntry)
	defer func() { _ = logEntry.Sync() }()

	srv := initialization.CreateServer("localhost:8090")
	go func() {
		defer cancel()
		err := srv.Start(ctx)
		if err != nil {
			log.Println("server start err:", err)
			return
		}
	}()
	<-ctx.Done()
}

func newStandardLogger(debug bool) *zap.Logger {
	cfg := zap.Config{
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.RFC3339TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:      false,
		DisableCaller:    false,
		Encoding:         "json",
	}

	if debug {
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		cfg.Development = true
	}

	logger, _ := cfg.Build()
	zap.ReplaceGlobals(logger)

	return logger
}
