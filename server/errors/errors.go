package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	descRegex = regexp.MustCompile(`target=(.*), reason=(.*), msg=(.*)`)
)

const (
	errTypeUnknown = "unknown"
	errTypeServer  = "server"
	errTypeClient  = "client"
)

// Response ...
type Response struct {
	// internal to server
	http    int
	grpc    codes.Code
	payload interface{}
	Err     error `json:"-"`
	// will be transported to client
	Target  string `json:"target"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

// JSONResponse is used only for marshalling to json
type JSONResponse struct {
	Target  string `json:"target"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

// MarshalJSON allows for custom rendering of an error payload
func (r *Response) MarshalJSON() ([]byte, error) {
	if r.payload != nil {
		switch payload := r.payload.(type) {
		case []byte:
			jResp := &JSONResponse{
				Target:  r.Target,
				Reason:  r.Reason,
				Message: string(payload),
			}
			return json.Marshal(jResp)
		default:
			return json.Marshal(payload)
		}
	}
	jResp := &JSONResponse{
		Target:  r.Target,
		Reason:  r.Reason,
		Message: r.Message,
	}
	return json.Marshal(jResp)
}

// Wrap stores a custom error payload. This will be used in MarshalJSON(), meaning it will be output in place of the
// standard error response when returned from the API.
func (r *Response) Wrap(payload interface{}) *Response {
	r.payload = payload
	return r
}

// SetHTTPCode overrides the HTTP status code
//
// This is intended for backwards-compatibility with server responses *only* and will be completely ignored by Grab-Kit
// clients and other functions like FromHTTPResponse.
func (r *Response) SetHTTPCode(code int) *Response {
	r.http = code
	return r
}

// ToHTTPCode ...
func (r *Response) ToHTTPCode() int {
	if r.http != 0 {
		return r.http
	}
	return httpStatusUnknown
}

// ToGRPCCode ...
func (r *Response) ToGRPCCode() codes.Code {
	// codes.OK is the default value, but shall never used for errors
	if r.grpc != codes.OK {
		return r.grpc
	}
	return codes.Unknown
}

// Payload returns the raw error response payload, if it was not a standard grab-kit error
func (r *Response) Payload() interface{} {
	return r.payload
}

// Error ...
func (r *Response) Error() string {
	return fmt.Sprintf("ServerError: target=%v, reason=%v, msg=%v", r.Target, r.Reason, r.Message)
}

// ToRPCError used for encoding server errors
func (r *Response) ToRPCError() error {
	return status.Errorf(r.grpc, r.rpcMessage())
}

// GRPCStatus returns the GRPC status
func (r *Response) GRPCStatus() *status.Status {
	return status.New(r.grpc, r.rpcMessage())
}

func (r *Response) rpcMessage() string {
	return fmt.Sprintf("target=%v, reason=%v, msg=%v", r.Target, r.Reason, r.Message)
}

// FromError attempts to convert a generic error back to a Grab-Kit error
func FromError(err error) error {
	// alias to FromRPCError for now
	return FromRPCError(err)
}

// grpcStatusFromErr attempts to convert an error to a GRPC error
func grpcStatusFromErr(err error) (s *status.Status, ok bool) {
	if se, ok := err.(interface {
		GRPCStatus() *status.Status
	}); ok {
		return se.GRPCStatus(), true
	}
	return status.New(codes.Unknown, err.Error()), false
}

// FromRPCError converts a gRPC error response back to a Grab-Kit error
//
// gRPC errors should *always* be in the Grab-Kit format unless there was a transport error, in which case we use the
// gRPC code directly. This assumes the gRPC code is reliable and not easily changed by the service.
//
// The default error if the reason code is present but doesn't match anything is Unknown
// If the reason code is not present, a default error is created with the original gRPC code and a HTTP code of Unknown.
func FromRPCError(err error) error {
	s, ok := status.FromError(err)
	if !ok {
		// TODO: remove once grpc-go is updated with the GRPCStatus() check
		s, _ = grpcStatusFromErr(err)
	}
	desc := s.Message()
	matched := descRegex.FindStringSubmatch(desc)
	// matched will return an array of [matched string, group 1, group 2]
	if len(matched) > 3 {
		errResp := reasonToErr(matched[2], matched[3])
		errResp.Target = matched[1]
		errResp.grpc = s.Code() // trust original grpc code for accuracy
		return errResp
	}
	httpStatus := httpStatusUnknown
	if hs, ok := err.(interface {
		ToHTTPCode() int
	}); ok {
		httpStatus = hs.ToHTTPCode()
	} else {
		httpStatus = grpcCodeToHTTPCode(s.Code()) // Try to convert to the relevant error if we can.
	}
	return &Response{
		grpc:    s.Code(),
		http:    httpStatus,
		Reason:  strings.ToLower(s.Code().String()),
		Message: desc,
	}
}

// FromHTTPResponse converts a HTTP error back to a Grab-Kit error
//
// HTTP responses are not guaranteed to be in the Grab-Kit format, since the response body can be overridden. In these
// cases, we rely on a GrabKit header to map to an appropriate Grab-Kit error.
//
// The default is always an Unknown error if the error can't be parsed or the status is unknown. We use this instead of
// the actual HTTP status code because it is assumed to be safer, since the HTTP status code could be overridden by the
// service.
func FromHTTPResponse(resp *http.Response) error {
	// Success response: no error
	if resp == nil || resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return nil
	}

	// No body: use status code
	if resp.Body == nil {
		return fromHTTPHeaders(resp, nil)
	}

	// Try parsing response
	bytesResponse, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return httpStatusCodeToErr(resp.StatusCode, fmt.Sprintf("(grab-kit: could not read response body: %v)", err))
	}

	// Try decoding response
	defer func() { _ = resp.Body.Close() }()

	errResp := &Response{}
	dec := json.NewDecoder(bytes.NewReader(bytesResponse))

	if err := dec.Decode(errResp); err != nil || errResp.Reason == "" {
		return fromHTTPHeaders(resp, bytesResponse)
	}

	// Got reason in response: convert to actual error
	reasonErr := reasonToErr(errResp.Reason, errResp.Message)
	reasonErr.Target = errResp.Target
	reasonErr.payload = bytesResponse
	return reasonErr
}

