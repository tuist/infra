package scheme

import (
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"

	infrav1 "github.com/tuist/infra/api/v1"
)

// Scheme is the central registry the controller manager and clients use to map
// YAML/JSON Kubernetes objects to Go structs and back.
var Scheme = runtime.NewScheme()

func init() {
	// Register built-in Kubernetes types first, then Infra's custom types.
	utilruntime.Must(clientgoscheme.AddToScheme(Scheme))
	utilruntime.Must(infrav1.AddToScheme(Scheme))
}

// AddToScheme copies the same registrations into another runtime.Scheme when a
// caller wants to build its own scheme instead of reusing the package-global one.
func AddToScheme(target *runtime.Scheme) error {
	if err := clientgoscheme.AddToScheme(target); err != nil {
		return err
	}

	return infrav1.AddToScheme(target)
}
