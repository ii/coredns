package sign

import (
	"fmt"
	"os"
	"time"

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
	// all types that should not be signed, have been dropped when reading the zone
	for _, rrs := range e.M() {
		for _, pair := range s.keys {
			rrsig, err := pair.signRRs(rrs, s.origin, s.ttl, s.inception, s.expiration)
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

	z, err := parse(rd, origin, s.dbfile)
	if err != nil {
		return err
	}

	s.inception, s.expiration = lifetime(time.Now().UTC())
	s.origin = origin

	s.ttl = z.Apex.SOA.Header().Ttl
	z.Apex.SOA.Serial = uint32(time.Now().Unix())

	for _, pair := range s.keys {
		z.Insert(pair.Public.ToDS(dns.SHA1))
		z.Insert(pair.Public.ToDS(dns.SHA256))
		z.Insert(pair.Public.ToCDNSKEY())
	}
	for _, pair := range s.keys {
		rrsig, err := pair.signRRs([]dns.RR{z.Apex.SOA}, s.origin, s.ttl, s.inception, s.expiration)
		if err != nil {
			return err
		}
		z.Insert(rrsig)
		rrsig, err = pair.signRRs(z.Apex.NS, s.origin, s.ttl, s.inception, s.expiration)
		if err != nil {
			return err
		}
		z.Insert(rrsig)
	}

	z.Tree.Do(s.signFunc)

	fmt.Println(z.Apex.SOA.String())
	for _, rr := range z.Apex.SIGSOA {
		fmt.Println(rr.String())
	}
	for _, rr := range z.Apex.NS {
		fmt.Println(rr.String())
	}
	for _, rr := range z.Apex.SIGNS {
		fmt.Println(rr.String())
	}
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
