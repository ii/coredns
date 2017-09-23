package proxy

import (
	"sync/atomic"
	"time"

	"github.com/coredns/coredns/plugin/pkg/healthcheck"

	"github.com/miekg/dns"
	"golang.org/x/net/context"
)

// checkDownFunc is the default function to use for CheckDown.
var checkDownFunc = func(upstream *staticUpstream) healthcheck.UpstreamHostDownFunc {
	return func(uh *healthcheck.UpstreamHost) bool {

		down := false

		uh.Lock()
		until := uh.OkUntil
		uh.Unlock()

		if !until.IsZero() && time.Now().After(until) {
			down = true
		}

		fails := atomic.LoadInt32(&uh.Fails)
		if fails >= upstream.MaxFails && upstream.MaxFails != 0 {
			down = true
		}
		return down
	}
}

func (d *dnsEx) HealthCheck(addr string) (*dns.Msg, error) {
	return d.Exchange(context.TODO(), addr, healthcheck.Payload())
}

func (g *google) HealthCheck(addr string) (*dns.Msg, error) {
	return g.Exchange(context.TODO(), addr, healthcheck.Payload())
}

func (g *grpcClient) HealthCheck(addr string) (*dns.Msg, error) {
	return g.Exchange(context.TODO(), addr, healthcheck.Payload())
}
