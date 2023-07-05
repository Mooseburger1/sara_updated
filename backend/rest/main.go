package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sara_updated/backend/rest/service"
	"sara_updated/backend/rest/service/authorization"
	photos_service "sara_updated/backend/rest/service/photos"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gopkg.in/boj/redistore.v1"
)

func photoServiceOptFuncBuilder(ps service.PhotosService) OptFunc {
	return func(o *Opts) {
		o.PhotoService = ps
	}
}

func authServiceOptFuncBuilder(as service.AuthorizationService) OptFunc {
	return func(o *Opts) {
		o.AuthService = as
	}
}

func main() {
	logger := log.New(os.Stdout, "server-manager", log.LstdFlags)

	store, err := redistore.NewRediStore(10, "tcp", "redis-server:6379", "", []byte("secret-key"))
	if err != nil {
		panic(err)
	}

	as := authorization.NewAuthMiddleware(store)

	ps, photosConnCloser := photos_service.NewPhotosClient()
	defer photosConnCloser()

	var api = NewApiRouter(photoServiceOptFuncBuilder(ps), authServiceOptFuncBuilder(as))

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
