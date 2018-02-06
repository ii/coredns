package kubernetes

import (
	"fmt"
	"github.com/coredns/coredns/plugin/pkg/watch"
)

func (k *Kubernetes) StartWatch(qname string, changes watch.WatchChan) error {
	fmt.Printf("starting watch for %s in k8s, %v", qname, changes)
	k.watchChan = changes
	return nil
}

func (k *Kubernetes) StopWatch(qname string) error { return nil }
