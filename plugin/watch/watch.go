package watch

import (
	"fmt"
	"io"

	"github.com/miekg/dns"

	"github.com/coredns/coredns/pb"
)

// Watcher is the interface watchable plugins should implement
// TODO: would Watchable be idiomatic? or Watchee?
type Watcher interface {
	StartWatch(qname string, changes WatchChan) error
	StopWatch(qname string) error
}

type WatchChan chan []string

// watch contains all the data needed to manage watches
type watch struct {
	changes WatchChan
	counter int64
	watches	map[string]watchlist
	watchers []Watcher
}

type watchlist map[int64]*watch
type watch struct {
	query *dns.Msg
	stream pb.DnsService_WatchServer
}

// NewWatch returns a structure for managing watches
func NewWatch() *watch {
	return &watch{changes; make(WatchChan), watches: make(map[string]watchlist)}
}

func (wc *watch) nextId() int64 {
	// TODO: should have a lock
	wc.counter += 1
	return wc.counter
}

// Watch is used to monitor the results of a given query. CoreDNS will push updated
// query responses down the stream.
func (wc *watch) Watch(stream pb.DnsService_WatchServer) error {
        for {
                in, err := stream.Recv()
                if err == io.EOF {
                        return nil
                }
                if err != nil {
                        return err
                }
		create := in.GetCreateRequest()
		if create != nil {
			msg := new(dns.Msg)
			err := msg.Unpack(create.Query.Msg)
			if err != nil {
				// TODO: should write back an error response not break the stream
				return err
			}

			id := wc.nextId()

			qname := msg.Question[0].Name
			if _, ok := wc.watches[qname]; !ok {
				wc.watches[qname] = make(watchlist)
			}
			wc.watches[msg.Question[0].Name][id] = &watch{query: msg, stream: stream}
			
			if err := stream.Send(&pb.WatchResponse{WatchId: id, Created: true}); err != nil {
				return err
			}

			fmt.Printf("watches: %v\n", wc.watches)
			continue
		}

		cancel := in.GetCancelRequest()
		if cancel != nil {
			//TODO: lock
			for qname, wl := range wc.watches {
				w, ok := wl[cancel.WatchId]
				if !ok {
					continue
				}
				if w.stream != stream {
					continue
				}
				delete(wl, cancel.WatchId)
				if len(wl) == 0 {
					delete(wc.watches, qname)
				}
				if err = stream.Send(&pb.WatchResponse{WatchId: cancel.WatchId, Canceled: true}); err != nil {
					return err
				}
			}
			fmt.Printf("watches: %v\n", wc.watches)
			continue
		}
        }
}
