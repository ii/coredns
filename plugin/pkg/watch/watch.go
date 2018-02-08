package watch

import (
	"github.com/coredns/coredns/pb"
)

// Watcher is the interface the Watch plugin implements
type Watcher interface {
	Watch(stream pb.DnsService_WatchServer) error
}

// Watchee is the interface watchable plugins should implement
type Watchee interface {
	StartWatch(qname string, changes NotifyChan) error
	StopWatch(qname string) error
}

// NotifyChan is used by plugins to inform the server when records have changed
type NotifyChan chan []string
