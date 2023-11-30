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
	listAlbumsReponse    func(*photos.AlbumListRequest) (*photos.AlbumsInfo, error)
	getMediaInfoResponse func(*photos.GetMediaRequest) (*photos.MediaInfo, error)
}

func (m *mockPhotosServer) UpdateListAlbumResponse(f func(*photos.AlbumListRequest) (*photos.AlbumsInfo, error)) {
	m.listAlbumsReponse = f
}

func (m *mockPhotosServer) UpdateGetMediaInfoResponse(f func(*photos.GetMediaRequest) (*photos.MediaInfo, error)) {
	m.getMediaInfoResponse = f
}

func (m *mockPhotosServer) ListAlbums(ctx context.Context, a *photos.AlbumListRequest) (*photos.AlbumsInfo, error) {
	return m.listAlbumsReponse(a)
}

func (m *mockPhotosServer) GetAlbumMedia(ctx context.Context, g *photos.GetMediaRequest) (*photos.MediaInfo, error) {
	return m.getMediaInfoResponse(g)
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

type protoExpectation interface {
	*photos.AlbumsInfo | *photos.MediaInfo
}

type protoRequest interface {
	*photos.AlbumListRequest | *photos.GetMediaRequest
}

type expectation[T protoExpectation] struct {
	value T
	err   error
}

type paramsCheckFuncBuilder[T protoRequest] func(*testing.T) paramsCheckFunc[T]
type paramsCheckFunc[T protoRequest] func(T)

func noCheckFunc[T protoRequest](t *testing.T) paramsCheckFunc[T] { return func(T) {} }

func TestPhotosClient_listAlbums(t *testing.T) {
	tests := map[string]struct {
		ctx context.Context
		e   expectation[*photos.AlbumsInfo]
		paramsCheckFuncBuilder[*photos.AlbumListRequest]
	}{
		"validResponseWithPageToken": {
			context.WithValue(context.Background(), service.OAUTH_CONFIG_KEY, &protoauth.OauthConfigInfo{}),
			expectation[*photos.AlbumsInfo]{
				value: &photos.AlbumsInfo{GooglePhotosAlbums: &photos.GooglePhotosAlbums{NextPageToken: "hello world"}},
				err:   nil,
			},
			noCheckFunc[*photos.AlbumListRequest],
		},
		"validResponseWithAlbumInfo": {
			context.WithValue(context.Background(), service.OAUTH_CONFIG_KEY, &protoauth.OauthConfigInfo{}),
			expectation[*photos.AlbumsInfo]{
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
			noCheckFunc[*photos.AlbumListRequest],
		},
		"validResponseWithMultipleAlbumInfo": {
			context.WithValue(context.Background(), service.OAUTH_CONFIG_KEY, &protoauth.OauthConfigInfo{}),
			expectation[*photos.AlbumsInfo]{
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
			noCheckFunc[*photos.AlbumListRequest],
		},
		"returnsError": {
			context.WithValue(context.Background(), service.OAUTH_CONFIG_KEY, &protoauth.OauthConfigInfo{}),
			expectation[*photos.AlbumsInfo]{
				value: nil,
				err:   errors.New("Failure"),
			},
			noCheckFunc[*photos.AlbumListRequest],
		},
		"pageSizeQueryParam": {
			context.WithValue(context.WithValue(context.Background(), service.OAUTH_CONFIG_KEY, &protoauth.OauthConfigInfo{}), service.ContextKey("queryParams"), &service.QueryParams{PageSize: 3}),
			expectation[*photos.AlbumsInfo]{
				value: nil,
				err:   errors.New("Failure"),
			},
			func(t *testing.T) paramsCheckFunc[*photos.AlbumListRequest] {
				return func(a *photos.AlbumListRequest) {
					if a.GoogleRequest.PageSize != 3 {
						t.Error("QueryParamsCheck: Expected 3\nActual: ", a.GoogleRequest.PageSize)
					}
				}
			},
		},
		"pageTokenQueryParam": {
			context.WithValue(context.WithValue(context.Background(), service.OAUTH_CONFIG_KEY, &protoauth.OauthConfigInfo{}), service.ContextKey("queryParams"), &service.QueryParams{PageToken: "Sparta"}),
			expectation[*photos.AlbumsInfo]{
				value: nil,
				err:   errors.New("Failure"),
			},
			func(t *testing.T) paramsCheckFunc[*photos.AlbumListRequest] {
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
		m.UpdateListAlbumResponse(func(a *photos.AlbumListRequest) (*photos.AlbumsInfo, error) {
			tt.paramsCheckFuncBuilder(t)(a)
			return tt.e.value, tt.e.err
		})
		t.Run(scenario, func(t *testing.T) {
			client, _ := NewPhotosClient(oFunc)

			response, err := client.ListAlbums(tt.ctx)

			if (tt.e.err != nil) && (tt.e.err.Error() != status.Convert(err).Message()) {
				t.Error("error: expected", tt.e.err, "received", status.Convert(err).Message())

			}

			if !proto.Equal(tt.e.value, response) {
				t.Errorf("\nTest %s \n Expected: %q\nActual: %q\n", scenario, tt.e.value, response)
			}
		})
	}

}

func TestPhotosClient_getMedia(t *testing.T) {
	tests := map[string]struct {
		ctx context.Context
		e   expectation[*photos.MediaInfo]
		paramsCheckFuncBuilder[*photos.GetMediaRequest]
	}{
		"validResponseWithPageToken": {
			context.WithValue(context.Background(), service.OAUTH_CONFIG_KEY, &protoauth.OauthConfigInfo{}),
			expectation[*photos.MediaInfo]{
				value: &photos.MediaInfo{GoogleMediaInfo: &photos.GooglePhotosMediaInfo{NextPageToken: "hello"}},
				err:   nil,
			},
			noCheckFunc[*photos.GetMediaRequest],
		},
		"validResponseWithMediaInfo": {
			context.WithValue(context.Background(), service.OAUTH_CONFIG_KEY, &protoauth.OauthConfigInfo{}),
			expectation[*photos.MediaInfo]{
				value: &photos.MediaInfo{GoogleMediaInfo: &photos.GooglePhotosMediaInfo{MediaItems: []*photos.Media{
					{
						Id:          "MediaId",
						ProductUrl:  "MediaURL",
						MimeType:    "MediaMimeType",
						Description: "MediaDescription",
					},
				}, NextPageToken: "hello"}},
				err: nil,
			},
			noCheckFunc[*photos.GetMediaRequest],
		},
		"validResponseWithMultipleMediaInfo": {
			context.WithValue(context.Background(), service.OAUTH_CONFIG_KEY, &protoauth.OauthConfigInfo{}),
			expectation[*photos.MediaInfo]{
				value: &photos.MediaInfo{GoogleMediaInfo: &photos.GooglePhotosMediaInfo{MediaItems: []*photos.Media{
					{
						Id:          "MediaId",
						ProductUrl:  "MediaURL",
						MimeType:    "MediaMimeType",
						Description: "MediaDescription",
					},
					{
						Id:          "MediaId2",
						ProductUrl:  "MediaURL2",
						MimeType:    "MediaMimeType2",
						Description: "MediaDescription2",
					},
				}, NextPageToken: "hello"}},
				err: nil,
			},
			noCheckFunc[*photos.GetMediaRequest],
		},
		"returnsError": {
			context.WithValue(context.Background(), service.OAUTH_CONFIG_KEY, &protoauth.OauthConfigInfo{}),
			expectation[*photos.MediaInfo]{
				value: nil,
				err:   errors.New("Failure"),
			},
			noCheckFunc[*photos.GetMediaRequest],
		},
		"pageSizeQueryParam": {
			context.WithValue(context.WithValue(context.Background(), service.OAUTH_CONFIG_KEY, &protoauth.OauthConfigInfo{}),
				service.ContextKey("mediaParams"), &service.GetAlbumMediaParams{Qp: &service.QueryParams{PageSize: 3}}),
			expectation[*photos.MediaInfo]{
				value: nil,
				err:   errors.New("Failure"),
			},
			func(t *testing.T) paramsCheckFunc[*photos.GetMediaRequest] {
				return func(a *photos.GetMediaRequest) {
					if a.GoogleRequest.PageSize != 3 {
						t.Error("QueryParamsCheck: Expected 3\nActual: ", a.GoogleRequest.PageSize)
					}
				}
			},
		},
		"pageTokenQueryParam": {
			context.WithValue(context.WithValue(context.Background(), service.OAUTH_CONFIG_KEY, &protoauth.OauthConfigInfo{}),
				service.ContextKey("mediaParams"), &service.GetAlbumMediaParams{Qp: &service.QueryParams{PageToken: "Sparta"}}),
			expectation[*photos.MediaInfo]{
				value: nil,
				err:   errors.New("Failure"),
			},
			func(t *testing.T) paramsCheckFunc[*photos.GetMediaRequest] {
				return func(a *photos.GetMediaRequest) {
					if a.GoogleRequest.PageToken != "Sparta" {
						t.Error("QueryParamsCheck: Expected 3\nActual: ", a.GoogleRequest.PageToken)
					}
				}
			},
		},
		"AlbumIdQueryParam": {
			context.WithValue(context.WithValue(context.Background(), service.OAUTH_CONFIG_KEY, &protoauth.OauthConfigInfo{}),
				service.ContextKey("mediaParams"), &service.GetAlbumMediaParams{AlbumId: "AlbumID", Qp: &service.QueryParams{}}),
			expectation[*photos.MediaInfo]{
				value: nil,
				err:   errors.New("Failure"),
			},
			func(t *testing.T) paramsCheckFunc[*photos.GetMediaRequest] {
				return func(a *photos.GetMediaRequest) {
					if a.GoogleRequest.AlbumId != "AlbumID" {
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
		m.UpdateGetMediaInfoResponse(func(a *photos.GetMediaRequest) (*photos.MediaInfo, error) {
			tt.paramsCheckFuncBuilder(t)(a)
			return tt.e.value, tt.e.err
		})
		t.Run(scenario, func(t *testing.T) {
			client, _ := NewPhotosClient(oFunc)

			response, err := client.GetAlbumMedia(tt.ctx)

			if (tt.e.err != nil) && (tt.e.err.Error() != status.Convert(err).Message()) {
				t.Error("error: expected", tt.e.err, "received", status.Convert(err).Message())

			}

			if !proto.Equal(tt.e.value, response) {
				t.Errorf("\nTest %s \n Expected: %q\nActual: %q\n", scenario, tt.e.value, response)
			}
		})
	}

}
