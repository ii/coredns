package forward

import (
	"time"

	"github.com/miekg/dns"
)

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
