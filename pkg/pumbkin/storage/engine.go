package storage

import (
	"context"
	"sync"
)

type InMemoryEngine struct {
	m    sync.RWMutex
	data map[string]string
}

func NewInMemoryEngine() *InMemoryEngine {
	return &InMemoryEngine{m: sync.RWMutex{}, data: make(map[string]string)}
}

func (i *InMemoryEngine) Get(ctx context.Context, key string) (string, error) {
	i.m.RLock()
	defer i.m.RUnlock()

	val := i.data[key]
	return val, nil
}

func (i *InMemoryEngine) Set(ctx context.Context, key, value string) error {
	i.m.Lock()
	defer i.m.Unlock()

	i.data[key] = value
	return nil
}

func (i *InMemoryEngine) Delete(ctx context.Context, key string) (bool, error) {
	i.m.Lock()
	defer i.m.Unlock()

	_, ok := i.data[key]
	if !ok {
		return false, nil
	}
	delete(i.data, key)
	return true, nil
}
