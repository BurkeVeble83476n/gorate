package policy

import (
	"testing"
)

func makeSortPolicies() []Policy {
	return []Policy{
		{Name: "gamma", Endpoint: "/c", Method: "POST", Limit: 30, Window: 60},
		{Name: "alpha", Endpoint: "/a", Method: "GET", Limit: 10, Window: 30},
		{Name: "beta", Endpoint: "/b", Method: "DELETE", Limit: 20, Window: 90},
	}
}

func TestSort_ByNameAsc(t *testing.T) {
	result := Sort(makeSortPolicies(), SortOptions{Field: SortByName, Order: SortAsc})
	if result[0].Name != "alpha" || result[1].Name != "beta" || result[2].Name != "gamma" {
		t.Errorf("unexpected order: %v", result)
	}
}

func TestSort_ByNameDesc(t *testing.T) {
	result := Sort(makeSortPolicies(), SortOptions{Field: SortByName, Order: SortDesc})
	if result[0].Name != "gamma" || result[1].Name != "beta" || result[2].Name != "alpha" {
		t.Errorf("unexpected order: %v", result)
	}
}

func TestSort_ByLimitAsc(t *testing.T) {
	result := Sort(makeSortPolicies(), SortOptions{Field: SortByLimit, Order: SortAsc})
	if result[0].Limit != 10 || result[1].Limit != 20 || result[2].Limit != 30 {
		t.Errorf("unexpected order by limit: %v", result)
	}
}

func TestSort_ByWindowDesc(t *testing.T) {
	result := Sort(makeSortPolicies(), SortOptions{Field: SortByWindow, Order: SortDesc})
	if result[0].Window != 90 || result[1].Window != 60 || result[2].Window != 30 {
		t.Errorf("unexpected order by window: %v", result)
	}
}

func TestSort_ByEndpointAsc(t *testing.T) {
	result := Sort(makeSortPolicies(), SortOptions{Field: SortByEndpoint, Order: SortAsc})
	if result[0].Endpoint != "/a" || result[2].Endpoint != "/c" {
		t.Errorf("unexpected order by endpoint: %v", result)
	}
}

func TestSort_UnknownField_PreservesOrder(t *testing.T) {
	original := makeSortPolicies()
	result := Sort(original, SortOptions{Field: SortField("unknown"), Order: SortAsc})
	for i, p := range result {
		if p.Name != original[i].Name {
			t.Errorf("order changed for unknown field at index %d", i)
		}
	}
}

func TestSort_DoesNotMutateOriginal(t *testing.T) {
	original := makeSortPolicies()
	originalFirst := original[0].Name
	Sort(original, SortOptions{Field: SortByName, Order: SortAsc})
	if original[0].Name != originalFirst {
		t.Error("Sort mutated the original slice")
	}
}
