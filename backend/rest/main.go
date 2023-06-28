package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	photos_service "sara_updated/backend/rest/service/photos"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	logger := log.New(os.Stdout, "server-manager", log.LstdFlags)

	var ps, photosConnCloser = photos_service.NewPhotosClient()
	defer photosConnCloser()

	var api = NewApiRouter(ps)

	//Serve Mux to replace the default ServeMux
	serveMux := mux.NewRouter()
	getRouter := serveMux.Methods(http.MethodGet).Subrouter()

	api.RegisterGetRoutes(getRouter)

	ch := handlers.CORS(handlers.AllowedOrigins([]string{"*"}))

	server := &http.Server{
		Addr:         ":9090",
		Handler:      ch(serveMux),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	//listen and server async
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan

	logger.Println("Received terminate signal, graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(tc)
}
