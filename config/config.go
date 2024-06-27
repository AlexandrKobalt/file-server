package config

import (
	"encoding/json"
	"os"

	fileservice "github.com/AlexandrKobalt/trip-track/backend/file-server/internal/file/service"
	"github.com/AlexandrKobalt/trip-track/backend/file-server/pkg/duration"
	grpcserver "github.com/AlexandrKobalt/trip-track/backend/file-server/pkg/grpc/server"
	"github.com/go-playground/validator"
)

const (
	path = "config/config.json"
)

type Config struct {
	StartTimeout duration.Seconds `validate:"required"`
	StopTimeout  duration.Seconds `validate:"required"`

	GRPC    grpcserver.Config
	Service struct {
		File fileservice.Config
	}
}

func LoadConfig() (cfg *Config, err error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(jsonFile).Decode(&cfg)
	if err != nil {
		return nil, err
	}

	err = validator.New().Struct(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
