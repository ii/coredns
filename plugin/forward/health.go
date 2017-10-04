package forward

import "github.com/miekg/dns"

// For HC we send to . IN NS +norec message to the upstream. Dial timeouts and empty
// replies are considered fails, basically anything else constitutes a healthy upstream.

func (h host) Check() {
	h.Lock()

	if h.Checking {
		h.Unlock()
		return
	}

	h.Checking = true
	h.Unlock()

	return
}

func (h host) send() (*dns.Msg, error) {
	hcping := new(dns.Msg)
	hcping.SetQuestion(".", dns.TypeNS)
	hcping.RecursionDesired = false

	m, e := hcclient.Exchange()
	return m, e
}

var hcclient = func() *dns.Client {
	c := new(dns.Client)
	c.net = "tcp"
	return c
}()
