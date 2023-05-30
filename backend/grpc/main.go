package main

import (
	"log"
	"os"
	"os/signal"

	"sara_updated/backend/common"
	photo_server "sara_updated/backend/grpc/photos/google"
)

func startServers(servers []common.SaraServer) {
	for _, server := range servers {
		go server.StartServer()
	}
}

func shutdownServers(servers []common.SaraServer) {
	for _, server := range servers {
		server.ShutdownServer()
	}
}

func main() {
	logger := log.New(os.Stdout, "grpc-server-manager", log.LstdFlags)

	// Initialize gRPC servers
	servers := []common.SaraServer{
		photo_server.NewPhotosServer(),
	}

	// Start Servers
	logger.Print("Starting RPC servers")
	startServers(servers)

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan

	logger.Printf("Received terminate signal of %s. Performing graceful shutdown", sig)
	shutdownServers(servers)

}
