// Package forward implements a forwarding proxy. It caches an upstream net.Conn for some time, so if the same
// client returns the upstream's Conn will be precached. Depending on how you benchmark this looks to be
// 50% faster than just openening a new connection for every client. It works with UDP and TCP and uses
// inband healthchecking.
package forward

import (
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"

	"github.com/miekg/dns"
	"golang.org/x/net/context"
)

// Forward represents a plugin instance that can proxy requests to another (DNS) server.
type Forward struct {
	proxies []*Proxy

	from string

	Next plugin.Handler
}

func (f Forward) Name() string { return "forward" }

func (f Forward) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {

	f.proxies[0].clientChan <- request.Request{W: w, Req: r}

	return 0, nil
}

/*
func (f Forwward) match(state request.Request) (u Upstream) {
	if p.Upstreams == nil {
		return nil
	}

	longestMatch := 0
	for _, upstream := range *p.Upstreams {
		from := upstream.From()

		if !plugin.Name(from).Matches(state.Name()) || !upstream.IsAllowedDomain(state.Name()) {
			continue
		}

		if lf := len(from); lf > longestMatch {
			longestMatch = lf
			u = upstream
		}
	}
	return u

}
*/
