package auth

import (
	"context"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func newTestLimiter(t *testing.T, limit rate.Limit, burst int) (*Limiter, *time.Time) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	clock := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	l := newLimiter(ctx, limit, burst, func() time.Time { return clock })
	return l, &clock
}
func TestResetGivesAFullBucket(t *testing.T) {
	l, _ := newTestLimiter(t, rate.Every(time.Minute), 2)

	l.Allow("a")
	l.Allow("a")
	if ok, _ := l.Allow("a"); ok {
		t.Fatal("third request should be blocked")
	}

	l.Reset("a")
	if ok, _ := l.Allow("a"); !ok {
		t.Error("after Reset the key should have a full bucket again")
	}
}
