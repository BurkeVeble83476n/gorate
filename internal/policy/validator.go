package policy

import (
	"fmt"
	"net/url"
	"strings"
)

// ValidationError holds a list of validation issues for a policy.
type ValidationError struct {
	Name   string
	Errors []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("policy %q has %d validation error(s): %s",
		e.Name, len(e.Errors), strings.Join(e.Errors, "; "))
}

// ValidateEndpoint checks that the endpoint is a valid relative or absolute path.
func ValidateEndpoint(endpoint string) error {
	if endpoint == "" {
		return fmt.Errorf("endpoint must not be empty")
	}
	if endpoint != "*" {
		_, err := url.ParseRequestURI(endpoint)
		if err != nil {
			return fmt.Errorf("endpoint %q is not a valid URI path: %w", endpoint, err)
		}
	}
	return nil
}

// ValidateMethod checks that the HTTP method is a known valid value or wildcard.
func ValidateMethod(method string) error {
	valid := map[string]bool{
		"GET": true, "POST": true, "PUT": true,
		"DELETE": true, "PATCH": true, "HEAD": true,
		"OPTIONS": true, "*": true,
	}
	upper := strings.ToUpper(method)
	if !valid[upper] {
		return fmt.Errorf("method %q is not a valid HTTP method or wildcard", method)
	}
	return nil
}

// ValidateLimit checks that limit and window values are positive.
func ValidateLimit(limit int, windowSeconds int) []string {
	var errs []string
	if limit <= 0 {
		errs = append(errs, fmt.Sprintf("limit must be > 0, got %d", limit))
	}
	if windowSeconds <= 0 {
		errs = append(errs, fmt.Sprintf("window_seconds must be > 0, got %d", windowSeconds))
	}
	return errs
}
