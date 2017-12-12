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
	nsec.TypeBitMap = filter(state.QType(), defaultBitmap)

	sigs, err := d.sign([]dns.RR{nsec}, zone, ttl, incep, expir)
	if err != nil {
		return nil, err
	}

	return append(sigs, nsec), nil
}

// defaultBitmap is the default bitmap every NSEC reply gets. We filter out the actual type ask for
// if it was in the map.
var defaultBitmap = [...]uint16{dns.TypeA, dns.TypeHINFO, dns.TypeTXT, dns.TypeAAAA, dns.TypeLOC, dns.TypeSRV, dns.TypeCERT, dns.TypeSSHFP, dns.TypeRRSIG, dns.TypeNSEC, dns.TypeTLSA, dns.TypeHIP, dns.TypeOPENPGPKEY, dns.TypeSPF}

// filter filters out t from bitmap (if it exists).
func filter(t uint16, bitmap [14]uint16) []uint16 { // 14 is len(defaultBitmap)
	for i := range bitmap {
		if bitmap[i] == t {
			return append(bitmap[:i], bitmap[i+1:]...)
		}
	}
	return defaultBitmap[:] // make a slice
}
