package photo_server

import (
	"log"
	"net"
	"os"

	google_oauth "github.com/Mooseburger1/sara_updated/backend/common"
	proto_photo "github.com/Mooseburger1/sara_updated/backend/grpc/proto/photos"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const PORT = "GOOGLE_PHOTO_SERVER_PORT"

type Opts struct {
	clientCreator google_oauth.ClientFunc
	listener      net.Listener
}

type OptFunc func(*Opts)

type photosServer struct {
	Server *grpc.Server
	logger *log.Logger
	Opts
}

func defaultOpts() Opts {
	return Opts{
		clientCreator: google_oauth.CreateClient,
		listener:      createDefaultListener(),
	}
}

func NewPhotosServer(opts ...OptFunc) *photosServer {
	o := defaultOpts()

	for _, fn := range opts {
		fn(&o)
	}

	ps := photosServer{
		Opts: o,
	}

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

	s.logger.Print("Google photos grpc listening on port 4000")
	s.Server.Serve(s.Opts.listener)

}

func (s *photosServer) ShutdownServer() {
	s.Server.GracefulStop()
}

func createDefaultListener() net.Listener {
	l, err := net.Listen("tcp", ":4000")
	if err != nil {
		panic(err)
	}

	return l
}
