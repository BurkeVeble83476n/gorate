package policy

import "strings"

// Policy tags allow grouping and categorization of policies.
// Tags are stored as a comma-separated string in the Tag field.

// AddTag adds a tag to the policy if it doesn't already exist.
func AddTag(p *Policy, tag string) {
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return
	}
	existing := GetTags(p)
	for _, t := range existing {
		if t == tag {
			return
		}
	}
	if p.Tags == "" {
		p.Tags = tag
	} else {
		p.Tags = p.Tags + "," + tag
	}
}

// RemoveTag removes a tag from the policy if it exists.
func RemoveTag(p *Policy, tag string) {
	tag = strings.TrimSpace(tag)
	existing := GetTags(p)
	filtered := make([]string, 0, len(existing))
	for _, t := range existing {
		if t != tag {
			filtered = append(filtered, t)
		}
	}
	p.Tags = strings.Join(filtered, ",")
}

// GetTags returns the list of tags for the policy.
func GetTags(p *Policy) []string {
	if p.Tags == "" {
		return []string{}
	}
	parts := strings.Split(p.Tags, ",")
	tags := make([]string, 0, len(parts))
	for _, t := range parts {
		t = strings.TrimSpace(t)
		if t != "" {
			tags = append(tags, t)
		}
	}
	return tags
}

// HasTag returns true if the policy has the given tag.
func HasTag(p *Policy, tag string) bool {
	tag = strings.TrimSpace(tag)
	for _, t := range GetTags(p) {
		if t == tag {
			return true
		}
	}
	return false
}

// FilterByTag returns policies that have the given tag.
func FilterByTag(policies []Policy, tag string) []Policy {
	result := make([]Policy, 0)
	for _, p := range policies {
		if HasTag(&p, tag) {
			result = append(result, p)
		}
	}
	return result
}
