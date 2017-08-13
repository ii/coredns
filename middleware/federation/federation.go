package federation

import (
	"github.com/coredns/coredns/middleware"
	"github.com/coredns/coredns/request"

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
	Federations FederationFunc
}

// FederationFunc needs to be implemented by any middleware that implements
// federation. Right now this is only the kubernetes middleware.
//type FederationFunc func(state request.Request) ([]msg.Service, []msg.Service)
type FederationFunc func(state request.Request) string // Use for testing

func New() Federation {
	return Federation{f: make(map[string]string)}
}

func (f Federation) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) {
	state := request.Request{W: w, Req: r}

	println(f.Federations(state))
}

func (f Federation) Name() string { return "federation" }
