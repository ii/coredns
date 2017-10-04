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

type connection struct {
	udp  *net.UDPConn
	w    dns.ResponseWriter
	used time.Time
}

type proxy struct {
	host         *host
	BufferSize   int
	ConnTimeout  time.Duration
	conns        map[string]connection
	closed       bool
	clientChan   chan (request.Request)
	upstreamChan chan (request.Request)
	sync.RWMutex
}

func newProxy(addr string) *proxy {
	proxy := &proxy{
		host:         newHost(addr),
		BufferSize:   udpBufSize,
		ConnTimeout:  connTimeout,
		closed:       false,
		conns:        make(map[string]connection),
		clientChan:   make(chan request.Request),
		upstreamChan: make(chan request.Request),
	}

	return proxy
}

func (p *proxy) setUsed(clientAddrString string) {
	p.Lock()
	if _, found := p.conns[clientAddrString]; found {
		connWrapper := p.conns[clientAddrString]
		connWrapper.used = time.Now()
		p.conns[clientAddrString] = connWrapper
	}
	p.Unlock()
}

func (p *proxy) clientRead(upstreamConn *net.UDPConn, w dns.ResponseWriter) {
	clientAddrString := w.RemoteAddr().String()
	for {
		buffer := make([]byte, p.BufferSize)
		size, _, err := upstreamConn.ReadFromUDP(buffer)
		if err != nil {
			p.Lock()
			upstreamConn.Close()
			delete(p.conns, clientAddrString)
			p.Unlock()
			return
		}

		ret := new(dns.Msg)
		ret.Unpack(buffer[:size])

		p.setUsed(clientAddrString)
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
		packetSourceString := pa.W.RemoteAddr().String()

		p.RLock()
		conn, found := p.conns[packetSourceString]
		p.RUnlock()

		buf, _ := pa.Req.Pack()

		if !found {
			c, err := net.DialTimeout("udp", p.host.addr, dialTimeout)
			if err != nil {
				// return err err
				// Wait return???
				return
			}
			conn := c.(*net.UDPConn)

			p.Lock()
			p.conns[packetSourceString] = connection{
				udp:  conn,
				w:    pa.W,
				used: time.Now(),
			}
			p.Unlock()

			conn.Write(buf)
			go p.clientRead(conn, pa.W)
		} else {
			conn.udp.Write(buf)

			p.RLock()
			updateUsed := false
			if _, ok := p.conns[packetSourceString]; ok {
				if p.conns[packetSourceString].used.Before(
					time.Now().Add(-p.ConnTimeout / 4)) {
					updateUsed = true
				}
			}
			p.RUnlock()

			if updateUsed {
				p.setUsed(packetSourceString)
			}
		}
	}
}

func (p *proxy) free() {
	for !p.closed {
		time.Sleep(p.ConnTimeout)
		var clientsToTimeout []string

		p.RLock()
		for client, conn := range p.conns {
			if conn.used.Before(time.Now().Add(-p.ConnTimeout)) {
				clientsToTimeout = append(clientsToTimeout, client)
			}
		}
		p.RUnlock()

		p.Lock()
		for _, client := range clientsToTimeout {
			p.conns[client].udp.Close()
			delete(p.conns, client)
		}
		p.Unlock()
	}
}

const (
	udpBufSize  = 4096
	dialTimeout = 1 * time.Second
	connTimeout = 2 * time.Second
)
