package common

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var httpClientCreationError = errors.New("Http client creation error")

func ClientCreationError() error {
	return httpClientCreationError
}

func RpcErrorResponse(statusCode int, desc string) *status.Status {
	switch statusCode {
	case 400:
		return status.New(codes.InvalidArgument, desc)
	default:
		return status.New(codes.InvalidArgument, desc)
	}
}
