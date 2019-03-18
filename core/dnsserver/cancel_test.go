package dnsserver

import (
	"context"
	"testing"
	"time"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/dnstest"
	"github.com/coredns/coredns/plugin/test"

	"github.com/miekg/dns"
)

type sleepPlugin struct{}

func (s sleepPlugin) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	i := 0
	m := new(dns.Msg)
	m.SetReply(r)
	for {
		if plugin.Done(ctx) {
			m.Rcode = dns.RcodeBadTime // use BadTime to return something time related
			w.WriteMsg(m)
			return 0, nil
		} else {
			time.Sleep(20 * time.Millisecond)
			i++
			if i > 2 {
				m.Rcode = dns.RcodeServerFailure
				w.WriteMsg(m)
				return 0, nil
			}
		}
	}
	return 0, nil
}

func (s sleepPlugin) Name() string { return "sleepplugin" }

func sleepConfig(p plugin.Handler) *Config {
	c := &Config{Zone: "example.com.", Transport: "dns", ListenHosts: []string{"127.0.0.1"}, Port: "53"}
	c.AddPlugin(func(next plugin.Handler) plugin.Handler { return p })
	return c
}

func TestContextCancel(t *testing.T) {
	s, err := NewServer("127.0.0.1:53", []*Config{sleepConfig(sleepPlugin{})})
	if err != nil {
		t.Fatalf("Expected no error for NewServer, got %s", err)
	}

	// have to make the context here, which is a bit unfortunate because it copies some code from server.go
	ctx := context.WithValue(context.TODO(), Key{}, s)
	ctx, cancel := context.WithTimeout(ctx, 20*time.Millisecond)
	defer cancel()

	w := dnstest.NewRecorder(&test.ResponseWriter{})
	m := new(dns.Msg)
	m.SetQuestion("aaa.example.com.", dns.TypeTXT)

	s.ServeDNS(ctx, w, m)
	if w.Rcode != dns.RcodeBadTime {
		t.Error("Expected ServeDNS to be canceled by context")
	}
}
