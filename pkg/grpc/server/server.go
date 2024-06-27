package server

import (
	"context"
	"log"
	"net"

	"github.com/AlexandrKobalt/trip-track/backend/file-server/pkg/duration"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type Config struct {
	Host              string           `validate:"required"`
	MaxConnectionIdle duration.Minutes `validate:"required"` // min
	Timeout           duration.Seconds `validate:"required"` // sec
	MaxConnectionAge  duration.Minutes `validate:"required"` // min
	Time              duration.Minutes `validate:"required"` // min
}

type Server struct {
	App      *grpc.Server
	listener net.Listener
}

func New(cfg Config) (server *Server, err error) {
	listener, err := net.Listen("tcp", cfg.Host)
	if err != nil {
		return nil, err
	}

	app := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: cfg.MaxConnectionIdle.Duration,
			Timeout:           cfg.Timeout.Duration,
			MaxConnectionAge:  cfg.MaxConnectionAge.Duration,
			Time:              cfg.Time.Duration,
		}),
	)

	return &Server{
		listener: listener,
		App:      app,
	}, nil
}

func (s *Server) Start(ctx context.Context) error {
	go func() {
		if err := s.App.Serve(s.listener); err != nil {
			log.Fatal(err.Error())
		}
	}()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.App.GracefulStop()

	return nil
}
