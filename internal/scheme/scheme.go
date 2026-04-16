package scheme

import (
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"

	kubeboxv1alpha1 "github.com/tuist/kubebox/api/v1alpha1"
)

var Scheme = runtime.NewScheme()

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(Scheme))
	utilruntime.Must(kubeboxv1alpha1.AddToScheme(Scheme))
}

func AddToScheme(target *runtime.Scheme) error {
	if err := clientgoscheme.AddToScheme(target); err != nil {
		return err
	}

	return kubeboxv1alpha1.AddToScheme(target)
}
