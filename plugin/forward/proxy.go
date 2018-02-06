package forward

import (
	"crypto/tls"
	"sync"
	"time"

	"github.com/miekg/dns"
)

// Proxy defines an upstream host.
type Proxy struct {
	host *host

	transport *transport

	// copied from Forward.
	forceTCP bool

	stop chan bool

	sync.RWMutex
}

// NewProxy returns a new proxy.
func NewProxy(addr string) *Proxy {
	host := newHost(addr)

	p := &Proxy{
		host:      host,
		stop:      make(chan bool),
		transport: newTransport(host),
	}
	return p
}

// SetTLSConfig sets the TLS config in the lower p.host.
func (p *Proxy) SetTLSConfig(cfg *tls.Config) { p.host.tlsConfig = cfg }

// SetExpire sets the expire duration in the lower p.host.
func (p *Proxy) SetExpire(expire time.Duration) { p.host.expire = expire }

// Dial connects to the host in p with the configured transport.
func (p *Proxy) Dial(proto string) (*dns.Conn, error) { return p.transport.Dial(proto) }

// Yield returns the connection to the pool.
func (p *Proxy) Yield(c *dns.Conn) { p.transport.Yield(c) }

// Down returns if this proxy is up or down.
func (p *Proxy) Down(maxfails uint32) bool { return p.host.down(maxfails) }

// close stops the health checking goroutine.
func (p *Proxy) close() { p.host.probe.Stop() }

// start starts the proxy's healthchecking.
func (p *Proxy) start() {
	p.host.SetClient()
	p.host.probe.Start(hcDuration)
}

// Healthcheck kicks of a round of health checks for this proxy.
func (p *Proxy) Healthcheck() { p.host.probe.Do(p.host.Check) }

const (
	dialTimeout = 4 * time.Second
	timeout     = 2 * time.Second
	hcDuration  = 500 * time.Millisecond
)
