package service

import (
	fileserverproto "github.com/AlexandrKobalt/trip-track/backend/proto/fileserver"
	"google.golang.org/protobuf/types/known/emptypb"
)

type IService interface {
	Upload(
		request *fileserverproto.UploadRequest,
	) (response *emptypb.Empty, err error)
	GetURL(
		request *fileserverproto.GetURLRequest,
	) (response *fileserverproto.GetURLResponse, err error)
}
