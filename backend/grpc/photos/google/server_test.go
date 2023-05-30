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

type mockServer struct {
	listResponse  *photos.AlbumsInfo
	albumResponse *photos.MediaInfo
}

func (m *mockServer) ListAlbums(ctx context.Context,
	rpc *photos.AlbumListRequest) (*photos.AlbumsInfo, error) {
	return m.listResponse, nil
}

func (m *mockServer) GetAlbumMedia(ctx context.Context,
	rpc *photos.FromAlbumRequest) (*photos.MediaInfo, error) {
	return m.albumResponse, nil
}

func listAlbumsApiCreator(lr *photos.AlbumsInfo) *mockServer {
	return apiServiceCreator(lr, &photos.MediaInfo{})
}

func getAlbumApiCreator(ar *photos.MediaInfo) *mockServer {
	return apiServiceCreator(&photos.AlbumsInfo{}, ar)
}

func apiServiceCreator(lr *photos.AlbumsInfo, ar *photos.MediaInfo) *mockServer {
	return &mockServer{
		listResponse:  lr,
		albumResponse: ar,
	}
}

func apiServiceOptFunc(lr *photos.AlbumsInfo) OptFunc {
	apiServer := listAlbumsApiCreator(lr)
	return func(options *Opts) {
		options.server_api = apiServer
	}
}

func opFuncListener(lis *bufconn.Listener) OptFunc {
	return func(options *Opts) {
		options.listener = lis
	}
}

func prepare(ctx context.Context, options []OptFunc) (photos.GooglePhotoServiceClient, func()) {
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

	client := photos.NewGooglePhotoServiceClient(conn)

	return client, closer
}

func TestNewPhotosServerListAlbums(t *testing.T) {

	type expectation struct {
		value string
		err   error
	}

	tests := map[string]struct {
		opts     []OptFunc
		in       *photos.AlbumListRequest
		expected expectation
	}{
		"Success": {
			opts: []OptFunc{apiServiceOptFunc(&photos.AlbumsInfo{NextPageToken: "hello"})},
			in:   &photos.AlbumListRequest{},
			expected: expectation{
				value: "hello",
				err:   nil,
			},
		},
	}

	for scenario, tt := range tests {
		t.Run(scenario, func(t *testing.T) {

			ctx := context.Background()
			client, closer := prepare(ctx, tt.opts)
			defer closer()

			req, err := client.ListAlbums(ctx, tt.in)

			if err != nil {
				if tt.expected.err.Error() != err.Error() {
					t.Errorf("\nTest %s\nExpected: %q\nActual: %q\n", scenario, tt.expected.err, err)
				}
				return
			}

			if tt.expected.value != req.NextPageToken {
				t.Errorf("\nTest %s\nExpected: %q\nActual: %q\n", scenario, tt.expected.value, req.NextPageToken)
			}

		})
	}
}
