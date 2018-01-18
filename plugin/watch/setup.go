package watch

import (
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"

	"github.com/mholt/caddy"
)

func init() {
	caddy.RegisterPlugin("watch", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	w := NewWatch()
	
        c.OnStartup(func() error {
                plugins := dnsserver.GetConfig(c).Handlers()
                for _, p := range plugins {
                        if x, ok := p.(Watcher); ok {
                                w.watchers = append(w.watchers, x)
                        }
                }
                return nil
        })

	// Do not do AddPlugin

	return nil
}
