package forward

import (
	"github.com/coredns/coredns/plugin"

	"github.com/miekg/dns"
	"golang.org/x/net/context"
)

// Proxy represents a plugin instance that can proxy requests to another (DNS) server.
type P struct {
	udp  *Proxy
	Next plugin.Handler
}

func (p P) Name() string { return "forward" }

func (p P) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {

	p.udp.clientChan <- packet{w: w, data: r}

	return 0, nil
}
