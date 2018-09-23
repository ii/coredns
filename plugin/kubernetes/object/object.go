package object

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Empty is an empty struct.
type Empty struct{}

// GetObjectKind implements the ObjectKind interface as a noop.
func (e *Empty) GetObjectKind() schema.ObjectKind { return schema.EmptyObjectKind }

// ToFunc converts one empty interface to another.
type ToFunc func(interface{}) interface{}
