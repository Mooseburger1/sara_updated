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
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/proto"
)

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

type paramsCheckFuncBuilder func(*testing.T) paramsCheckFunc
type paramsCheckFunc func(*photos.AlbumListRequest)

func noCheckFunc(t *testing.T) paramsCheckFunc { return func(a *photos.AlbumListRequest) {} }

func TestPhotosClient_listAlbums(t *testing.T) {
	tests := map[string]struct {
		ctx context.Context
		e   expectation
		paramsCheckFuncBuilder
	}{
		"validResponseWithPageToken": {
			context.Background(),
			expectation{
				value: &photos.AlbumsInfo{GooglePhotosAlbums: &photos.GooglePhotosAlbums{NextPageToken: "hello world"}},
				err:   nil,
			},
			noCheckFunc,
		},
		"validResponseWithAlbumInfo": {
			context.Background(),
			expectation{
				value: &photos.AlbumsInfo{GooglePhotosAlbums: &photos.GooglePhotosAlbums{Albums: []*photos.GoogleAlbumInfo{{
					Id:                    "abc",
					Title:                 "123",
					ProductUrl:            "www.thisisatest.com",
					MediaItemsCount:       int32(4),
					CoverPhotoBaseUrl:     "somePhotoUrl",
					CoverPhotoMediaItemId: "someMediaItem",
				}}}},
				err: nil,
			},
			noCheckFunc,
		},
		"validResponseWithMultipleAlbumInfo": {
			context.Background(),
			expectation{
				value: &photos.AlbumsInfo{GooglePhotosAlbums: &photos.GooglePhotosAlbums{Albums: []*photos.GoogleAlbumInfo{{
					Id:                    "album1",
					Title:                 "children",
					ProductUrl:            "www.thisisatest.com",
					MediaItemsCount:       int32(4),
					CoverPhotoBaseUrl:     "somePhotoUrl",
					CoverPhotoMediaItemId: "someMediaItem",
				}, {
					Id:         "album2",
					Title:      "family",
					ProductUrl: "www.family.com",
				}}}},
				err: nil,
			},
			noCheckFunc,
		},
		"returnsError": {
			context.Background(),
			expectation{
				value: nil,
				err:   errors.New("Failure"),
			},
			noCheckFunc,
		},
		"pageSizeQueryParam": {
			context.WithValue(context.Background(), ContextKey("queryParams"), &service.QueryParams{PageSize: 3}),
			expectation{
				value: nil,
				err:   errors.New("Failure"),
			},
			func(t *testing.T) paramsCheckFunc {
				return func(a *photos.AlbumListRequest) {
					if a.GoogleRequest.PageSize != 3 {
						t.Error("QueryParamsCheck: Expected 3\nActual: ", a.GoogleRequest.PageSize)
					}
				}
			},
		},
		"pageTokenQueryParam": {
			context.WithValue(context.Background(), ContextKey("queryParams"), &service.QueryParams{PageToken: "Sparta"}),
			expectation{
				value: nil,
				err:   errors.New("Failure"),
			},
			func(t *testing.T) paramsCheckFunc {
				return func(a *photos.AlbumListRequest) {
					if a.GoogleRequest.PageToken != "Sparta" {
						t.Error("QueryParamsCheck: Expected 3\nActual: ", a.GoogleRequest.PageToken)
					}
				}
			},
		},
	}

	m := &mockPhotosServer{}
	conn, err := grpc.DialContext(context.Background(), "", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(dialer(m)))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	oFunc := makeConnOptFunc(conn)

	for scenario, tt := range tests {
		m.UpdateResponse(func(a *photos.AlbumListRequest) (*photos.AlbumsInfo, error) {
			tt.paramsCheckFuncBuilder(t)(a)
			return tt.e.value, tt.e.err
		})
		t.Run(scenario, func(t *testing.T) {
			client, _ := NewPhotosClient(oFunc)

			response, err := client.ListAlbums(tt.ctx, &protoauth.OauthConfigInfo{})

			if (tt.e.err != nil) && (tt.e.err.Error() != status.Convert(err).Message()) {
				t.Error("error: expected", tt.e.err, "received", status.Convert(err).Message())

			}

			if !proto.Equal(tt.e.value, response) {
				t.Errorf("\nTest %s \n Expected: %q\nActual: %q\n", scenario, tt.e.value, response)
			}
		})
	}

}
