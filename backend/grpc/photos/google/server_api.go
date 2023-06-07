package photo_server

import (
	"context"

	"sara_updated/backend/common"
	"sara_updated/backend/grpc/proto/photos"
)

// GPhotosAPI is the implementation of the
// Google Photo RPC server. It implements the
// * ListAlbums service
// * GetAlbumMedia service
type GPhotosAPI struct {
	clientCreator common.ClientFunc
}

// Builder for instantiating a GPhotosAPI
func NewGPhotosApiStub(cf common.ClientFunc) *GPhotosAPI {
	return &GPhotosAPI{
		clientCreator: cf,
	}
}

// ListAlbums is a RPC service endpoint. It receives an AlbumListRequest
// proto and returns an AlbumsInfo proto. Internally it makes an Oauth2
// authorized REST request to the Google Photos API server for listing albums.
func (g *GPhotosAPI) ListAlbums(ctx context.Context,
	rpc *photos.AlbumListRequest) (*photos.AlbumsInfo, error) {

	//TODO
	return nil, nil

}

// GetAlbumMedia is a RPC service endpoint. It receives a
// GetMediaRequest proto and returns a MediaInfo proto. Internally
// it makes an Oauth2 authorized REST request to the Google Photos API
// server for listing photos from a specific album or directory.
func (g *GPhotosAPI) GetAlbumMedia(ctx context.Context,
	rpc *photos.GetMediaRequest) (*photos.MediaInfo, error) {

	//TODO
	return nil, nil
}
