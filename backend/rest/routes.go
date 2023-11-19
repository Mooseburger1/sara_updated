package main

import (
	"context"
	"net/http"
	"sara_updated/backend/common"
	"sara_updated/backend/rest/service"

	"github.com/gorilla/mux"
)

type AuthMiddleWare interface {
	EnsureAuthorized(http.Handler) http.Handler
	GoogleRedirectCallback(http.ResponseWriter, *http.Request)
}

// router registers and maintains all API endpoints for the publicly exposed URLs.
type router struct {
	Opts
}

// Opts holds the configurations for the router struct and sets the services in which router should
// use to fulfill requests
type Opts struct {
	PhotoService service.PhotosService
	Auth         AuthMiddleWare
}

// defaultOpts creates the default services the router will use if no Opts configurations are provided
// during router creation
func defaultOpts() Opts {
	return Opts{}
}

// OptFunc is a function that takes in a pointer to the router Options and sets a new options field.
// The router object uses the services set in Opts to fullfill API requests
type OptFunc func(*Opts)

// NewApiRouter receives zero or more OptFuncs to configure the API router services. It returns a
// pointer to a new router object to be utilized to register mux routers to handle GET, POST, etc
// requests.
func NewApiRouter(opts ...OptFunc) *router {
	o := defaultOpts()

	for _, fn := range opts {
		fn(&o)
	}

	r := router{
		Opts: o,
	}

	return &r
}

// RegisterGetRoutes expectes a GET mux subrouter. It utilizes the subrouter to handled all incoming
// GET requests for the specified routes.
func (r *router) RegisterGetRoutes(get *mux.Router) {
	// Register the auth Middleware
	get.Use(r.Opts.Auth.EnsureAuthorized)
	get.HandleFunc("/api/v1/ListAlbums", r.listAlbumsRouter(r.Opts.PhotoService.ListAlbums))
}

func (r *router) listAlbumsRouter(rpcHandler service.RpcAlbumsHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		pageToken := req.URL.Query().Get(service.PAGE_TOKEN)
		pagesize := req.URL.Query().Get(service.PAGE_SIZE)

		var qp *service.QueryParams
		qp.PageToken = pageToken

		if pagesize != "" {
			i, err := common.Str2Int32(pagesize)
			if err != nil {
				panic(err)
			}
			qp.PageSize = i
		}

		ctx := context.WithValue(req.Context(), service.ContextKey("queryParams"), qp)
		rpcHandler(ctx)
	}

}
