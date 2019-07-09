package sign

import (
	"crypto"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/miekg/dns"
)

type Pair struct {
	Public  *dns.DNSKEY
	Private crypto.PrivateKey
}

// readKeyPair read the public and private key from disk.
func readKeyPair(public, private string) (Pair, error) {
	rk, err := os.Open(public)
	if err != nil {
		return Pair{}, err
	}
	b, err := ioutil.ReadAll(rk)
	if err != nil {
		return Pair{}, err
	}
	dnskey, err := dns.NewRR(string(b))
	if err != nil {
		return Pair{}, err
	}
	if _, ok := dnskey.(*dns.DNSKEY); !ok {
		return Pair{}, fmt.Errorf("file %s, doesn't parse as a DNSKEY: %d", public, dnskey.Header().Rrtype)
	}

	rp, err := os.Open(private)
	if err != nil {
		return Pair{}, err
	}
	privkey, err := dnskey.(*dns.DNSKEY).ReadPrivateKey(rp, private)
	if err != nil {
		return Pair{}, err
	}

	return Pair{Public: dnskey.(*dns.DNSKEY), Private: privkey}, nil
}

// more needed, find keypars for zone names in a directory; pick them up and parse them.
