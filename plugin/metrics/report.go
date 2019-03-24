package metrics

import (
	"context"
	"time"

	"github.com/coredns/coredns/plugin/metrics/vars"
	"github.com/coredns/coredns/request"

	"github.com/miekg/dns"
)

// report reports the metrics data associated with request.
func report(ctx context.Context, req request.Request, zone, rcode string, size int, start time.Time) {
	// Proto and Family.
	net := req.Proto()
	fam := "1"
	if req.Family() == 2 {
		fam = "2"
	}

	server := WithServer(ctx)

	typ := req.QType()
	vars.RequestCount.WithLabelValues(server, zone, net, fam).Inc()
	vars.RequestDuration.WithLabelValues(server, zone).Observe(time.Since(start).Seconds())

	if req.Do() {
		vars.RequestDo.WithLabelValues(server, zone).Inc()
	}

	if _, known := monitorType[typ]; known {
		vars.RequestType.WithLabelValues(server, zone, dns.Type(typ).String()).Inc()
	} else {
		vars.RequestType.WithLabelValues(server, zone, other).Inc()
	}

	vars.ResponseSize.WithLabelValues(server, zone, net).Observe(float64(size))
	vars.RequestSize.WithLabelValues(server, zone, net).Observe(float64(req.Len()))

	vars.ResponseRcode.WithLabelValues(server, zone, rcode).Inc()
}

var monitorType = map[uint16]struct{}{
	dns.TypeAAAA:   struct{}{},
	dns.TypeA:      struct{}{},
	dns.TypeCNAME:  struct{}{},
	dns.TypeDNSKEY: struct{}{},
	dns.TypeDS:     struct{}{},
	dns.TypeMX:     struct{}{},
	dns.TypeNSEC3:  struct{}{},
	dns.TypeNSEC:   struct{}{},
	dns.TypeNS:     struct{}{},
	dns.TypePTR:    struct{}{},
	dns.TypeRRSIG:  struct{}{},
	dns.TypeSOA:    struct{}{},
	dns.TypeSRV:    struct{}{},
	dns.TypeTXT:    struct{}{},
	// Meta Qtypes
	dns.TypeIXFR: struct{}{},
	dns.TypeAXFR: struct{}{},
	dns.TypeANY:  struct{}{},
}

const other = "other"
