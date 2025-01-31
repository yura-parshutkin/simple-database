package network

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net"
	"simpledatabase/pkg/pumbkin/concurrency"
	"simpledatabase/pkg/pumbkin/handler"
	"time"
)

type TcpServer struct {
	address        string
	maxMessageSize int
	idleTimeout    time.Duration
	handler        *handler.Handler
	semaphore      *concurrency.Semaphore
	listen         net.Listener
	logger         *zap.Logger
}

func NewTcpServer(
	address string,
	maxConnections int,
	maxMessageSize int,
	idleTimeout time.Duration,
	handler *handler.Handler,
	logger *zap.Logger,
) *TcpServer {
	if maxConnections == 0 {
		maxConnections = 100
	}
	if maxMessageSize == 0 {
		maxMessageSize = 1024 * 4
	}
	return &TcpServer{
		address:        address,
		maxMessageSize: maxMessageSize,
		idleTimeout:    idleTimeout,
		handler:        handler,
		semaphore:      concurrency.NewSemaphore(maxConnections),
		logger:         logger,
	}
}

func (s *TcpServer) Start(ctx context.Context) error {
	var err error
	s.listen, err = net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	defer func() { _ = s.listen.Close() }()
	s.logger.Info("server listening", zap.String("address", s.address))

	for {
		conn, errAc := s.listen.Accept()
		if errAc != nil {
			s.logger.Error("error accepting:", zap.Error(errAc))
			continue
		}
		s.semaphore.Acquire()
		go func(c net.Conn) {
			defer s.semaphore.Release()
			s.handle(ctx, conn)
		}(conn)
	}
}

func (s *TcpServer) handle(ctx context.Context, conn net.Conn) {
	defer func() {
		if v := recover(); v != nil {
			s.logger.Error("recover from panic", zap.Any("panic", v))
		}
	}()
	defer func() { _ = conn.Close() }()

	buff := make([]byte, s.maxMessageSize)
	if s.idleTimeout != 0 {
		errRe := conn.SetReadDeadline(time.Now().Add(s.idleTimeout))
		if errRe != nil {
			s.logger.Error("error set count read deadline:", zap.Error(errRe))
			return
		}
	}
	if s.idleTimeout != 0 {
		errWr := conn.SetWriteDeadline(time.Now().Add(s.idleTimeout))
		if errWr != nil {
			s.logger.Error("error set count write deadline:", zap.Error(errWr))
			return
		}
	}
	for {
		count, errRe := conn.Read(buff)
		if errRe == io.EOF {
			return
		} else if errRe != nil {
			s.logger.Error("error read:", zap.Error(errRe))
			return
		}
		if count == len(buff) {
			s.logger.Error("error max buffer size reached", zap.Int("count", count))
			return
		}
		out, errHand := s.handler.Handle(ctx, string(buff[:count]))
		if errHand != nil {
			s.logger.Error("connection count error:", zap.Error(errHand))
			return
		}
		_, errWr := conn.Write([]byte("[ok] " + out))
		if errWr != nil {
			s.logger.Error("connection write error:", zap.Error(errWr))
			return
		}
	}
}

func (s *TcpServer) Stop() error {
	return s.listen.Close()
}
