package retailcrm

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSingleKeyLimiter(t *testing.T) {
	limiter := NewSingleKeyLimiter()
	require.NotNil(t, limiter, "NewSingleKeyLimiter returned nil")

	skl, ok := limiter.(*SingleKeyLimiter)
	require.True(t, ok, "NewSingleKeyLimiter did not return *SingleKeyLimiter")
	assert.NotNil(t, skl.regularLimiter, "regularLimiter should not be nil")
	assert.NotNil(t, skl.telephonyLimiter, "telephonyLimiter should not be nil")
}

func TestSingleKeyLimiter_Limit_RegularRoute(t *testing.T) {
	limiter := NewSingleKeyLimiter().(*SingleKeyLimiter)
	ctx := context.Background()

	start := time.Now()
	require.NoError(t, limiter.Limit(ctx, "/api/orders", "test-key"))
	require.WithinDuration(t, start, time.Now(), 10*time.Millisecond, "First request should not be delayed")

	start = time.Now()
	require.NoError(t, limiter.Limit(ctx, "/api/customers", "test-key"))
	elapsed := time.Since(start)

	expectedDelay := regularDelay
	tolerance := 50 * time.Millisecond
	assert.GreaterOrEqual(t, elapsed, expectedDelay-tolerance,
		"Request completed too quickly: %v (expected ~%v)", elapsed, expectedDelay)
}

func TestSingleKeyLimiter_Limit_TelephonyRoute(t *testing.T) {
	limiter := NewSingleKeyLimiter().(*SingleKeyLimiter)
	ctx := context.Background()

	start := time.Now()
	require.NoError(t, limiter.Limit(ctx, "/telephony/calls", "test-key"))
	require.WithinDuration(t, start, time.Now(), 10*time.Millisecond,
		"First telephony request should not be delayed")

	start = time.Now()
	require.NoError(t, limiter.Limit(ctx, "/telephony/status", "test-key"))
	elapsed := time.Since(start)

	expectedDelay := telephonyDelay
	tolerance := 20 * time.Millisecond
	assert.GreaterOrEqual(t, elapsed, expectedDelay-tolerance,
		"Telephony request completed too quickly: %v (expected ~%v)", elapsed, expectedDelay)
}

func TestSingleKeyLimiter_Limit_SeparateLimiters(t *testing.T) {
	limiter := NewSingleKeyLimiter().(*SingleKeyLimiter)
	ctx := context.Background()

	require.NoError(t, limiter.Limit(ctx, "/api/orders", "test-key"))

	start := time.Now()
	require.NoError(t, limiter.Limit(ctx, "/telephony/calls", "test-key"))
	elapsed := time.Since(start)

	assert.Less(t, elapsed, 50*time.Millisecond,
		"Telephony request was delayed by regular limiter: %v", elapsed)
}

func TestSingleKeyLimiter_Limit_ContextCancellation(t *testing.T) {
	limiter := NewSingleKeyLimiter().(*SingleKeyLimiter)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	require.NoError(t, limiter.Limit(context.Background(), "/api/test", "test-key"))

	err := limiter.Limit(ctx, "/api/test", "test-key")
	assert.ErrorIs(t, err, context.Canceled, "Expected context.Canceled error")
}

func TestSingleKeyLimiter_Limit_ContextTimeout(t *testing.T) {
	limiter := NewSingleKeyLimiter().(*SingleKeyLimiter)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	require.NoError(t, limiter.Limit(context.Background(), "/api/test", "test-key"))

	time.Sleep(5 * time.Millisecond)

	err := limiter.Limit(ctx, "/api/test", "test-key")
	assert.Error(t, err, "Expected an error with timed out context")
	assert.ErrorIs(t, err, context.DeadlineExceeded,
		"Expected error to be or contain context.DeadlineExceeded")
}

func TestSingleKeyLimiter_Limit_URIPrefixMatching(t *testing.T) {
	tests := []struct {
		name               string
		uri                string
		shouldUseTelephony bool
	}{
		{"Exact telephony prefix", "/telephony", true},
		{"Telephony with path", "/telephony/calls/123", true},
		{"Regular route", "/api/orders", false},
		{"Similar but not telephony", "/telephone", false},
		{"Empty URI", "", false},
		{"Root path", "/", false},
		{"Telephony uppercase", "/TELEPHONY/calls", false}, // Case sensitive
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limiter := NewSingleKeyLimiter().(*SingleKeyLimiter)
			ctx := context.Background()

			require.NoError(t, limiter.Limit(ctx, tt.uri, "test-key"))

			start := time.Now()
			require.NoError(t, limiter.Limit(ctx, tt.uri, "test-key"))
			elapsed := time.Since(start)

			if tt.shouldUseTelephony {
				assert.Less(t, elapsed, 40*time.Millisecond,
					"Expected telephony rate limit (~25ms), got %v", elapsed)
			} else {
				assert.GreaterOrEqual(t, elapsed, 50*time.Millisecond,
					"Expected regular rate limit (~100ms), got %v", elapsed)
			}
		})
	}
}

