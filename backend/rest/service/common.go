package service

import (
	"context"
	protos "sara_updated/backend/grpc/proto/photos"
	"sara_updated/backend/grpc/proto/protoauth"
)

// All accepted URL query params
const (
	PAGE_TOKEN = "pageToken"
	PAGE_SIZE  = "pageSize"
)

// QueryParams is used to pass along any params extracted from API URL during
// an API request. It is passed through the layers in the context of each call
type QueryParams struct {
	PageSize  int32
	PageToken string
}

type PhotosService interface {
	ListAlbums(context.Context, *protoauth.OauthConfigInfo) (*protos.AlbumsInfo, error)
}

// Use this to set context params so as to avoid collisions
type ContextKey string
