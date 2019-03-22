// Package cancel implements a plugin adds a canceling context to each request.
package cancel

import (
	"context"
	"time"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"

	"github.com/mholt/caddy"
	"github.com/miekg/dns"
)

func init() {
	caddy.RegisterPlugin("cancel", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return Cancel{Next: next}
	})

	return nil
}

// Cancel is a plugin that adds a canceling context to each request's context.
type Cancel struct {
	Next plugin.Handler
}

// ServeDNS implements the plugin.Handler interface.
func (c Cancel) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)

	code, err := plugin.NextOrFailure(c.Name(), c.Next, ctx, w, r)

	cancel()

	return code, err
}

// Name implements the Handler interface.
func (c Cancel) Name() string { return "cancel" }

// timeout made a variable so we can override in a test.
var timeout = 5001 * time.Millisecond
