/*
Package external implements external names for kubernetes clusters. The following
name are allowed:

* name1
* name2

queries for the zone's apex will ...

Basically any plugin can implement this if implement the ExternalFunc.

I.e:

func (m Plugin) External(state request.Request) ([]msg.Service, error) {
	return nil, nil
}
*/
package external

import (
	"context"
	"fmt"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/etcd/msg"
	"github.com/coredns/coredns/request"

	"github.com/miekg/dns"
)

// Func defines the function a plugin should implement to return set of services.
type Func func(request.Request) ([]msg.Service, error)

// Externaler defines the interface that a plugin should implement in order to be
// used by External.
type Externaler interface {
	External(request.Request) ([]msg.Service, error)
}

// External resolves Ingress and Loadbalance IPs from kubernetes clusters
type External struct {
	Next  plugin.Handler
	Zones []string

	externalFunc Func
}

// ServeDNS implements the plugin.Handle interface.
func (e *External) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}

	zone := plugin.Zones(e.Zones).Matches(state.Name())
	if zone == "" {
		return plugin.NextOrFailure(e.Name(), e.Next, ctx, w, r)
	}

	if e.externalFunc == nil {
		return plugin.NextOrFailure(e.Name(), e.Next, ctx, w, r)
	}

	svc, err := e.externalFunc(state)
	if err != nil {
		return dns.RcodeServerFailure, nil
	}
	// TODO, actually make some DNS data out of this. Do we need to full backend_lookup.go machinery? Probably not, because
	// these things are never CNAMEs that need resolver. Hard to say without actually documentation regarding this feature.
	// We can call the various NewSRV and friends on these services, and we'll probably need to handle/copy some of that code
	// over here.
	// Apex query, zone transfers and other bits can all be implemented later.
	fmt.Printf("%v\n", svc)

	return dns.RcodeServerFailure, nil
}

// Name implements the Handler interface.
func (e *External) Name() string { return "external" }
