package main

import (
	"log"
	"os"

	"github.com/Mooseburger1/sara_updated/backend/common"
	photo_server "github.com/Mooseburger1/sara_updated/backend/grpc/photos/google"
)

func startServers(servers ...common.SaraServer) {
	for _, server := range servers {
		go server.StartServer()
	}
}

func main() {
	logger := log.New(os.Stdout, "grpc-server-manager", log.LstdFlags)

	// Initialize gRPC servers
	photoServer := photo_server.NewPhotosServer()

	// Start Servers
	logger.Print("Starting RPC servers")
	startServers(photoServer)
}
