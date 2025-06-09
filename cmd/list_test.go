package cmd

import "testing"

func TestGetKubeClient_InvalidPath(t *testing.T) {
	_, err := getKubeClient("/invalid/path")
	if err == nil {
		t.Error("expected error for invalid kubeconfig path")
	}
}
