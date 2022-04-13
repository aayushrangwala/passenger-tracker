package errors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	InvalidArgumentError = status.Error(codes.InvalidArgument, "Invalid argument passed")
	UnImplementedError   = status.Error(codes.Unimplemented, "method not implemented")
	InternalServerError  = status.Error(codes.Internal, "failed to serve the endpoint")
)

// IsInvalidArgument returns true if the error is NotFound error.
func IsInvalidArgument(err error) bool {
	return codes.InvalidArgument == status.Code(err)
}

// IsInternalServer returns true if the error is NotFound error.
func IsInternalServer(err error) bool {
	return codes.Internal == status.Code(err)
}

// IsUnImplementedError returns true if the error passed is UnImplemented error.
func IsUnImplementedError(err error) bool {
	return codes.Unimplemented == status.Code(err)
}
