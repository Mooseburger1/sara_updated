package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	mw "sara_updated/backend/rest/middleware"
	"sara_updated/backend/rest/service"
	photos_service "sara_updated/backend/rest/service/photos"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gopkg.in/boj/redistore.v1"
)

func photoServiceOptFuncBuilder(ps service.PhotosService) OptFunc {
	return func(o *Opts) {
		o.PhotoService = ps
	}
}

func authServiceOptFuncBuilder(as AuthMiddleWare) OptFunc {
	return func(o *Opts) {
		o.Auth = as
	}
}

func main() {
	logger := log.New(os.Stdout, "server-manager", log.LstdFlags)

	// store, err := redistore.NewRediStore(10, "tcp", "192.168.224.2:6379", "", []byte("secret-key"))
	store, err := redistore.NewRediStore(10, "tcp", "redis-server:6379", "", []byte("secret-key"))
	if err != nil {
		panic(err)
	}

	// TODO - move the oauthconfig to a better spot
	am := mw.NewAuthMiddleware(mw.WithSessionStore(store),
		mw.WithOAuthConfig(&oauth2.Config{
			ClientID:     os.Getenv("GOOGLE_API_ID"),
			ClientSecret: os.Getenv("GOOGLE_API_SECRET"),
			RedirectURL:  "http://localhost:9090/oauth-callback",
			Scopes: []string{"https://www.googleapis.com/auth/photoslibrary.readonly",
				"https://www.googleapis.com/auth/calendar.readonly",
				"https://www.googleapis.com/auth/calendar.events.readonly",
			},
			Endpoint: google.Endpoint,
		}))

	ps, photosConnCloser := photos_service.NewPhotosClient()
	defer photosConnCloser()

	var api = NewApiRouter(photoServiceOptFuncBuilder(ps), authServiceOptFuncBuilder(am))

	//Serve Mux to replace the default ServeMux
	serveMux := mux.NewRouter()
	serveMux.HandleFunc("/oauth-callback", am.GoogleRedirectCallback)

	getRouter := serveMux.Methods(http.MethodGet).Subrouter()

	api.RegisterGetRoutes(getRouter)

	ch := handlers.CORS(handlers.AllowedOrigins([]string{"*"}))

	server := &http.Server{
		Addr:         ":9090",
		Handler:      ch(serveMux),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}

	//listen and server async
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()

	logger.Println("Listening on port 9090")
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan

	logger.Println("Received terminate signal, graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(tc)
}
