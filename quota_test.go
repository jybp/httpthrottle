package httpthrottle_test

import (
	"context"
	"testing"
	"time"

	"github.com/jybp/httpthrottle"
)

func TestQuota(t *testing.T) {
	ctx := context.Background()
	q := httpthrottle.NewQuota(time.Millisecond*10, 2)
	assertFn := func() {
		if err := q.Wait(ctx); err != nil {
			t.Fatal(err)
		}
		if err := q.Wait(ctx); err != nil {
			t.Fatal(err)
		}
		if err := q.Wait(ctx); err != httpthrottle.ErrQuotaExceeded {
			t.Fatal("ErrQuotaExceeded expected", err)
		}
		if err := q.Wait(ctx); err != httpthrottle.ErrQuotaExceeded {
			t.Fatal("ErrQuotaExceeded expected", err)
		}
	}
	assertFn()
	time.Sleep(time.Millisecond * 10)
	assertFn()
	time.Sleep(time.Millisecond * 20)
	assertFn()
}
