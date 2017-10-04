package forward

import "sync"

type host struct {
	fails int // use atomic here?

	addr string

	sync.RWMutex
	checking bool
}

func newHost(addr string) *host { return &host{addr: addr} }

// Proxies has Select() returns in random order, but if conn is known that is first
// check Down function on healthyness of each upstream, use when healthy
