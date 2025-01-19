package initialization

import (
	"simpledatabase/pkg/pumbkin"
	"simpledatabase/pkg/pumbkin/compute"
	"simpledatabase/pkg/pumbkin/storage"
)

func CreateServer(host string) *pumbkin.Server {
	handler := pumbkin.NewHandler(storage.NewInMemoryEngine(), compute.NewParser())
	return pumbkin.NewServer(host, handler)
}
