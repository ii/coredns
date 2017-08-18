package kubernetes

import (
	"fmt"

	"github.com/coredns/coredns/middleware/pkg/dnsutil"
	"github.com/coredns/coredns/request"

	"github.com/miekg/dns"
)

type recordRequest struct {
	// The named port from the kubernetes DNS spec, this is the service part (think _https) from a well formed
	// SRV record.
	port string
	// The protocol is usually _udp or _tcp (if set), and comes from the protocol part of a well formed
	// SRV record.
	protocol string
	endpoint string
	// The servicename used in Kubernetes.
	service string
	// The namespace used in Kubernetes.
	namespace string
	// A each name can be for a pod or a service, here we track what we've seen, either "pod" or "service".
	podOrSvc string
}

// parseRequest parses the qname to find all the elements we need for querying k8s. Anything
// that is not parsed will have the wildcard "*" value. Potential underscores are stripped
// from _port and _protocol.
func (k *Kubernetes) parseRequest(state request.Request) (r recordRequest, err error) {
	// 3 Possible cases:
	// o SRV Request: _port._protocol.service.namespace.pod|svc.zone
	// o A Request (endpoint): endpoint.service.namespace.pod|svc.zone
	// o A Request (service): service.namespace.pod|svc.zone
	//
	// Federations are handled in the federation middleware.

	base, _ := dnsutil.TrimZone(state.Name(), state.Zone)
	segs := dns.SplitDomainName(base)

	r.port = "*"
	r.service = "*"
	r.endpoint = "*"
	r.namespace = "*"

	defer func() {
		fmt.Printf("%#v\n", r)
	}()

	// start at the right and fill out recordRequest with the bits we find, so we look for
	// pod|svc.namespace.service and then either
	// * endpoint
	// *_protocol._port

	last := len(segs) - 1
	r.podOrSvc = segs[last]
	if r.podOrSvc != Pod && r.podOrSvc != Svc {
		return r, errInvalidRequest
	}
	last--
	if last < 0 {
		return r, nil
	}

	r.namespace = segs[last]
	last--
	if last < 0 {
		return r, nil
	}

	r.service = segs[last]
	last--
	if last < 0 {
		return r, nil
	}

	if segs[last][0] == '_' {
		r.protocol = segs[last][1:]
	} else {
		r.endpoint = segs[last]
	}
	last--
	if last < 0 {
		return r, nil
	}

	if segs[last][0] == '_' {
		r.port = segs[last][1:]
	}

	if last > 0 { // Too long, so NXDOMAIN these.
		return r, errInvalidRequest

	}
	return r, nil
}

// String return a string representation of r, it just returns all fields concatenated with dots.
// This is mostly used in tests.
func (r recordRequest) String() string {
	s := r.port
	s += "." + r.protocol
	s += "." + r.endpoint
	s += "." + r.service
	s += "." + r.namespace
	s += "." + r.podOrSvc
	return s
}
