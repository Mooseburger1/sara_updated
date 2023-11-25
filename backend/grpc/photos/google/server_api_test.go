package photo_server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sara_updated/backend/common"
	"sara_updated/backend/grpc/proto/photos"
	"sara_updated/backend/grpc/proto/protoauth"
	"strings"
	"testing"

	"google.golang.org/protobuf/proto"
)

var (
	DEFAULT_ALBUM_LIST_REQUEST  = &photos.AlbumListRequest{}
	DEFAULT_ALBUM_LIST_RESPONSE = &photos.AlbumsInfo{GooglePhotosAlbums: &photos.GooglePhotosAlbums{}}
	DEFAULT_GET_MEDIA_REQUEST   = &photos.GetMediaRequest{}
	DEFAULT_MEDIA_INFO          = &photos.MediaInfo{GoogleMediaInfo: &photos.GooglePhotosMediaInfo{}}
)

type reqChecksFunc func(req *http.Request)
type testingClientFunc func(r *http.Response, c reqChecksFunc) common.ClientFunc

type albumsExpectation struct {
	value *photos.AlbumsInfo
	err   error
}

type mediaExpectation struct {
	value *photos.MediaInfo
	err   error
}

type albumsHttpResponse struct {
	value *http.Response
	err   error
}

func createClientFunc(r *http.Response, c reqChecksFunc) common.ClientFunc {
	return func(o *protoauth.OauthConfigInfo) (*http.Client, error) {

		return &http.Client{
			Transport: common.RoundTripFunc(func(req *http.Request) (*http.Response, error) {
				if c != nil {
					c(req)
				}
				return r, nil
			}),
		}, nil
	}
}

