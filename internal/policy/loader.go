package policy

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// policyFile represents the top-level structure of a policy YAML file.
type policyFile struct {
	Policies []Policy `yaml:"policies"`
}

// LoadFromFile reads a YAML file at the given path and returns a slice of
// validated Policy values. An error is returned if the file cannot be read,
// parsed, or if any policy fails validation.
func LoadFromFile(path string) ([]Policy, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading policy file %q: %w", path, err)
	}

	var pf policyFile
	if err := yaml.Unmarshal(data, &pf); err != nil {
		return nil, fmt.Errorf("parsing policy file %q: %w", path, err)
	}

	for i := range pf.Policies {
		if err := Validate(&pf.Policies[i]); err != nil {
			return nil, fmt.Errorf("policy %d in %q: %w", i, path, err)
		}
	}

	return pf.Policies, nil
}
