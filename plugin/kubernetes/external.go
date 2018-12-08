package kubernetes

import (
	"github.com/coredns/coredns/plugin/etcd/msg"
	"github.com/coredns/coredns/plugin/kubernetes/object"
	"github.com/coredns/coredns/plugin/pkg/dnsutil"
	"github.com/coredns/coredns/request"

	"github.com/miekg/dns"
)

// External implements the ExternalFunc call from the external plugin.
// It returns a per-query search path or nil indicating no searchpathing should happen.
func (k *Kubernetes) External(state request.Request) ([]msg.Service, error) {
	// get a name like a1.<namespace>.<zone>
	base, _ := dnsutil.TrimZone(state.Name(), state.Zone)
	if base == "" {
		return nil, nil
	}
	segs := dns.SplitDomainName(base)
	ls := len(segs)
	if ls < 3 {
		return nil, nil
	}
	namespace := segs[ls-2]
	service := segs[ls-3]
	//port/name
	// What do with longer name, nx domain?

	// don't support wildcards for this, except for the port and protocol, because we allow those *not* to be specified.

	idx := object.ServiceKey(service, namespace)
	serviceList := k.APIConn.SvcIndex(idx)

	services := []msg.Service{}
	for _, svc := range serviceList {
		if namespace != svc.Namespace {
			continue
		}
		if service != svc.Name {
			continue
		}

		/*
			for _, ip := range svc.ExternalIPs {
				for _, p := range svc.Ports {
					if !(match(r.port, p.Name) && match(r.protocol, string(p.Protocol))) {
						continue
					}
					err = nil
					s := msg.Service{Host: ip, Port: int(p.Port), TTL: k.ttl}
					s.Key = strings.Join([]string{zonePath, svc.Namespace, svc.Name}, "/")
					services = append(services, s)
				}
			}
		*/
	}
	// TODO set err and errNoItems
	return services, nil
}
