package policy

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func makeRequest(method, path string) *http.Request {
	req := httptest.NewRequest(method, path, nil)
	return req
}

func TestMatch_ExactEndpointAndMethod(t *testing.T) {
	p := Policy{Name: "p", Endpoint: "/api/v1", Method: "GET", Limit: 10, Window: "1m"}
	r := makeRequest("GET", "/api/v1")
	if !Match(r, p) {
		t.Error("expected match")
	}
}

func TestMatch_WrongMethod(t *testing.T) {
	p := Policy{Name: "p", Endpoint: "/api/v1", Method: "POST", Limit: 10, Window: "1m"}
	r := makeRequest("GET", "/api/v1")
	if Match(r, p) {
		t.Error("expected no match on wrong method")
	}
}

func TestMatch_WrongEndpoint(t *testing.T) {
	p := Policy{Name: "p", Endpoint: "/api/v2", Method: "GET", Limit: 10, Window: "1m"}
	r := makeRequest("GET", "/api/v1")
	if Match(r, p) {
		t.Error("expected no match on wrong endpoint")
	}
}

func TestMatch_WildcardMethod(t *testing.T) {
	p := Policy{Name: "p", Endpoint: "/api/v1", Method: "*", Limit: 10, Window: "1m"}
	for _, method := range []string{"GET", "POST", "DELETE", "PUT"} {
		r := makeRequest(method, "/api/v1")
		if !Match(r, p) {
			t.Errorf("expected wildcard method to match %s", method)
		}
	}
}

func TestMatch_WildcardEndpointPrefix(t *testing.T) {
	p := Policy{Name: "p", Endpoint: "/api/*", Method: "GET", Limit: 10, Window: "1m"}
	r := makeRequest("GET", "/api/users")
	if !Match(r, p) {
		t.Error("expected wildcard prefix endpoint to match")
	}
}

func TestMatch_WildcardEndpointNoMatch(t *testing.T) {
	p := Policy{Name: "p", Endpoint: "/api/*", Method: "GET", Limit: 10, Window: "1m"}
	r := makeRequest("GET", "/other/path")
	if Match(r, p) {
		t.Error("expected wildcard prefix endpoint not to match /other/path")
	}
}

func TestFindMatch_ReturnsFirst(t *testing.T) {
	policies := []Policy{
		{Name: "first", Endpoint: "/api/v1", Method: "GET", Limit: 5, Window: "1m"},
		{Name: "second", Endpoint: "/api/v1", Method: "*", Limit: 100, Window: "1m"},
	}
	r := makeRequest("GET", "/api/v1")
	match := FindMatch(r, policies)
	if match == nil {
		t.Fatal("expected a match")
	}
	if match.Name != "first" {
		t.Errorf("expected first policy, got %s", match.Name)
	}
}

func TestFindMatch_NoMatch(t *testing.T) {
	policies := []Policy{
		{Name: "p", Endpoint: "/api/v2", Method: "POST", Limit: 5, Window: "1m"},
	}
	r := makeRequest("GET", "/api/v1")
	if FindMatch(r, policies) != nil {
		t.Error("expected no match")
	}
}
