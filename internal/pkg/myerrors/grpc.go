package myerrors

import "google.golang.org/grpc/status"

type GrpcErr struct{}

func (g *GrpcErr) WriteError(err error) error {
	
	unWrappedErr := unWrap(err)

	return status.New(unWrappedErr.GrpcStatus(), err.Error()).Err()
}