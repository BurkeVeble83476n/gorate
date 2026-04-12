package policy

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// rawPolicy is used for JSON unmarshalling with string-based duration.
type rawPolicy struct {
	Name     string `json:"name"`
	Endpoint string `json:"endpoint"`
	Method   string `json:"method"`
	Limit    int    `json:"limit"`
	Window   string `json:"window"`
}

// LoadFromFile reads and parses a JSON policy file, returning validated policies.
func LoadFromFile(path string) ([]*Policy, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading policy file: %w", err)
	}

	var raws []rawPolicy
	if err := json.Unmarshal(data, &raws); err != nil {
		return nil, fmt.Errorf("parsing policy file: %w", err)
	}

	policies := make([]*Policy, 0, len(raws))
	for i, r := range raws {
		window, err := time.ParseDuration(r.Window)
		if err != nil {
			return nil, fmt.Errorf("policy[%d] invalid window %q: %w", i, r.Window, err)
		}
		p := &Policy{
			Name:     r.Name,
			Endpoint: r.Endpoint,
			Method:   r.Method,
			Limit:    r.Limit,
			Window:   window,
		}
		if err := p.Validate(); err != nil {
			return nil, fmt.Errorf("policy[%d] validation failed: %w", i, err)
		}
		policies = append(policies, p)
	}
	return policies, nil
}
