package photo_server

import (
	"log"
	"net"
	"os"

	google_oauth "sara_updated/backend/common"
	proto_photo "sara_updated/backend/grpc/proto/photos"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// photoServer maintains the functionality and state of the grpc server for the
// photo service gRPC server/endpoint
type photosServer struct {
	Server *grpc.Server
	logger *log.Logger
	Opts
}

// Opts represents the different options, settings and configurations the photos gRPC server can be
// configured with
type Opts struct {
	clientCreator google_oauth.ClientFunc
	listener      net.Listener
	server_api    proto_photo.GooglePhotoServiceServer
}

// An OptFunc is a function meant to modify an existing Opts struct with new configurations
// prior to starting a photos gRPC server. It receives a pointer to an Opt struct and sets
// one or more fields of the struct. The modifications are in place, thus an OptFunc should
// not return anything.
type OptFunc func(*Opts)

func defaultOpts() Opts {
	return Opts{
		listener:   createDefaultListener(),
		server_api: NewGPhotosApiStub(google_oauth.CreateClient),
	}
}

func createDefaultListener() net.Listener {
	l, err := net.Listen("tcp", ":4000")
	if err != nil {
		panic(err)
	}

	return l
}

// NewPhotosServer takes in 0 or more OptFuncs to instantiate and configure the photos
// gRPC server. If no OptFuncs are supplied, the defaults will be utilized.
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

	proto_photo.RegisterGooglePhotoServiceServer(grpcServer, ps.Opts.server_api)
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
