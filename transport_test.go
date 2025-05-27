package httpthrottle_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jybp/httpthrottle"
	"golang.org/x/time/rate"
)

func TestTransport(t *testing.T) {
	l := rate.NewLimiter(rate.Every(time.Millisecond*10), 3)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !l.Allow() {
			w.WriteHeader(http.StatusTooManyRequests)
		}
	}))
	client := &http.Client{Transport: &httpthrottle.Transport{Limiter: rate.NewLimiter(rate.Every(time.Millisecond*10), 2)}}
	assertStatusFn := func(expected int) {
		resp, err := client.Get(srv.URL)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != expected {
			t.Fatalf("expected:%d\tactual:%d\n", expected, resp.StatusCode)
		}
	}
	assertStatusFn(http.StatusOK)
	assertStatusFn(http.StatusOK)
	assertStatusFn(http.StatusOK)
	assertStatusFn(http.StatusOK)
}

func TestMultiLimiter(t *testing.T) {
	l := httpthrottle.MultiLimiters(
		httpthrottle.NewQuota(time.Second, 101),
		rate.NewLimiter(rate.Every(time.Millisecond), 1),
	)
	start := time.Now()
	for i := 0; i < 101; i++ {
		if err := l.Wait(context.Background()); err != nil {
			t.Fatal(err)
		}
	}
	if elapsed := time.Since(start); elapsed < time.Millisecond*100 {
		t.Fatalf("101 wait took %v", elapsed)
	}
	if err := l.Wait(context.Background()); err != httpthrottle.ErrQuotaExceeded {
		t.Fatal(err)
	}
}

// Wait won't block if a Quota is reached.
func TestMultiLimiterQuota(t *testing.T) {
	l := httpthrottle.MultiLimiters(
		rate.NewLimiter(1, 1),
		httpthrottle.NewQuota(time.Second, 1),
	)
	start := time.Now()
	if err := l.Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := l.Wait(context.Background()); err != httpthrottle.ErrQuotaExceeded {
		t.Fatal(err)
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Fatalf("100 wait took %v", elapsed)
	}
}

func Example() {
	client := &http.Client{
		Transport: httpthrottle.Default(
			// Returns ErrQuotaExceeded if more than 36000 requests occured within an hour.
			httpthrottle.NewQuota(time.Hour, 36000),
			// Blocks to never exceed 99 requests per second.
			rate.NewLimiter(99, 1),
		),
	}
	resp, err := client.Get("https://golang.org/")
	if err == httpthrottle.ErrQuotaExceeded {
		// Handle err.
	}
	if err != nil {
		// Handle err.
	}
	_ = resp // Do something with resp.
}
