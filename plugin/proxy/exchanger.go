package proxy

import (
	"context"

	"github.com/coredns/coredns/request"

	"github.com/miekg/dns"
)

// Exchanger is an interface that specifies a type implementing a DNS resolver that
// can use whatever transport it likes.
type Exchanger interface {
	// Exchange performs the exchange with the upstream.
	Exchange(ctx context.Context, addr string, state request.Request) (*dns.Msg, error)

	// Protocol returns the protocol as specified in the config, grpc, https_google etc.
	Protocol() string

	// Transport returns the only transport protocol used by this Exchanger or "".
	// If the return value is "", Exchange must use `state.Proto()`.
	Transport() string

	HealthCheck(string) (*dns.Msg, error)

	OnStartup(*Proxy) error
	OnShutdown(*Proxy) error
}
