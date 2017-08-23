// Package health implements an HTTP handler that responds to health checks.
package health

import (
	"io"
	"log"
	"net"
	"net/http"
	"sync"
)

var once sync.Once

type health struct {
	Addr string

	ln  net.Listener
	mux *http.ServeMux

	// A slice of Healthers that the health middleware will poll every second for their
	// health status.
	h []Healther
	sync.RWMutex
	ok bool // ok is the global boolean indicating all healthy middleware stack
}

// Ok returns the global health status of all middleware configured in this server.
func (h *health) Ok() bool {
	h.RLock()
	defer h.RUnlock()
	return h.ok
}

// SetOk sets the global health status of all middleware configured in this server.
func (h *health) SetOk(ok bool) {
	h.Lock()
	defer h.Unlock()
	h.ok = ok
}

func (h *health) Startup() error {
	if h.Addr == "" {
		h.Addr = defAddr
	}

	once.Do(func() {
		ln, err := net.Listen("tcp", h.Addr)
		if err != nil {
			log.Printf("[ERROR] Failed to start health handler: %s", err)
			return
		}

		h.ln = ln

		h.mux = http.NewServeMux()

		h.mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			if !h.Ok() {
				w.WriteHeader(http.StatusOK)
				io.WriteString(w, ok)
				return
			}
			w.WriteHeader(http.StatusServiceUnavailable)
		})

		go func() {
			http.Serve(h.ln, h.mux)
		}()
	})
	return nil
}

func (h *health) Shutdown() error {
	if h.ln != nil {
		return h.ln.Close()
	}
	return nil
}

// Poll polls all healthers and sets the global state.
func (h *health) Poll() {
	for _, m := range h.h {
		if !m.Health() {
			h.SetOk(false)
			return
		}
	}
	h.SetOk(true)
}

// Middleware that implements the Healther interface.
var healthers = map[string]bool{
	"erratic": true,
}

const (
	ok      = "OK"
	nok     = "NOT OK"
	defAddr = ":8080"
	path    = "/health"
)
