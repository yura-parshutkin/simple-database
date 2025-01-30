package test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net"
	"simpledatabase/pkg/pumbkin/compute"
	"simpledatabase/pkg/pumbkin/handler"
	"simpledatabase/pkg/pumbkin/network"
	"simpledatabase/pkg/pumbkin/storage"
	"testing"
	"time"
)

func TestTcpServer_SetAndGetData(t *testing.T) {
	h := handler.NewHandler(storage.NewInMemoryEngine(), compute.NewParser())

	host := "localhost:8090"
	srv := network.NewTcpServer(host, 0, 0, 0, h)
	ctx := context.Background()

	go func() {
		err := srv.Start(ctx)
		require.NoError(t, err)
	}()
	// waiting for a server
	time.Sleep(1 * time.Second)

	buf := make([]byte, 1024)
	var err error
	var count int

	client, err := net.Dial("tcp", host)
	require.NoError(t, err)

	_, err = client.Write([]byte("SET foo bar"))
	require.NoError(t, err)

	count, err = client.Read(buf)
	require.NoError(t, err)
	assert.Equal(t, "[ok] true", string(buf[:count]))

	_, err = client.Write([]byte("GET foo"))
	require.NoError(t, err)

	count, err = client.Read(buf)
	require.NoError(t, err)
	assert.Equal(t, "[ok] bar", string(buf[:count]))

	srv.Stop()
}

func TestTcpServer_RateConnections(t *testing.T) {
	n := handler.NewHandler(storage.NewInMemoryEngine(), compute.NewParser())

	host := "localhost:8090"
	srv := network.NewTcpServer(host, 2, 0, 1*time.Second, n)
	ctx := context.Background()

	go func() {
		err := srv.Start(ctx)
		require.NoError(t, err)
	}()
	// waiting for a server
	time.Sleep(1 * time.Second)

	var err error
	buf := make([]byte, 1024)

	c1, err := net.Dial("tcp", host)
	require.NoError(t, err)
	time.Sleep(1 * time.Second)
	_, _ = c1.Write([]byte("GET foo1"))
	_, _ = c1.Read(buf)

	c2, err := net.Dial("tcp", host)
	require.NoError(t, err)
	_, _ = c2.Write([]byte("GET foo2"))
	_, _ = c2.Read(buf)

	c3, err := net.Dial("tcp", host)
	require.NoError(t, err)
	_, _ = c3.Write([]byte("GET foo3"))
	_, _ = c3.Read(buf)

	c4, err := net.Dial("tcp", host)
	require.NoError(t, err)
	_, _ = c4.Write([]byte("GET foo4"))
	_, _ = c4.Read(buf)

	srv.Stop()
}
