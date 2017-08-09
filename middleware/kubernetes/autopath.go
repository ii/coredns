package kubernetes

import (
	"fmt"

	"github.com/coredns/coredns/middleware"
	"github.com/coredns/coredns/request"

	"k8s.io/client-go/1.5/pkg/api"
)

func (k *Kubernetes) AutoPath(state request.Request) ([]string, error) {
	// Check if the query falls in a zone we are actually authoriative for and thus if we want autopath.
	zone := middleware.Zones(k.Zones).Matches(state.Name())
	if zone == "" {
		return nil, fmt.Errorf("kubernetes: no authoriative for %s", state.Name())
	}

	ip := state.IP()

	pod := k.PodWithIP(ip)
	if pod == nil {
		return nil, fmt.Errorf("kubernetes: no pod found for %s", ip)
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

// PodWithIP return the api.Pod for source IP ip. It return nil if nothing can be found.
func (k *Kubernetes) PodWithIP(ip string) (p *api.Pod) {
	objList := k.APIConn.PodIndex(ip)
	for _, o := range objList {
		p, ok := o.(*api.Pod)
		if !ok {
			return nil
		}
		return p
	}
	return nil
}
