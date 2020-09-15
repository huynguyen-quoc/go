package metadata

// Getter is an interface that gets endpoint metadata given a key
type Getter interface {
	Get(key string, protocol Protocol) Metadata
}

// NoOpGetter can be used for default / tests
type NoOpGetter struct{}

// Get ...
func (n *NoOpGetter) Get(key string, protocol Protocol) Metadata { return Metadata{} }

// Protocol for which md is requested
type Protocol string

const (
	// ProtocolHTTP ...
	ProtocolHTTP Protocol = "http"

	// ProtocolGRPC ...
	ProtocolGRPC Protocol = "grpc"
)
