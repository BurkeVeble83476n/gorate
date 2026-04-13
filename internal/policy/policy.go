package policy

import "fmt"

// Policy defines a rate-limit rule for an HTTP endpoint.
type Policy struct {
	Name     string `yaml:"name"`
	Endpoint string `yaml:"endpoint"`
	Method   string `yaml:"method"`
	Limit    int    `yaml:"limit"`
	Window   int    `yaml:"window"` // seconds
}

// Validate checks that a Policy has the required fields and valid values.
func Validate(p Policy) error {
	if p.Name == "" {
		return fmt.Errorf("policy missing required field: name")
	}
	if p.Endpoint == "" {
		return fmt.Errorf("policy %q missing required field: endpoint", p.Name)
	}
	if p.Limit <= 0 {
		return fmt.Errorf("policy %q: limit must be greater than 0", p.Name)
	}
	if p.Window <= 0 {
		return fmt.Errorf("policy %q: window must be greater than 0", p.Name)
	}
	return nil
}

// ApplyDefaults fills in optional fields with sensible defaults.
func ApplyDefaults(p *Policy) {
	if p.Method == "" {
		p.Method = "*"
	}
}
