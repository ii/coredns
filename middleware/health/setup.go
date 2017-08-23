package health

import (
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/middleware"

	"github.com/mholt/caddy"
)

func init() {
	caddy.RegisterPlugin("health", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	addr, err := healthParse(c)
	if err != nil {
		return middleware.Error("health", err)
	}

	h := &health{Addr: addr}

	c.OnStartup(func() error {
		for he := range healthers {
			m := dnsserver.GetConfig(c).Handler(he)
			if x, ok := m.(Healther); ok {
				h.h = append(h.h, x)
			}
		}
		return nil
	})

	c.OnStartup(func() error {
		// Kick of Go-func that every second calls h.Poll()
		h.Poll()
		return nil
	})

	c.OnStartup(h.Startup)
	c.OnShutdown(h.Shutdown)

	// Don't do AddMiddleware, as health is not *really* a middleware just a separate webserver running.
	return nil
}

func healthParse(c *caddy.Controller) (string, error) {
	addr := ""
	for c.Next() {
		args := c.RemainingArgs()

		switch len(args) {
		case 0:
		case 1:
			addr = args[0]
		default:
			return "", c.ArgErr()
		}
	}
	return addr, nil
}
