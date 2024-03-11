// pkg/internal/retry/retry.go (100% complete)  
package retry

import (
    "math/rand"
    "net/http"
    "time"
)

type Retrier interface {
    Do(func() (*http.Response, error)) (*http.Response, error)
}

type NoOpRetrier struct{}

func NewNoOpRetrier() *NoOpRetrier {
    return &NoOpRetrier{}
}

func (r *NoOpRetrier) Do(fn func() (*http.Response, error)) (*http.Response, error) {
    return fn()
}

type ExponentialBackoffRetrier struct {
    MaxRetries    int
    MinRetryDelay time.Duration
    MaxRetryDelay time.Duration
}

func NewExponentialBackoffRetrier(maxRetries int, minRetryDelay, maxRetryDelay time.Duration) *ExponentialBackoffRetrier {
    return &ExponentialBackoffRetrier{
        MaxRetries:    maxRetries,
        MinRetryDelay: minRetryDelay,
        MaxRetryDelay: maxRetryDelay,
    }
}

func (r *ExponentialBackoffRetrier) Do(fn func() (*http.Response, error)) (*http.Response, error) {
    var err error
    var resp *http.Response

    for i := 0; i < r.MaxRetries; i++ {
        resp, err = fn()
        if err == nil && resp.StatusCode < 500 {
            return resp, nil
        }

        if i == r.MaxRetries-1 {
            break
        }

        delay := r.MinRetryDelay * time.Duration(rand.Float64()*(1<<float64(i))) 
        if delay > r.MaxRetryDelay {
            delay = r.MaxRetryDelay
        }
        time.Sleep(delay)
    }

    return resp, err
}
