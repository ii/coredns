package object

import (
	"time"

	api "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// Pod is a stripped down api.Pod with only the items we need for CoreDNS.
type Pod struct {
	PodIP             string
	DeletionTimeStamp time.Time
	*Empty
}

// ToPod converts an api.Pod to a *Pod.
func ToPod(pod *api.Pod) *Pod {
	p := &Pod{}
	p.PodIP = pod.Status.PodIP

	t := pod.ObjectMeta.DeletionTimestamp
	if t != nil {
		p.DeletionTimeStamp = (*t).Time
	}

	return p
}

var _ runtime.Object = &Pod{}

// DeepCopyObject implements the runtime.Object interface.
func (p *Pod) DeepCopyObject() runtime.Object {
	return &Pod{PodIP: p.PodIP, DeletionTimeStamp: p.DeletionTimeStamp}
}
