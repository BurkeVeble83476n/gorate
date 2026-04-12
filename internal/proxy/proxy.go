package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/user/gorate/internal/policy"
)

// RateLimiter tracks request counts per policy.
type RateLimiter struct {
	counts map[string]*bucket
}

type bucket struct {
	count    int
	resetAt  time.Time
}

// NewRateLimiter creates a new RateLimiter.
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{counts: make(map[string]*bucket)}
}

// Allow checks whether a request is permitted under the given policy.
func (rl *RateLimiter) Allow(p policy.Policy) bool {
	now := time.Now()
	b, ok := rl.counts[p.Name]
	if !ok || now.After(b.resetAt) {
		rl.counts[p.Name] = &bucket{
			count:   1,
			resetAt: now.Add(time.Duration(p.WindowSecs) * time.Second),
		}
		return true
	}
	if b.count >= p.Limit {
		return false
	}
	b.count++
	return true
}

// Handler returns an http.Handler that enforces the policy and proxies allowed requests.
func Handler(p policy.Policy, target string, rl *RateLimiter) (http.Handler, error) {
	u, err := url.Parse(target)
	if err != nil {
		return nil, fmt.Errorf("invalid target URL %q: %w", target, err)
	}
	rp := httputil.NewSingleHostReverseProxy(u)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != p.Method {
			rp.ServeHTTP(w, r)
			return
		}
		if !rl.Allow(p) {
			w.Header().Set("X-RateLimit-Policy", p.Name)
			http.Error(w, "429 Too Many Requests", http.StatusTooManyRequests)
			return
		}
		rp.ServeHTTP(w, r)
	}), nil
}
