/*
The federation package implements kubernetes federation. It checks if the qname matches
a possible federation. If this is the case and the captured answer is an NXDOMAIN,
federation is performed. If this is not the case the next middleware in the chain
is called.

Federation is only useful in conjunction with the kubernetes middleware, without it is a noop.
*/
package federation

import (
	"github.com/coredns/coredns/middleware"
	"github.com/coredns/coredns/middleware/etcd/msg"
	"github.com/coredns/coredns/request"

	"github.com/miekg/dns"
	"golang.org/x/net/context"
)

// Federation contains the name to zone mapping used for federation
// in kubernetes.
type Federation struct {
	f     map[string]string
	zones []string

	Next        middleware.Handler
	Fallthrough bool
	Federations FederationFunc
}

// FederationFunc needs to be implemented by any middleware that implements
// federation. Right now this is only the kubernetes middleware.
type FederationFunc func(state request.Request) ([]msg.Service, error)

func New() *Federation {
	return &Federation{f: make(map[string]string)}
}

func (f *Federation) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	// Kubernetes is not loaded, we can't get to the data.
	if f.Federations == nil {
		if f.Fallthrough {
			return middleware.NextOrFailure(f.Name(), f.Next, ctx, w, r)
		}
		return dns.RcodeServerFailure, nil
	}

	state := request.Request{W: w, Req: r}
	zone := middleware.Zones(f.zones).Matches(state.Name())
	if zone == "" {
		if f.Fallthrough {
			return middleware.NextOrFailure(f.Name(), f.Next, ctx, w, r)
		}
		return dns.RcodeServerFailure, nil
	}

	state.Zone = zone
	services, err := f.Federations(state)

	if err != nil {
		if f.Fallthrough {
			return middleware.NextOrFailure(f.Name(), f.Next, ctx, w, r)
		}
		return dns.RcodeServerFailure, nil
	}
	if err != nil {
		return dns.RcodeServerFailure, err
	}

	var records []dns.RR
	for _, serv := range services {
		records = append(records, serv.NewCNAME(state.QName(), serv.Host))
	}

	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative, m.RecursionAvailable, m.Compress = true, true, true
	m.Answer = append(m.Answer, records...)

	state.SizeAndDo(m)
	m, _ = state.Scrub(m)
	w.WriteMsg(m)

	return dns.RcodeSuccess, nil
}

func (f *Federation) Name() string { return "federation" }
