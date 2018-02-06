package watch

import (
	"github.com/coredns/coredns/core/dnsserver"
	"log"

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
			log.Printf("Checking if %s is a Watchee\n", p.Name())
			if x, ok := p.(watch.Watchee); ok {
				log.Printf("Yes\n")
				w.watchees = append(w.watchees, x)
			} else {
				log.Printf("No\n")
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
