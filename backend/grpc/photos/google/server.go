package photo_server

import (
	"log"
	"net"
	"net/http"
	"os"

	proto_photo "github.com/Mooseburger1/sara_updated/backend/grpc/proto/photos"
	"github.com/Mooseburger1/sara_updated/backend/grpc/proto/protoauth"
	google_oauth "github.com/Mooseburger1/sara_updated/backend/common"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const PORT = "GOOGLE_PHOTO_SERVER_PORT"

type ClientFunc func(*protoauth.OauthConfigInfo) (*http.Client, error)

type Opts struct {
	clientCreator ClientFunc
}

type photosServer struct {
	Server  *grpc.Server
	logger  *log.Logger
	Opts
}

func DefaultOpts() Opts {
	return Opts {
		clientCreator: google_oauth.CreateClient,
	}
}

func NewPhotosServer() *photosServer {
	ps := photosServer{}
	ps.initServer()
	return &ps
}

func (ps *photosServer) initServer() {
	logger := log.New(os.Stdout, "photos-rpc-server", log.LstdFlags)
	ps.logger = logger
	grpcServer := grpc.NewServer()
	photoServer := NewGphotoStub(logger)

	proto_photo.RegisterGooglePhotoServiceServer(grpcServer, photoServer)
	ps.Server = grpcServer
}

func (s *photosServer) StartServer() {
	reflection.Register(s.Server)
	l, err := net.Listen("tcp", os.Getenv("GOOGLE_PHOTO_SERVER_PORT"))
	if err != nil {
		s.logger.Fatal(err)
		os.Exit(1)
	}
	s.logger.Printf("Google photos grpc listening on %s", os.Getenv("GOOGLE_PHOTO_SERVER_PORT"))
	s.Server.Serve(l)

}

func (s *photosServer) ShutdownServer() {
	s.Server.GracefulStop()
}
