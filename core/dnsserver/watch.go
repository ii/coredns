package dnsserver

import (
	"io"
	"log"
	"sync"

	"github.com/miekg/dns"

	"github.com/coredns/coredns/pb"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/watch"
)

// watcher contains all the data needed to manage watches
type watcher struct {
	changes watch.Chan
	counter int64
	watches map[string]watchlist
	plugins []watch.Watchable
	mutex   sync.Mutex
}

type watchlist map[int64]pb.DnsService_WatchServer

func newWatcher(zones map[string]*Config) *watcher {
	w := &watcher{changes: make(watch.Chan), watches: make(map[string]watchlist)}

	for _, config := range zones {
		plugins := config.Handlers()
		for _, p := range plugins {
			if x, ok := p.(watch.Watchable); ok {
				x.SetWatchChan(w.changes)
				w.plugins = append(w.plugins, x)
			}
		}
	}

	//TODO: maybe a stop channel, work properly with reloads?
	go w.processWatches()
	return w
}

func (w *watcher) nextID() int64 {
	w.mutex.Lock()

	w.counter++
	id := w.counter

	w.mutex.Unlock()
	return id
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

			if err := stream.Send(&pb.WatchResponse{WatchId: id, Created: true}); err != nil {
				return err
			}

			qname := msg.Question[0].Name
			if _, ok := w.watches[qname]; !ok {
				w.watches[qname] = make(watchlist)
			}

			w.watches[qname][id] = stream

			for _, p := range w.plugins {
				err := p.Watch(qname)
				if err != nil {
					log.Printf("[WARNING] Failed to start watch for %s in plugin %s: %s\n", qname, p.Name(), err)
				}
			}

			log.Printf("watches: %v\n", w.watches)
			continue
		}

		cancel := in.GetCancelRequest()
		if cancel != nil {
			//TODO: lock
			for qname, wl := range w.watches {
				ws, ok := wl[cancel.WatchId]
				if !ok {
					continue
				}

				// only allow cancels from the client that started it
				if ws != stream {
					continue
				}

				delete(wl, cancel.WatchId)

				// if there are no more watches for this qname, we should tell the plugins
				if len(wl) == 0 {
					for _, p := range w.plugins {
						p.StopWatching(qname)
					}
					delete(w.watches, qname)
				}

				// let the client know we canceled the watch
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
				for id, stream := range wl {
					wr := pb.WatchResponse{WatchId: id, Qname: qname}
					log.Printf("Sending %v over %v\n", wr, stream)
					err := stream.Send(&wr)
					log.Printf("Sent, err = %s", err)
					if err != nil {
						log.Printf("Error sending to watch %d: %s\n", id, err)
					}
				}
			}
		}
	}
}
