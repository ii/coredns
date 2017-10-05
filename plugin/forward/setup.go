package forward

import (
	"strconv"
	"time"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/dnsutil"

	"github.com/mholt/caddy"
)

func init() {
	caddy.RegisterPlugin("forward", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	f, err := parseForward(c)
	if err != nil {
		return err
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		f.Next = next
		return f
	})

	c.OnStartup(func() error {
		return f.Start()
	})
	c.OnShutdown(func() error {
		return f.Close()
	})

	return nil
}

// Start starts a goroutines for all proxies.
func (f *Forward) Start() (err error) {
	for _, p := range f.proxies {
		if p.ConnTimeout.Nanoseconds() > 0 {
			go p.free()
		}

		go p.handlerUpstreamPackets()
		go p.handleClientPackets()
		go p.healthCheck()
	}
	return nil
}

// Close stops all configures proxies.
func (f *Forward) Close() error {
	for _, p := range f.proxies {
		p.Lock()
		p.closed = true
		for _, conn := range p.conns {
			conn.c.Close()
		}
		p.Unlock()
	}
	return nil
}

func parseForward(c *caddy.Controller) (Forward, error) {
	f := Forward{}
	for c.Next() {
		if !c.Args(&f.from) {
			return f, c.ArgErr()
		}
		f.from = plugin.Host(f.from).Normalize()

		to := c.RemainingArgs()
		if len(to) == 0 {
			return f, c.ArgErr()
		}
		toHosts, err := dnsutil.ParseHostPortOrFile(to...)
		if err != nil {
			return f, err
		}
		for _, h := range toHosts {
			p := newProxy(h)
			f.proxies = append(f.proxies, p)

		}

		for c.NextBlock() {
			if err := parseBlock(c, &f); err != nil {
				return f, err
			}
		}
	}
	return f, nil
}

func parseBlock(c *caddy.Controller, f *Forward) error {
	switch c.Val() {
	case "except":
		ignore := c.RemainingArgs()
		if len(ignore) == 0 {
			return c.ArgErr()
		}
		for i := 0; i < len(ignore); i++ {
			ignore[i] = plugin.Host(ignore[i]).Normalize()
		}
		f.ignored = ignore
	case "max_fails":
		if !c.NextArg() {
			return c.ArgErr()
		}
		n, err := strconv.Atoi(c.Val())
		if err != nil {
			return err
		}
		f.maxfails = uint32(n)
	case "health_check":
		if !c.NextArg() {
			return c.ArgErr()
		}
		dur, err := time.ParseDuration(c.Val())
		if err != nil {
			return err
		}
		for i := range f.proxies {
			f.proxies[i].hcInterval = dur
		}
	default:
		return c.Errf("unknown property '%s'", c.Val())
	}
	return nil
}
