package policy

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func makeRequest(method, path string) *http.Request {
	return httptest.NewRequest(method, path, nil)
}

func TestMatch_ExactEndpointAndMethod(t *testing.T) {
	p := Policy{Name: "test", Endpoint: "/api/v1", Method: "GET", Limit: 10, Window: 60}
	if !Match(p, makeRequest("GET", "/api/v1")) {
		t.Error("expected match for exact endpoint and method")
	}
}

func TestMatch_WrongMethod(t *testing.T) {
	p := Policy{Name: "test", Endpoint: "/api/v1", Method: "POST", Limit: 10, Window: 60}
	if Match(p, makeRequest("GET", "/api/v1")) {
		t.Error("expected no match for wrong method")
	}
}

func TestMatch_WrongEndpoint(t *testing.T) {
	p := Policy{Name: "test", Endpoint: "/api/v1", Method: "GET", Limit: 10, Window: 60}
	if Match(p, makeRequest("GET", "/api/v2")) {
		t.Error("expected no match for wrong endpoint")
	}
}

func TestMatch_WildcardMethod(t *testing.T) {
	p := Policy{Name: "test", Endpoint: "/health", Method: "*", Limit: 5, Window: 30}
	for _, m := range []string{"GET", "POST", "DELETE"} {
		if !Match(p, makeRequest(m, "/health")) {
			t.Errorf("expected wildcard method to match %s", m)
		}
	}
}

func TestMatch_WildcardEndpoint(t *testing.T) {
	p := Policy{Name: "test", Endpoint: "/api/*", Method: "GET", Limit: 10, Window: 60}
	for _, path := range []string{"/api/users", "/api/orders", "/api/v2/items"} {
		if !Match(p, makeRequest("GET", path)) {
			t.Errorf("expected wildcard endpoint to match path %s", path)
		}
	}
}

func TestMatch_WildcardEndpoint_NoMatch(t *testing.T) {
	p := Policy{Name: "test", Endpoint: "/api/*", Method: "GET", Limit: 10, Window: 60}
	if Match(p, makeRequest("GET", "/other/path")) {
		t.Error("expected wildcard endpoint not to match /other/path")
	}
}

func TestFindMatch_ReturnsFirst(t *testing.T) {
	policies := []Policy{
		{Name: "p1", Endpoint: "/foo", Method: "GET", Limit: 5, Window: 10},
		{Name: "p2", Endpoint: "/foo", Method: "GET", Limit: 20, Window: 10},
	}
	r := makeRequest("GET", "/foo")
	got := FindMatch(policies, r)
	if got == nil {
		t.Fatal("expected a match, got nil")
	}
	if got.Name != "p1" {
		t.Errorf("expected first matching policy 'p1', got '%s'", got.Name)
	}
}

func TestFindMatch_NoMatch(t *testing.T) {
	policies := []Policy{
		{Name: "p1", Endpoint: "/bar", Method: "POST", Limit: 5, Window: 10},
	}
	r := makeRequest("GET", "/foo")
	if got := FindMatch(policies, r); got != nil {
		t.Errorf("expected nil, got policy '%s'", got.Name)
	}
}
