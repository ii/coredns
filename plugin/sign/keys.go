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
	Tag     uint16
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
		return Pair{}, fmt.Errorf("RR in %q is not a DNSKEY: %d", public, dnskey.Header().Rrtype)
	}
	ksk := dnskey.(*dns.DNSKEY).Flags&(1<<8) == (1<<8) && dnskey.(*dns.DNSKEY).Flags&1 == 1
	if !ksk {
		return Pair{}, fmt.Errorf("DNSKEY in %q, DNSKEY is not a CSK/KSK", public)
	}

	rp, err := os.Open(private)
	if err != nil {
		return Pair{}, err
	}
	privkey, err := dnskey.(*dns.DNSKEY).ReadPrivateKey(rp, private)
	if err != nil {
		return Pair{}, err
	}

	return Pair{Public: dnskey.(*dns.DNSKEY), Tag: dnskey.(*dns.DNSKEY).KeyTag(), Private: privkey}, nil
}

// more needed, find keypars for zone names in a directory; pick them up and parse them.
