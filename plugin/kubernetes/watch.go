package kubernetes

import (
	"fmt"
	"github.com/coredns/coredns/plugin/pkg/watch"
)

func (k *Kubernetes) StartWatch(qname string, changes watch.WatchChan) {
	fmt.Printf("starting watch for %s in k8s, %v", qname, changes)
	k.watchChan = changes
}

func (k *Kubernetes) StopWatch(qname string) {}
