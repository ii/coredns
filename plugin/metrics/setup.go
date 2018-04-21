package metrics

import (
	"net"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"

	"github.com/mholt/caddy"
)

func init() {
	caddy.RegisterPlugin("prometheus", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})

	uniqAddr = addrs{a: make(map[string]int)}
}

func setup(c *caddy.Controller) error {
	m, err := prometheusParse(c)
	if err != nil {
		return plugin.Error("prometheus", err)
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		m.Next = next
		return m
	})

	c.OnStartup(func() error {
		for a, v := range uniqAddr.a {
			println("setting up for", a)
			if v == todo {
				c.OncePerServerBlock(m.OnStartup)
			}
			uniqAddr.a[a] = done
		}
		return nil
	})

	c.OnShutdown(m.OnFinalShutdown)
	c.OnRestart(m.OnRestart)

	return nil
}

func prometheusParse(c *caddy.Controller) (*Metrics, error) {
	var met = New(defaultAddr)

	defer func() {
		uniqAddr.SetAddress(met.Addr)
	}()

	i := 0
	for c.Next() {
		if i > 0 {
			return nil, plugin.ErrOnce
		}
		i++

		for _, z := range c.ServerBlockKeys {
			met.AddZone(plugin.Host(z).Normalize())
		}
		args := c.RemainingArgs()

		switch len(args) {
		case 0:
		case 1:
			met.Addr = args[0]
			_, _, e := net.SplitHostPort(met.Addr)
			if e != nil {
				return met, e
			}
		default:
			return met, c.ArgErr()
		}
	}
	return met, nil
}

// addrs keeps track on which addrs we listen, so we only start one listener, is
// prometheus is used in multiple Server Blocks.
type addrs struct {
	a map[string]int
}

var uniqAddr addrs

func (a *addrs) SetAddress(addr string) {
	// If already there and set to done, we've already started this listener.
	if a.a[addr] == done {
		return
	}
	a.a[addr] = todo
}

// defaultAddr is the address the where the metrics are exported by default.
const defaultAddr = "localhost:9153"

const (
	todo = 1
	done = 2
)
