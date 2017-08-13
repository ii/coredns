package federation

import (
	"github.com/coredns/coredns/middleware"

	"github.com/miekg/dns"
	"golang.org/x/net/context"
)

// Federation contains the name to zone mapping used for federation
// in kubernetes.
type Federation struct {
	f     map[string]string
	zones []string

	Next        middleware.Handler
	Fallthrough bool
}

func New() Federation {
	return Federation{f: make(map[string]string)}
}

func (f Federation) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) {}

func (f Federation) Name() string { return "federation" }
