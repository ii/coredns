package forward

import (
	"net"
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
	upstream, err := parseForward(c)
	if err != nil {
		return err
	}

	timeout := time.Second
	udp := New(upstream.to[0].addr, 4096, timeout, timeout)
	p := P{udp: udp}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		p.Next = next
		return p
	})

	c.OnStartup(func() error {
		return p.Start(upstream.to[0].addr)
	})
	c.OnShutdown(func() error {
		return p.Close()
	})

	return nil
}

func (p *P) Start(addr string) (err error) {
	p.udp.upstream, err = net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}
	if p.udp.ConnTimeout.Nanoseconds() > 0 {
		go p.udp.freeIdleSocketsLoop()
	}
	if p.udp.ResolveTTL.Nanoseconds() > 0 {
		go p.udp.resolveUpstreamLoop()
	}
	go p.udp.handlerUpstreamPackets()
	go p.udp.handleClientPackets()
	return nil
}

func (p *P) Close() error {
	p.udp.Lock()
	p.udp.closed = true
	for _, conn := range p.udp.connsMap {
		conn.udp.Close()
	}
	p.udp.Unlock()
	return nil
}

func parseForward(c *caddy.Controller) (upstream, error) {
	u := upstream{}
	for c.Next() {
		if !c.Args(&u.from) {
			return u, c.ArgErr()
		}
		u.from = plugin.Host(u.from).Normalize()

		to := c.RemainingArgs()
		if len(to) == 0 {
			return u, c.ArgErr()
		}
		toHosts, err := dnsutil.ParseHostPortOrFile(to...)
		if err != nil {
			return u, err
		}
		u.to = toHost(toHosts)

		for c.NextBlock() {
			if err := parseBlock(c, &u); err != nil {
				return u, err
			}
		}
	}
	return u, nil
}

func parseBlock(c *caddy.Controller, u *upstream) error {
	switch c.Val() {
	default:
		return c.Errf("unknown property '%s'", c.Val())
	}
	return nil
}
