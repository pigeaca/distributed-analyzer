// Package ratelimit provides rate-limiting functionality for HTTP and gRPC services.
package ratelimit

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// HTTPLimiter is a rate limiter for HTTP requests.
type HTTPLimiter struct {
	// limiter is the underlying rate limiter.
	limiter *rate.Limiter

	// mu protects the limiter map.
	mu sync.RWMutex

	// limiters is a map of IP addresses to rate limiters.
	limiters map[string]*rate.Limiter

	// config is the configuration for the rate limiter.
	config HTTPConfig
}

// HTTPConfig holds configuration for the HTTP rate limiter.
type HTTPConfig struct {
	// Rate is the maximum number of requests per second.
	Rate float64

	// Burst is the maximum number of requests that can be made in a burst.
	Burst int

	// IPLookup is a function that extracts the IP address from a request.
	// If nil, the default implementation is used, which gets the IP from
	// the X-Forwarded-For header or the RemoteAddr.
	IPLookup func(*http.Request) string

	// ExcludedPaths is a list of paths that are excluded from rate limiting.
	ExcludedPaths []string

	// Message is the message to return when the rate limit is exceeded.
	Message string

	// StatusCode is the HTTP status code to return when the rate limit is exceeded.
	StatusCode int

	// KeyFunc is a function that extracts a key from a request.
	// If nil, the IP address is used as the key.
	KeyFunc func(*http.Request) string

	// CleanupInterval is the interval at which to clean up old limiters.
	// If 0, cleanup is disabled.
	CleanupInterval time.Duration

	// MaxKeys is the maximum number of keys to track.
	// If 0, there is no limit.
	MaxKeys int
}

// DefaultHTTPConfig returns a default configuration for the HTTP rate limiter.
func DefaultHTTPConfig() HTTPConfig {
	return HTTPConfig{
		Rate:            100,
		Burst:           200,
		Message:         "Too many requests",
		StatusCode:      http.StatusTooManyRequests,
		CleanupInterval: 10 * time.Minute,
		MaxKeys:         10000,
	}
}

// NewHTTPLimiter creates a new HTTP rate limiter with the given configuration.
func NewHTTPLimiter(config HTTPConfig) *HTTPLimiter {
	limiter := &HTTPLimiter{
		limiter:  rate.NewLimiter(rate.Limit(config.Rate), config.Burst),
		limiters: make(map[string]*rate.Limiter),
		config:   config,
	}

	// Start the cleanup goroutine if cleanup is enabled
	if config.CleanupInterval > 0 {
		go limiter.cleanup()
	}

	return limiter
}

// cleanup periodically cleans up old limiters.
func (l *HTTPLimiter) cleanup() {
	ticker := time.NewTicker(l.config.CleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		l.mu.Lock()
		// If we have too many keys, remove the oldest ones
		if l.config.MaxKeys > 0 && len(l.limiters) > l.config.MaxKeys {
			// This is a simple implementation that just clears all limiters
			// A more sophisticated implementation would track the last used time
			// and remove the oldest ones
			l.limiters = make(map[string]*rate.Limiter)
		}
		l.mu.Unlock()
	}
}

// getIPAddress extracts the IP address from a request.
func (l *HTTPLimiter) getIPAddress(r *http.Request) string {
	if l.config.IPLookup != nil {
		return l.config.IPLookup(r)
	}

	// Get IP from X-Forwarded-For header
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		// If no X-Forwarded-For header, use RemoteAddr
		ip = r.RemoteAddr
	}
	return ip
}

// getKey extracts a key from a request.
func (l *HTTPLimiter) getKey(r *http.Request) string {
	if l.config.KeyFunc != nil {
		return l.config.KeyFunc(r)
	}
	return l.getIPAddress(r)
}

// getLimiter gets or creates a rate limiter for the given key.
func (l *HTTPLimiter) getLimiter(key string) *rate.Limiter {
	l.mu.RLock()
	limiter, exists := l.limiters[key]
	l.mu.RUnlock()

	if !exists {
		l.mu.Lock()
		// Check again in case another goroutine created the limiter
		limiter, exists = l.limiters[key]
		if !exists {
			limiter = rate.NewLimiter(rate.Limit(l.config.Rate), l.config.Burst)
			l.limiters[key] = limiter
		}
		l.mu.Unlock()
	}

	return limiter
}

// there isExcluded checks if the request path is excluded from rate limiting.
func (l *HTTPLimiter) isExcluded(path string) bool {
	for _, excludedPath := range l.config.ExcludedPaths {
		if path == excludedPath {
			return true
		}
	}
	return false
}

// Middleware returns a middleware that rate limits HTTP requests.
func (l *HTTPLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the path is excluded
		if l.isExcluded(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		// Get the key for the request
		key := l.getKey(r)

		// Get the rate limiter for the key
		limiter := l.getLimiter(key)

		// Check if the request is allowed
		if !limiter.Allow() {
			// Return a 429 Too Many Requests responses
			http.Error(w, l.config.Message, l.config.StatusCode)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// Handler returns a middleware function that rate limits HTTP requests.
func (l *HTTPLimiter) Handler(next http.Handler) http.Handler {
	return l.Middleware(next)
}

// HandlerFunc returns a middleware function that rate limits HTTP requests.
func (l *HTTPLimiter) HandlerFunc(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l.Middleware(next).ServeHTTP(w, r)
	}
}
