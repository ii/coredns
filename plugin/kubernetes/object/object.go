package object

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

// Empty is an empty struct.
type Empty struct{}

// GetObjectKind implements the ObjectKind interface as a noop.
func (e *Empty) GetObjectKind() schema.ObjectKind { return schema.EmptyObjectKind }

// ToFunc converts one empty interface to another.
type ToFunc func(interface{}) interface{}

func (e *Empty) GetGenerateName() string                       { return "" }
func (e *Empty) SetGenerateName(name string)                   {}
func (e *Empty) GetUID() types.UID                             { return "" }
func (e *Empty) SetUID(uid types.UID)                          {}
func (e *Empty) GetResourceVersion() string                    { return "" }
func (e *Empty) SetResourceVersion(version string)             {}
func (e *Empty) GetGeneration() int64                          { return 0 }
func (e *Empty) SetGeneration(generation int64)                {}
func (e *Empty) GetSelfLink() string                           { return "" }
func (e *Empty) SetSelfLink(selfLink string)                   {}
func (e *Empty) GetCreationTimestamp() v1.Time                 { return v1.Time{} }
func (e *Empty) SetCreationTimestamp(timestamp v1.Time)        {}
func (e *Empty) GetDeletionTimestamp() *v1.Time                { return &v1.Time{} }
func (e *Empty) SetDeletionTimestamp(timestamp *v1.Time)       {}
func (e *Empty) GetDeletionGracePeriodSeconds() *int64         { return nil }
func (e *Empty) SetDeletionGracePeriodSeconds(*int64)          {}
func (e *Empty) GetLabels() map[string]string                  { return nil }
func (e *Empty) SetLabels(labels map[string]string)            {}
func (e *Empty) GetAnnotations() map[string]string             { return nil }
func (e *Empty) SetAnnotations(annotations map[string]string)  {}
func (e *Empty) GetInitializers() *v1.Initializers             { return nil }
func (e *Empty) SetInitializers(initializers *v1.Initializers) {}
func (e *Empty) GetFinalizers() []string                       { return nil }
func (e *Empty) SetFinalizers(finalizers []string)             {}
func (e *Empty) GetOwnerReferences() []v1.OwnerReference       { return nil }
func (e *Empty) SetOwnerReferences([]v1.OwnerReference)        {}
func (e *Empty) GetClusterName() string                        { return "" }
func (e *Empty) SetClusterName(clusterName string)             {}
