package photo_server

import (
	"context"
	"io/ioutil"
	"net/http"
	"sara_updated/backend/common"
	"sara_updated/backend/grpc/proto/photos"
	"sara_updated/backend/grpc/proto/protoauth"
	"strings"
	"testing"
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
			in:       &photos.AlbumListRequest{},
			expected: &photos.AlbumsInfo{AlbumsInfo: []*photos.AlbumInfo{{Id: "123"}}},
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
					  },
					  {
						"id": "someId",
						"productUrl": "someProductUrl"
					  }
					]
				  }`)),
			},
			checks: nil},
	}

	for scenario, tt := range tests {
		ctx := context.Background()
		g := NewGPhotosApiStub(createClientFunc(tt.resp, tt.checks))

		value, err := g.ListAlbums(ctx, tt.in)

		if err != nil {
			t.Errorf("Got error of %v", err)
		}

		if value.AlbumsInfo[0].Id != tt.expected.AlbumsInfo[0].Id {
			t.Errorf("\nTest %s\nExpected: %q\nActual: %q\n", scenario, tt.expected, value)
		}
	}
}
