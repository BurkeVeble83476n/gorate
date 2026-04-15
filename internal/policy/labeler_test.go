package policy

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeLabelerPolicies() []Policy {
	return []Policy{
		{Name: "alpha", Endpoint: "/api/alpha", Method: "GET", Limit: 10, Window: 60, Labels: []string{"core"}},
		{Name: "beta", Endpoint: "/api/beta", Method: "POST", Limit: 5, Window: 30, Labels: []string{}},
		{Name: "gamma", Endpoint: "/api/gamma", Method: "*", Limit: 20, Window: 120, Labels: []string{"core", "public"}},
	}
}

func TestAddLabel_NewLabel(t *testing.T) {
	policies := makeLabelerPolicies()
	updated, err := AddLabel(policies, "beta", "internal")
	require.NoError(t, err)
	labels, _ := GetLabels(updated, "beta")
	assert.Contains(t, labels, "internal")
}

func TestAddLabel_DuplicateIsIdempotent(t *testing.T) {
	policies := makeLabelerPolicies()
	updated, err := AddLabel(policies, "alpha", "core")
	require.NoError(t, err)
	labels, _ := GetLabels(updated, "alpha")
	count := 0
	for _, l := range labels {
		if l == "core" {
			count++
		}
	}
	assert.Equal(t, 1, count)
}

func TestAddLabel_PolicyNotFound(t *testing.T) {
	policies := makeLabelerPolicies()
	_, err := AddLabel(policies, "nonexistent", "tag")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestAddLabel_EmptyName(t *testing.T) {
	policies := makeLabelerPolicies()
	_, err := AddLabel(policies, "", "tag")
	assert.Error(t, err)
}

func TestAddLabel_EmptyLabel(t *testing.T) {
	policies := makeLabelerPolicies()
	_, err := AddLabel(policies, "alpha", "")
	assert.Error(t, err)
}

func TestRemoveLabel_Existing(t *testing.T) {
	policies := makeLabelerPolicies()
	updated, err := RemoveLabel(policies, "alpha", "core")
	require.NoError(t, err)
	labels, _ := GetLabels(updated, "alpha")
	assert.NotContains(t, labels, "core")
}

func TestRemoveLabel_NonExistentLabelNoError(t *testing.T) {
	policies := makeLabelerPolicies()
	_, err := RemoveLabel(policies, "alpha", "missing")
	assert.NoError(t, err)
}

func TestRemoveLabel_PolicyNotFound(t *testing.T) {
	policies := makeLabelerPolicies()
	_, err := RemoveLabel(policies, "ghost", "core")
	assert.Error(t, err)
}

func TestGetLabels_ReturnsCorrectLabels(t *testing.T) {
	policies := makeLabelerPolicies()
	labels, err := GetLabels(policies, "gamma")
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"core", "public"}, labels)
}

func TestGetLabels_PolicyNotFound(t *testing.T) {
	policies := makeLabelerPolicies()
	_, err := GetLabels(policies, "unknown")
	assert.Error(t, err)
}

func TestFilterByLabel_ReturnsMatching(t *testing.T) {
	policies := makeLabelerPolicies()
	result := FilterByLabel(policies, "core")
	assert.Len(t, result, 2)
}

func TestFilterByLabel_NoMatches(t *testing.T) {
	policies := makeLabelerPolicies()
	result := FilterByLabel(policies, "nonexistent")
	assert.Empty(t, result)
}
