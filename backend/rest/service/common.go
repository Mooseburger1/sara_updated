package service

import (
	"context"
	"sara_updated/backend/grpc/proto/photos"
	"sara_updated/backend/grpc/proto/protoauth"
)

type QueryParams struct {
	PageSize  int32
	PageToken string
}

type PhotosService interface {
	ListAlbums(context.Context, *protoauth.OauthConfigInfo) (*photos.AlbumsInfo, error)
}
