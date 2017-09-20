package dnstest

import (
	"net"

	"github.com/miekg/dns"
)

// A Server is an DNS server listening on a system-chosen port on the local
// loopback interface, for use in end-to-end DNS tests.
type Server struct {
	Addr string // Address where the server listening.

	s *dns.Server
}

// NewServer starts and returns a new Server. The caller should call Close when
// finished, to shut it down.
func NewServer(f dns.HandlerFunc) *Server {
	dns.HandleFunc(".", f)

	p, _ := net.ListenPacket("udp", ":0")
	l, _ := net.Listen("tcp", p.LocalAddr().String())

	s := &dns.Server{PacketConn: p, Listener: l}
	go s.ActivateAndServe()

	return &Server{s: s, Addr: p.LocalAddr().String()}
}

// Close shuts down the server.
func (s *Server) Close() { s.s.Shutdown() }
