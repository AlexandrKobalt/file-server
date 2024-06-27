package service

import (
	fileserverproto "github.com/AlexandrKobalt/trip-track_proto/fileserver"
)

type IService interface {
	Upload(
		request *fileserverproto.UploadRequest,
	) (response *fileserverproto.UploadResponse, err error)
	GetURL(
		request *fileserverproto.GetURLRequest,
	) (response *fileserverproto.GetURLResponse, err error)
}
