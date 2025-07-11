package mcp

import "testing"

func TestGetMcpContextUnset(t *testing.T) {
	ResetMcpContext()
	instance, cluster, org := GetMcpContext()
	if instance != "" || cluster != "" || org != "" {
		t.Fatalf("expected empty context when unset, got %q %q %q", instance, cluster, org)
	}
}

func TestGetMcpContextSet(t *testing.T) {
	ResetMcpContext()
	SetMcpContext("inst", "cluster", "org")
	instance, cluster, org := GetMcpContext()
	if instance != "inst" || cluster != "cluster" || org != "org" {
		t.Fatalf("expected context to be returned, got %q %q %q", instance, cluster, org)
	}
}
