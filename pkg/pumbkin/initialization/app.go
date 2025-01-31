package initialization

import (
	"context"
	"fmt"
	"simpledatabase/pkg/pumbkin/config"
	"simpledatabase/pkg/pumbkin/network"
)

type App struct {
	server *network.TcpServer
}

func CreateApp(config *config.Config) (*App, error) {
	logger, err := CreateLogger(&config.Logging)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}
	srv, err := CreateServer(&config.Network, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create server: %w", err)
	}
	return &App{
		server: srv,
	}, nil
}

func (app *App) Start(ctx context.Context) error {
	return app.server.Start(ctx)
}

func (app *App) Stop() error {
	return app.server.Stop()
}
