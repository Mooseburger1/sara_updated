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

	"google.golang.org/protobuf/proto"
)

type reqChecksFunc func(req *http.Request)

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
		in       *photos.AlbumListRequest
		expected *photos.AlbumsInfo
		resp     *http.Response
		checks   reqChecksFunc
	}{
		"NoQueryParams": {
			in: &photos.AlbumListRequest{},
			expected: &photos.AlbumsInfo{Albums: []*photos.AlbumInfo{{
				Id:                    "foo",
				Title:                 "bar",
				ProductUrl:            "baz",
				MediaItemsCount:       200,
				CoverPhotoBaseUrl:     "someUrl",
				CoverPhotoMediaItemId: "someOtherUrl"}}},
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
			}},
		"QueryParams": {
			in:       &photos.AlbumListRequest{PageSize: 10, PageToken: "Foo"},
			expected: &photos.AlbumsInfo{},
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
			}},
	}

	for scenario, tt := range tests {
		ctx := context.Background()
		g := NewGPhotosApiStub(createClientFunc(tt.resp, tt.checks))

		value, err := g.ListAlbums(ctx, tt.in)

		if err != nil {
			t.Errorf("Got error of %v", err)
		}

		if !proto.Equal(value, tt.expected) {
			t.Errorf("\nTest %s\nExpected: %q\nActual: %q\n", scenario, tt.expected, value)
		}
	}
}