func fromHTTPHeaders(resp *http.Response, bytesResponse []byte) error {
	var httpErr *Response
	// Unable to decode or reason was empty. Check error header first.
	if errHeader := resp.Header.Get(errorsHeader); errHeader != "" {
		httpErr = reasonToErr(errHeader, "")
	} else {
		// Couldn't read header. Use status code instead
		httpErr = httpStatusCodeToErr(resp.StatusCode, "")
	}
	httpErr.payload = bytesResponse
	return httpErr
}

// reasonToErr converts a reason code to an actual error Response using the constructors, which guarantees all internal
// fields are properly set.
func reasonToErr(reason, msg string) *Response {
	switch reason {
	case reasonUnknown:
		return Unknown(msg)
	case reasonInternal:
		return Internal(msg, nil) // Err is unused
	case reasonInvalidArgument:
		return InvalidArgument(msg)
	case reasonAlreadyExists:
		return AlreadyExists(msg)
	case reasonNotFound:
		return NotFound(msg)
	case reasonConflict:
		return Conflict(msg)
	case reasonUnimplemented:
		return Unimplemented(msg)
	case reasonForbidden:
		return Forbidden(msg)
	case reasonRateExceeded:
		return RateExceeded(msg)
	case reasonUnauthorised:
		return Unauthorised(msg)
	case reasonUnavailable:
		return Unavailable(msg)
	case reasonUnavailableForLegalReasons:
		return UnavailableForLegalReasons(msg)
	case reasonCanceled:
		return Canceled(msg)
	case reasonDeadlineExceeded:
		return DeadlineExceeded(msg)
	}

	// Default: unknown error
	return Unknown(msg)
}

func httpStatusCodeToErr(statusCode int, msg string) *Response {
	switch statusCode {
	case httpStatusUnknown:
		return Unknown(msg)
	case http.StatusInternalServerError:
		return Internal(msg, nil)
	case http.StatusBadRequest:
		return InvalidArgument(msg)
	case http.StatusConflict:
		return AlreadyExists(msg)
	case http.StatusNotFound:
		return NotFound(msg)
	case http.StatusNotImplemented:
		return Unimplemented(msg)
	case http.StatusForbidden:
		return Forbidden(msg)
	case http.StatusTooManyRequests:
		return RateExceeded(msg)
	case http.StatusUnauthorized:
		return Unauthorised(msg)
	case http.StatusServiceUnavailable:
		return Unavailable(msg)
	case http.StatusUnavailableForLegalReasons:
		return UnavailableForLegalReasons(msg)
	case httpStatusClientClosedRequest:
		return Canceled(msg)
	case http.StatusGatewayTimeout:
		return DeadlineExceeded(msg)
	}
	return Unknown(msg)
}

// GetErrorType identifies error type as server/client error based on http spec 4xx - client error & 5xx - server error
// http code used is based on the closest mapping of gRPC error codes to http status codes
func (r *Response) GetErrorType() string {
	if r.ToGRPCCode() == codes.Unknown || r.ToHTTPCode() == httpStatusUnknown {
		return errTypeUnknown
	}
	httpStatus := r.ToHTTPCode()
	if httpStatus >= 400 && httpStatus <= 499 {
		return errTypeClient
	}
	if httpStatus >= 500 && httpStatus <= 599 {
		return errTypeServer
	}
	return errTypeUnknown
}

// Unwrap implements xerrors.Wrapper to support Go 1.13 error values
func (r *Response) Unwrap() error {
	return r.Err
}
