package liveness

import (
	"sync"
	"time"
)

// Probe is used to run a single UpFunc until it returns true (indicating a target is healthy). If an UpFunc
// is already in progress no new one will be added, i.e. there is always a maximum of 1 checks in flight.
type Probe struct {
	do   chan UpFunc
	stop chan bool

	interval time.Duration
	target   string

	sync.Mutex
	inprogress bool
}

// UpFunc is used to determine if a target is alive. If so this function must return true.
type UpFunc func(target string) bool

// New returns a pointer to an intialized Probe.
func New(interval time.Duration) *Probe {
	return &Probe{stop: make(chan bool), do: make(chan UpFunc), interval: interval}
}

// Do will probe target, if a probe is already in progress this is a noop.
func (p *Probe) Do(f UpFunc) { p.do <- f }

// Stop stops the probing.
func (p *Probe) Stop() { p.stop <- true }

// Start will start the probe manager, after which probes can be initialized with Do.
func (p *Probe) Start(target string) { go p.start(target) }

func (p *Probe) start(target string) {
	for {
		select {
		case <-p.stop:
			return
		case f := <-p.do:
			p.Lock()
			if p.inprogress {
				p.Unlock()
				continue
			}
			p.inprogress = true
			p.Unlock()

			// Passed the lock. Now run f for as long it returns false. If a true is returned
			// we return from the goroutine and we can accept another UpFunc to run.
			go func() {
				for {
					if ok := f(target); ok {
						break
					}
					time.Sleep(p.interval)
				}
				p.Lock()
				p.inprogress = false
				p.Unlock()
			}()
		}
	}
}
