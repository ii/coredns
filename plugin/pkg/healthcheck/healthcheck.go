package healthcheck

import (
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

// UpstreamHostDownFunc can be used to customize how Down behaves.
type UpstreamHostDownFunc func(*UpstreamHost) bool

// UpstreamHost represents a single proxy upstream
type UpstreamHost struct {
	Conns       int64  // must be first field to be 64-bit aligned on 32-bit systems
	Name        string // IP address (and port) of this upstream host
	Fails       int32
	FailTimeout time.Duration
	OkUntil     time.Time
	CheckDown   UpstreamHostDownFunc
	Checking    bool
	CheckMu     sync.Mutex
}

// Down checks whether the upstream host is down or not. Down will try to use
// uh.CheckDown first, and will fall back to some default criteria if
// necessary.
func (uh *UpstreamHost) Down() bool {
	if uh.CheckDown == nil {
		// Default settings
		fails := atomic.LoadInt32(&uh.Fails)
		after := false

		uh.CheckMu.Lock()
		until := uh.OkUntil
		uh.CheckMu.Unlock()

		if !until.IsZero() && time.Now().After(until) {
			after = true
		}

		return after || fails > 0
	}
	return uh.CheckDown(uh)
}

// HostPool is a collection of UpstreamHosts.
type HostPool []*UpstreamHost

// HealthCheck is used for performing healthcheck
// on a collection of upstream hosts and select
// one based on the policy.
type HealthCheck struct {
	wg   sync.WaitGroup // Used to wait for running goroutines to stop.
	stop chan struct{}  // Signals running goroutines to stop.

	Hosts       HostPool
	Policy      Policy
	Spray       Policy
	FailTimeout time.Duration
	MaxFails    int32
	Future      time.Duration
	Interval    time.Duration
}

// Start starts the healthcheck
func (u *HealthCheck) Start() {
	u.stop = make(chan struct{})
	u.wg.Add(1)
	go func() {
		defer u.wg.Done()
		u.healthCheckWorker(u.stop)
	}()
}

// Stop sends a signal to all goroutines started by this staticUpstream to exit
// and waits for them to finish before returning.
func (u *HealthCheck) Stop() error {
	close(u.stop)
	u.wg.Wait()
	return nil
}

// This was moved into a thread so that each host could throw a health
// check at the same time.  The reason for this is that if we are checking
// 3 hosts, and the first one is gone, and we spend minutes timing out to
// fail it, we would not have been doing any other health checks in that
// time.  So we now have a per-host lock and a threaded health check.
//
// We use the Checking bool to avoid concurrent checks against the same
// host; if one is taking a long time, the next one will find a check in
// progress and simply return before trying.
//
// We are carefully avoiding having the mutex locked while we check,
// otherwise checks will back up, potentially a lot of them if a host is
// absent for a long time.  This arrangement makes checks quickly see if
// they are the only one running and abort otherwise.
func Check(f Func, okUntil time.Time, host *UpstreamHost) {
	// lock for our bool check.  We don't just defer the unlock because
	// we don't want the lock held while the check runs.
	host.CheckMu.Lock()

	// Are we mid check?  Don't run another one.
	if host.Checking {
		host.CheckMu.Unlock()
		return
	}
	host.Checking = true

	host.CheckMu.Unlock()

	// Exchange the payload package. This has been moved into a go func
	// because when the remote host is not merely not serving, but actually
	// absent, then tcp syn timeouts can be very long, and so one fetch
	// could last several check intervals
	if r, err := f(host.Name); err != nil {
		log.Printf("[WARNING] Host %s health check probe failed: %v\n", host.Name, err)
		okUntil = time.Unix(0, 0)
	} else {
		if r.Rcode == dns.RcodeSuccess {
			// add this??
		}
		atomic.StoreInt32(&host.Fails, 0) // reset Fails as well
	}

	host.CheckMu.Lock()
	host.Checking = false
	host.OkUntil = okUntil
	host.CheckMu.Unlock()
}

func (u *HealthCheck) healthCheck() {
	for _, host := range u.Hosts {

		// calculate this before the get
		okUntil := time.Now().Add(u.Future)

		// locks/bools should prevent requests backing up
		go Check(nil, okUntil, host)
	}
}

func (u *HealthCheck) healthCheckWorker(stop chan struct{}) {
	ticker := time.NewTicker(u.Interval)
	u.healthCheck()
	for {
		select {
		case <-ticker.C:
			u.healthCheck()
		case <-stop:
			ticker.Stop()
			return
		}
	}
}

// Select selects an upstream host based on the policy
// and the healthcheck result.
func (u *HealthCheck) Select() *UpstreamHost {
	pool := u.Hosts
	if len(pool) == 1 {
		if pool[0].Down() && u.Spray == nil {
			return nil
		}
		return pool[0]
	}
	allDown := true
	for _, host := range pool {
		if !host.Down() {
			allDown = false
			break
		}
	}
	if allDown {
		if u.Spray == nil {
			return nil
		}
		return u.Spray.Select(pool)
	}

	if u.Policy == nil {
		h := (&Random{}).Select(pool)
		if h != nil {
			return h
		}
		if h == nil && u.Spray == nil {
			return nil
		}
		return u.Spray.Select(pool)
	}

	h := u.Policy.Select(pool)
	if h != nil {
		return h
	}

	if u.Spray == nil {
		return nil
	}
	return u.Spray.Select(pool)
}

// payload is the default payload we send to the upstream. This is using tcp for the checking.
var payload = func() request.Request {
	m := new(dns.Msg)
	m.SetQuestion(".", dns.TypeNS)
	m.RecursionDesired = false

	return request.Request{Req: m, W: &responseWriter{}}
}()

// Func is the function that gets called to perform healthchecks, it uses
// payload is the default payload.
type Func func(addr string) (*dns.Msg, error)

// responseWriter is a health specific one that defaults to using a TCP transport.
type responseWriter struct{ dns.ResponseWriter }

func (r *responseWriter) LocalAddr() net.Addr { return &net.TCPAddr{IP: ip, Port: 53, Zone: ""} }

var ip = net.ParseIP("127.0.0.1")
