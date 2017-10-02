package fastproxy

import (
	"fmt"
	"net"
	"time"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"

	"github.com/mholt/caddy"
)

func init() {
	caddy.RegisterPlugin("fastproxy", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	timeout := time.Second
	udp := New("8.8.8.8", 53, 4096, timeout, timeout)
	p := P{udp: udp}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		p.Next = next
		return p
	})

	c.OnStartup(func() error {
		return p.Start("8.8.8.8", 53)
	})
	c.OnShutdown(func() error {
		return p.Close()
	})

	return nil
}

func (p *P) Start(addr string, port int) (err error) {
	p.udp.upstream, err = net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", addr, port))
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
	p.udp.connectionsLock.Lock()
	p.udp.closed = true
	for _, conn := range p.udp.connsMap {
		conn.udp.Close()
	}
	if p.udp.listenerConn != nil {
		p.udp.listenerConn.Close()
	}
	p.udp.connectionsLock.Unlock()
	return nil
}
