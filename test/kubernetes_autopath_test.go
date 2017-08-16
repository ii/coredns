// +build k8s

package test

import (
	"testing"

	"github.com/miekg/dns"
)

func TestLookupAutoPathKubernetes(t *testing.T) {
	corefile := `cluster.local:0 {
		kubernetes {
                endpoint http://localhost:8080
		namespaces test-1
		pods insecure
	}
	autopath @kubernetes
	proxy . 8.8.8.8:53
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

	m := new(dns.Msg)
	m.SetQuestion("google.com", dns.TypeA)

	r, err := dns.Exchange(m, udp)
	if err != nil {
		t.Fatalf("Failed to sent query: %q", err)
	}
	t.Logf("%v\n", r)

}
