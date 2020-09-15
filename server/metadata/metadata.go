package metadata

import (
	"net/textproto"
	"strings"
)

const (
	maxHeaderSize = 1000
)

// Metadata ...
type Metadata map[string][]string

// Len returns the number of items in meta.
func (meta Metadata) Len() int {
	return len(meta)
}

// Get attempts to do a case-insensitive lookup of key. It returns the first value found.
// If you need a case-sensitive lookup or multiple values, access the map directly.
func (meta Metadata) Get(key string) string {
	// Short-circuit for common header casing first
	for _, k := range []string{key, textproto.CanonicalMIMEHeaderKey(key), strings.ToLower(key)} {
		if v, ok := meta[k]; ok && len(v) >= 1 {
			return v[0]
		}
	}

	if meta.Len() > maxHeaderSize {
		return ""
	}

	// Iterate over the array. We cannot follow the net/http approach of using Set/Get for backwards-compatibility, and net/http headers are
	// not compatible with gRPC (http2)
	for k, v := range meta {
		if strings.EqualFold(k, key) && len(v) >= 1 {
			return v[0]
		}
	}
	return ""
}

// Copy returns a copy of meta.
func (meta Metadata) Copy() Metadata {
	return Join(meta)
}

// Join joins any number of meta into a single Metadata.
func Join(metas ...Metadata) Metadata {
	out := Metadata{}
	for _, meta := range metas {
		for k, v := range meta {
			out[k] = append(out[k], v...)
		}
	}
	return out
}
