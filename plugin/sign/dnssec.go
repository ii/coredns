package sign

import (
	"github.com/miekg/dns"
)

func (p Pair) signRRs(rrs []dns.RR, signerName string, origTTL, ttl, incep, expir uint32) (*dns.RRSIG, error) {
	rrsig := p.newRRSIG(signerName, origTTL, ttl, incep, expir)
	e := rrsig.Sign(p.Private, rrs)
	return rrsig, e
}

func (p Pair) newRRSIG(signerName string, origTTL, ttl, incep, expir uint32) *dns.RRSIG {
	return &dns.RRSIG{
		Hdr:        dns.RR_Header{Rrtype: dns.TypeRRSIG, Ttl: ttl},
		Algorithm:  p.Public.Algorithm,
		KeyTag:     p.Tag,
		SignerName: signerName,
		OrigTtl:    origTTL,
		Inception:  incep,
		Expiration: expir,
	}
}
