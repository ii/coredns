package dnsserver

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/miekg/dns"
)

type zoneAddr struct {
	Zone      string
	Port      string
	Transport string // dns, tls or grpc
}

// String return the string representation of z.
func (z zoneAddr) String() string { return z.Transport + "://" + z.Zone + ":" + z.Port }

// Transport returns the protocol of the string s
func Transport(s string) string {
	switch {
	case strings.HasPrefix(s, TransportTLS+"://"):
		return TransportTLS
	case strings.HasPrefix(s, TransportDNS+"://"):
		return TransportDNS
	case strings.HasPrefix(s, TransportGRPC+"://"):
		return TransportGRPC
	}
	return TransportDNS
}

// normalizeZone parses an zone string into a structured format with separate
// host, and port portions, as well as the original input string.
func normalizeZone(str string) (zoneAddr, error) {
	var err error

	// Default to DNS if there isn't a transport protocol prefix.
	trans := TransportDNS

	switch {
	case strings.HasPrefix(str, TransportTLS+"://"):
		trans = TransportTLS
		str = str[len(TransportTLS+"://"):]
	case strings.HasPrefix(str, TransportDNS+"://"):
		trans = TransportDNS
		str = str[len(TransportDNS+"://"):]
	case strings.HasPrefix(str, TransportGRPC+"://"):
		trans = TransportGRPC
		str = str[len(TransportGRPC+"://"):]
	}

	// If there is: :[0-9]+ on the end we assume this is the port. This works for (ascii) domain
	// names and our reverse syntax, which always needs a /mask *before* the port.
	// So from the back, find first colon, and then check if its a number.
	host := str
	port := ""

	colon := strings.LastIndex(str, ":")
	if colon == len(str)-1 {
		return zoneAddr{}, fmt.Errorf("expecting data after last colon: %q", str)
	}
	if colon != -1 {
		if p, err := strconv.Atoi(str[colon+1:]); err == nil {
			port = strconv.Itoa(p)
			host = str[:colon]
		}
	}

	// TODO(miek): this should take escaping into account.
	if len(host) > 255 {
		return zoneAddr{}, fmt.Errorf("specified zone is too long: %d > 255", len(host))
	}

	_, d := dns.IsDomainName(host)
	if !d {
		return zoneAddr{}, fmt.Errorf("zone is not a valid domain name: %s", host)
	}

	// Check if it parses as a reverse zone, if so we use that. Must be fully
	// specified IP and mask and mask % 8 = 0.
	ip, net, err := net.ParseCIDR(host)
	if err == nil {
		if rev, e := dns.ReverseAddr(ip.String()); e == nil {
			ones, bits := net.Mask.Size()
			if (bits-ones)%8 == 0 {
				offset, end := 0, false
				for i := 0; i < (bits-ones)/8; i++ {
					offset, end = dns.NextLabel(rev, offset)
					if end {
						break
					}
				}
				host = rev[offset:]
			}
		}
	}

	if port == "" {
		if trans == TransportDNS {
			port = Port
		}
		if trans == TransportTLS {
			port = TLSPort
		}
		if trans == TransportGRPC {
			port = GRPCPort
		}
	}

	return zoneAddr{Zone: strings.ToLower(dns.Fqdn(host)), Port: port, Transport: trans}, nil
}

// Supported transports.
const (
	TransportDNS  = "dns"
	TransportTLS  = "tls"
	TransportGRPC = "grpc"
)
