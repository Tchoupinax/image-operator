package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var (
	prometheusInstalledByTests bool
	certManagerInstalledByTests bool
)

// EnsureSafeE2ECluster blocks e2e tests from mutating unexpected clusters.
// Set E2E_TESTS_ENABLED=true and either use a kind context or set E2E_KUBE_CONTEXT
// to the exact disposable cluster context name.
func EnsureSafeE2ECluster() error {
	if os.Getenv("E2E_TESTS_ENABLED") != "true" {
		return fmt.Errorf(
			"e2e tests are disabled: set E2E_TESTS_ENABLED=true to run against a disposable cluster",
		)
	}

	currentContext, err := currentKubeContext()
	if err != nil {
		return err
	}

	if allowedContext := strings.TrimSpace(os.Getenv("E2E_KUBE_CONTEXT")); allowedContext != "" {
		if currentContext != allowedContext {
			return fmt.Errorf(
				"refusing to run e2e tests: current context %q does not match E2E_KUBE_CONTEXT %q",
				currentContext,
				allowedContext,
			)
		}
		return nil
	}

	if !isKindContext(currentContext) {
		return fmt.Errorf(
			"refusing to run e2e tests: context %q is not a kind cluster (set E2E_KUBE_CONTEXT to explicitly allow a disposable cluster)",
			currentContext,
		)
	}

	return nil
}

func currentKubeContext() (string, error) {
	cmd := exec.Command("kubectl", "config", "current-context")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("unable to read kubectl context: %w", err)
	}

	context := strings.TrimSpace(string(output))
	if context == "" {
		return "", fmt.Errorf("kubectl current-context is empty")
	}

	return context, nil
}

func isKindContext(context string) bool {
	return context == "kind" || strings.HasPrefix(context, "kind-")
}
