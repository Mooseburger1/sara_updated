package photos

import (
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const PORT = "GOOGLE_PHOTO_SERVER_PORT"


type photosServer struct {
	Server *grpc.Server
	Logger *log.Logger
}

func (s *photosServer) StartServer() {
	reflection.Register(s.Server)
	l, err := net.Listen("tcp", os.Getenv("GOOGLE_PHOTO_SERVER_PORT"))
	if err != nil {
		s.Logger.Fatal(err)
		os.Exit(1)
	}
	s.Logger.Printf("Google photos grpc listening on %s", os.Getenv("GOOGLE_PHOTO_SERVER_PORT"))
	s.Server.Serve(l)

}

func (s *photosServer) ShutdownServer() {
	s.Server.GracefulStop()
}