package main

import (
	"net/http"
	"sara_updated/backend/common"
	"sara_updated/backend/grpc/proto/protoauth"
	"sara_updated/backend/rest/service"

	"github.com/gorilla/mux"
)

type router struct {
	ps service.PhotosService
	as service.AuthorizationService
}

func NewApiRouter(ps service.PhotosService, as service.AuthorizationService) *router {
	return &router{ps: ps, as: as}
}

func (r *router) RegisterGetRoutes(get *mux.Router) {
	get.HandleFunc("/api/v1/ListAlbums", r.as.IsAuthorized(r.ps.ListAlbums))
}

func (r *router) listAlbums(rpcHandler service.RpcHandlerFunc) service.OauthHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, o *protoauth.OauthConfigInfo) {
		pageToken := r.URL.Query().Get(service.PAGE_TOKEN)
		pagesize := r.URL.Query().Get(service.PAGE_SIZE)

		var qp service.QueryParams

		qp.PageToken = pageToken

		if pagesize != "" {
			i, err := common.Str2Int32(pagesize)
			if err != nil {
				panic(err)
			}
			qp.PageSize = i
		}
	}

}
