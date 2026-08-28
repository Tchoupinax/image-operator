package utils

import (
	"os"
	"testing"
)

func TestEnsureSafeE2EClusterRequiresOptIn(t *testing.T) {
	t.Setenv("E2E_TESTS_ENABLED", "")
	t.Setenv("E2E_KUBE_CONTEXT", "kind-kind")

	err := EnsureSafeE2ECluster()
	if err == nil {
		t.Fatal("expected error when E2E_TESTS_ENABLED is unset")
	}
}

func TestEnsureSafeE2EClusterRequiresExplicitContextForNonKind(t *testing.T) {
	t.Setenv("E2E_TESTS_ENABLED", "true")
	t.Setenv("E2E_KUBE_CONTEXT", "")

	if isKindContext("production") {
		t.Fatal("production should not be treated as kind context")
	}
}

func TestIsKindContext(t *testing.T) {
	tests := map[string]bool{
		"kind-kind":     true,
		"kind-dev":      true,
		"kind":          true,
		"prod":          false,
		"production":    false,
		"gke_prod_us_1": false,
	}

	for context, expected := range tests {
		if isKindContext(context) != expected {
			t.Fatalf("isKindContext(%q) = %v, want %v", context, !expected, expected)
		}
	}
}

func TestEnsureSafeE2EClusterAllowsExplicitContext(t *testing.T) {
	if os.Getenv("KUBECONFIG") == "" {
		t.Skip("kubectl context unavailable")
	}

	currentContext, err := currentKubeContext()
	if err != nil {
		t.Skip("kubectl context unavailable")
	}

	t.Setenv("E2E_TESTS_ENABLED", "true")
	t.Setenv("E2E_KUBE_CONTEXT", currentContext)

	if err := EnsureSafeE2ECluster(); err != nil {
		t.Fatalf("expected explicit context match to pass: %v", err)
	}
}
