package forward

import (
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

func (f *Forward) Start() (err error) {
	if f.proxies[0].ConnTimeout.Nanoseconds() > 0 {
		go f.proxies[0].free()
	}

	go f.proxies[0].handlerUpstreamPackets()
	go f.proxies[0].handleClientPackets()
	return nil
}

func (f *Forward) Close() error {
	f.proxies[0].Lock()
	f.proxies[0].closed = true
	for _, conn := range f.proxies[0].conns {
		conn.udp.Close()
	}
	f.proxies[0].Unlock()
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
			p := NewProxy(h)
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
	default:
		return c.Errf("unknown property '%s'", c.Val())
	}
	return nil
}
