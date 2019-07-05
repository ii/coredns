package sign

import (
	"os"

	"github.com/miekg/dns"
)

type Sign struct {
	dnskeys []*dns.DNSKEY
	expiration
	dbfile string
}

func (s Sign) Sign(origin string) error {
	rd, err := os.Open(dbfile)
	if err != nil {
		return err
	}

	z, err := Parse(rd, origin, s.dbfile, 0)
	if err != nil {
		return err
	}

}
