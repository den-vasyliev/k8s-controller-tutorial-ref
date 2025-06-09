// +kubebuilder:object:generate=true
// +groupName=frontendpage.alex0m.io

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	// SchemeGroupVersion is group version used to register these objects
	SchemeGroupVersion = schema.GroupVersion{Group: "frontendpage.alex0m.io", Version: "v1alpha1"}
)

// AddToScheme registers FrontendPage and FrontendPageList with the given scheme.
func AddToScheme(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion, &FrontendPage{}, &FrontendPageList{})
	return nil
}
