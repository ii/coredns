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

package fastproxy

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type connection struct {
	udp          *net.UDPConn
	lastActivity time.Time
}

type packet struct {
	src  *net.UDPAddr
	data []byte
}

type Proxy struct {
	BindPort               int
	BindAddress            string
	UpstreamAddress        string
	UpstreamPort           int
	listenerConn           *net.UDPConn
	client                 *net.UDPAddr
	upstream               *net.UDPAddr
	BufferSize             int
	ConnTimeout            time.Duration
	ResolveTTL             time.Duration
	connsMap               map[string]connection
	connectionsLock        *sync.RWMutex
	closed                 bool
	clientMessageChannel   chan (packet)
	upstreamMessageChannel chan (packet)
}

func New(bindPort int, bindAddress string, upstreamAddress string, upstreamPort int, bufferSize int, connTimeout time.Duration, resolveTTL time.Duration) *Proxy {
	proxy := &Proxy{
		BindPort:               bindPort,
		BindAddress:            bindAddress,
		BufferSize:             bufferSize,
		ConnTimeout:            connTimeout,
		UpstreamAddress:        upstreamAddress,
		UpstreamPort:           upstreamPort,
		connectionsLock:        new(sync.RWMutex),
		connsMap:               make(map[string]connection),
		closed:                 false,
		ResolveTTL:             resolveTTL,
		clientMessageChannel:   make(chan packet),
		upstreamMessageChannel: make(chan packet),
	}

	return proxy
}

func (p *Proxy) updateClientLastActivity(clientAddrString string) {
	p.connectionsLock.Lock()
	if _, found := p.connsMap[clientAddrString]; found {
		connWrapper := p.connsMap[clientAddrString]
		connWrapper.lastActivity = time.Now()
		p.connsMap[clientAddrString] = connWrapper
	}
	p.connectionsLock.Unlock()
}

func (p *Proxy) clientConnectionReadLoop(clientAddr *net.UDPAddr, upstreamConn *net.UDPConn) {
	clientAddrString := clientAddr.String()
	for {
		buffer := make([]byte, p.BufferSize)
		size, _, err := upstreamConn.ReadFromUDP(buffer)
		if err != nil {
			p.connectionsLock.Lock()
			upstreamConn.Close()
			delete(p.connsMap, clientAddrString)
			p.connectionsLock.Unlock()
			return
		}
		p.updateClientLastActivity(clientAddrString)
		p.upstreamMessageChannel <- packet{
			src:  clientAddr,
			data: buffer[:size],
		}
	}
}

func (p *Proxy) handlerUpstreamPackets() {
	for pa := range p.upstreamMessageChannel {
		p.listenerConn.WriteTo(pa.data, pa.src)
	}
}

func (p *Proxy) handleClientPackets() {
	for pa := range p.clientMessageChannel {
		packetSourceString := pa.src.String()

		p.connectionsLock.RLock()
		conn, found := p.connsMap[packetSourceString]
		p.connectionsLock.RUnlock()

		if !found {
			conn, err := net.DialUDP("udp", p.client, p.upstream)
			if err != nil {
				return
			}

			p.connectionsLock.Lock()
			p.connsMap[packetSourceString] = connection{
				udp:          conn,
				lastActivity: time.Now(),
			}
			p.connectionsLock.Unlock()

			conn.Write(pa.data)
			go p.clientConnectionReadLoop(pa.src, conn)
		} else {
			conn.udp.Write(pa.data)
			p.connectionsLock.RLock()
			shouldUpdateLastActivity := false
			if _, found := p.connsMap[packetSourceString]; found {
				if p.connsMap[packetSourceString].lastActivity.Before(
					time.Now().Add(-p.ConnTimeout / 4)) {
					shouldUpdateLastActivity = true
				}
			}
			p.connectionsLock.RUnlock()
			if shouldUpdateLastActivity {
				p.updateClientLastActivity(packetSourceString)
			}
		}
	}
}

func (p *Proxy) readLoop() {
	for !p.closed {
		buffer := make([]byte, p.BufferSize)
		size, srcAddress, err := p.listenerConn.ReadFromUDP(buffer)
		if err != nil {
			continue
		}
		p.clientMessageChannel <- packet{
			src:  srcAddress,
			data: buffer[:size],
		}
	}
}

func (p *Proxy) resolveUpstreamLoop() {
	for !p.closed {
		time.Sleep(p.ResolveTTL)
		upstreamAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", p.UpstreamAddress, p.UpstreamPort))
		if err != nil {
			continue
		}
		if p.upstream.String() != upstreamAddr.String() {
			p.upstream = upstreamAddr
		}
	}
}

func (p *Proxy) freeIdleSocketsLoop() {
	for !p.closed {
		time.Sleep(p.ConnTimeout)
		var clientsToTimeout []string

		p.connectionsLock.RLock()
		for client, conn := range p.connsMap {
			if conn.lastActivity.Before(time.Now().Add(-p.ConnTimeout)) {
				clientsToTimeout = append(clientsToTimeout, client)
			}
		}
		p.connectionsLock.RUnlock()

		p.connectionsLock.Lock()
		for _, client := range clientsToTimeout {
			p.connsMap[client].udp.Close()
			delete(p.connsMap, client)
		}
		p.connectionsLock.Unlock()
	}
}

func (p *Proxy) Close() {
	p.connectionsLock.Lock()
	p.closed = true
	for _, conn := range p.connsMap {
		conn.udp.Close()
	}
	if p.listenerConn != nil {
		p.listenerConn.Close()
	}
	p.connectionsLock.Unlock()
}

func (p *Proxy) Start() {
	ProxyAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", p.BindAddress, p.BindPort))
	if err != nil {
		return
	}
	p.upstream, err = net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", p.UpstreamAddress, p.UpstreamPort))
	p.client = &net.UDPAddr{
		IP:   ProxyAddr.IP,
		Port: 0,
		Zone: ProxyAddr.Zone,
	}
	p.listenerConn, err = net.ListenUDP("udp", ProxyAddr)
	if err != nil {
		return
	}
	if p.ConnTimeout.Nanoseconds() > 0 {
		go p.freeIdleSocketsLoop()
	}
	if p.ResolveTTL.Nanoseconds() > 0 {
		go p.resolveUpstreamLoop()
	}
	go p.handlerUpstreamPackets()
	go p.handleClientPackets()
	go p.readLoop()
}
