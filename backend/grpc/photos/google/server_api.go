package photo_server

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"sara_updated/backend/common"
	"sara_updated/backend/grpc/proto/photos"

	"google.golang.org/protobuf/encoding/protojson"
)

const (
	ALBUMS_ENDPOINT = "https://photoslibrary.googleapis.com/v1/albums"
)

// GPhotosAPI is the implementation of the
// Google Photo RPC server. It implements the
// * ListAlbums service
// * GetAlbumMedia service
type GPhotosAPI struct {
	clientCreator common.ClientFunc
	logger        *log.Logger
}

// Builder for instantiating a GPhotosAPI
func NewGPhotosApiStub(cf common.ClientFunc) *GPhotosAPI {
	return &GPhotosAPI{
		clientCreator: cf,
		logger:        log.New(os.Stdout, "google-photos-grpc-server", log.LstdFlags),
	}
}

// ListAlbums is a RPC service endpoint. It receives an AlbumListRequest
// proto and returns an AlbumsInfo proto. Internally it makes an Oauth2
// authorized REST request to the Google Photos API server for listing albums.
func (g *GPhotosAPI) ListAlbums(ctx context.Context,
	rpc *photos.AlbumListRequest) (*photos.AlbumsInfo, error) {

	client, err := g.clientCreator(rpc.GetOauthInfo())
	if err != nil {
		g.logger.Printf("Error creating oauth http client for grpc photos request: %v", err)
		return nil, common.ClientCreationError()
	}

	req, err := http.NewRequest("GET", ALBUMS_ENDPOINT, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Accept", "application/json")
	addQueryParams(req, rpc)

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		g.logger.Printf("Call to List Albums returned status code %v, not %v", resp.StatusCode, http.StatusOK)
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		return nil, common.RpcErrorResponse(resp.StatusCode, string(bodyBytes)).Err()
	}

	respRpc := &photos.AlbumsInfo{}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	err = protojson.Unmarshal(bytes, respRpc)
	if err != nil {
		panic(err)
	}
	return respRpc, nil
}

// GetAlbumMedia is a RPC service endpoint. It receives a
// FromAlbumRequest proto and returns a PhotosInfo proto. Internally
// it makes an Oauth2 authorized rest request to the Google Photos API
// server for listing photos from a specific album
func (g *GPhotosAPI) GetAlbumMedia(ctx context.Context,
	rpc *photos.FromAlbumRequest) (*photos.MediaInfo, error) {

	//TODO
	return nil, nil
}

func addQueryParams(req *http.Request, rpc *photos.AlbumListRequest) {
	query := req.URL.Query()
	query.Add("pageToken", rpc.GetPageToken())
	query.Add("pageSize", strconv.Itoa(int(rpc.GetPageSize())))
	req.URL.RawQuery = query.Encode()
}
