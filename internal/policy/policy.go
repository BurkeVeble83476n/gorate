package policy

import (
	"errors"
	"time"
)

// Policy defines a rate-limit policy for an HTTP endpoint.
type Policy struct {
	Name     string        `json:"name"`
	Endpoint string        `json:"endpoint"`
	Method   string        `json:"method"`
	Limit    int           `json:"limit"`
	Window   time.Duration `json:"window"`
}

// Validate checks that the policy fields are valid.
func (p *Policy) Validate() error {
	if p.Name == "" {
		return errors.New("policy name is required")
	}
	if p.Endpoint == "" {
		return errors.New("policy endpoint is required")
	}
	if p.Limit <= 0 {
		return errors.New("policy limit must be greater than zero")
	}
	if p.Window <= 0 {
		return errors.New("policy window must be greater than zero")
	}
	if p.Method == "" {
		p.Method = "*"
	}
	return nil
}

// RatePerSecond returns the effective rate in requests per second.
func (p *Policy) RatePerSecond() float64 {
	return float64(p.Limit) / p.Window.Seconds()
}
