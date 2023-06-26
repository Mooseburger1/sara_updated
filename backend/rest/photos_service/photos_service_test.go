package photos_service

import (
	"context"
	"log"
	"net"
	"sara_updated/backend/grpc/proto/photos"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

type mockPhotosServer struct {
	photos.UnimplementedPhotoServiceServer
	funcResponse
}

type funcResponse struct {
	listAlbumsFunc *photos.AlbumsInfo
}

func (m *mockPhotosServer) ListAlbums(ctx context.Context, a *photos.AlbumListRequest) (photos.AlbumsInfo, error) {
	return m.mockFuncs.listAlbumsFunc(ctx, a)
}

func (m *mockPhotosServer) GetAlbumMedia(ctx context.Context, g *photos.GetMediaRequest) (photos.MediaInfo, error) {
	return photos.MediaInfo{}, nil
}

func dialer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()

	photos.RegisterPhotoServiceServer(server, &)

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func TestPhotosClient_listAlbums(t *testing.T) {
	tests := []struct {
		name string
		updater mockFuncsUpdater
		err  error
	}{
		{
			"valid response",
			func(m *mockFuncsfunc){
				m.
			}
		},
	}

}
