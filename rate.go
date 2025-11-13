package retailcrm

import (
	"context"
	"net/http"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

const (
	// RPS for the single key.
	RPS = 10
	// TelephonyRPS is used for telephony routes instead of RPS.
	TelephonyRPS   = 40
	regularDelay   = time.Second / RPS          // Delay between regular requests.
	telephonyDelay = time.Second / TelephonyRPS // Delay between telephony requests.
)

// Limiter describes basic rate limiter.
type Limiter interface {
	// Limit the request. Returned error will be received from API method which is being called.
	Limit(ctx context.Context, uri, key string) error
}

// ResponseAware can be used to make rate limiter which will be able to read every response.
// Body is not guaranteed for this method and may be already read.
type ResponseAware interface {
	ProcessResponse(resp *http.Response)
}

// RateLimiter is the alias for SingleKeyLimiter.
// Deprecated: use SingleKeyLimiter and NewSingleKeyLimiter instead.
type RateLimiter SingleKeyLimiter

// SingleKeyLimiter manages API request rates to prevent hitting rate limits. Works for only one key.
type SingleKeyLimiter struct {
	regularLimiter   *rate.Limiter
	telephonyLimiter *rate.Limiter
}

// NewSingleKeyLimiter instantiates new SingleKeyLimiter.
func NewSingleKeyLimiter() Limiter {
	return &SingleKeyLimiter{
		regularLimiter:   rate.NewLimiter(rate.Limit(RPS), 1),
		telephonyLimiter: rate.NewLimiter(rate.Limit(TelephonyRPS), 1),
	}
}

// Limit the request.
func (r *SingleKeyLimiter) Limit(ctx context.Context, uri, _ string) error {
	if strings.HasPrefix(uri, "/telephony") {
		return r.telephonyLimiter.Wait(ctx)
	}

	return r.regularLimiter.Wait(ctx)
}
