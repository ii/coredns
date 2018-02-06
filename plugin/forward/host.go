package forward

import (
	"crypto/tls"
	"time"

	"github.com/coredns/coredns/plugin/pkg/up"

	"github.com/miekg/dns"
)

type host struct {
	addr   string
	client *dns.Client

	tlsConfig *tls.Config
	expire    time.Duration
	probe     *up.Probe

	fails uint32
}

// newHost returns a new host.
func newHost(addr string) *host {
	return &host{addr: addr, fails: 0, expire: defaultExpire, probe: up.New()}
}

// setClient sets and configures the dns.Client in host.
func (h *host) SetClient() {
	c := new(dns.Client)
	c.Net = "udp"
	c.ReadTimeout = 2 * time.Second
	c.WriteTimeout = 2 * time.Second

	if h.tlsConfig != nil {
		c.Net = "tcp-tls"
		c.TLSConfig = h.tlsConfig
	}

	h.client = c
}

const defaultExpire = 10 * time.Second
