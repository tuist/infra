package scheme

import (
	"testing"

	infrav1 "github.com/tuist/infra/api/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestSchemeIncludesBuiltInAndInfraTypes(t *testing.T) {
	t.Parallel()

	assertSchemeCreates[*infrav1.HostPool](t, Scheme, infrav1.SchemeGroupVersion.WithKind("HostPool"))
	assertSchemeCreates[*infrav1.HostImage](t, Scheme, infrav1.SchemeGroupVersion.WithKind("HostImage"))
	assertSchemeCreates[*corev1.Pod](t, Scheme, corev1.SchemeGroupVersion.WithKind("Pod"))
}

func TestAddToSchemeRegistersBuiltInAndInfraTypes(t *testing.T) {
	t.Parallel()

	target := runtime.NewScheme()
	if err := AddToScheme(target); err != nil {
		t.Fatalf("AddToScheme() returned error: %v", err)
	}

	assertSchemeCreates[*infrav1.Sandbox](t, target, infrav1.SchemeGroupVersion.WithKind("Sandbox"))
	assertSchemeCreates[*infrav1.HostImage](t, target, infrav1.SchemeGroupVersion.WithKind("HostImage"))
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
