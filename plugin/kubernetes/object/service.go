package object

import (
	api "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// Service is a stripped down api.Service with only the items we need for CoreDNS.
type Service struct {
	Name         string
	Namespace    string
	Index        string
	ClusterIP    string
	Type         api.ServiceType
	ExternalName string
	Ports        []api.ServicePort

	*Empty
}

// ToService converts an api.Service to a *Service.
func ToService(obj interface{}) interface{} {
	svc, ok := obj.(*api.Service)
	if !ok {
		return nil
	}

	s := &Service{
		Name:         svc.GetName(),
		Namespace:    svc.GetNamespace(),
		Index:        svc.GetName() + "." + svc.GetNamespace(),
		ClusterIP:    svc.Spec.ClusterIP,
		Type:         svc.Spec.Type,
		ExternalName: svc.Spec.ExternalName,
	}
	copy(s.Ports, svc.Spec.Ports)

	return s
}

var _ runtime.Object = &Service{}

// DeepCopyObject implements the runtime.Object interface.
func (s *Service) DeepCopyObject() runtime.Object {
	s1 := &Service{
		Name:         s.GetName(),
		Namespace:    s.GetNamespace(),
		Index:        s.GetName() + "." + s.GetNamespace(),
		ClusterIP:    s.ClusterIP,
		Type:         s.Type,
		ExternalName: s.ExternalName,
	}
	copy(s1.Ports, s.Ports)
	return s1
}

// GetNamespace implements the metav1.Object interface.
func (s *Service) GetNamespace() string { return s.Namespace }

// SetNamespace implements the metav1.Object interface.
func (s *Service) SetNamespace(namespace string) {}

// GetName implements the metav1.Object interface.
func (s *Service) GetName() string { return s.Name }

// SetName implements the metav1.Object interface.
func (s *Service) SetName(name string) {}
