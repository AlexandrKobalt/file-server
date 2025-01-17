package grpc

import (
	"context"

	"github.com/AlexandrKobalt/trip-track_file-server/internal/file/service"
	fileserverproto "github.com/AlexandrKobalt/trip-track_proto/fileserver"
)

type Server struct {
	service service.IService
	fileserverproto.FileServer
}

func New(service service.IService) fileserverproto.FileServer {
	return &Server{service: service}
}

func (s *Server) Upload(
	_ context.Context,
	request *fileserverproto.UploadRequest,
) (response *fileserverproto.UploadResponse, err error) {
	return s.service.Upload(request)
}

func (s *Server) GetURL(
	_ context.Context,
	request *fileserverproto.GetURLRequest,
) (response *fileserverproto.GetURLResponse, err error) {
	return s.service.GetURL(request)
}