func TestSingleKeyLimiter_Limit_KeyIgnored(t *testing.T) {
	limiter := NewSingleKeyLimiter().(*SingleKeyLimiter)
	ctx := context.Background()

	require.NoError(t, limiter.Limit(ctx, "/api/test", "key1"))

	start := time.Now()
	require.NoError(t, limiter.Limit(ctx, "/api/test", "key2"))
	elapsed := time.Since(start)

	assert.GreaterOrEqual(t, elapsed, 50*time.Millisecond,
		"Different keys should still share the same rate limit")
}

func TestSingleKeyLimiter_Limit_Concurrent(t *testing.T) {
	limiter := NewSingleKeyLimiter().(*SingleKeyLimiter)
	ctx := context.Background()

	results := make(chan time.Duration, 5)
	start := time.Now()

	for i := 0; i < 5; i++ {
		go func() {
			reqStart := time.Now()
			err := limiter.Limit(ctx, "/api/test", "test-key")
			require.NoError(t, err, "Request should not fail")
			results <- time.Since(reqStart)
		}()
	}

	for i := 0; i < 5; i++ {
		<-results
	}
	totalElapsed := time.Since(start)

	expectedMin := 350 * time.Millisecond
	expectedMax := 550 * time.Millisecond

	assert.GreaterOrEqual(t, totalElapsed, expectedMin,
		"Concurrent requests completed too quickly: %v", totalElapsed)
	assert.LessOrEqual(t, totalElapsed, expectedMax,
		"Concurrent requests took too long: %v", totalElapsed)
}

func TestSingleKeyLimiter_Limit_RateEnforcement(t *testing.T) {
	limiter := NewSingleKeyLimiter().(*SingleKeyLimiter)
	ctx := context.Background()

	count := 5
	start := time.Now()

	for i := 0; i < count; i++ {
		require.NoError(t, limiter.Limit(ctx, "/api/orders", "key1"))
	}

	elapsed := time.Since(start)
	expectedMin := regularDelay * time.Duration(count-1)

	assert.GreaterOrEqual(t, elapsed, expectedMin,
		"Rate not enforced: expected at least %v, got %v", expectedMin, elapsed)
}

func TestSingleKeyLimiter_Limit_TelephonyRateEnforcement(t *testing.T) {
	limiter := NewSingleKeyLimiter().(*SingleKeyLimiter)
	ctx := context.Background()

	count := 5
	start := time.Now()

	for i := 0; i < count; i++ {
		require.NoError(t, limiter.Limit(ctx, "/telephony/call", "key1"))
	}

	elapsed := time.Since(start)
	expectedMin := telephonyDelay * time.Duration(count-1)

	assert.GreaterOrEqual(t, elapsed, expectedMin,
		"Telephony rate not enforced: expected at least %v, got %v", expectedMin, elapsed)
}

func TestSingleKeyLimiter_Limit_CanReturnAnError(t *testing.T) {
	limiter := NewSingleKeyLimiter().(*SingleKeyLimiter)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	assert.ErrorIs(t, limiter.Limit(ctx, "/telephony/call", "key1"), context.Canceled)
}

func TestConstants(t *testing.T) {
	assert.Equal(t, 10, RPS, "RPS should be 10")
	assert.Equal(t, 40, TelephonyRPS, "TelephonyRPS should be 40")
	assert.Equal(t, time.Second/RPS, regularDelay, "regularDelay should be 100ms")
	assert.Equal(t, time.Second/TelephonyRPS, telephonyDelay, "telephonyDelay should be 25ms")
}

func BenchmarkSingleKeyLimiter_Limit_Regular(b *testing.B) {
	limiter := NewSingleKeyLimiter()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = limiter.Limit(ctx, "/api/test", "test-key")
	}
}

func BenchmarkSingleKeyLimiter_Limit_Telephony(b *testing.B) {
	limiter := NewSingleKeyLimiter()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = limiter.Limit(ctx, "/telephony/test", "test-key")
	}
}
