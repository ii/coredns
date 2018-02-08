package kubernetes

import (
	"fmt"
	"github.com/coredns/coredns/plugin/pkg/watch"
)

// StartWatch is called when a watch is started for a name.
func (k *Kubernetes) StartWatch(qname string, changes watch.NotifyChan) error {
	fmt.Printf("starting watch for %s in k8s, %v", qname, changes)
	k.watchers = changes
	return nil
}

// StopWatch is called when a watch is stopped for a name.
func (k *Kubernetes) StopWatch(qname string) error { return nil }
