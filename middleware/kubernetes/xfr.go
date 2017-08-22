package kubernetes

import (
	"fmt"
	"sync"
	"time"

	"github.com/coredns/coredns/middleware/pkg/dnsutil"

	"github.com/miekg/dns"
	"k8s.io/client-go/1.5/pkg/api"
)

type Xfr struct {
	*Kubernetes
	sync.RWMutex
	epoch time.Time
}

func NewXfr(k *Kubernetes) *Xfr {
	return &Xfr{Kubernetes: k, epoch: time.Now().UTC()}
}

// All returns all kubernetes records with a SOA at the start.
func (x *Xfr) All(zone string) []dns.RR {
	res := []dns.RR{}

	serviceList := x.APIConn.ServiceList()
	for _, svc := range serviceList {

		name := dnsutil.Join([]string{svc.Name, svc.Namespace, zone})

		// Endpoint query or headless service
		if svc.Spec.ClusterIP == api.ClusterIPNone {

			endpointsList := x.APIConn.EndpointsList()
			for _, ep := range endpointsList.Items {
				if ep.ObjectMeta.Name != svc.Name || ep.ObjectMeta.Namespace != svc.Namespace {
					continue
				}
				for _, eps := range ep.Subsets {
					for _, addr := range eps.Addresses {
						for _, p := range eps.Ports {

							fmt.Printf("%s IN A %s\n", name, addr.IP)
							fmt.Printf("_%s._%s.%s IN SRV %d %s.%s\n", p.Name, p.Protocol, name, p.Port, endpointHostname(addr), name)
							fmt.Printf("%s.%s IN A %s\n", endpointHostname(addr), name, addr.IP)
						}
					}
				}
			}
			continue
		}

		// External service
		if svc.Spec.ExternalName != "" {
			fmt.Printf("%s IN CNAME %s", name, svc.Spec.ExternalName)
			continue
		}

		// ClusterIP service
		fmt.Printf("%s IN A %s\n", name, svc.Spec.ClusterIP)
		for _, p := range svc.Spec.Ports {
			fmt.Printf("_%s._%s.%s IN SRV %s\n", p.Name, p.Protocol, name, name)
		}
	}
	return res
}

func (x *Xfr) serial() uint32 {
	x.RLock()
	defer x.RUnlock()
	return uint32(x.epoch.Unix())
}

// Give these to dnscontroller via the options, so these functions get exectuted and the SOA's serial gets updated.
func (x *Xfr) AddDeleteXfrHandler(a interface{}) {
	x.Lock()
	defer x.Unlock()
	x.epoch = time.Now().UTC()
}

func (x *Xfr) UpdateXfrHandler(a, b interface{}) {
	x.Lock()
	defer x.Unlock()
	x.epoch = time.Now().UTC()
}

/*
cache.ResourceEventHandlerFuncs{
    AddFunc: x.AddDeleteXfrHandler,
    DeleteFunc: x.AddDeleteXfrHandler,
    UpdateFunc: x.UpdateXfrHandler,
}

// set to nil? Or noop functions as defaults?
*/
