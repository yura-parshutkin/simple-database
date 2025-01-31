package initialization

import (
	"go.uber.org/zap"
	"simpledatabase/pkg/pumbkin/compute"
	"simpledatabase/pkg/pumbkin/config"
	"simpledatabase/pkg/pumbkin/handler"
	"simpledatabase/pkg/pumbkin/network"
	"simpledatabase/pkg/pumbkin/storage"
)

func CreateServer(config *config.Network, logger *zap.Logger) (*network.TcpServer, error) {
	return network.NewTcpServer(
		config.Address,
		config.MaxConnections,
		config.MaxMessageSize,
		config.IdleTimeout,
		handler.NewHandler(storage.NewInMemoryEngine(), compute.NewParser()),
		logger,
	), nil
}
