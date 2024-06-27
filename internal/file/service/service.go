package service

import (
	"fmt"
	"os"
	"path/filepath"

	fileserverproto "github.com/AlexandrKobalt/trip-track/backend/proto/fileserver"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Config struct {
	UploadDirectory string `validate:"required"`
	BaseURL         string `validate:"required"`
}

type service struct {
	cfg Config
}

func New(cfg Config) IService {
	return &service{cfg: cfg}
}

func (s *service) Upload(
	request *fileserverproto.UploadRequest,
) (response *emptypb.Empty, err error) {
	key := uuid.New().String()
	filePath := filepath.Join(s.cfg.UploadDirectory, key)

	out, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("error on file create: %w", err)
	}
	defer out.Close()

	_, err = out.Write(request.GetFile())
	if err != nil {
		return nil, fmt.Errorf("error on file write: %w", err)
	}

	return response, nil
}

func (s *service) GetURL(
	request *fileserverproto.GetURLRequest,
) (response *fileserverproto.GetURLResponse, err error) {
	filePath := filepath.Join(s.cfg.UploadDirectory, request.GetKey())
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %w", err)
	}

	url := fmt.Sprintf("%s%s", s.cfg.BaseURL, request.GetKey())

	return &fileserverproto.GetURLResponse{
		Url: url,
	}, nil
}
