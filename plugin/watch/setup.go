package watch

import (
	"github.com/coredns/coredns/core/dnsserver"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/watch"

	"github.com/mholt/caddy"
)

func init() {
	caddy.RegisterPlugin("watch", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	w := NewWatcher()

	c.OnStartup(func() error {
		plugins := dnsserver.GetConfig(c).Handlers()
		for _, p := range plugins {
			if x, ok := p.(watch.Watchee); ok {
				w.watchees = append(w.watchees, x)
			}
		}
		return nil
	})

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		w.Next = next
		return w
	})

	return nil
}
