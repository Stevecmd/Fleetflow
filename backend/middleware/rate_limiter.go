package middleware

import (
	"net/http"
	"strings"
	"sync"

	"golang.org/x/time/rate"
)

// Global limiter instance
// 100 requests per second per IP with a burst of 200
var globalLimiter = NewIPRateLimiter(rate.Limit(100), 200)

type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit
	b   int
}

// NewIPRateLimiter creates a new IPRateLimiter instance.
// It returns a new pointer to an IPRateLimiter struct.
// The struct contains a map of IP addresses to rate limiters,
// a read-write mutex to protect the map, the rate limit, and the burst size.
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	return &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}
}

// GetLimiter retrieves the rate limiter associated with the given IP address.
// If no rate limiter exists for the IP, a new one is created with the specified rate and burst size.
// This function is thread-safe and locks the map for writing to ensure consistency.

func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(i.r, i.b)
		i.ips[ip] = limiter
	}

	return limiter
}

// getIP determines the IP address of the request.
// It first checks the X-Forwarded-For header and returns the first IP.
// If the header is not set, it checks the X-Real-IP header.
// If the header is not set, it falls back to the RemoteAddr.
func getIP(r *http.Request) string {
	// Check X-Forwarded-For header
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		return strings.Split(forwarded, ",")[0]
	}

	// Check X-Real-IP header
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}

	// Fall back to RemoteAddr
	return strings.Split(r.RemoteAddr, ":")[0]
}

// RateLimitMiddleware is a middleware function that applies rate limiting
// to incoming HTTP requests based on the client's IP address. It retrieves
// the IP address from the request, checks the associated rate limiter, and
// allows the request to proceed if the rate limit has not been exceeded.
// If the request exceeds the rate limit, a 429 Too Many Requests response
// is returned with a Retry-After header indicating when the client can retry.

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getIP(r)
		limiter := globalLimiter.GetLimiter(ip)

		if !limiter.Allow() {
			w.Header().Set("Retry-After", "1")
			http.Error(w, "Too many requests. Please try again later.", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
