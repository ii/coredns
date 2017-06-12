package cache

import (
	"context"
	"testing"
	"time"

	"github.com/coredns/coredns/middleware/pkg/cache"
	"github.com/coredns/coredns/middleware/pkg/dnsrecorder"

	"github.com/coredns/coredns/middleware/test"
	"github.com/miekg/dns"
)

func TestPrefetch(t *testing.T) {
	c := &Cache{Zones: []string{"."}, pcap: defaultCap, ncap: defaultCap, pttl: maxTTL, nttl: maxTTL}
	c.pcache = cache.New(c.pcap)
	c.ncache = cache.New(c.ncap)
	c.prefetch = 1
	c.duration = 1 * time.Second

	ctx := context.TODO()

	req := new(dns.Msg)
	req.SetQuestion("lowttl.example.org", dns.TypeA)
	c.Next = test.NextHandler(dns.RcodeSuccess, nil)

	rec := dnsrecorder.New(&test.ResponseWriter{})
	code, _ := c.ServeDNS(ctx, rec, req)
	println(code)

}
