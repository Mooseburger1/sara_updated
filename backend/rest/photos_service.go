package main

import (
	"log"
	"net/http"
	"os"
	protos "sara_updated/backend/grpc/proto/photos"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var logger = log.New(os.Stdout, "rest-server-photos", log.LstdFlags)

type photoService struct {
	pc protos.PhotoServiceClient
}

func NewPhotoService() *photoService {
	ps := &photoService{}
	return ps
}

// InitServiceAndReturnCloseFunc intializes any connections to backend services while creating any
// needed dependencies. It returns a function to be deferred on in order to close any open connections
// and clean up and resources as necessary
func (ps *photoService) InitServiceAndReturnCloseFunc() func() {

	// Initialze the connection to the photos grpc server
	conn, err := grpc.Dial("grpc_backend:4000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	// Create the photos grpc client and persist
	ps.pc = protos.NewPhotoServiceClient(conn)

	return func() { conn.Close() }
}

func (ps *photoService) RegisterGetRoutes(getRouter *mux.Router) {
	//route for listing albums - optional params {pageSize | pageToken}
	getRouter.HandleFunc("/photos/albums/list", listAlbums)
}

func listAlbums(w http.ResponseWriter, r *http.Request) {
	//TODO
}
