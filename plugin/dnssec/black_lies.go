package dnssec

import (
	"github.com/coredns/coredns/request"

	"github.com/miekg/dns"
)

// nsec returns an NSEC useful for NXDOMAIN respsones.
// See https://tools.ietf.org/html/draft-valsorda-dnsop-black-lies-00
// For example, a request for the non-existing name a.example.com would
// cause the following NSEC record to be generated:
//	a.example.com. 3600 IN NSEC \000.a.example.com. ( RRSIG NSEC )
// This inturn makes every NXDOMAIN answer a NODATA one, don't forget to flip
// the header rcode to NOERROR.
func (d Dnssec) nsec(state request.Request, zone string, ttl, incep, expir uint32) ([]dns.RR, error) {
	nsec := &dns.NSEC{}
	nsec.Hdr = dns.RR_Header{Name: state.QName(), Ttl: ttl, Class: dns.ClassINET, Rrtype: dns.TypeNSEC}
	nsec.NextDomain = "\\000." + state.QName()
	if state.Name() == zone {
		nsec.TypeBitMap = filter18(state.QType(), apexBitmap)
	} else {
		nsec.TypeBitMap = filter14(state.QType(), zoneBitmap)
	}

	sigs, err := d.sign([]dns.RR{nsec}, zone, ttl, incep, expir)
	if err != nil {
		return nil, err
	}

	return append(sigs, nsec), nil
}

// The nsec bit maps we return
var (
	zoneBitmap = [...]uint16{dns.TypeA, dns.TypeHINFO, dns.TypeTXT, dns.TypeAAAA, dns.TypeLOC, dns.TypeSRV, dns.TypeCERT, dns.TypeSSHFP, dns.TypeRRSIG, dns.TypeNSEC, dns.TypeTLSA, dns.TypeHIP, dns.TypeOPENPGPKEY, dns.TypeSPF}
	apexBitmap = [...]uint16{dns.TypeA, dns.TypeNS, dns.TypeSOA, dns.TypeHINFO, dns.TypeMX, dns.TypeTXT, dns.TypeAAAA, dns.TypeLOC, dns.TypeSRV, dns.TypeCERT, dns.TypeSSHFP, dns.TypeRRSIG, dns.TypeNSEC, dns.TypeDNSKEY, dns.TypeTLSA, dns.TypeHIP, dns.TypeOPENPGPKEY, dns.TypeSPF}
)

// filter14 filters out t from bitmap (if it exists).
func filter14(t uint16, bitmap [14]uint16) []uint16 {
	for i := range bitmap {
		if bitmap[i] == t {
			return append(bitmap[:i], bitmap[i+1:]...)
		}
	}
	return zoneBitmap[:] // make a slice
}

// filter18 filters out t from bitmap (if it exists).
func filter18(t uint16, bitmap [18]uint16) []uint16 {
	for i := range bitmap {
		if bitmap[i] == t {
			return append(bitmap[:i], bitmap[i+1:]...)
		}
	}
	return apexBitmap[:] // make a slice
}
