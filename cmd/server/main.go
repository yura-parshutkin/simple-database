package main

import (
	"context"
	"go.uber.org/zap"
	"log"
	"os/signal"
	"simpledatabase/pkg/pumbkin/config"
	"simpledatabase/pkg/pumbkin/initialization"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, errConf := config.Parse("./config.yaml")
	if errConf != nil {
		log.Fatal("failed to parse config %w", errConf.Error())
	}
	app, err := initialization.CreateApp(cfg)
	if err != nil {
		log.Fatal("failed to create app %w", err)
	}
	go func() {
		defer cancel()
		errSrv := app.Start(ctx)
		if errSrv != nil {
			log.Fatal("server error", zap.Error(errSrv))
		}
	}()
	<-ctx.Done()
}
