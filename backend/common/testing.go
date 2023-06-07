package common

import "net/http"

// RoundTripFunc is to be used in unit tests that inject their own http client.
// The injected http.Client struct should set its "Transport" member with an
// implementation of a RoundTripFunc. Doing so allows testing of a component that
// makes an http request to receive the *http.Response as returned from the
// RoundTripFunc implementation
type RoundTripFunc func(req *http.Request) (*http.Response, error)

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
