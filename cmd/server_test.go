package cmd

import (
	"testing"
)

func TestServerCommandDefined(t *testing.T) {
	if serverCmd == nil {
		t.Fatal("serverCmd should be defined")
	}
	if serverCmd.Use != "server" {
		t.Errorf("expected command use 'server', got %s", serverCmd.Use)
	}
	portFlag := serverCmd.Flags().Lookup("port")
	if portFlag == nil {
		t.Error("expected 'port' flag to be defined")
	}
}

func TestGetServerKubeClient_InvalidPath(t *testing.T) {
	_, err := getServerKubeClient("/invalid/path", false)
	if err == nil {
		t.Error("expected error for invalid kubeconfig path")
	}
}
