package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/sample-crd/pkg/apis/samplecrd"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime"
)

var SchemeGroupVersion = schema.GroupVersion{Group: samplecrd.GroupName, Version: "v1"}

func Kind(kind string) schema.GroupKind {
	return SchemeGroupVersion.WithKind(kind).GroupKind()
}

func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&SampleCRD{},
		&SampleCRDList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}