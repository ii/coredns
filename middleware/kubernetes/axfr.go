package kubernetes

import (
	"fmt"

	"github.com/coredns/coredns/middleware/pkg/dnsutil"

	"k8s.io/client-go/1.5/pkg/api"
)

type Xfr struct {
	*Kubernetes
}

func (x *Xfr) services(zone string) ([]kService, error) {
	res := []kService{}

	serviceList := x.APIConn.ServiceList()
	for _, svc := range serviceList {
		s := kService{name: svc.Name, namespace: svc.Namespace}
		fmt.Printf("name %v %v\n", svc.Name, svc.Namespace)
		suffix := dnsutil.Join([]string{svc.Name, svc.Namespace, zone})

		// Endpoint query or headless service
		if svc.Spec.ClusterIP == api.ClusterIPNone {
			s.addr = svc.Spec.ClusterIP

			endpointsList := x.APIConn.EndpointsList()
			for _, ep := range endpointsList.Items {
				if ep.ObjectMeta.Name != svc.Name || ep.ObjectMeta.Namespace != svc.Namespace {
					continue
				}
				for _, eps := range ep.Subsets {
					for _, addr := range eps.Addresses {
						for _, p := range eps.Ports {
							s.endpoints = append(s.endpoints, endpoint{addr: addr, port: p})
							fmt.Printf("%s IN A %s\n", suffix, addr.IP)
							fmt.Printf("_%s_%s.%s IN SRV %s\n", p.Name, p.Protocol, suffix, addr.IP)
							fmt.Printf("%s.%s IN A %s\n", endpointHostname(addr), suffix, addr.IP)
						}
					}
				}
			}
			if len(s.endpoints) > 0 {
				res = append(res, s)
			}
			continue
		}
		println("clusterip", svc.Spec.ClusterIP)

		// External service
		if svc.Spec.ExternalName != "" {
			s.addr = svc.Spec.ExternalName
			println("external", s.addr)
			res = append(res, s)
			continue
		}

		// ClusterIP service
		s.addr = svc.Spec.ClusterIP
		for _, p := range svc.Spec.Ports {
			fmt.Printf("SRV record _%s._%s %s\n", p.Name, p.Protocol, suffix)
			s.ports = append(s.ports, p)
			// srv target is suffix which gets cluster IP
			fmt.Printf("A record %s %s\n", suffix, svc.Spec.ClusterIP)
		}

		res = append(res, s)
	}
	return res, nil
}
