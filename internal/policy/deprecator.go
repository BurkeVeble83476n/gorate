package policy

import (
	"fmt"
	"time"
)

const (
	AnnotationDeprecated = "deprecated"
	AnnotationDeprecatedSince = "deprecated-since"
	AnnotationDeprecatedReason = "deprecated-reason"
	AnnotationDeprecatedReplacement = "deprecated-replacement"
)

// DeprecationInfo holds structured info about a deprecated policy.
type DeprecationInfo struct {
	Name        string
	Since       string
	Reason      string
	Replacement string
}

// Deprecate marks a policy as deprecated with optional metadata.
func Deprecate(policies []Policy, name, reason, replacement string) ([]Policy, error) {
	if name == "" {
		return nil, fmt.Errorf("policy name must not be empty")
	}
	found := false
	for i, p := range policies {
		if p.Name != name {
			continue
		}
		found = true
		if policies[i].Annotations == nil {
			policies[i].Annotations = map[string]string{}
		}
		policies[i].Annotations[AnnotationDeprecated] = "true"
		policies[i].Annotations[AnnotationDeprecatedSince] = time.Now().UTC().Format(time.RFC3339)
		if reason != "" {
			policies[i].Annotations[AnnotationDeprecatedReason] = reason
		}
		if replacement != "" {
			policies[i].Annotations[AnnotationDeprecatedReplacement] = replacement
		}
	}
	if !found {
		return nil, fmt.Errorf("policy %q not found", name)
	}
	return policies, nil
}

// Undeprecate removes deprecation markers from a policy.
func Undeprecate(policies []Policy, name string) ([]Policy, error) {
	if name == "" {
		return nil, fmt.Errorf("policy name must not be empty")
	}
	found := false
	for i, p := range policies {
		if p.Name != name {
			continue
		}
		found = true
		for _, key := range []string{AnnotationDeprecated, AnnotationDeprecatedSince, AnnotationDeprecatedReason, AnnotationDeprecatedReplacement} {
			delete(policies[i].Annotations, key)
		}
	}
	if !found {
		return nil, fmt.Errorf("policy %q not found", name)
	}
	return policies, nil
}

// IsDeprecated returns true if the policy is marked deprecated.
func IsDeprecated(p Policy) bool {
	return p.Annotations[AnnotationDeprecated] == "true"
}

// GetDeprecationInfo returns structured deprecation info for a policy.
func GetDeprecationInfo(p Policy) (DeprecationInfo, bool) {
	if !IsDeprecated(p) {
		return DeprecationInfo{}, false
	}
	return DeprecationInfo{
		Name:        p.Name,
		Since:       p.Annotations[AnnotationDeprecatedSince],
		Reason:      p.Annotations[AnnotationDeprecatedReason],
		Replacement: p.Annotations[AnnotationDeprecatedReplacement],
	}, true
}

// ListDeprecated returns all deprecated policies.
func ListDeprecated(policies []Policy) []Policy {
	var result []Policy
	for _, p := range policies {
		if IsDeprecated(p) {
			result = append(result, p)
		}
	}
	return result
}
