package forward

import "sync"

type upstream struct {
	to []host
	// hc + downfunction
}

type host struct {
	fails int // use atomic here?
	addr  string
	sync.RWMutex
	checking bool
}

func toHost(addr []string) []host {
	h := make([]host, len(addr))
	for i := range addr {
		h[i].addr = addr[i]
	}
	return h
}
