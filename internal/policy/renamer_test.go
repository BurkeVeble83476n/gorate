package policy

import (
	"testing"
)

func makeRenamerPolicies() []Policy {
	return []Policy{
		{Name: "alpha", Endpoint: "/a", Method: "GET", Limit: 10, Window: 60},
		{Name: "beta", Endpoint: "/b", Method: "POST", Limit: 20, Window: 60},
		{Name: "gamma", Endpoint: "/c", Method: "PUT", Limit: 30, Window: 60},
	}
}

func TestRename_Success(t *testing.T) {
	policies := makeRenamerPolicies()
	result, err := Rename(policies, RenameOptions{OldName: "alpha", NewName: "zeta"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[0].Name != "zeta" {
		t.Errorf("expected name %q, got %q", "zeta", result[0].Name)
	}
}

func TestRename_OldNameNotFound(t *testing.T) {
	policies := makeRenamerPolicies()
	_, err := Rename(policies, RenameOptions{OldName: "missing", NewName: "zeta"})
	if err == nil {
		t.Fatal("expected error for missing old name")
	}
}

func TestRename_NewNameConflict(t *testing.T) {
	policies := makeRenamerPolicies()
	_, err := Rename(policies, RenameOptions{OldName: "alpha", NewName: "beta"})
	if err == nil {
		t.Fatal("expected error when new name already exists")
	}
}

func TestRename_SameName(t *testing.T) {
	policies := makeRenamerPolicies()
	_, err := Rename(policies, RenameOptions{OldName: "alpha", NewName: "alpha"})
	if err == nil {
		t.Fatal("expected error when old and new names are identical")
	}
}

func TestRename_EmptyOldName(t *testing.T) {
	policies := makeRenamerPolicies()
	_, err := Rename(policies, RenameOptions{OldName: "", NewName: "zeta"})
	if err == nil {
		t.Fatal("expected error for empty old name")
	}
}

func TestRename_EmptyNewName(t *testing.T) {
	policies := makeRenamerPolicies()
	_, err := Rename(policies, RenameOptions{OldName: "alpha", NewName: ""})
	if err == nil {
		t.Fatal("expected error for empty new name")
	}
}

func TestRename_OtherPoliciesUnchanged(t *testing.T) {
	policies := makeRenamerPolicies()
	result, err := Rename(policies, RenameOptions{OldName: "beta", NewName: "delta"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[0].Name != "alpha" {
		t.Errorf("expected alpha unchanged, got %q", result[0].Name)
	}
	if result[2].Name != "gamma" {
		t.Errorf("expected gamma unchanged, got %q", result[2].Name)
	}
}
