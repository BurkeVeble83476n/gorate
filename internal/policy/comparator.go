package policy

import (
	"fmt"
	"sort"
	"strings"
)

// CompareResult holds the result of comparing two policy sets.
type CompareResult struct {
	OnlyInA  []Policy
	OnlyInB  []Policy
	InBoth   []Policy
	Conflicts []CompareConflict
}

// CompareConflict describes a policy present in both sets but with differing fields.
type CompareConflict struct {
	Name   string
	FieldA string
	FieldB string
	Field  string
}

// Compare compares two slices of policies by name and returns a CompareResult.
func Compare(a, b []Policy) CompareResult {
	result := CompareResult{}

	aMap := make(map[string]Policy, len(a))
	for _, p := range a {
		aMap[p.Name] = p
	}

	bMap := make(map[string]Policy, len(b))
	for _, p := range b {
		bMap[p.Name] = p
	}

	for _, p := range a {
		if q, ok := bMap[p.Name]; ok {
			result.InBoth = append(result.InBoth, p)
			if conflicts := detectConflicts(p, q); len(conflicts) > 0 {
				result.Conflicts = append(result.Conflicts, conflicts...)
			}
		} else {
			result.OnlyInA = append(result.OnlyInA, p)
		}
	}

	for _, p := range b {
		if _, ok := aMap[p.Name]; !ok {
			result.OnlyInB = append(result.OnlyInB, p)
		}
	}

	return result
}

func detectConflicts(a, b Policy) []CompareConflict {
	var conflicts []CompareConflict
	if a.Limit != b.Limit {
		conflicts = append(conflicts, CompareConflict{
			Name:   a.Name,
			Field:  "limit",
			FieldA: fmt.Sprintf("%d", a.Limit),
			FieldB: fmt.Sprintf("%d", b.Limit),
		})
	}
	if a.Window != b.Window {
		conflicts = append(conflicts, CompareConflict{
			Name:   a.Name,
			Field:  "window",
			FieldA: a.Window,
			FieldB: b.Window,
		})
	}
	if !strings.EqualFold(a.Method, b.Method) {
		conflicts = append(conflicts, CompareConflict{
			Name:   a.Name,
			Field:  "method",
			FieldA: a.Method,
			FieldB: b.Method,
		})
	}
	return conflicts
}

// FormatCompare returns a human-readable summary of a CompareResult.
func FormatCompare(r CompareResult) string {
	var sb strings.Builder

	sort.Slice(r.OnlyInA, func(i, j int) bool { return r.OnlyInA[i].Name < r.OnlyInA[j].Name })
	sort.Slice(r.OnlyInB, func(i, j int) bool { return r.OnlyInB[i].Name < r.OnlyInB[j].Name })

	for _, p := range r.OnlyInA {
		sb.WriteString(fmt.Sprintf("< only in A: %s\n", p.Name))
	}
	for _, p := range r.OnlyInB {
		sb.WriteString(fmt.Sprintf("> only in B: %s\n", p.Name))
	}
	for _, c := range r.Conflicts {
		sb.WriteString(fmt.Sprintf("~ conflict [%s] field=%s A=%s B=%s\n", c.Name, c.Field, c.FieldA, c.FieldB))
	}
	if sb.Len() == 0 {
		sb.WriteString("no differences found\n")
	}
	return sb.String()
}
