package proxy

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/gorate/internal/policy"
)

func makePolicy(limit int, window int) policy.Policy {
	return policy.Policy{
		Name:       "test",
		Endpoint:   "/api",
		Method:     "GET",
		Limit:      limit,
		WindowSecs: window,
	}
}

func TestAllow_UnderLimit(t *testing.T) {
	rl := NewRateLimiter()
	p := makePolicy(3, 60)
	for i := 0; i < 3; i++ {
		if !rl.Allow(p) {
			t.Fatalf("expected request %d to be allowed", i+1)
		}
	}
}

func TestAllow_ExceedsLimit(t *testing.T) {
	rl := NewRateLimiter()
	p := makePolicy(2, 60)
	rl.Allow(p)
	rl.Allow(p)
	if rl.Allow(p) {
		t.Fatal("expected third request to be denied")
	}
}

func TestAllow_WindowReset(t *testing.T) {
	rl := NewRateLimiter()
	p := makePolicy(1, 0) // 0-second window expires immediately
	rl.Allow(p)
	time.Sleep(10 * time.Millisecond)
	if !rl.Allow(p) {
		t.Fatal("expected request to be allowed after window reset")
	}
}

func TestHandler_RateLimited(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	p := makePolicy(1, 60)
	rl := NewRateLimiter()
	h, err := Handler(p, backend.URL, rl)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	rr2 := httptest.NewRecorder()
	h.ServeHTTP(rr2, req)
	if rr2.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d", rr2.Code)
	}
}
