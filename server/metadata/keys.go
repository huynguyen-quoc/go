package metadata

import (
	"fmt"
)

const (
	// Be careful to follow http2 convention here, and note that uppercase chars will fail to parse, leading to
	// fatal errors on the client
	// See: https://gitlab.myteksi.net/gophers/go/blob/2420f08092f8debaa5a240b1cd841b62eb618de4/vendor/golang.org/x/net/http2/http2.go#L189-L189

	// KeyGRPCCacheTTL in seconds
	KeyGRPCCacheTTL = "cache-ttl"
)

var (
	// These are to standardize metadata/header keys between HTTP and gRPC

	// KeyRequestAuthorization is the Authorization header of the incoming request
	KeyRequestAuthorization = makeKey("auth")

	// KeyCacheTTL is the cache ttl
	KeyCacheTTL = makeKey("cache-ttl")

	// KeyCacheControl is the cache control header (equivalent of HTTP)
	KeyCacheControl = makeKey("cache-control")

	// KeyGrabServiceID is a user authentication field
	KeyGrabServiceID = makeKey("grab-serviceid")

	// KeyGrabServiceUserID is a user authentication field
	KeyGrabServiceUserID = makeKey("grab-service-userid")

	// KeyClientID is the id of the calling service
	KeyClientID = makeKey("clientid")

	// KeyClientIP is the ip address of the calling service
	KeyClientIP = makeKey("clientip")

	// KeyUserAgent is a HTTP request header string that allows to identify the application type,
	// operating system, software vendor or software version of the requesting software user agent.
	// Example - Grab/5.14.0 (iPhone; iOS 12.0; Scale/2.0)
	KeyUserAgent = "user-agent"

	// KeyAcceptLanguage is a HTTP request header. It advertises which language the client is able to understand
	KeyAcceptLanguage = "accept-language"

	// KeyForwardedUserAgent is populated by the HTTP to GRPC proxy (API Gateway) from the original HTTP request.
	// This is necessary as GRPC will automatically override the standard HTTP user-agent header with its own value
	KeyForwardedUserAgent = "x-forwarded-user-agent"
)

func makeKey(key string) string {
	return fmt.Sprintf("x-grabkit-%s", key)
}
