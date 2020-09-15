package errors

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

const (
	httpStatusUnknown             = 520
	httpStatusClientClosedRequest = 499

	reasonUnknown                    = "unknown"
	reasonInternal                   = "internal"
	reasonInvalidArgument            = "invalid_argument"
	reasonAlreadyExists              = "already_exists"
	reasonNotFound                   = "not_found"
	reasonConflict                   = "conflict"
	reasonLocked                     = "locked"
	reasonUnimplemented              = "unimplemented"
	reasonForbidden                  = "forbidden"
	reasonUnauthorised               = "unauthorised"
	reasonRateExceeded               = "rate_exceeded"
	reasonUnavailable                = "unavailable"
	reasonCanceled                   = "canceled"
	reasonDeadlineExceeded           = "deadline_exceeded"
	reasonUnavailableForLegalReasons = "unavailable_for_legal_reasons"
	errorsHeader                     = "X-GrabKit-Error-Reason"
)

// Unknown represents an error returned from the server that had no code or did not map to any known code.
//
// This should never be returned by Grab-Kit services.
func Unknown(msg string) *Response {
	return &Response{
		http:    httpStatusUnknown,
		grpc:    codes.Unknown,
		Reason:  reasonUnknown,
		Message: msg,
	}
}

// CustomHTTP returns an error with a custom HTTP code for backwards-compatibility purposes only. The Grab-Kit client
// treats all CustomHTTP errors as internal server errors.
func CustomHTTP(msg string, code int) *Response {
	return &Response{
		http:    code,
		grpc:    codes.Internal,
		Reason:  reasonInternal,
		Message: msg,
	}
}

// Internal ...
func Internal(msg string, err error) *Response {
	return &Response{
		http:    http.StatusInternalServerError,
		grpc:    codes.Internal,
		Reason:  reasonInternal,
		Message: msg,
		Err:     err,
	}
}

// InvalidArgument ...
func InvalidArgument(msg string) *Response {
	return &Response{
		http:    http.StatusBadRequest,
		grpc:    codes.InvalidArgument,
		Reason:  reasonInvalidArgument,
		Message: msg,
	}
}

// AlreadyExists ...
func AlreadyExists(msg string) *Response {
	return &Response{
		http:    http.StatusConflict,
		grpc:    codes.AlreadyExists,
		Reason:  reasonAlreadyExists,
		Message: msg,
	}
}

// NotFound ...
func NotFound(msg string) *Response {
	return &Response{
		http:    http.StatusNotFound,
		grpc:    codes.NotFound,
		Reason:  reasonNotFound,
		Message: msg,
	}
}

// Conflict ...
func Conflict(msg string) *Response {
	return &Response{
		http:    http.StatusConflict,
		grpc:    codes.AlreadyExists,
		Reason:  reasonConflict,
		Message: msg,
	}
}

// Locked ...
func Locked(msg string) *Response {
	return &Response{
		http:    http.StatusLocked,
		grpc:    codes.FailedPrecondition,
		Reason:  reasonLocked,
		Message: msg,
	}
}

// UnavailableForLegalReasons ...
func UnavailableForLegalReasons(msg string) *Response {
	return &Response{
		http:    http.StatusUnavailableForLegalReasons,
		grpc:    codes.PermissionDenied,
		Reason:  reasonUnavailableForLegalReasons,
		Message: msg,
	}
}

// Unimplemented ...
func Unimplemented(msg string) *Response {
	return &Response{
		http:    http.StatusNotImplemented,
		grpc:    codes.Unimplemented,
		Reason:  reasonUnimplemented,
		Message: msg,
	}
}

// Forbidden ...
func Forbidden(msg string) *Response {
	return &Response{
		http:    http.StatusForbidden,
		grpc:    codes.PermissionDenied,
		Reason:  reasonForbidden,
		Message: msg,
	}
}

// RateExceeded ...
func RateExceeded(msg string) *Response {
	return &Response{
		http:    http.StatusTooManyRequests,
		grpc:    codes.ResourceExhausted,
		Reason:  reasonRateExceeded,
		Message: msg,
	}
}

// Unauthorised ...
func Unauthorised(msg string) *Response {
	return &Response{
		http:    http.StatusUnauthorized,
		grpc:    codes.Unauthenticated,
		Reason:  reasonUnauthorised,
		Message: msg,
	}
}

// Unavailable ...
func Unavailable(msg string) *Response {
	return &Response{
		http:    http.StatusServiceUnavailable,
		grpc:    codes.Unavailable,
		Reason:  reasonUnavailable,
		Message: msg,
	}
}

// DeadlineExceeded ...
func DeadlineExceeded(msg string) *Response {
	return &Response{
		http:    http.StatusGatewayTimeout,
		grpc:    codes.DeadlineExceeded,
		Reason:  reasonDeadlineExceeded,
		Message: msg,
	}
}

// Canceled ...
func Canceled(msg string) *Response {
	return &Response{
		http:    httpStatusClientClosedRequest, // gRPC status code CANCELED, Closest HTTP Mapping - 499 Client Closed Request
		grpc:    codes.Canceled,
		Reason:  reasonCanceled,
		Message: msg,
	}
}

func grpcCodeToHTTPCode(code codes.Code) int {
	switch code {
	case codes.Canceled:
		return httpStatusClientClosedRequest
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.PermissionDenied:
		return http.StatusForbidden // also http.StatusUnavailableForLegalReasons, but ignoring that.
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.NotFound:
		return http.StatusNotFound
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.Unknown:
		fallthrough
	default:
		return httpStatusUnknown
	}
}
