package constant

type (
	// IdempotencyLevel determines the idempotency for this endpoint
	IdempotencyLevel string
)

const (
	// IdempotencyUnknown means the idempotency is not specified; assume not idempotent
	IdempotencyUnknown IdempotencyLevel = "unknown"

	// IdempotencyIdempotent means the endpoint is idempotent, but may have side-effects
	IdempotencyIdempotent IdempotencyLevel = "idempotent"

	// IdempotencySafe means the endpoint is idempotent and is guaranteed not to have side-effects
	IdempotencySafe IdempotencyLevel = "safe"

	// IdempotencyDefault is an alias for the default idempotency level
	IdempotencyDefault = IdempotencyUnknown

)
