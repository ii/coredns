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
		Index:        svc.GetName() + "." + svc.GetNamespace(), // Used as index key
		ClusterIP:    svc.Spec.ClusterIP,
		Type:         svc.Spec.Type,
		ExternalName: svc.Spec.ExternalName,
		Ports:        svc.Spec.Ports,
	}

	return s
}

var _ runtime.Object = &Service{}

// DeepCopyObject implements the runtime.Object interface.
func (s *Service) DeepCopyObject() runtime.Object {
	return nil
}
