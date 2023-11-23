package photo_server

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"sara_updated/backend/common"
	"sara_updated/backend/grpc/proto/photos"

	"google.golang.org/protobuf/encoding/protojson"
)

const (
	// Endpoint for Google photos REST API
	// TODO - move constants like these to a central location for scalability/supportability
	ALBUMS_ENDPOINT = "https://photoslibrary.googleapis.com/v1/albums"
	PHOTOS_ENDPOINT = "https://photoslibrary.googleapis.com/v1/mediaItems:search"
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

// GetAlbumMedia is a RPC service endpoint. It receives a
// GetMediaRequest proto and returns a MediaInfo proto. Internally
// it makes an Oauth2 authorized REST request to the Google Photos API
// server for listing photos from a specific album or directory.
func (g *GPhotosAPI) GetAlbumMedia(ctx context.Context,
	rpc *photos.GetMediaRequest) (*photos.MediaInfo, error) {

	client, err := g.clientCreator(rpc.GetOauthInfo())
	if err != nil {
		g.logger.Printf("Error creating oauth http client for grpc photos request: %v", err)
		return nil, common.ClientCreationError()
	}

	requestBody := []byte(fmt.Sprintf(`{"albumId":"%v", "pageSize":"%v", "pageToken":"%v"}`,
		rpc.GoogleRequest.GetAlbumId(), rpc.GoogleRequest.GetPageSize(), rpc.GoogleRequest.GetPageToken()))

	req, err := http.NewRequest(http.MethodPost, PHOTOS_ENDPOINT, bytes.NewBuffer(requestBody))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		g.logger.Printf("Call to get media returned status code %v, not %v", resp.StatusCode, http.StatusOK)
		bodyBytes, err := io.ReadAll(resp.Body)

		if err != nil {
			panic(err)
		}

		return nil, common.RpcErrorResponse(resp.StatusCode, string(bodyBytes)).Err()
	}

	r := &photos.GooglePhotosMediaInfo{}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		g.logger.Printf("%s", err)
		panic(err)
	}
	err = protojson.Unmarshal(bytes, r)
	if err != nil {
		g.logger.Printf("%s", err)
		panic(err)
	}

	return &photos.MediaInfo{GoogleMediaInfo: r}, nil
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

	req, err := http.NewRequest(http.MethodGet, ALBUMS_ENDPOINT, nil)
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

	r := &photos.GooglePhotosAlbums{}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		g.logger.Printf("%s", err)
		panic(err)
	}
	err = protojson.Unmarshal(bytes, r)
	if err != nil {
		g.logger.Printf("%s", err)
		panic(err)
	}
	return &photos.AlbumsInfo{GooglePhotosAlbums: r}, nil
}

// addQueryParams receives a pre-created http URL and appends
// the appropriate URL query params required with respect to the to
// the incoming RPC that triggered the api endpoint.
func addQueryParams(req *http.Request, rpc *photos.AlbumListRequest) {
	query := req.URL.Query()
	query.Add("pageToken", rpc.GetGoogleRequest().GetPageToken())
	query.Add("pageSize", strconv.Itoa(int(rpc.GetGoogleRequest().GetPageSize())))
	req.URL.RawQuery = query.Encode()
}
