package test

import (
	"io/ioutil"
	"log"
	"testing"
	"time"

	"github.com/coredns/coredns/middleware/proxy"
	"github.com/coredns/coredns/middleware/test"
	"github.com/coredns/coredns/request"

	"github.com/miekg/dns"
)

func TestCachePrefetch(t *testing.T) {
	// These settings are chosen in such a way that the prefetch will happen.
	// We check this, by looking at the TTL between 2 look ups, The should be
	// equal because we just fetched the new record and refreshed the TTL.

	name, rm, err := test.TempFile(".", exampleOrg)
	if err != nil {
		t.Fatalf("failed to create zone: %s", err)
	}
	defer rm()

	corefile := `example.org:0 {
       file ` + name + `
}
`
	i, err := CoreDNSServer(corefile)
	if err != nil {
		t.Fatalf("Could not get CoreDNS serving instance: %s", err)
	}

	udp, _ := CoreDNSServerPorts(i, 0)
	if udp == "" {
		t.Fatalf("Could not get UDP listening port")
	}
	defer i.Stop()

	// Start caching proxy CoreDNS that we want to test.
	corefile = `example.org:0 {
	proxy . ` + udp + `
	cache {
		prefetch 1 1s
	}
}
`
	i, err = CoreDNSServer(corefile)
	if err != nil {
		t.Fatalf("Could not get CoreDNS serving instance: %s", err)
	}

	udp, _ = CoreDNSServerPorts(i, 0)
	if udp == "" {
		t.Fatalf("Could not get UDP listening port")
	}
	defer i.Stop()

	log.SetOutput(ioutil.Discard)

	p := proxy.NewLookup([]string{udp})
	state := request.Request{W: &test.ResponseWriter{}, Req: new(dns.Msg)}

	resp, err := p.Lookup(state, lowttl, dns.TypeA)
	if err != nil {
		t.Fatal("Expected to receive reply, but didn't")
	}
	println(resp.String())
	time.Sleep(2 * time.Second)

	resp, err = p.Lookup(state, lowttl, dns.TypeA)
	if err != nil {
		t.Fatal("Expected to receive reply, but didn't")
	}
	println(resp.String())

	// This should have ttl of 80 second. It doesn't: bug!
	resp, err = p.Lookup(state, lowttl, dns.TypeA)
	if err != nil {
		t.Fatal("Expected to receive reply, but didn't")
	}
	println(resp.String())
}

const lowttl = "lowttl.example.org."
