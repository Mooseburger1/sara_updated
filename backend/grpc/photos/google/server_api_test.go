package photo_server

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"sara_updated/backend/common"
	"sara_updated/backend/grpc/proto/photos"
	"sara_updated/backend/grpc/proto/protoauth"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/proto"
)

var DEFAULT_ALBUM_LIST_REQUEST = &photos.AlbumListRequest{}
var DEFAULT_ALBUM_LIST_RESPONSE = &photos.AlbumsInfo{}

type reqChecksFunc func(req *http.Request)
type testingClientFunc func(r *http.Response, c reqChecksFunc) common.ClientFunc

type albumsExpectation struct {
	value *photos.AlbumsInfo
	err   error
}

func createClientFunc(r *http.Response, c reqChecksFunc) common.ClientFunc {
	return func(o *protoauth.OauthConfigInfo) (*http.Client, error) {

		return &http.Client{
			Transport: common.RoundTripFunc(func(req *http.Request) *http.Response {
				if c != nil {
					c(req)
				}
				return r
			}),
		}, nil
	}
}

func TestListAlbums(t *testing.T) {
	tests := map[string]struct {
		in         *photos.AlbumListRequest
		expected   albumsExpectation
		resp       *http.Response
		checks     reqChecksFunc
		clientFunc testingClientFunc
	}{
		"NoQueryParams": {
			in: DEFAULT_ALBUM_LIST_REQUEST,
			expected: albumsExpectation{value: &photos.AlbumsInfo{Albums: []*photos.AlbumInfo{{
				Id:                    "foo",
				Title:                 "bar",
				ProductUrl:            "baz",
				MediaItemsCount:       200,
				CoverPhotoBaseUrl:     "someUrl",
				CoverPhotoMediaItemId: "someOtherUrl"}}},
				err: nil},
			resp: &http.Response{
				StatusCode: http.StatusOK,
				Body: ioutil.NopCloser(strings.NewReader(`{
					"albums": [
					  {
						"id": "foo",
						"title": "bar",
						"productUrl": "baz",
						"mediaItemsCount": "200",
						"coverPhotoBaseUrl": "someUrl",
						"coverPhotoMediaItemId": "someOtherUrl"
					  }
					]
				  }`)),
			},
			checks: func(req *http.Request) {
				host := req.Host
				path := req.URL.Path
				url := fmt.Sprintf("https://%s%s", host, path)
				if url != ALBUMS_ENDPOINT {
					t.Errorf("Expected URL: %q\nActual: %q\n", ALBUMS_ENDPOINT, url)
				}

				if req.Method != "GET" {
					t.Errorf("Expected Verb: %q\nActual: %q\n", "GET", req.Method)
				}
			},
			clientFunc: nil},
		"QueryParams": {
			in:       &photos.AlbumListRequest{PageSize: 10, PageToken: "Foo"},
			expected: albumsExpectation{value: DEFAULT_ALBUM_LIST_RESPONSE, err: nil},
			resp: &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(strings.NewReader(`{}`)),
			},
			checks: func(req *http.Request) {
				if pageSize := req.URL.Query().Get("pageSize"); pageSize != "10" {
					t.Errorf("Expected pageSize: %q\nActual: %q\n", "10", pageSize)
				}

				if pageToken := req.URL.Query().Get("pageToken"); pageToken != "Foo" {
					t.Errorf("Expected pageToken: %q\nActual: %q\n", "Foo", pageToken)
				}
			},
			clientFunc: nil},
		"ClientCreationError": {
			in:       DEFAULT_ALBUM_LIST_REQUEST,
			expected: albumsExpectation{value: nil, err: common.CreateClientCreationError(fmt.Errorf("Unused")).Err()},
			resp: &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(strings.NewReader(`{}`)),
			},
			checks: nil,
			clientFunc: func(r *http.Response, c reqChecksFunc) common.ClientFunc {
				return func(o *protoauth.OauthConfigInfo) (*http.Client, error) {

					return nil, common.CreateClientCreationError(fmt.Errorf("Unused")).Err()
				}
			}},
	}

	for scenario, tt := range tests {
		var g *GPhotosAPI

		ctx := context.Background()

		if clientCreator := tt.clientFunc; clientCreator != nil {
			g = NewGPhotosApiStub(clientCreator(tt.resp, tt.checks))
		} else {
			g = NewGPhotosApiStub(createClientFunc(tt.resp, tt.checks))
		}

		value, err := g.ListAlbums(ctx, tt.in)

		if err != nil {
			if cmp.Equal(err, tt.expected.err, cmpopts.EquateErrors()) {
				t.Errorf("\nTest %s\nExpected error: %q\nActual error: %q", scenario, tt.expected.err, err)
			}

		}

		if !proto.Equal(value, tt.expected.value) {
			t.Errorf("\nTest %s\nExpected: %q\nActual: %q\n", scenario, tt.expected.value, value)
		}
	}
}
