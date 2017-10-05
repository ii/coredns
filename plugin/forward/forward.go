/*
 * Copyright (c) 2016 Felipe Cavalcanti <fjfcavalcanti@gmail.com>
 * Author: Felipe Cavalcanti <fjfcavalcanti@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of
 * this software and associated documentation files (the "Software"), to deal in
 * the Software without restriction, including without limitation the rights to
 * use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
 * the Software, and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
 * FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
 * IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package forward

import (
	"net"
	"sync"
	"time"

	"github.com/coredns/coredns/request"

	"github.com/miekg/dns"
)

type conn struct {
	c    net.Conn
	w    dns.ResponseWriter
	used time.Time
}

type proxy struct {
	host         *host
	BufferSize   int
	ConnTimeout  time.Duration
	conns        map[string]conn
	closed       bool
	clientChan   chan request.Request
	upstreamChan chan request.Request
	sync.RWMutex
}

func newProxy(addr string) *proxy {
	proxy := &proxy{
		host:         newHost(addr),
		BufferSize:   udpBufSize,
		ConnTimeout:  connTimeout,
		closed:       false,
		conns:        make(map[string]conn),
		clientChan:   make(chan request.Request),
		upstreamChan: make(chan request.Request),
	}

	return proxy
}

func (p *proxy) setUsed(clientID string) {
	p.Lock()
	if _, found := p.conns[clientID]; found {
		connWrapper := p.conns[clientID]
		connWrapper.used = time.Now()
		p.conns[clientID] = connWrapper
	}
	p.Unlock()
}

func (p *proxy) clientRead(upstreamConn net.Conn, w dns.ResponseWriter) {
	clientID, _ := clientID(w)
	for {
		buffer := make([]byte, p.BufferSize)
		size, err := upstreamConn.Read(buffer)
		if err != nil {
			p.Lock()
			upstreamConn.Close()
			delete(p.conns, clientID)
			p.Unlock()
			return
		}

		ret := new(dns.Msg)
		ret.Unpack(buffer[:size])

		p.setUsed(clientID)
		p.upstreamChan <- request.Request{Req: ret, W: w}
	}
}

func (p *proxy) handlerUpstreamPackets() {
	for pa := range p.upstreamChan {
		pa.W.WriteMsg(pa.Req)
	}
}

func (p *proxy) handleClientPackets() {
	for pa := range p.clientChan {
		clientID, proto := clientID(pa.W)

		p.RLock()
		c, found := p.conns[clientID]
		p.RUnlock()

		buf, _ := pa.Req.Pack()

		if !found {
			c, err := net.DialTimeout(proto, p.host.addr, dialTimeout)
			if err != nil {
				continue
			}

			p.Lock()
			p.conns[clientID] = conn{
				c:    c,
				w:    pa.W,
				used: time.Now(),
			}
			p.Unlock()

			c.Write(buf)

			go p.clientRead(c, pa.W)
			continue
		}

		c.c.Write(buf)

		p.RLock()
		if _, ok := p.conns[clientID]; ok {
			if p.conns[clientID].used.Before(
				time.Now().Add(-p.ConnTimeout / 4)) {
				p.setUsed(clientID)
			}
		}
		p.RUnlock()
	}
}

func (p *proxy) free() {
	for !p.closed {
		time.Sleep(p.ConnTimeout)

		p.RLock()
		for client, conn := range p.conns {
			if conn.used.Before(time.Now().Add(-p.ConnTimeout)) {
				delete(p.conns, client)
			}
		}
		p.RUnlock()
	}
}

// clientID returns a string that identifies this particular client's 3-tuple.
func clientID(w dns.ResponseWriter) (id, proto string) {
	if _, ok := w.RemoteAddr().(*net.UDPAddr); ok {
		return w.RemoteAddr().String() + "udp", "udp"
	}
	return w.RemoteAddr().String() + "tcp", "tcp"
}

const (
	udpBufSize  = 4096
	dialTimeout = 1 * time.Second
	connTimeout = 2 * time.Second
)
