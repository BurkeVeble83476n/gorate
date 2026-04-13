package policy

import (
	"errors"
	"fmt"
	"strings"
)

// Template defines a reusable policy template with placeholder values.
type Template struct {
	Name     string            `yaml:"name"`
	Endpoint string            `yaml:"endpoint"`
	Method   string            `yaml:"method"`
	Limit    int               `yaml:"limit"`
	Window   int               `yaml:"window"`
	Vars     map[string]string `yaml:"vars,omitempty"`
}

// ApplyTemplate creates a Policy from a Template by substituting variable
// placeholders in the form {{VAR_NAME}} with values from the provided vars map.
func ApplyTemplate(tmpl Template, vars map[string]string) (Policy, error) {
	if tmpl.Name == "" {
		return Policy{}, errors.New("template name is required")
	}

	merged := make(map[string]string)
	for k, v := range tmpl.Vars {
		merged[k] = v
	}
	for k, v := range vars {
		merged[k] = v
	}

	name := interpolate(tmpl.Name, merged)
	endpoint := interpolate(tmpl.Endpoint, merged)
	method := interpolate(tmpl.Method, merged)

	p := Policy{
		Name:     name,
		Endpoint: endpoint,
		Method:   method,
		Limit:    tmpl.Limit,
		Window:   tmpl.Window,
	}

	ApplyDefaults(&p)

	if err := Validate(p); err != nil {
		return Policy{}, fmt.Errorf("template %q produced invalid policy: %w", tmpl.Name, err)
	}

	return p, nil
}

// interpolate replaces all {{KEY}} occurrences in s with values from vars.
func interpolate(s string, vars map[string]string) string {
	for k, v := range vars {
		s = strings.ReplaceAll(s, "{"+"{"+k+"}}", v)
	}
	return s
}
