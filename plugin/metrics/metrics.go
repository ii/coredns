// Package metrics implement a handler and plugin that provides Prometheus metrics.
package metrics

import (
	"log"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/metrics/vars"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics holds the prometheus configuration. The metrics' path is fixed to be /metrics
type Metrics struct {
	Next plugin.Handler
	Addr string
	Reg  *prometheus.Registry
	ln   net.Listener
	mux  *http.ServeMux

	zoneNames []string
	zoneMap   map[string]bool
	zoneMu    sync.RWMutex
}

// New returns a new instance of Metrics with the given address
func New(addr string) *Metrics {
	met := &Metrics{
		Addr:    addr,
		Reg:     prometheus.NewRegistry(),
		zoneMap: make(map[string]bool),
	}
	// Add the default collectors
	met.Reg.MustRegister(prometheus.NewGoCollector())
	met.Reg.MustRegister(prometheus.NewProcessCollector(os.Getpid(), ""))

	// Add all of our collectors
	met.Reg.MustRegister(vars.RequestCount)
	met.Reg.MustRegister(vars.RequestDuration)
	met.Reg.MustRegister(vars.RequestSize)
	met.Reg.MustRegister(vars.RequestDo)
	met.Reg.MustRegister(vars.RequestType)
	met.Reg.MustRegister(vars.ResponseSize)
	met.Reg.MustRegister(vars.ResponseRcode)
	return met
}

// AddZone adds zone z to m.
func (m *Metrics) AddZone(z string) {
	m.zoneMu.Lock()
	m.zoneMap[z] = true
	m.zoneNames = keys(m.zoneMap)
	m.zoneMu.Unlock()
}

// RemoveZone remove zone z from m.
func (m *Metrics) RemoveZone(z string) {
	m.zoneMu.Lock()
	delete(m.zoneMap, z)
	m.zoneNames = keys(m.zoneMap)
	m.zoneMu.Unlock()
}

// ZoneNames returns the zones of m.
func (m *Metrics) ZoneNames() []string {
	m.zoneMu.RLock()
	s := m.zoneNames
	m.zoneMu.RUnlock()
	return s
}

// OnStartup sets up the metrics on startup.
func (m *Metrics) OnStartup() error {
	ln, err := net.Listen("tcp", m.Addr)
	if err != nil {
		log.Printf("[ERROR] Failed to start metrics handler: %s", err)
		return err
	}

	m.ln = ln
	m.mux = http.NewServeMux()
	m.mux.Handle("/metrics", promhttp.HandlerFor(m.Reg, promhttp.HandlerOpts{}))

	go func() {
		http.Serve(m.ln, m.mux)
	}()
	return nil
}

// OnShutdown tears down the metrics on shutdown and restart.
func (m *Metrics) OnShutdown() error {
	if m.ln != nil {
		return m.ln.Close()
	}
	return nil
}

func keys(m map[string]bool) []string {
	sx := []string{}
	for k := range m {
		sx = append(sx, k)
	}
	return sx
}
