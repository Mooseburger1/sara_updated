package service

import (
	"context"
	protos "sara_updated/backend/grpc/proto/photos"
	"sara_updated/backend/grpc/proto/protoauth"
)

type QueryParams struct {
	PageSize  int32
	PageToken string
}

type PhotosService interface {
	ListAlbums(context.Context, *protoauth.OauthConfigInfo) (*protos.AlbumsInfo, error)
}
