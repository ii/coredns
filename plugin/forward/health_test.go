package forward

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/coredns/coredns/plugin/pkg/dnstest"
	"github.com/coredns/coredns/plugin/test"
	"golang.org/x/net/context"

	"github.com/miekg/dns"
)

func TestHealth(t *testing.T) {
	const expected = 1
	i := uint32(0)
	s := dnstest.NewServer(func(w dns.ResponseWriter, r *dns.Msg) {
		if r.Question[0].Name == "." {
			// health check, answer
			ret := new(dns.Msg)
			atomic.AddUint32(&i, 1)
			ret.SetReply(r)
			w.WriteMsg(ret)
		}
		// not a health check, do a timeout
	})
	defer s.Close()

	p := NewProxy(s.Addr)
	f := New()
	f.SetProxy(p)
	defer f.Close()

	req := new(dns.Msg)
	req.SetQuestion("example.org.", dns.TypeA)

	f.ServeDNS(context.TODO(), &test.ResponseWriter{}, req)

	time.Sleep(1 * time.Second)
	i1 := atomic.LoadUint32(&i)
	if i1 != expected {
		t.Errorf("Expected number of health checks to be %d, got %d", expected, i1)
	}
}

func TestNoHealth(t *testing.T) {
	const expected = 0
	i := uint32(0)
	s := dnstest.NewServer(func(w dns.ResponseWriter, r *dns.Msg) {
		if r.Question[0].Name == "." {
			atomic.AddUint32(&i, 1)
		}
		ret := new(dns.Msg)
		ret.SetReply(r)
		w.WriteMsg(ret)
	})
	defer s.Close()

	p := NewProxy(s.Addr)
	f := New()
	f.SetProxy(p)
	defer f.Close()

	req := new(dns.Msg)
	req.SetQuestion("example.org.", dns.TypeA)

	f.ServeDNS(context.TODO(), &test.ResponseWriter{}, req)

	time.Sleep(1 * time.Second)
	i1 := atomic.LoadUint32(&i)
	if i1 != expected {
		t.Errorf("Expected number of health checks to be %d, got %d", expected, i1)
	}
}
