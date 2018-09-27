package object

import (
	api "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
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

func (s *Service) GetNamespace() string          { return s.Namespace }
func (s *Service) SetNamespace(namespace string) {}
func (s *Service) GetName() string               { return s.Name }
func (s *Service) SetName(name string)           {}

type Object interface {
	GetNamespace() string
	SetNamespace(namespace string)
	GetName() string
	SetName(name string)
	GetGenerateName() string
	SetGenerateName(name string)
	GetUID() types.UID
	SetUID(uid types.UID)
	GetResourceVersion() string
	SetResourceVersion(version string)
	GetGeneration() int64
	SetGeneration(generation int64)
	GetSelfLink() string
	SetSelfLink(selfLink string)
	GetCreationTimestamp() v1.Time
	SetCreationTimestamp(timestamp v1.Time)
	GetDeletionTimestamp() *v1.Time
	SetDeletionTimestamp(timestamp *v1.Time)
	GetDeletionGracePeriodSeconds() *int64
	SetDeletionGracePeriodSeconds(*int64)
	GetLabels() map[string]string
	SetLabels(labels map[string]string)
	GetAnnotations() map[string]string
	SetAnnotations(annotations map[string]string)
	GetInitializers() *v1.Initializers
	SetInitializers(initializers *v1.Initializers)
	GetFinalizers() []string
	SetFinalizers(finalizers []string)
	GetOwnerReferences() []v1.OwnerReference
	SetOwnerReferences([]v1.OwnerReference)
	GetClusterName() string
	SetClusterName(clusterName string)
}
