package proxy

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/user/gorate/internal/policy"
)

// Server wraps an http.Server with rate-limiting proxy handlers.
type Server struct {
	addr   string
	httpSv *http.Server
}

// NewServer builds a Server that enforces all provided policies against target.
func NewServer(addr, target string, policies []policy.Policy) (*Server, error) {
	rl := NewRateLimiter()
	mux := http.NewServeMux()

	for _, p := range policies {
		h, err := Handler(p, target, rl)
		if err != nil {
			return nil, fmt.Errorf("policy %q: %w", p.Name, err)
		}
		log.Printf("registering policy %q on %s %s (limit %d/%ds)",
			p.Name, p.Method, p.Endpoint, p.Limit, p.WindowSecs)
		mux.Handle(p.Endpoint, h)
	}

	return &Server{
		addr: addr,
		httpSv: &http.Server{
			Addr:         addr,
			Handler:      mux,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}, nil
}

// Start begins listening and blocks until the context is cancelled.
func (s *Server) Start(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		log.Printf("gorate proxy listening on %s", s.addr)
		if err := s.httpSv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()
	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return s.httpSv.Shutdown(shutCtx)
	}
}
