# httpthrottle

[![GoDoc](https://godoc.org/github.com/jybp/httpthrottle?status.svg)](https://godoc.org/github.com/jybp/httpthrottle)

Package httpthrottle provides a http.RoundTripper to rate limit HTTP requests.

## Usage

```go
package example

import (
    "errors"
    "net/http"
    "github.com/jybp/httpthrottle"
    "golang.org/x/time/rate"
)

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
    if errors.Is(err, httpthrottle.ErrQuotaExceeded) {
        // Handle err.
    }
    if err != nil {
        // Handle err.
    }
    _ = resp // Do something with resp.
}
```
