package main

import (
	"context"
	"net/http"
	"sara_updated/backend/grpc/proto/protoauth"
	"sara_updated/backend/rest/service"

	"github.com/gorilla/mux"
)

type router struct {
	ps service.PhotosService
}

func NewApiRouter(ps service.PhotosService) *router {
	return &router{ps: ps}
}

func (r *router) RegisterGetRoutes(get *mux.Router) {
	get.HandleFunc("/api/v1/ListAlbums", r.listAlbums)
}

func (r *router) listAlbums(w http.ResponseWriter, req *http.Request) {
	temp := &protoauth.OauthConfigInfo{}
	r.ps.ListAlbums(context.Background(), temp)
}
