package cmd

import (
	"bytes"
	"os"
	"testing"
)

func TestAddNewUser(t *testing.T) {
	k8s := Kubernetes{
		Name:    "test-cluster",
		Version: "1.0",
		Users:   []string{"alice", "bob"},
	}

	k8s.AddNewUser("charlie")

	expected := []string{"alice", "bob", "charlie"}
	for i, user := range expected {
		if k8s.Users[i] != user {
			t.Errorf("expected user %s at index %d, got %s", user, i, k8s.Users[i])
		}
	}
}

func TestGetUsers(t *testing.T) {
	k8s := Kubernetes{
		Name:    "test-cluster",
		Version: "1.0",
		Users:   []string{"alice", "bob"},
	}

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	k8s.GetUsers()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, err := buf.ReadFrom(r)
	if err != nil {
		t.Fatalf("failed to read from pipe: %v", err)
	}

	output := buf.String()
	expected := "alice\nbob\n"
	if output != expected {
		t.Errorf("expected output %q, got %q", expected, output)
	}
}
