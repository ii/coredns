package object

import (
	api "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// Endpoints is a stripped down api.Endpoints with only the items we need for CoreDNS.
type Endpoints struct {
	Name      string
	Namespace string

	Addresses []string
	Ports     []int32

	*Empty
}

// ToEndpint converts an api.Service to a *Service.
func ToEndpoints(obj interface{}) interface{} {
	end, ok := obj.(*api.Endpoints)
	if !ok {
		return nil
	}

	e := &Endpoints{
		Name:      end.GetName(),
		Namespace: end.GetNamespace(),
	}
	for _, eps := range end.Subsets {
		for _, a := range eps.Addresses {
			e.Addresses = append(e.Addresses, a.IP)
		}
		for _, p := range eps.Ports {
			e.Ports = append(e.Ports, p.Port)
		}
	}

	return e
}

var _ runtime.Object = &Endpoints{}

// GetNamespace implements the metav1.Object interface.
func (e *Endpoints) GetNamespace() string { return e.Namespace }

// SetNamespace implements the metav1.Object interface.
func (e *Endpoints) SetNamespace(namespace string) {}

// GetName implements the metav1.Object interface.
func (e *Endpoints) GetName() string { return e.Name }

// SetName implements the metav1.Object interface.
func (e *Endpoints) SetName(name string) {}
