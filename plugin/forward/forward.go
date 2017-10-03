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

	"github.com/miekg/dns"
)

type connection struct {
	udp          *net.UDPConn
	w            dns.ResponseWriter
	lastActivity time.Time
}

type packet struct {
	w    dns.ResponseWriter
	data *dns.Msg
}

type Proxy struct {
	addr         string
	BufferSize   int // can go
	ConnTimeout  time.Duration
	connsMap     map[string]connection
	closed       bool
	clientChan   chan (packet)
	upstreamChan chan (packet)
	sync.RWMutex
}

func New(addr string, bufferSize int, connTimeout time.Duration) *Proxy {
	proxy := &Proxy{
		BufferSize:   bufferSize,
		ConnTimeout:  connTimeout,
		addr:         addr,
		connsMap:     make(map[string]connection),
		closed:       false,
		clientChan:   make(chan packet),
		upstreamChan: make(chan packet),
	}

	return proxy
}

func (p *Proxy) updateClientLastActivity(clientAddrString string) {
	p.Lock()
	if _, found := p.connsMap[clientAddrString]; found {
		connWrapper := p.connsMap[clientAddrString]
		connWrapper.lastActivity = time.Now()
		p.connsMap[clientAddrString] = connWrapper
	}
	p.Unlock()
}

func (p *Proxy) clientConnectionReadLoop(upstreamConn *net.UDPConn, w dns.ResponseWriter) {
	clientAddrString := w.RemoteAddr().String()
	for {
		buffer := make([]byte, p.BufferSize)
		size, _, err := upstreamConn.ReadFromUDP(buffer)
		if err != nil {
			p.Lock()
			upstreamConn.Close()
			delete(p.connsMap, clientAddrString)
			p.Unlock()
			return
		}

		ret := new(dns.Msg)
		ret.Unpack(buffer[:size])

		p.updateClientLastActivity(clientAddrString)
		p.upstreamChan <- packet{
			data: ret,
			w:    w,
		}
	}
}

func (p *Proxy) handlerUpstreamPackets() {
	for pa := range p.upstreamChan {
		pa.w.WriteMsg(pa.data)
	}
}

func (p *Proxy) handleClientPackets() {
	for pa := range p.clientChan {
		packetSourceString := pa.w.RemoteAddr().String()

		p.RLock()
		conn, found := p.connsMap[packetSourceString]
		p.RUnlock()

		buf, _ := pa.data.Pack()

		if !found {
			c, err := net.Dial("udp", p.addr)
			if err != nil {
				return
			}
			conn := c.(*net.UDPConn)

			p.Lock()
			p.connsMap[packetSourceString] = connection{
				udp:          conn,
				w:            pa.w, // We're setting this once for this socket, is that safe?
				lastActivity: time.Now(),
			}
			p.Unlock()

			conn.Write(buf)
			go p.clientConnectionReadLoop(conn, pa.w)
		} else {
			conn.udp.Write(buf)
			p.RLock()
			shouldUpdateLastActivity := false
			if _, found := p.connsMap[packetSourceString]; found {
				if p.connsMap[packetSourceString].lastActivity.Before(
					time.Now().Add(-p.ConnTimeout / 4)) {
					shouldUpdateLastActivity = true
				}
			}
			p.RUnlock()
			if shouldUpdateLastActivity {
				p.updateClientLastActivity(packetSourceString)
			}
		}
	}
}

func (p *Proxy) freeIdleSocketsLoop() {
	for !p.closed {
		time.Sleep(p.ConnTimeout)
		var clientsToTimeout []string

		p.RLock()
		for client, conn := range p.connsMap {
			if conn.lastActivity.Before(time.Now().Add(-p.ConnTimeout)) {
				clientsToTimeout = append(clientsToTimeout, client)
			}
		}
		p.RUnlock()

		p.Lock()
		for _, client := range clientsToTimeout {
			p.connsMap[client].udp.Close()
			delete(p.connsMap, client)
		}
		p.Unlock()
	}
}
