package proxy

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/gorate/internal/policy"
)

func TestNewServer_InvalidTarget(t *testing.T) {
	policies := []policy.Policy{makePolicy(5, 60)}
	_, err := NewServer(":0", "://bad url", policies)
	if err == nil {
		t.Fatal("expected error for invalid target URL")
	}
}

func TestNewServer_ValidPolicies(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	policies := []policy.Policy{makePolicy(10, 60)}
	s, err := NewServer(":0", backend.URL, policies)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil server")
	}
}

func TestServer_StartsAndShutdown(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	policies := []policy.Policy{makePolicy(10, 60)}
	s, err := NewServer(":19876", backend.URL, policies)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	errCh := make(chan error, 1)
	go func() { errCh <- s.Start(ctx) }()

	time.Sleep(50 * time.Millisecond)
	resp, err := http.Get(fmt.Sprintf("http://localhost:19876%s", "/api"))
	if err == nil {
		resp.Body.Close()
	}

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("server returned error on shutdown: %v", err)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("server did not shut down in time")
	}
}
