package retailcrm

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type contextTestTransport func(*http.Request) (*http.Response, error)

func (f contextTestTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

type contextTestBody struct {
	io.Reader
	closed int
}

func (b *contextTestBody) Close() error { b.closed++; return nil }

func TestWithContextClonesClient(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	original := New("https://crm.test", "key").EnableRateLimiter(3)
	original.Debug = true
	clone := original.WithContext(ctx)
	require.NotSame(t, original, clone)
	assert.Nil(t, original.ctx)
	assert.Same(t, ctx, clone.ctx)
	assert.Equal(t, original.URL, clone.URL)
	assert.Equal(t, original.Key, clone.Key)
	assert.Equal(t, original.Debug, clone.Debug)
	assert.Same(t, original.httpClient, clone.httpClient)
	assert.Same(t, original.limiter, clone.limiter)
	assert.Equal(t, uint(3), clone.maxAttempts)
	// Cloning must not copy a locked mutex or mutate the original's settings.
	clone.EnableCustomRateLimiter(nil, 0)
	assert.NotNil(t, original.limiter)
	assert.Equal(t, uint(3), original.maxAttempts)
	cancel()
	assert.NoError(t, original.requestContext().Err())
	assert.ErrorIs(t, clone.requestContext().Err(), context.Canceled)
	reset := clone.WithContext(nil)
	assert.NoError(t, reset.requestContext().Err())
	assert.ErrorIs(t, clone.requestContext().Err(), context.Canceled)
}

func TestWithContextConcurrentClones(t *testing.T) {
	type contextKey struct{}
	original := New("http://crm.test", "key").EnableRateLimiter(1)
	original.WithHTTPClient(&http.Client{Transport: contextTestTransport(func(r *http.Request) (*http.Response, error) {
		assert.Equal(t, r.URL.Query().Get("id"), r.Context().Value(contextKey{}))
		return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(`{"success":true}`)), Request: r}, nil
	})})
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			ctx := context.WithValue(context.Background(), contextKey{}, id)
			_, _, err := original.WithContext(ctx).GetRequest("/test?id=" + id)
			assert.NoError(t, err)
		}(fmt.Sprint(i))
	}
	wg.Wait()
	assert.Nil(t, original.ctx)
}

var contextTestCalls = []struct {
	name string
	call func(*Client) error
}{
	{"GET", func(c *Client) error { _, _, err := c.GetRequest("/test"); return err }},
	{"POST", func(c *Client) error { _, _, err := c.PostRequest("/test", url.Values{}); return err }},
	{"plate", func(c *Client) error {
		body, _, err := c.GetOrderPlate(ByID, "1", "main", 1)
		if body != nil {
			body.Close()
		}
		return err
	}},
}

func TestWithContextCancelsHTTPRequests(t *testing.T) {
	for _, call := range contextTestCalls {
		t.Run(call.name, func(t *testing.T) {
			for _, duringBody := range []bool{false, true} {
				t.Run(fmt.Sprint("duringBody=", duringBody), func(t *testing.T) {
					entered := make(chan struct{})
					server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						if duringBody {
							w.WriteHeader(200)
							w.(http.Flusher).Flush()
						}
						close(entered)
						<-r.Context().Done()
					}))
					defer server.Close()
					original := New(server.URL, "key")
					ctx, cancel := context.WithCancel(context.Background())
					defer cancel()
					done := make(chan error, 1)
					go func() {
						// Read streaming responses too: the context must cover their lifetime.
						if call.name == "plate" && duringBody {
							body, _, err := original.WithContext(ctx).GetOrderPlate(ByID, "1", "main", 1)
							if err == nil {
								_, err = io.ReadAll(body)
								body.Close()
							}
							done <- err
						} else {
							done <- call.call(original.WithContext(ctx))
						}
					}()
					select {
					case <-entered:
					case <-time.After(time.Second):
						t.Fatal("request did not start")
					}
					cancel()
					select {
					case err := <-done:
						require.ErrorIs(t, err, context.Canceled)
					case <-time.After(time.Second):
						t.Fatal("request was not canceled")
					}
					assert.Nil(t, original.ctx)
				})
			}
		})
	}
}

