package proxy

import (
	"github.com/coredns/coredns/plugin/pkg/fuzz"

	"github.com/mholt/caddy"
)

func Fuzz(data []byte) int {
	c := caddy.NewTestController("dns", "proxy . 8.8.8.8:53")
	up, err := NewStaticUpstreams(&c.Dispenser)
	if err != nil {
		panic("fuzz: should not panic:" + err.Error())
	}
	p := &Proxy{Upstreams: &up}

	return fuzz.Do(p, data)
}
