package plugin

import "context"

// Done returns true if the context has been canceled, usually because the timeout has expired.
func Done(ctx context.Context) bool {
	for {
		select {
		case <-ctx.Done():
			return true
		default:
			return false
		}
	}
	return false
}
