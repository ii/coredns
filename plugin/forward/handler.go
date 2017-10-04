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
	proxies []*proxy

	from string
	// ignoredNames here

	Next plugin.Handler
}

// Name implements plugin.Handler.
func (f Forward) Name() string { return "forward" }

// ServeDNS implements plugin.Handler.
func (f Forward) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {

	state := request.Request{W: w, Req: r}
	if !f.match(state) {
		return plugin.NextOrFailure(f.Name(), f.Next, ctx, w, r)
	}

	f.proxies[0].clientChan <- state
	// inspect w is something is written?? If that is hard because way may use an old w in the forward code...

	return 0, nil
}

func (f Forward) match(state request.Request) bool {
	from := f.from

	if !plugin.Name(from).Matches(state.Name()) { // || !upstream.IsAllowedDomain(state.Name()) {
		return false
	}

	return true
}
