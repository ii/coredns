package sign

import (
	"fmt"
	"os"
	"time"

	"github.com/coredns/coredns/plugin/file"
	"github.com/coredns/coredns/plugin/file/tree"

	"github.com/miekg/dns"
)

type Sign struct {
	keys       []Pair
	expiration time.Duration
	inception  time.Duration
	dbfile     string
}

func (s Sign) signFunc(e *tree.Elem) bool {
	for qtype, rrs := range e.M() {
		println(qtype)
		println(rrs[0].String())
	}
	return false
}

func (s Sign) Sign(origin string) error {
	rd, err := os.Open(s.dbfile)
	if err != nil {
		return err
	}

	z, err := file.Parse(rd, origin, s.dbfile, 0)
	if err != nil {
		return err
	}
	// use SOA TTL?

	incep, expir := lifetime(time.Now().UTC())
	// CDS, CDNSKEY records.
	cds := []dns.RR{}
	cdnskeys := []dns.RR{}
	for _, pair := range s.keys {
		cds = append(cds, pair.Public.ToDS(dns.SHA1))
		cds = append(cds, pair.Public.ToDS(dns.SHA256))
		cdnskeys = append(cdnskeys, pair.Public.ToCDNSKEY())
	}
	for _, pair := range s.keys {
		rrsig, err := pair.signRRs(cds, origin, 3600, incep, expir)
		if err != nil {
			return err
		}
		println(rrsig.String())
		rrsig, err = pair.signRRs(cdnskeys, origin, 3600, incep, expir)
		if err != nil {
			return err
		}
		println(rrsig.String())
	}

	// sign it
	z.Tree.Do(s.signFunc)

	// print it
	z.Tree.Do(func(e *tree.Elem) bool {
		for _, r := range e.All() {
			fmt.Println(r.String())
		}
		return false
	})

	return nil
}

func lifetime(now time.Time) (uint32, uint32) {
	incep := uint32(now.Add(-3 * time.Hour).Unix()) // -(2+1) hours, be sure to catch daylight saving time and such
	expir := uint32(now.Add(threeWeeks).Unix())     // sign for 21 days
	return incep, expir
}

const threeWeeks = 3 * 7 * 24 * time.Hour
