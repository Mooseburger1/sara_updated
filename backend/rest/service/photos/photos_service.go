package photos_service

import (
	"context"
	"net/http"
	protos "sara_updated/backend/grpc/proto/photos"
	"sara_updated/backend/grpc/proto/protoauth"
	"sara_updated/backend/rest/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OptFunc func(*Opts)

// Opts persists all options set for the photos REST service
type Opts struct {
	ConnFunc       func() (*grpc.ClientConn, error)
	ListAlbumsFunc http.HandlerFunc
}

func defaultOpts() Opts {
	return Opts{
		ConnFunc: defaultConnFunc,
	}
}

func defaultConnFunc() (*grpc.ClientConn, error) {
	return grpc.Dial("grpc_backend:4000", grpc.WithTransportCredentials(insecure.NewCredentials()))
}

type photosClient struct {
	Opts
	pc protos.PhotoServiceClient
}

// NewPhotoClient intializes any connections to backend services while creating any
// needed dependencies. It returns a photoService object and a function to be deferred
// on in order to close any open connections and clean up and resources as necessary
func NewPhotosClient(opts ...OptFunc) (*photosClient, func()) {
	o := defaultOpts()

	for _, fn := range opts {
		fn(&o)
	}

	ps := &photosClient{
		Opts: o,
	}

	conn, err := o.ConnFunc()
	if err != nil {
		panic(err)
	}

	ps.pc = protos.NewPhotoServiceClient(conn)

	return ps, func() { conn.Close() }
}

func (p *photosClient) ListAlbums(ctx context.Context) (*protos.AlbumsInfo, error) {
	oa, ok := ctx.Value(service.OAUTH_CONFIG_KEY).(*protoauth.OauthConfigInfo)
	if !ok {
		panic("could not extract OauthConfig")
	}
	qp, ok := extractQueryParams(ctx)

	greq := &protos.GooglePhotosAlbumsRequest{}
	if ok {
		greq.PageSize = qp.PageSize
		greq.PageToken = qp.PageToken
	}

	req := &protos.AlbumListRequest{
		GoogleRequest: greq,
		OauthInfo:     oa,
	}

	return p.pc.ListAlbums(ctx, req)
}

func (p *photosClient) GetAlbumMedia(ctx context.Context) (*protos.MediaInfo, error) {
	oa, ok := ctx.Value(service.OAUTH_CONFIG_KEY).(*protoauth.OauthConfigInfo)
	if !ok {
		panic("could not extract OauthConfig")
	}

	mp, ok := extractMediaParams(ctx)
	greq := &protos.GooglePhotosMediaRequest{}
	if ok {
		greq.AlbumId = mp.AlbumId
		greq.PageSize = mp.Qp.PageSize
		greq.PageToken = mp.Qp.PageToken
	}

	req := &protos.GetMediaRequest{
		GoogleRequest: greq,
		OauthInfo:     oa,
	}

	return p.pc.GetAlbumMedia(ctx, req)
}

func extractQueryParams(ctx context.Context) (*service.QueryParams, bool) {
	v := ctx.Value(service.ContextKey("queryParams"))
	if v == nil {
		return nil, false
	}

	qp, ok := v.(*service.QueryParams)
	if !ok {
		panic("Recevied unexpected object in queryParams context")
	}

	return qp, true
}

func extractMediaParams(ctx context.Context) (*service.GetAlbumMediaParams, bool) {
	v := ctx.Value(service.ContextKey("mediaParams"))
	if v == nil {
		return nil, false
	}

	mp, ok := v.(*service.GetAlbumMediaParams)
	if !ok {
		panic("Recevied unexpected object in mediaParams context")
	}

	return mp, true
}
