package service

import (
	"context"
	"net/http"
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

type GetAlbumMediaParams struct {
	Qp      *QueryParams
	AlbumId string
}

// PhotoService defines the capabilities any service should have in order to be
// classified as a photos service. It exposes methods to list any and all albums
// provided by the service.
type PhotosService interface {
	ListAlbums(context.Context) (*protos.AlbumsInfo, error)
	GetAlbumMedia(context.Context) (*protos.MediaInfo, error)
}

// AuthorizationService is reponsible for ensuring all API calls to provided services
// has the proper credentials and permissions to do so.
type AuthorizationService interface {
	IsAuthorized(OauthHandlerFunc) http.HandlerFunc
	Authenticate()
	RedirectCallback()
}

// OauthHandlerFunc is an amended http.HandlerFunc that takes in the typical params of
// a http.HandlerFunc plus a *OauthConfigInfo proto. It is the expected input handler for the
// OAuth middleware.
type OauthHandlerFunc func(http.ResponseWriter, *http.Request, *protoauth.OauthConfigInfo)

type RpcAlbumsHandlerFunc func(context.Context) (*protos.AlbumsInfo, error)
type RpcGetMediaHandlerFunc func(context.Context) (*protos.MediaInfo, error)

// Use this to set context params so as to avoid collisions
type ContextKey string

// Use this to get/set context params for get album media requests to avoid collisions
type GetMediaKey string

var OAUTH_CONFIG_KEY = ContextKey("OAUTH_CONFIG_KEY")
