package kubernetes

import (
	"fmt"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"

	api "k8s.io/client-go/pkg/api/v1"
)

// AutoPath implements the AutoPathFunc call from the autopath plugin.
// It returns a per-query search path or nil indicating no searchpathing should happen.
func (k *Kubernetes) AutoPath(state request.Request, namespace string) ([]string, error) {
	// Check if the query falls in a zone we are actually authoriative for and thus if we want autopath.
	zone := plugin.Zones(k.Zones).Matches(state.Name())
	if zone == "" {
		return nil, fmt.Errorf("not authoriative: %q", state.Name())
	}

	ip := state.IP()

	pod := k.podWithIP(ip)
	if pod == nil {
		return nil, fmt.Errorf("no pod found with IP: %q", ip)
	}

	if namespace != "" && namespace != pod.Namespace {
		return nil, fmt.Errorf("pod namespace %q not allowed to do autopath for %q", pod.Namespace, namespace)
	}

	search := make([]string, 3)
	if zone == "." {
		search[0] = pod.Namespace + ".svc."
		search[1] = "svc."
		search[2] = "."
	} else {
		search[0] = pod.Namespace + ".svc." + zone
		search[1] = "svc." + zone
		search[2] = zone
	}

	search = append(search, k.autoPathSearch...)
	search = append(search, "") // sentinal
	return search, nil
}

// podWithIP return the api.Pod for source IP ip. It returns nil if nothing can be found.
func (k *Kubernetes) podWithIP(ip string) *api.Pod {
	ps := k.APIConn.PodIndex(ip)
	if len(ps) == 0 {
		return nil
	}
	return ps[0]
}
