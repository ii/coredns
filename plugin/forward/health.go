package forward

import "github.com/miekg/dns"

// For HC we send to . IN NS +norec message to the upstream. Dial timeouts and empty
// replies are considered fails, basically anything else constitutes a healthy upstream.

func (h *host) Check() {
	h.Lock()

	if h.checking {
		h.Unlock()
		return
	}

	h.checking = true
	h.Unlock()

	return
}

func (h *host) send() (*dns.Msg, error) {
	hcping := new(dns.Msg)
	hcping.SetQuestion(".", dns.TypeNS)
	hcping.RecursionDesired = false

	// track rtt as well?
	m, _, e := hcclient.Exchange(hcping, h.addr)
	return m, e
}

var hcclient = func() *dns.Client {
	c := new(dns.Client)
	c.Net = "tcp"
	return c
}()
