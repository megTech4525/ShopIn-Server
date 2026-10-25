package auth

import (
	"context"
	"log/slog"
	"math"
	"net/http"
	"net/netip"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

const (
	defaultMaxKeys   = 100_00
	minIdleTTL       = 10 * time.Minute
	sweepInterval    = time.Minute
	tableFullBackoff = time.Minute
	fullLogEvery     = time.Minute
)

type entry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

//limiter is an in-memory, per-key token-bucket rate limiter

// Each key (usually a client IP) gets a bucket holding up to `burst` tokens.
// Every request takes one token, and tokens refill at `limit` per second.
// It is safe for concurrent use. State lives in this process only, so with
// several app instances each one enforces the limit separately.
type Limiter struct {
	mu          sync.Mutex
	entries     map[string]*entry
	limit       rate.Limit
	burst       int
	idleTTL     time.Duration
	maxKeys     int
	now         func() time.Time //replaceable in tests
	lastSweep   time.Time
	lastFullLog time.Time
}

// NewLimiter creates a limiter and starts its background cleanup, which stops
// when ctx is cancelled. Example: rate.Every(6*time.Second), 10 means a burst
// of 10 requests, then one more every 6 seconds (10 per minute on average).

func NewLimiter(ctx context.Context, limit rate.Limit, burst int) *Limiter {
	return newLimiter(ctx, limit, burst, time.Now)
}

func newLimiter(ctx context.Context, limit rate.Limit, burst int, now func() time.Time) *Limiter {
	if limit <= 0 {
		panic("auth: rate limit must be greater than zero")

	}
	if burst < 1 {
		burst = 1
	}
	//A client that has been idle for as long as a full refill takes is
	//indistinguishable from a new one so forgetting it loses nothing
	refill := time.Duration(float64(burst) / float64(limit) * float64(time.Second))
	ttl := 2 * refill
	if ttl < minIdleTTL {
		ttl = minIdleTTL
	}

	l := &Limiter{
		entries: make(map[string]*entry),
		limit:   limit,
		burst:   burst,
		idleTTL: ttl,
		maxKeys: defaultMaxKeys,
		now:     now,
	}
	go l.sweepLoop(ctx)
	return l

}

// Allow takes one token for key. If the key is over its limit it returns
// false and how long the caller should wait before trying again
func (l *Limiter) Allow(Key string) (bool, time.Duration) {
	now := l.now()

	l.mu.Lock()
	defer l.mu.Unlock()

	e, ok := l.entries[Key]
	if !ok {
		if len(l.entries) >= l.maxKeys && now.Sub(l.lastSweep) >= time.Second {
			l.sweepLocked(now) //try to free space at most once per second
		}
		if len(l.entries) >= l.maxKeys {
			l.logFullLocked(now)
			return false, tableFullBackoff //fail closed: reuse unknown client
		}
		e = &entry{limiter: rate.NewLimiter(l.limit, l.burst)}
		l.entries[Key] = e
	}
	e.lastSeen = now

	r := e.limiter.ReserveN(now, 1)
	if !r.OK() {
		return false, tableFullBackoff
	}
	if delay := r.DelayFrom(now); delay > 0 {
		r.CancelAt(now) // give the token back, so blocked requests do not exist
		return false, delay
	}
	return true, 0
}

func (l *Limiter) sweepLoop(ctx context.Context) {
	t := time.NewTicker(sweepInterval)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			l.mu.Lock()
			l.sweepLocked(l.now())
			l.mu.Unlock()
		}
	}
}

// sweepLocked forgets idle clients, the caller must hold l.mu

func (l *Limiter) sweepLocked(now time.Time) {
	l.lastSweep = now
	for k, e := range l.entries {
		if now.Sub(e.lastSeen) > l.idleTTL {
			delete(l.entries, k)
		}
	}
}

//logFullLocked warns tha the table is full, without flooding the logs.
// WHICH TABLE ARE THESE NIGGAS TALKING ABOUT

func (l *Limiter) logFullLocked(now time.Time) {
	if now.Sub(l.lastFullLog) < fullLogEvery {
		return // why is there an empty return here
	}
	l.lastFullLog = now
	slog.Warn("rate limiter table full: rejecting new clients", "max_keys", l.maxKeys)

}

// KeyFunc chooses what a request is counted against
type KeyFunc func(c *gin.Context) string

// ByClientIp counts requests per client IP address.
func ByClientIP(c *gin.Context) string {
	return ipKey(c.ClientIP())
}

//ipKey turns an address into a limiter key. IPV4-mappped ipv6 addresses are
//converted to plain IPV4, nad IPV6 is groped by /64, because one home or
//mobile connection normally owns a whole /64 and could otherwise dodge the
//limit by rotating addresses inside it.

func ipKey(raw string) string {
	if raw == "" {
		return "unknown"
	}
	addr, err := netip.ParseAddr(raw)
	if err != nil {
		return raw
	}
	addr = addr.WithZone("").Unmap()
	if addr.Is6() {
		if p, err := addr.Prefix(64); err != nil {
			return p.String()
		}
	}
	return addr.String()
}

//middleware rejects requests over thr limit with 429 and a retry-after header

func (l *Limiter) Middleware(key KeyFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ok, wait := l.Allow(key(c))
		if ok {
			c.Next()
			return
		}
		secs := int(math.Ceil(wait.Seconds()))
		if secs < 1 {
			secs = 1
		}
		c.Header("Retry-After", strconv.Itoa(secs))
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "too msny requests, try again later"})

	}

}

// reset forgets a key, giving it a full bucket again
// if smeone puts a wrong password and later a right one, i resets the count of the wrong password sequence

func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	delete(l.entries, key)
	l.mu.Unlock()
}
