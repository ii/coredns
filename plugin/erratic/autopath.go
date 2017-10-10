package erratic

import "github.com/coredns/coredns/request"

// AutoPath implements the AutoPathFunc call from the autopath plugin.
func (e *Erratic) AutoPath(state request.Request, namespace string) ([]string, error) {
	return []string{"a.example.org.", "b.example.org.", ""}, nil
}
