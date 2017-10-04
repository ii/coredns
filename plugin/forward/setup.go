package forward

import (
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

	timeout := time.Second

	udp := New(f.udp.to[0].addr, timeout)
	f.udp = udp

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		f.Next = next
		return f
	})

	c.OnStartup(func() error {
		return f.Start(upstream.to[0].addr)
	})
	c.OnShutdown(func() error {
		return f.Close()
	})

	return nil
}

func (f *Forward) Start(addr string) (err error) {
	f.udp.addr = addr
	if err != nil {
		return err
	}
	if f.udp.ConnTimeout.Nanoseconds() > 0 {
		go f.udp.free()
	}

	go f.udp.handlerUpstreamPackets()
	go f.udp.handleClientPackets()
	return nil
}

func (f *Forward) Close() error {
	f.udp.Lock()
	f.udp.closed = true
	for _, conn := range f.udp.conns {
		conn.udp.Close()
	}
	f.udp.Unlock()
	return nil
}

func parseForward(c *caddy.Controller) (Forward, error) {
	f := Forward{}
	u := upstream{}
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
		u.to = toHost(toHosts)

		for c.NextBlock() {
			if err := parseBlock(c, &f); err != nil {
				return f, err
			}
		}
	}
	f.udp = u
	return f, nil
}

func parseBlock(c *caddy.Controller, f *Forward) error {
	switch c.Val() {
	default:
		return c.Errf("unknown property '%s'", c.Val())
	}
	return nil
}
