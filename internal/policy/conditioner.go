package policy

import (
	"fmt"
	"strings"
)

// Condition represents a rule that must be satisfied for a policy to be active.
type Condition struct {
	Field    string `json:"field" yaml:"field"`
	Operator string `json:"operator" yaml:"operator"`
	Value    string `json:"value" yaml:"value"`
}

// validOperators lists all supported condition operators.
var validOperators = map[string]bool{
	"eq":       true,
	"neq":      true,
	"contains": true,
	"prefix":   true,
}

// SetCondition attaches a condition to a named policy.
func SetCondition(policies []Policy, name, field, operator, value string) ([]Policy, error) {
	if name == "" {
		return nil, fmt.Errorf("policy name must not be empty")
	}
	if !validOperators[operator] {
		return nil, fmt.Errorf("unsupported operator %q: must be one of eq, neq, contains, prefix", operator)
	}
	for i, p := range policies {
		if p.Name == name {
			if policies[i].Annotations == nil {
				policies[i].Annotations = map[string]string{}
			}
			policies[i].Annotations["condition.field"] = field
			policies[i].Annotations["condition.operator"] = operator
			policies[i].Annotations["condition.value"] = value
			return policies, nil
		}
	}
	return nil, fmt.Errorf("policy %q not found", name)
}

// GetCondition retrieves the condition attached to a named policy, if any.
func GetCondition(policies []Policy, name string) (*Condition, error) {
	for _, p := range policies {
		if p.Name == name {
			if p.Annotations == nil {
				return nil, nil
			}
			field, ok1 := p.Annotations["condition.field"]
			op, ok2 := p.Annotations["condition.operator"]
			val, ok3 := p.Annotations["condition.value"]
			if !ok1 || !ok2 || !ok3 {
				return nil, nil
			}
			return &Condition{Field: field, Operator: op, Value: val}, nil
		}
	}
	return nil, fmt.Errorf("policy %q not found", name)
}

// EvaluateCondition checks whether the given attributes satisfy the policy's condition.
// Returns true if no condition is set or if the condition is met.
func EvaluateCondition(p Policy, attrs map[string]string) bool {
	if p.Annotations == nil {
		return true
	}
	field := p.Annotations["condition.field"]
	op := p.Annotations["condition.operator"]
	expected := p.Annotations["condition.value"]
	if field == "" || op == "" {
		return true
	}
	actual := attrs[field]
	switch op {
	case "eq":
		return actual == expected
	case "neq":
		return actual != expected
	case "contains":
		return strings.Contains(actual, expected)
	case "prefix":
		return strings.HasPrefix(actual, expected)
	}
	return false
}
