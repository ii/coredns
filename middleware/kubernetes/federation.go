package kubernetes

import (
	"strings"

	"github.com/coredns/coredns/middleware/etcd/msg"
	"github.com/coredns/coredns/request"
)

const (
	// TODO: Do not hardcode these labels. Pull them out of the API instead.
	//
	// We can get them via ....
	//   import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//     metav1.LabelZoneFailureDomain
	//     metav1.LabelZoneRegion
	//
	// But importing above breaks coredns with flag collision of 'log_dir'

	labelZone   = "failure-domain.beta.kubernetes.io/zone"
	labelRegion = "failure-domain.beta.kubernetes.io/region"
)

type FederationFunc func(state request.Request) (msg.Service, error)

// Federations is used from the federations middleware to return the service that should be
// returned as a CNAME for federation to work.
func (k *Kubernetes) Federations(state request.Request) (msg.Service, error) {
	nodeName := k.localNodeName()
	node, err := k.APIConn.GetNodeByName(nodeName)
	if err != nil {
		return msg.Service{}, err
	}
	r, err := k.parseRequest(state)

	if r.endpoint == "" {
		s := msg.Service{
			Key:  strings.Join([]string{msg.Path(r.zone, "coredns"), r.podOrSvc, state.Zone, r.namespace, r.service}, "/"),
			Host: strings.Join([]string{r.service, r.namespace, state.Zone, r.podOrSvc, node.Labels[labelZone], node.Labels[labelRegion], state.Zone}, "."),
		}
		return s, nil
	}
	s := msg.Service{
		Key:  strings.Join([]string{msg.Path(r.zone, "coredns"), r.podOrSvc, state.Zone, r.namespace, r.service, r.endpoint}, "/"),
		Host: strings.Join([]string{r.endpoint, r.service, r.namespace, state.Zone, r.podOrSvc, node.Labels[labelZone], node.Labels[labelRegion], state.Zone}, "."),
	}
	return s, nil
}