func TestWithContextDeadlineAndOriginalRemainIndependent(t *testing.T) {
	original := New("http://crm.test", "key").WithHTTPClient(&http.Client{Transport: contextTestTransport(func(r *http.Request) (*http.Response, error) {
		if r.Context().Done() != nil {
			<-r.Context().Done()
			return nil, r.Context().Err()
		}
		return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(`{"success":true}`)), Request: r}, nil
	})})
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	_, _, err := original.WithContext(ctx).Customers(CustomersRequest{})
	require.ErrorIs(t, err, context.DeadlineExceeded)
	response, _, err := original.Customers(CustomersRequest{})
	require.NoError(t, err)
	assert.True(t, response.Success)
}

func TestWithContextCancelsRetryWait(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var requests int
	body := &contextTestBody{Reader: strings.NewReader(`{"success":false}`)}
	c := New("http://crm.test", "key").EnableRateLimiter(0).WithHTTPClient(&http.Client{Transport: contextTestTransport(func(r *http.Request) (*http.Response, error) {
		requests++
		cancel()
		return &http.Response{StatusCode: 503, Header: make(http.Header), Body: body, Request: r}, nil
	})}).WithContext(ctx)
	started := time.Now()
	_, _, err := c.GetRequest("/test")
	require.ErrorIs(t, err, context.Canceled)
	assert.Less(t, time.Since(started), regularDelay)
	assert.Equal(t, 1, requests)
	assert.Equal(t, 1, body.closed)
}

func TestHTTPResponseBodiesClosed(t *testing.T) {
	for _, call := range contextTestCalls {
		for _, status := range []int{200, 400, 429, 500, 502, 503} {
			t.Run(fmt.Sprintf("%s/%d", call.name, status), func(t *testing.T) {
				body := &contextTestBody{Reader: strings.NewReader(`{"success":true}`)}
				c := New("http://crm.test", "key").EnableRateLimiter(1).WithHTTPClient(&http.Client{Transport: contextTestTransport(func(r *http.Request) (*http.Response, error) {
					return &http.Response{StatusCode: status, Header: make(http.Header), Body: body, Request: r}, nil
				})})
				_ = call.call(c)
				assert.Equal(t, 1, body.closed)
			})
		}
	}
}

func TestPlateBodyClosedBeforeRetry(t *testing.T) {
	var bodies []*contextTestBody
	c := New("http://crm.test", "key").EnableRateLimiter(2).WithHTTPClient(&http.Client{Transport: contextTestTransport(func(r *http.Request) (*http.Response, error) {
		if len(bodies) > 0 {
			require.Equal(t, 1, bodies[0].closed)
		}
		body := &contextTestBody{Reader: strings.NewReader("plate")}
		bodies = append(bodies, body)
		status := 503
		if len(bodies) == 2 {
			status = 200
		}
		return &http.Response{StatusCode: status, Header: make(http.Header), Body: body, Request: r}, nil
	})})
	body, status, err := c.GetOrderPlate(ByID, "1", "main", 1)
	require.NoError(t, err)
	require.Equal(t, 200, status)
	require.Len(t, bodies, 2)
	assert.Equal(t, 0, bodies[1].closed, "successful stream belongs to the caller")
	_, err = io.ReadAll(body)
	require.NoError(t, err)
	require.NoError(t, body.Close())
	assert.Equal(t, 1, bodies[1].closed)
}

func TestHTTPResponseBodyClosedOnReadError(t *testing.T) {
	for _, call := range contextTestCalls[:2] {
		t.Run(call.name, func(t *testing.T) {
			body := &contextTestBody{Reader: contextErrorReader{}}
			c := New("http://crm.test", "key").WithHTTPClient(&http.Client{Transport: contextTestTransport(func(r *http.Request) (*http.Response, error) {
				return &http.Response{StatusCode: 200, Header: make(http.Header), Body: body, Request: r}, nil
			})})
			require.Error(t, call.call(c))
			assert.Equal(t, 1, body.closed)
		})
	}
}

type contextErrorReader struct{}

func (contextErrorReader) Read([]byte) (int, error) { return 0, errors.New("read failed") }
