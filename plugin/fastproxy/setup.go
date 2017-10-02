package fastproxy

import (
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
	//	udp := New(
	//func New(bindPort int, bindAddress string, upstreamAddress string, upstreamPort int, bufferSize int, connTimeout time.Duration, resolveTTL time.Duration) *Proxy {
	//)
	p := P{}
	//	p = P{udp: New()}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		p.Next = next
		return p
	})
	return nil
}
