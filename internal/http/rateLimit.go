package http

import (
	"container/list"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

const rateLimitMaxEntries = 10_000

type lruEntry struct {
	limiter *rate.Limiter
	elem    *list.Element
}

type ipRateLimiter struct {
	mu      sync.Mutex
	entries map[string]*lruEntry
	order   *list.List
	r       rate.Limit
	burst   int
}

func newIPRateLimiter(count int, window time.Duration) *ipRateLimiter {
	return &ipRateLimiter{
		entries: make(map[string]*lruEntry, rateLimitMaxEntries),
		order:   list.New(),
		r:       rate.Limit(float64(count) / window.Seconds()),
		burst:   count,
	}
}

func (l *ipRateLimiter) get(ip string) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()

	if e, ok := l.entries[ip]; ok {
		l.order.MoveToFront(e.elem)
		return e.limiter
	}

	if len(l.entries) >= rateLimitMaxEntries {
		back := l.order.Back()
		if back != nil {
			l.order.Remove(back)
			delete(l.entries, back.Value.(string))
		}
	}

	lim := rate.NewLimiter(l.r, l.burst)
	elem := l.order.PushFront(ip)
	l.entries[ip] = &lruEntry{limiter: lim, elem: elem}
	return lim
}

func realIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}

	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		if i := strings.IndexByte(fwd, ','); i >= 0 {
			return strings.TrimSpace(fwd[:i])
		}
		return strings.TrimSpace(fwd)
	}

	return host
}

// RateLimit returns a middleware that allows at most count requests per window per IP.
// Per-IP limiters are kept in an LRU capped at rateLimitMaxEntries to bound memory.
func rateLimit(count int, window time.Duration) gin.HandlerFunc {
	limiter := newIPRateLimiter(count, window)

	return func(c *gin.Context) {
		ip := realIP(c.Request)
		if !limiter.get(ip).Allow() {
			log.Printf("Error: Rate limit exceeded for IP %s\n", ip)
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}
		c.Next()
	}
}
