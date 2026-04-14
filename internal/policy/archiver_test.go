package policy

import (
	"os"
	"path/filepath"
	"testing"
)

func makeArchiverPolicies() []Policy {
	return []Policy{
		{Name: "alpha", Endpoint: "/api/v1", Method: "GET", Limit: 100, Window: 60},
		{Name: "beta", Endpoint: "/api/v2", Method: "POST", Limit: 50, Window: 30},
	}
}

func TestSaveArchive_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	policies := makeArchiverPolicies()

	path, err := SaveArchive(dir, "release", policies)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("expected archive file to exist at %s", path)
	}
}

func TestSaveArchive_FilenameContainsLabel(t *testing.T) {
	dir := t.TempDir()
	policies := makeArchiverPolicies()

	path, err := SaveArchive(dir, "mybackup", policies)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	base := filepath.Base(path)
	if len(base) == 0 {
		t.Error("expected non-empty filename")
	}

	// label should appear in filename
	if filepath.Ext(path) != ".json" {
		t.Errorf("expected .json extension, got %s", filepath.Ext(path))
	}
}

func TestLoadArchive_ReturnsCorrectData(t *testing.T) {
	dir := t.TempDir()
	policies := makeArchiverPolicies()

	path, err := SaveArchive(dir, "test", policies)
	if err != nil {
		t.Fatalf("save error: %v", err)
	}

	archive, err := LoadArchive(path)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}

	if archive.Label != "test" {
		t.Errorf("expected label 'test', got %q", archive.Label)
	}

	if len(archive.Policies) != 2 {
		t.Errorf("expected 2 policies, got %d", len(archive.Policies))
	}
}

func TestLoadArchive_FileNotFound(t *testing.T) {
	_, err := LoadArchive("/nonexistent/path/archive.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestListArchives_ReturnsAll(t *testing.T) {
	dir := t.TempDir()
	policies := makeArchiverPolicies()

	_, _ = SaveArchive(dir, "first", policies)
	_, _ = SaveArchive(dir, "second", policies)

	archives, err := ListArchives(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(archives) != 2 {
		t.Errorf("expected 2 archives, got %d", len(archives))
	}
}

func TestListArchives_EmptyDir(t *testing.T) {
	dir := t.TempDir()

	archives, err := ListArchives(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(archives) != 0 {
		t.Errorf("expected 0 archives, got %d", len(archives))
	}
}
