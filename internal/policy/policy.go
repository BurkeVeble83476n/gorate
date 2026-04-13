package policy

import (
	"errors"
	"time"
)

// Policy defines a rate-limit rule for a specific HTTP endpoint.
type Policy struct {
	Name     string `yaml:"name"`
	Endpoint string `yaml:"endpoint"`
	Method   string `yaml:"method"`
	Limit    int    `yaml:"limit"`
	Window   string `yaml:"window"`
}

// ParsedWindow returns the parsed duration from the Window field.
func (p Policy) ParsedWindow() (time.Duration, error) {
	return time.ParseDuration(p.Window)
}

// Validate checks that a Policy has all required fields and valid values.
func Validate(p Policy) error {
	if p.Name == "" {
		return errors.New("policy name is required")
	}
	if p.Endpoint == "" {
		return errors.New("policy endpoint is required")
	}
	if p.Limit <= 0 {
		return errors.New("policy limit must be greater than zero")
	}
	if _, err := p.ParsedWindow(); err != nil {
		return errors.New("policy window is invalid: " + err.Error())
	}
	if p.Method == "" {
		p.Method = "*"
	}
	return nil
}

// ApplyDefaults fills in optional fields with sensible defaults.
func ApplyDefaults(p *Policy) {
	if p.Method == "" {
		p.Method = "*"
	}
}
