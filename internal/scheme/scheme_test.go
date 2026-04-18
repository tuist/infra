package scheme

import (
	"testing"

	kubeboxv1 "github.com/tuist/kubebox/api/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestSchemeIncludesBuiltInAndKubeboxTypes(t *testing.T) {
	t.Parallel()

	assertSchemeCreates[*kubeboxv1.HostPool](t, Scheme, kubeboxv1.SchemeGroupVersion.WithKind("HostPool"))
	assertSchemeCreates[*corev1.Pod](t, Scheme, corev1.SchemeGroupVersion.WithKind("Pod"))
}

func TestAddToSchemeRegistersBuiltInAndKubeboxTypes(t *testing.T) {
	t.Parallel()

	target := runtime.NewScheme()
	if err := AddToScheme(target); err != nil {
		t.Fatalf("AddToScheme() returned error: %v", err)
	}

	assertSchemeCreates[*kubeboxv1.Sandbox](t, target, kubeboxv1.SchemeGroupVersion.WithKind("Sandbox"))
	assertSchemeCreates[*corev1.Service](t, target, corev1.SchemeGroupVersion.WithKind("Service"))
}

func assertSchemeCreates[T runtime.Object](t *testing.T, scheme *runtime.Scheme, kind schema.GroupVersionKind) {
	t.Helper()

	object, err := scheme.New(kind)
	if err != nil {
		t.Fatalf("scheme.New(%s) returned error: %v", kind, err)
	}

	if _, ok := object.(T); !ok {
		t.Fatalf("scheme.New(%s) returned %T, want %T", kind, object, *new(T))
	}
}
