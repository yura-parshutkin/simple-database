package pumbkin

import (
	"bufio"
	"context"
	"fmt"
	"go.uber.org/zap"
	"net"
)

type Server struct {
	host    string
	handler *Handler
}

func NewServer(host string, handler *Handler) *Server {
	return &Server{host: host, handler: handler}
}

func (s *Server) Start(ctx context.Context) error {
	listen, err := net.Listen("tcp", s.host)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	defer func() { _ = listen.Close() }()
	zap.L().Info("server listening", zap.String("host", s.host))

	for {
		conn, errAc := listen.Accept()
		if errAc != nil {
			zap.L().Error("error accepting:", zap.Error(errAc))
			continue
		}
		go func() {
			defer func() { _ = conn.Close() }()
			sc := bufio.NewScanner(conn)
			for sc.Scan() {
				out, errHand := s.handler.Handle(ctx, sc.Text())
				if errHand != nil {
					zap.L().Error("connection read error:", zap.Error(errHand))
					continue
				}
				_, errWe := conn.Write([]byte(out + "\r\n"))
				if errWe != nil {
					zap.L().Error("connection write error:", zap.Error(errWe))
					continue
				}
			}
		}()
	}
}
