package main

import (
	"log"
	"os"

	"sara_updated/backend/common"
	photo_server "sara_updated/backend/grpc/photos/google"
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

	select {}
}
