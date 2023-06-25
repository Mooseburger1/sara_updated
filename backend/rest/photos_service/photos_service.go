package photos_service

import (
	"fmt"
	"log"
	"net/http"
	"os"
	protos "sara_updated/backend/grpc/proto/photos"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var logger = log.New(os.Stdout, "rest-server-photos", log.LstdFlags)

type OptFunc func(*Opts)

// Opts persists all options set for the photos REST service
type Opts struct {
	service        string
	port           int16
	creds          credentials.TransportCredentials
	listAlbumsFunc http.HandlerFunc
}

func defaultOpts() Opts {
	return Opts{
		service:        "grpc_backend",
		port:           4000,
		creds:          insecure.NewCredentials(),
		listAlbumsFunc: listAlbums,
	}
}

type photoService struct {
	Opts
	pc protos.PhotoServiceClient
}

// NewPhotoService intializes any connections to backend services while creating any
// needed dependencies. It returns a photoService object and a function to be deferred
// on in order to close any open connections and clean up and resources as necessary
func NewPhotoService(opts ...OptFunc) (*photoService, func()) {
	o := defaultOpts()

	for _, fn := range opts {
		fn(&o)
	}

	ps := &photoService{
		Opts: o,
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", o.service, o.port), grpc.WithTransportCredentials(o.creds))
	if err != nil {
		panic(err)
	}

	ps.pc = protos.NewPhotoServiceClient(conn)

	return ps, func() { conn.Close() }
}

func (ps *photoService) RegisterGetRoutes(getRouter *mux.Router) {
	//route for listing albums - optional params {pageSize | pageToken}
	getRouter.HandleFunc("/photos/albums/list", listAlbums)
}

func listAlbums(w http.ResponseWriter, r *http.Request) {
	//TODO
}
