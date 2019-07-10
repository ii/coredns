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
	expiration uint32
	inception  uint32
	ttl        uint32
	dbfile     string
	origin     string
}

func (s Sign) signFunc(e *tree.Elem) bool {
	for qtype, rrs := range e.M() {
		if qtype == dns.TypeRRSIG {
			// delete
			continue
		}
		for _, pair := range s.keys {
			rrsig, err := pair.signRRs(rrs, s.origin, 3600, s.inception, s.expiration)
			if err != nil {
				return true
			}
			e.Insert(rrsig)
		}
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

	s.inception, s.expiration = lifetime(time.Now().UTC())
	s.origin = origin
	for _, pair := range s.keys {
		z.Insert(pair.Public.ToDS(dns.SHA1))
		z.Insert(pair.Public.ToDS(dns.SHA256))
		z.Insert(pair.Public.ToCDNSKEY())
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
