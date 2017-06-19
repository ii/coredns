package test

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/coredns/coredns/middleware/proxy"
	"github.com/coredns/coredns/middleware/test"
	"github.com/coredns/coredns/request"

	"github.com/miekg/dns"
)

func TestEmptySecondaryZone(t *testing.T) {
	// Corefile that fails to transfer example.org.
	corefile := `example.org:0 {
		secondary {
			transfer from 127.0.0.1:1717
		}
	}
`

	i, err := CoreDNSServer(corefile)
	if err != nil {
		t.Fatalf("Could not get CoreDNS serving instance: %s", err)
	}

	udp, _ := CoreDNSServerPorts(i, 0)
	if udp == "" {
		t.Fatal("Could not get UDP listening port")
	}
	defer i.Stop()

	log.SetOutput(ioutil.Discard)

	p := proxy.NewLookup([]string{udp})
	state := request.Request{W: &test.ResponseWriter{}, Req: new(dns.Msg)}

	resp, err := p.Lookup(state, "www.example.org.", dns.TypeA)
	if err != nil {
		t.Fatal("Expected to receive reply, but didn't")
	}
	if resp.Rcode != dns.RcodeServerFailure {
		t.Fatalf("Expected reply to be a SERVFAIL, got %d", resp.Rcode)
	}
}

func TestSecondaryZoneTransfer(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	name, rm, err := test.TempFile(".", exampleCom)
	if err != nil {
		t.Fatalf("failed to create zone: %s", err)
	}
	defer rm()

	corefile := `example.com:32054 {
       file ` + name + ` {
	       transfer to 127.0.0.1:32053
       }
}
`

	prim, err := CoreDNSServer(corefile)
	if err != nil {
		t.Fatalf("Could not get CoreDNS serving instance: %s", err)
	}

	defer prim.Stop()

	corefile = `example.com:32053 {
	secondary {
		transfer from 127.0.0.1:32054
       }
}
`

	sec, err := CoreDNSServer(corefile)
	if err != nil {
		t.Fatalf("Could not get CoreDNS serving instance: %s", err)
	}

	defer sec.Stop()

	m := new(dns.Msg)
	m.SetQuestion("example.com.", dns.TypeSOA)

	r, err := dns.Exchange(m, "127.0.0.1:32054")
	if err != nil {
		t.Fatalf("Expected to receive reply, but didn't: %s", err)
	}
	if r.Answer[0].(*dns.SOA).Serial != 2017042730 {
		t.Fatalf("Expected serial of %d, got %d", 2017042730, r.Answer[0].(*dns.SOA).Serial)
	}
}

const exampleCom = `
example.com.		3600	IN	SOA	sns.dns.icann.org. noc.dns.icann.org. 2017042730 7200 3600 1209600 3600

example.com.		65118	IN	NS	a.iana-servers.net.
example.com.		65118	IN	NS	b.iana-servers.net.
cname.example.com.       434334  IN      CNAME   a.miek.nl.
`
