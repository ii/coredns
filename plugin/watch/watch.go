package watch

import (
	"context"
	"fmt"
	"io"

	"github.com/miekg/dns"

	"github.com/coredns/coredns/pb"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/watch"
)

// watch contains all the data needed to manage watches
type watcher struct {
	Next    plugin.Handler
	changes watch.WatchChan
	counter int64
	watches	map[string]watchlist
	watchees []watch.Watchee
}

type watchlist map[int64]*watchquery
type watchquery struct {
	query *dns.Msg
	stream pb.WatchService_WatchServer
}

// Newwatcher returns the watcher plugin
func NewWatcher() *watcher {
	w := &watcher{changes: make(watch.WatchChan), watches: make(map[string]watchlist)}
	go w.processWatches()
	return w
}

func (w *watcher) Name() string { return "watch" }

func (w *watcher) nextId() int64 {
	// TODO: should have a lock
	w.counter += 1
	return w.counter
}

// Watch is used to monitor the results of a given query. CoreDNS will push updated
// query responses down the stream.
func (w *watcher) Watch(stream pb.WatchService_WatchServer) error {
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

			id := w.nextId()

			qname := msg.Question[0].Name
			if _, ok := w.watches[qname]; !ok {
				w.watches[qname] = make(watchlist)
			}
			w.watches[qname][id] = &watchquery{query: msg, stream: stream}
			for wee := range w.watchees {
				w.watchees[wee].StartWatch(qname, w.changes)
			}
			if err := stream.Send(&pb.WatchResponse{WatchId: id, Created: true}); err != nil {
				return err
			}

			fmt.Printf("watches: %v\n", w.watches)
			continue
		}

		cancel := in.GetCancelRequest()
		if cancel != nil {
			//TODO: lock
			for qname, wl := range w.watches {
				ww, ok := wl[cancel.WatchId]
				if !ok {
					continue
				}
				if ww.stream != stream {
					continue
				}
				delete(wl, cancel.WatchId)
				if len(wl) == 0 {
					for wee := range w.watchees {
						w.watchees[wee].StopWatch(qname)
					}
					delete(w.watches, qname)
				}
				if err = stream.Send(&pb.WatchResponse{WatchId: cancel.WatchId, Canceled: true}); err != nil {
					return err
				}
			}
			fmt.Printf("watches: %v\n", w.watches)
			continue
		}
        }
}

func (w *watcher) processWatches() {
	for { 
		select {
		case changed := <-w.changes:
			fmt.Printf("A change: %v\n", changed)
		}
	}
}

func (w *watcher) ServeDNS(ctx context.Context, rw dns.ResponseWriter, r *dns.Msg) (int, error) {
	return plugin.NextOrFailure(w.Name(), w.Next, ctx, rw, r)
}

