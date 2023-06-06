package common

import (
	"encoding/json"
	"fmt"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func CreateClientCreationError(err error) *status.Status {
	st := status.New(codes.InvalidArgument, "Client creation error")
	desc := fmt.Sprintf("Error creating client for making REST calls to Google photos API: %s", err)
	v := &errdetails.ErrorInfo{Reason: desc}
	st, err = st.WithDetails(v)
	if err != nil {
		panic(fmt.Sprintf("Unexpected error attaching metadata: %v", err))
	}
	return st
}

func CreateErrorResponseError(statusCode int, response []byte) *status.Status {
	var rpcErrCode codes.Code
	var desc errResponse

	json.Unmarshal(response, &desc)

	switch statusCode {
	case 400:
		rpcErrCode = codes.InvalidArgument
	default:
		rpcErrCode = codes.InvalidArgument
	}

	return status.New(rpcErrCode, desc.Error.Message)
}

type errResponse struct {
	Error errDetails `json:"error"`
}

type errDetails struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
