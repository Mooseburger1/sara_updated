package photo_server

import (
	"context"
	"log"
	"net"
	"sara_updated/backend/grpc/proto/photos"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

var (
	DEFAULT_LIST_REQUEST   = &photos.AlbumListRequest{}
	DEFAULT_MEDIA_REQUEST  = &photos.GetMediaRequest{}
	DEFAULT_LIST_RESPONSE  = &photos.AlbumsInfo{}
	DEFAULT_MEDIA_RESPONSE = &photos.MediaInfo{}
)

type mockServer struct {
	listResponse  *photos.AlbumsInfo
	albumResponse *photos.MediaInfo
}

func (m *mockServer) ListAlbums(ctx context.Context,
	rpc *photos.AlbumListRequest) (*photos.AlbumsInfo, error) {
	return m.listResponse, nil
}

func (m *mockServer) GetAlbumMedia(ctx context.Context,
	rpc *photos.GetMediaRequest) (*photos.MediaInfo, error) {
	return m.albumResponse, nil
}

func apiServiceOptFunc(lr *photos.AlbumsInfo, ar *photos.MediaInfo) OptFunc {
	apiServer := &mockServer{
		listResponse:  lr,
		albumResponse: ar,
	}
	return func(options *Opts) {
		options.server_api = apiServer
	}
}

func opFuncListener(lis *bufconn.Listener) OptFunc {
	return func(options *Opts) {
		options.listener = lis
	}
}

func prepareTest(ctx context.Context, options []OptFunc) (*photos.PhotoServiceClient, func()) {
	buffer := 101024 * 1024
	listener := bufconn.Listen(buffer)

	options = append(options, opFuncListener(listener))
	server := NewPhotosServer(options...)
	go server.StartServer()

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("error connecting to server: %v", err)
	}

	closer := func() {
		err := listener.Close()
		if err != nil {
			log.Printf("error closing listener %v", err)
		}

		server.ShutdownServer()
	}

	client := photos.NewPhotoServiceClient(conn)

	return &client, closer
}

type expectation struct {
	value string
	err   error
}

type testFunc func(context.Context, *photos.PhotoServiceClient) (string, error)

func TestNewPhotosServerListAlbums(t *testing.T) {

	tests := map[string]struct {
		opts     []OptFunc
		test     testFunc
		expected expectation
	}{
		"listAlbums": {
			opts: []OptFunc{apiServiceOptFunc(&photos.AlbumsInfo{
				PageTokens: &photos.PageTokens{
					GoogleToken: &photos.GooglePageToken{
						NextPageToken: "hello"}}},
				DEFAULT_MEDIA_RESPONSE)},
			test: func(ctx context.Context, c *photos.PhotoServiceClient) (string, error) {
				resp, err := (*c).ListAlbums(ctx, &photos.AlbumListRequest{})
				return resp.GetPageTokens().GoogleToken.GetNextPageToken(), err
			},
			expected: expectation{
				value: "hello",
				err:   nil,
			},
		},
		"getMedia": {
			opts: []OptFunc{apiServiceOptFunc(DEFAULT_LIST_RESPONSE,
				&photos.MediaInfo{
					PageToken: &photos.PageTokens{
						GoogleToken: &photos.GooglePageToken{
							NextPageToken: "world"}}})},
			test: func(ctx context.Context, c *photos.PhotoServiceClient) (string, error) {
				resp, err := (*c).GetAlbumMedia(ctx, &photos.GetMediaRequest{})
				return resp.PageToken.GetGoogleToken().GetNextPageToken(), err
			},
			expected: expectation{
				value: "world",
				err:   nil,
			},
		},
	}

	for scenario, tt := range tests {
		t.Run(scenario, func(t *testing.T) {

			ctx := context.Background()
			client, closer := prepareTest(ctx, tt.opts)
			defer closer()

			val, err := tt.test(ctx, client)

			if err != nil {
				if tt.expected.err.Error() != err.Error() {
					t.Errorf("\nTest %s\nExpected: %q\nActual: %q\n", scenario, tt.expected.err, err)
				}
				return
			}

			if tt.expected.value != val {
				t.Errorf("\nTest %s\nExpected: %q\nActual: %q\n", scenario, tt.expected.value, val)
			}

		})
	}
}
