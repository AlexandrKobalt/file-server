package client

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	Host string
}

func New(cfg Config) (connection *grpc.ClientConn, err error) {
	return grpc.Dial(
		cfg.Host,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
}
