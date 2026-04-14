package policy

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Archive represents a saved archive of policies with metadata.
type Archive struct {
	CreatedAt time.Time `json:"created_at"`
	Label     string    `json:"label"`
	Policies  []Policy  `json:"policies"`
}

// SaveArchive writes policies to a named archive file in the given directory.
func SaveArchive(dir, label string, policies []Policy) (string, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create archive directory: %w", err)
	}

	timestamp := time.Now().UTC().Format("20060102T150405Z")
	filename := fmt.Sprintf("%s_%s.json", timestamp, label)
	path := filepath.Join(dir, filename)

	archive := Archive{
		CreatedAt: time.Now().UTC(),
		Label:     label,
		Policies:  policies,
	}

	data, err := json.MarshalIndent(archive, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal archive: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write archive file: %w", err)
	}

	return path, nil
}

// LoadArchive reads an archive file and returns its contents.
func LoadArchive(path string) (*Archive, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read archive file: %w", err)
	}

	var archive Archive
	if err := json.Unmarshal(data, &archive); err != nil {
		return nil, fmt.Errorf("failed to parse archive file: %w", err)
	}

	return &archive, nil
}

// ListArchives returns all archive files found in the given directory.
func ListArchives(dir string) ([]Archive, error) {
	entries, err := filepath.Glob(filepath.Join(dir, "*.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to list archives: %w", err)
	}

	var archives []Archive
	for _, entry := range entries {
		a, err := LoadArchive(entry)
		if err != nil {
			continue
		}
		archives = append(archives, *a)
	}

	return archives, nil
}
