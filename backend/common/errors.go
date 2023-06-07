package common

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// httpClientCreationError represents an error experienced while creating
// an Oauth2 capable http.Client.
var httpClientCreationError = errors.New("Http client creation error")

// ClientCreationError returns a default httpClientCreationError
// TODO - Create a more robust error propagation design
func ClientCreationError() error {
	return httpClientCreationError
}

// RpcErrorResponse returns an error to be returned by a gRPC API endpoint.
// The exact error is dependent on a received http standard error (e.g. 404).
// The appropriate RPC that correlates to the http error is returned.
func RpcErrorResponse(statusCode int, desc string) *status.Status {
	switch statusCode {
	case 400:
		return status.New(codes.InvalidArgument, desc)
	default:
		return status.New(codes.InvalidArgument, desc)
	}
}
