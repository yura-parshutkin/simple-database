package initialization

import (
	"simpledatabase/pkg/pumbkin/compute"
	"simpledatabase/pkg/pumbkin/handler"
	"simpledatabase/pkg/pumbkin/network"
	"simpledatabase/pkg/pumbkin/storage"
)

func CreateServer(host string) *network.TcpServer {
	handler := handler.NewHandler(storage.NewInMemoryEngine(), compute.NewParser())
	return network.NewTcpServer(host, 0, 0, 0, handler)
}
