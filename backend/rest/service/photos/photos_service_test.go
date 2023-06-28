package photos_service

import (
	"context"
	"errors"
	"log"
	"net"
	"sara_updated/backend/grpc/proto/photos"
	"sara_updated/backend/grpc/proto/protoauth"
	"sara_updated/backend/rest/service"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/proto"
)

// TODO - SET QUERY PARAMS IN MOCKPHOTO SERVER FOR VALIDATION!

type mockPhotosServer struct {
	photos.UnimplementedPhotoServiceServer
	listAlbumsReponse func(*photos.AlbumListRequest) (*photos.AlbumsInfo, error)
}

func (m *mockPhotosServer) UpdateResponse(f func(*photos.AlbumListRequest) (*photos.AlbumsInfo, error)) {
	m.listAlbumsReponse = f
}

func (m *mockPhotosServer) ListAlbums(ctx context.Context, a *photos.AlbumListRequest) (*photos.AlbumsInfo, error) {
	return m.listAlbumsReponse(a)
}

func (m *mockPhotosServer) GetAlbumMedia(ctx context.Context, g *photos.GetMediaRequest) (*photos.MediaInfo, error) {
	return &photos.MediaInfo{}, nil
}

func makeConnOptFunc(conn *grpc.ClientConn) OptFunc {
	return func(o *Opts) {
		o.ConnFunc = func() (*grpc.ClientConn, error) {
			return conn, nil
		}
	}
}

func dialer(m *mockPhotosServer) func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()

	photos.RegisterPhotoServiceServer(server, m)

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

type expectation struct {
	value *photos.AlbumsInfo
	err   error
}

func TestPhotosClient_listAlbums(t *testing.T) {
	tests := map[string]struct {
		ctx context.Context
		e   expectation
	}{
		"validResponse": {
			context.Background(),
			expectation{
				value: &photos.AlbumsInfo{GooglePhotosAlbums: &photos.GooglePhotosAlbums{NextPageToken: "hello world"}},
				err:   nil,
			},
		},
		"returnsError": {
			context.Background(),
			expectation{
				value: nil,
				err:   errors.New("Failure"),
			},
		},
		"pageSizeQueryParam": {
			context.WithValue(context.Background(), "queryParams", service.QueryParams{PageSize: 3}),
			expectation{
				value: nil,
				err:   errors.New("Failure"),
			},
		},
	}

	m := &mockPhotosServer{}
	conn, err := grpc.DialContext(context.Background(), "", grpc.WithInsecure(), grpc.WithContextDialer(dialer(m)))
	defer conn.Close()
	if err != nil {
		panic(err)
	}

	oFunc := makeConnOptFunc(conn)

	for scenario, tt := range tests {
		m.UpdateResponse(func(*photos.AlbumListRequest) (*photos.AlbumsInfo, error) { return tt.e.value, tt.e.err })
		t.Run(scenario, func(t *testing.T) {
			client, _ := NewPhotosClient(oFunc)

			response, err := client.ListAlbums(tt.ctx, &protoauth.OauthConfigInfo{})

			if (tt.e.err != nil) && (tt.e.err.Error() != grpc.ErrorDesc(err)) {
				t.Error("error: expected", tt.e.err, "received", grpc.ErrorDesc(err))

			}

			if !proto.Equal(tt.e.value, response) {
				t.Errorf("\nTest %s \n Expected: %q\nActual: %q\n", scenario, tt.e.value, response)
			}
		})
	}

}
