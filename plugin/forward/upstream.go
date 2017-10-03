package forward

type host struct {
	addr string
}

func toHost(addr []string) []host {
	h := make([]host, len(addr))
	for i := range addr {
		h[i].addr = addr[i]
	}
	return h
}

type upstream struct {
	from string
	to   []host
}
