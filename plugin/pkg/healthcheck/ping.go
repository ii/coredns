package healthcheck

import "github.com/miekg/dns"

// newPingMsg returns a dns.Msg that is used for the health check ping.
// The message's RD flag is not set. The name, type and class or
// ".", NS, IN repectively.
func newPingMsg() *dns.Msg {
	m := new(dns.Msg)
	m.SetQuestion(".", dns.TypeNS)
	m.RecursionDesired = false
	return m
}