func TestListAlbums(t *testing.T) {
	tests := map[string]struct {
		in          *photos.AlbumListRequest
		expected    albumsExpectation
		resp        *http.Response
		checks      reqChecksFunc
		clientFunc  testingClientFunc
		shouldPanic bool
	}{
		"NoQueryParams": {
			in: DEFAULT_ALBUM_LIST_REQUEST,
			expected: albumsExpectation{value: &photos.AlbumsInfo{
				GooglePhotosAlbums: &photos.GooglePhotosAlbums{
					Albums: []*photos.GoogleAlbumInfo{
						{
							Id:                    "foo",
							Title:                 "bar",
							ProductUrl:            "baz",
							MediaItemsCount:       200,
							CoverPhotoBaseUrl:     "someUrl",
							CoverPhotoMediaItemId: "someOtherUrl"}}}},
				err: nil},
			resp: &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(strings.NewReader(`{
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
			clientFunc:  createClientFunc,
			shouldPanic: false},
		"MultipleAlbums": {
			in: DEFAULT_ALBUM_LIST_REQUEST,
			expected: albumsExpectation{value: &photos.AlbumsInfo{
				GooglePhotosAlbums: &photos.GooglePhotosAlbums{
					Albums: []*photos.GoogleAlbumInfo{
						{
							Id:                    "foo",
							Title:                 "bar",
							ProductUrl:            "baz",
							MediaItemsCount:       200,
							CoverPhotoBaseUrl:     "someUrl",
							CoverPhotoMediaItemId: "someOtherUrl"},
						{
							Id:                    "abc",
							Title:                 "def",
							ProductUrl:            "ghi",
							MediaItemsCount:       -100,
							CoverPhotoBaseUrl:     "sparta.com",
							CoverPhotoMediaItemId: "answer.to.life.is.42"}}}},
				err: nil},
			resp: &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(strings.NewReader(`{
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
						"id": "abc",
						"title": "def",
						"productUrl": "ghi",
						"mediaItemsCount": "-100",
						"coverPhotoBaseUrl": "sparta.com",
						"coverPhotoMediaItemId": "answer.to.life.is.42"
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
			clientFunc:  createClientFunc,
			shouldPanic: false},
		"QueryParams": {
			in: &photos.AlbumListRequest{
				GoogleRequest: &photos.GooglePhotosAlbumsRequest{
					PageSize: 10, PageToken: "Foo"}},
			expected: albumsExpectation{value: DEFAULT_ALBUM_LIST_RESPONSE, err: nil},
			resp: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{}`)),
			},
			checks: func(req *http.Request) {
				if pageSize := req.URL.Query().Get("pageSize"); pageSize != "10" {
					t.Errorf("Expected pageSize: %q\nActual: %q\n", "10", pageSize)
				}

				if pageToken := req.URL.Query().Get("pageToken"); pageToken != "Foo" {
					t.Errorf("Expected pageToken: %q\nActual: %q\n", "Foo", pageToken)
				}
			},
			clientFunc:  createClientFunc,
			shouldPanic: false},
		"ClientCreationErrorReturnsError": {
			in:       DEFAULT_ALBUM_LIST_REQUEST,
			expected: albumsExpectation{value: nil, err: common.ClientCreationError()},
			resp:     nil,
			checks:   nil,
			clientFunc: func(r *http.Response, c reqChecksFunc) common.ClientFunc {
				return func(o *protoauth.OauthConfigInfo) (*http.Client, error) {

					return nil, common.ClientCreationError()
				}
			},
			shouldPanic: false},
		"ClientRequestErrorShouldPanic": {
			in:       DEFAULT_ALBUM_LIST_REQUEST,
			expected: albumsExpectation{value: nil, err: nil},
			resp:     nil,
			checks:   nil,
			clientFunc: func(r *http.Response, c reqChecksFunc) common.ClientFunc {
				return func(o *protoauth.OauthConfigInfo) (*http.Client, error) {

					return &http.Client{
						Transport: common.RoundTripFunc(func(req *http.Request) (*http.Response, error) {

							return nil, errors.New("Unused")
						}),
					}, nil
				}
			},
			shouldPanic: true},
		"Non200StatusReturnsError": {
			in:       DEFAULT_ALBUM_LIST_REQUEST,
			expected: albumsExpectation{value: nil, err: common.RpcErrorResponse(http.StatusBadRequest, "Unused").Err()},
			resp: &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(strings.NewReader("Unused")),
			},
			checks:      nil,
			clientFunc:  createClientFunc,
			shouldPanic: false},
		"MalFormedJsonShouldPanic": {
			in:       DEFAULT_ALBUM_LIST_REQUEST,
			expected: albumsExpectation{value: nil, err: nil},
			resp: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`"albums":[],`)),
			},
			checks:      nil,
			clientFunc:  createClientFunc,
			shouldPanic: true},
	}

	for scenario, tt := range tests {
		t.Run(scenario, func(t *testing.T) {
			defer func() {
				r := recover()

				if (r == nil) && (tt.shouldPanic) {
					t.Errorf("%s should have panicked but did not!", scenario)
				}

				if (r != nil) && (!tt.shouldPanic) {
					t.Errorf("%s Test should not have panicked but did!", scenario)
				}
			}()

			var g *GPhotosAPI

			ctx := context.Background()
			g = NewGPhotosApiStub(tt.clientFunc(tt.resp, tt.checks))
			value, err := g.ListAlbums(ctx, tt.in)

			if err != nil {
				if !strings.Contains(tt.expected.err.Error(), err.Error()) {
					t.Errorf("\nTest %s\nExpected error: %v\nActual error: %v", scenario, tt.expected.err.Error(), err.Error())
				}
			}

			if !proto.Equal(value, tt.expected.value) {
				t.Errorf("\nTest %s\nExpected: %q\nActual: %q\n", scenario, tt.expected.value, value)
			}
		})
	}
}

func TestGetAlbumMedia(t *testing.T) {
	tests := map[string]struct {
		in          *photos.GetMediaRequest
		expected    mediaExpectation
		resp        *http.Response
		checks      reqChecksFunc
		clientFunc  testingClientFunc
		shouldPanic bool
	}{
		"NoQueryParams": {
			in: DEFAULT_GET_MEDIA_REQUEST,
			expected: mediaExpectation{
				value: &photos.MediaInfo{
					GoogleMediaInfo: &photos.GooglePhotosMediaInfo{
						MediaItems: []*photos.Media{
							{
								Id:          "hello",
								ProductUrl:  "http://www.someurl.com",
								MimeType:    "mp4",
								Description: "",
							},
						},
						NextPageToken: "page_token",
					},
				},
				err: nil},
			resp: &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(strings.NewReader(`{
					"mediaItems": [
					  {
						"id": "hello",
						"productUrl": "http://www.someurl.com",
						"baseUrl": "https://lh3.googleusercontent.com/somealbum",
						"mimeType": "mp4",
						"mediaMetadata": {
						  "creationTime": "2020-05-05T20:17:32Z",
						  "width": "4032",
						  "height": "3024",
						  "photo": {
							"cameraMake": "samsung",
							"cameraModel": "SM-N975U",
							"focalLength": 4.3,
							"apertureFNumber": 2.4,
							"isoEquivalent": 50,
							"exposureTime": "0.000703729s"
						  }
						},
						"filename": "20200505_161732.jpg"
					  }
					],
					"nextPageToken": "page_token"
				  }`)),
			},
			checks: func(req *http.Request) {
				host := req.Host
				path := req.URL.Path
				url := fmt.Sprintf("https://%s%s", host, path)
				if url != PHOTOS_ENDPOINT {
					t.Errorf("Expected URL: %q\nActual: %q\n", PHOTOS_ENDPOINT, url)
				}

				if req.Method != http.MethodPost {
					t.Errorf("Expected Verb: %q\nActual: %q\n", http.MethodPost, req.Method)
				}
			},
			clientFunc:  createClientFunc,
			shouldPanic: false,
		},
		"MultipleMedia": {
			in: DEFAULT_GET_MEDIA_REQUEST,
			expected: mediaExpectation{
				value: &photos.MediaInfo{
					GoogleMediaInfo: &photos.GooglePhotosMediaInfo{
						MediaItems: []*photos.Media{
							{
								Id:          "hello",
								ProductUrl:  "http://www.someurl.com",
								MimeType:    "mp4",
								Description: "",
							},
							{
								Id:          "goodbye",
								ProductUrl:  "http://www.someurl2.com",
								MimeType:    "jpg",
								Description: "",
							},
						},
						NextPageToken: "page_token",
					},
				},
				err: nil},
			resp: &http.Response{
				StatusCode: http.StatusOK,
				Body: io.NopCloser(strings.NewReader(`{
						"mediaItems": [
						  {
							"id": "hello",
							"productUrl": "http://www.someurl.com",
							"baseUrl": "https://lh3.googleusercontent.com/lr/",
							"mimeType": "mp4",
							"mediaMetadata": {
							  "creationTime": "2020-05-05T20:17:32Z",
							  "width": "4032",
							  "height": "3024",
							  "photo": {
								"cameraMake": "samsung",
								"cameraModel": "SM-N975U",
								"focalLength": 4.3,
								"apertureFNumber": 2.4,
								"isoEquivalent": 50,
								"exposureTime": "0.000703729s"
							  }
							},
							"filename": "20200505_161732.jpg"
						  },
						  {
							"id": "goodbye",
							"productUrl": "http://www.someurl2.com",
							"baseUrl": "https://lh3.googleusercontent.com/lr/",
							"mimeType": "jpg",
							"mediaMetadata": {
							  "creationTime": "2020-05-07T20:48:18Z",
							  "width": "2208",
							  "height": "2944",
							  "photo": {
								"cameraMake": "samsung",
								"cameraModel": "SM-N975U",
								"focalLength": 3.3,
								"apertureFNumber": 1.9,
								"isoEquivalent": 500,
								"exposureTime": "0.033333335s"
							  }
							},
							"filename": "20200507_164818.jpg"
						  }
						],
						"nextPageToken": "page_token"
					  }`)),
			},
			checks: func(req *http.Request) {
				host := req.Host
				path := req.URL.Path
				url := fmt.Sprintf("https://%s%s", host, path)
				if url != PHOTOS_ENDPOINT {
					t.Errorf("Expected URL: %q\nActual: %q\n", PHOTOS_ENDPOINT, url)
				}

				if req.Method != http.MethodPost {
					t.Errorf("Expected Verb: %q\nActual: %q\n", http.MethodPost, req.Method)
				}
			},
			clientFunc:  createClientFunc,
			shouldPanic: false},
		"QueryParams": {
			in: &photos.GetMediaRequest{
				GoogleRequest: &photos.GooglePhotosMediaRequest{
					PageSize: 10, PageToken: "Foo"}},
			expected: mediaExpectation{value: DEFAULT_MEDIA_INFO, err: nil},
			resp: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{}`)),
			},
			checks: func(req *http.Request) {
				body, _ := io.ReadAll(req.Body)
				var b getMediaRequestBody

				err := json.Unmarshal(body, &b)
				if err != nil {
					t.Errorf("Error parsing json in test: %v", err)
					return
				}

				if pageSize := b.PageSize; pageSize != 10 {
					t.Errorf("Expected pageSize: %q\nActual: %q\n", "10", pageSize)
				}

				if pageToken := b.PageToken; pageToken != "Foo" {
					t.Errorf("Expected pageToken: %q\nActual: %q\n", "Foo", pageToken)
				}
			},
			clientFunc:  createClientFunc,
			shouldPanic: false},
		"ClientCreationErrorReturnsError": {
			in:       DEFAULT_GET_MEDIA_REQUEST,
			expected: mediaExpectation{value: nil, err: common.ClientCreationError()},
			resp:     nil,
			checks:   nil,
			clientFunc: func(r *http.Response, c reqChecksFunc) common.ClientFunc {
				return func(o *protoauth.OauthConfigInfo) (*http.Client, error) {

					return nil, common.ClientCreationError()
				}
			},
			shouldPanic: false},
		"ClientRequestErrorShouldPanic": {
			in:       DEFAULT_GET_MEDIA_REQUEST,
			expected: mediaExpectation{value: nil, err: nil},
			resp:     nil,
			checks:   nil,
			clientFunc: func(r *http.Response, c reqChecksFunc) common.ClientFunc {
				return func(o *protoauth.OauthConfigInfo) (*http.Client, error) {

					return &http.Client{
						Transport: common.RoundTripFunc(func(req *http.Request) (*http.Response, error) {

							return nil, errors.New("Unused")
						}),
					}, nil
				}
			},
			shouldPanic: true},
		"Non200StatusReturnsError": {
			in:       DEFAULT_GET_MEDIA_REQUEST,
			expected: mediaExpectation{value: nil, err: common.RpcErrorResponse(http.StatusBadRequest, "Unused").Err()},
			resp: &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(strings.NewReader("Unused")),
			},
			checks:      nil,
			clientFunc:  createClientFunc,
			shouldPanic: false},
		"MalFormedJsonShouldPanic": {
			in:       DEFAULT_GET_MEDIA_REQUEST,
			expected: mediaExpectation{value: nil, err: nil},
			resp: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`"mediaItems":[],`)),
			},
			checks:      nil,
			clientFunc:  createClientFunc,
			shouldPanic: true},
	}

	for scenario, tt := range tests {
		t.Run(scenario, func(t *testing.T) {
			defer func() {
				r := recover()

				if (r == nil) && (tt.shouldPanic) {
					t.Errorf("%s should have panicked but did not!", scenario)
				}

				if (r != nil) && (!tt.shouldPanic) {
					t.Errorf("%s Test should not have panicked but did!", scenario)
				}
			}()

			var g *GPhotosAPI

			ctx := context.Background()
			g = NewGPhotosApiStub(tt.clientFunc(tt.resp, tt.checks))
			value, err := g.GetAlbumMedia(ctx, tt.in)

			if err != nil {
				if !strings.Contains(tt.expected.err.Error(), err.Error()) {
					t.Errorf("\nTest %s\nExpected error: %v\nActual error: %v", scenario, tt.expected.err.Error(), err.Error())
				}
			}

			if !proto.Equal(value, tt.expected.value) {
				t.Errorf("\nTest %s\nExpected: %q\nActual: %q\n", scenario, tt.expected.value, value)
			}
		})
	}
}
