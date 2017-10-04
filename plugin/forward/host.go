package forward

import "sync"

type host struct {
	fails int // use atomic here?

	addr string

	sync.RWMutex
	checking bool
}

func newHost(addr string) host { return host{addr: addr} }

func toHost(addr []string) []host {
	h := make([]host, len(addr))
	for i := range addr {
		h[i].addr = addr[i]
	}
	return h
}

// Proxies has Select() returns in random order, but if conn is known that is first
// check Down function on healthyness of each upstream, use when healthy

// Down function?
// Select - not down - round robin fashsion
// Select return all up hosts in randomized
/*
can't really do this in ServeDNS because then we might select one for which we don't have a socket,
otoh if the upstream is down, we can use it, but do to randomization I don't want to ranmize I need to
now where I sent this client before
*/
