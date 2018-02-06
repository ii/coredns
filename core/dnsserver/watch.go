package dnsserver

import (
	"io"
	"log"

	"github.com/miekg/dns"

	"github.com/coredns/coredns/pb"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/watch"
)

// watch contains all the data needed to manage watches
type watcher struct {
	changes  watch.WatchChan
	counter  int64
	watches  map[string]watchlist
	watchees []watch.Watchee
}

type watchlist map[int64]*watchquery
type watchquery struct {
	qtype uint32
	stream pb.DnsService_WatchServer
}

func newWatcher() *watcher {
	w := &watcher{changes: make(watch.WatchChan), watches: make(map[string]watchlist)}
	go w.processWatches()
	return w
}

func (w *watcher) Name() string { return "watch" }

func (w *watcher) nextID() int64 {
	// TODO: should have a lock
	w.counter += 1
	return w.counter
}

// watch is used to monitor the results of a given query. CoreDNS will push updated
// query responses down the stream.
func (w *watcher) watch(stream pb.DnsService_WatchServer) error {
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

			id := w.nextID()

			qname := msg.Question[0].Name
			if _, ok := w.watches[qname]; !ok {
				w.watches[qname] = make(watchlist)
			}
			w.watches[qname][id] = &watchquery{stream: stream, qtype: uint32(msg.Question[0].Qtype)}
			for wee := range w.watchees {
				w.watchees[wee].StartWatch(qname, w.changes)
			}
			if err := stream.Send(&pb.WatchResponse{WatchId: id, Created: true}); err != nil {
				return err
			}

			log.Printf("watches: %v\n", w.watches)
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
			log.Printf("watches: %v\n", w.watches)
			continue
		}
	}
}

func (w *watcher) processWatches() {
	for {
		select {
		case changed := <-w.changes:
			log.Printf("Change: %v, checking watches in %v\n", changed, w.watches)
			for qname, wl := range w.watches {
				log.Printf("Checking %s against %s\n", changed, qname)
				if plugin.Zones(changed).Matches(qname) == "" {
					continue
				}
				log.Printf("Matches %s\n", qname)
				for id, wq := range wl {
					wr := pb.WatchResponse{WatchId: id, Qname: qname, Qtype: wq.qtype}
					log.Printf("Sending %v over %v\n", wr, wq.stream)
					err := wq.stream.Send(&wr)
					log.Printf("Sent, err = %s", err)
					if err != nil {
						log.Printf("Error sending to watch %d: %s\n", id, err)
					}
				}
			}
		}
	}
}
