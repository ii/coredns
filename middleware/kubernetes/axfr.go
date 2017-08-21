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

func (x *Xfr) services(zone string) []dns.RR {
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

/*
From https://stackoverflow.com/questions/35192712/kubernetes-watch-pod-events-with-api

similar to what is in controller.go ? May extend that when transfers are enabled?
grab xfr lock and epoch = time.Now().UTC() on every event watch

    watchlist := cache.NewListWatchFromClient(clientset.Core().RESTClient(), "pods", v1.NamespaceDefault,
       fields.Everything())
    _, controller := cache.NewInformer(
        watchlist,
        &v1.Pod{},
        time.Second * 0,
        cache.ResourceEventHandlerFuncs{
            AddFunc: func(obj interface{}) {
                fmt.Printf("add: %s \n", obj)
            },
            DeleteFunc: func(obj interface{}) {
                fmt.Printf("delete: %s \n", obj)
            },
            UpdateFunc:func(oldObj, newObj interface{}) {
                fmt.Printf("old: %s, new: %s \n", oldObj, newObj)
            },
        },
    )
    stop := make(chan struct{})
    go controller.Run(stop)
}
*/
