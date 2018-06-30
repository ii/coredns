package metadata

import (
	"context"

	"github.com/coredns/coredns/request"
)

// Provider interface needs to be implemented by each plugin willing to provide
// metadata information for other plugins.
type Provider interface {
	// Metadata adds metadata to the context and returns a (potentially) new context.
	// Note: this method should work quickly, because it is called for every request
	// from the metadata plugin.
	Metadata(ctx context.Context, state request.Request) context.Context
}

// Labels returns all metadata keys stored in the context. These label names should be named
// as: plugin/NAME, where NAME is something descriptive.
func Labels(ctx context.Context) []string {
	if metadata := ctx.Value(key{}); metadata != nil {
		if m, ok := metadata.(md); ok {
			return keys(m)
		}
	}
	return nil
}

// Value returns the value of label. If none can be found the empty string is returned.
func Value(ctx context.Context, label string) string {
	if metadata := ctx.Value(key{}); metadata != nil {
		if m, ok := metadata.(md); ok {
			return m[label]
		}
	}
	return ""
}

// SetValue set the metadata label to value. If no metadata can be found this is a noop and
// false is returned. Any existing value is overwritten.
func SetValue(ctx context.Context, label, value string) bool {
	if metadata := ctx.Value(key{}); metadata != nil {
		if m, ok := metadata.(md); ok {
			m[label] = value
			return true
		}
	}
	return false
}

// md is metadata information storage.
type md map[string]string

// key defines the type of key that is used to save metadata into the context.
type key struct{}

func keys(m map[string]string) []string {
	s := make([]string, len(m))
	i := 0
	for k, _ := range m {
		s[i] = k
		i++
	}
	return s
}
