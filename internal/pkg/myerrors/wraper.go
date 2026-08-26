package myerrors

import (
	"google.golang.org/grpc/codes"
)

const (
	Internal      = 500
	NotFound      = 404
	BadRequest    = 400
	Unauthorized  = 401
	AlreadyExists = 409
)

func Wrap(err error, code int) error {
	return &wrappedErr{
		message: err.Error(),
		code:    code,
	}
}

func unWrap(err error) *wrappedErr {
	unWrappedErr, ok := err.(*wrappedErr)
	if !ok {
		return nil
	}
	return unWrappedErr
}

type wrappedErr struct {
	message string
	code    int
}

func (w *wrappedErr) Error() string {
	return w.message
}

func (w *wrappedErr) HttpCode() int {
	return w.code
}

func (w *wrappedErr) GrpcStatus() codes.Code {
	switch w.code {
	case Internal:
		return codes.Internal
	case NotFound:
		return codes.NotFound
	case BadRequest:
		return codes.InvalidArgument
	case Unauthorized:
		return codes.Unauthenticated
	case AlreadyExists:
		return codes.AlreadyExists
	}

	return codes.Unknown
}
