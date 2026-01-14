# httpx

<p>
<a href="https://github.com/goapt/httpx/actions"><img src="https://github.com/goapt/httpx/workflows/build/badge.svg" alt="Build Status"></a>
<a href="https://codecov.io/gh/goapt/httpx"><img src="https://codecov.io/gh/goapt/httpx/branch/master/graph/badge.svg" alt="codecov"></a>
<a href="https://goreportcard.com/report/github.com/goapt/httpx"><img src="https://goreportcard.com/badge/github.com/goapt/httpx" alt="Go Report Card
"></a>
<a href="https://godoc.org/github.com/goapt/httpx"><img src="https://godoc.org/github.com/goapt/httpx?status.svg" alt="GoDoc"></a>
<a href="https://opensource.org/licenses/mit-license.php" rel="nofollow"><img src="https://badges.frapsoft.com/os/mit/mit.svg?v=103"></a>
</p>

Simple, composable HTTP client with middleware for Go.

## Features

- Option-based client configuration (timeout, transport, middlewares)
- Composable, chainable middlewares (logging, debugging, tracing, mocking)
- Pluggable transport layer (custom TLS, connection tuning)
- Friendly for unit tests via an HTTP mock middleware

## Installation

```bash
go get github.com/goapt/httpx
```

## Quick Start

```go
package main

import (
	"net/http"
	"time"

	"github.com/goapt/httpx"
)

func main() {
	client := httpx.NewClient(
		httpx.WithTimeout(5*time.Second),
		httpx.WithMiddleware(httpx.Debug()),
	)

	req, _ := http.NewRequest(http.MethodGet, "https://httpbin.org/get", nil)
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
}
```

## Middlewares

The middleware model wraps a `http.RoundTripper` in a chain, allowing cross-cutting features without coupling to request logic.

- AccessLog: structured access logging
- Debug: human-readable request/response dump for local development
- Trace: OpenTelemetry tracing for HTTP client requests
- Mock: programmable mock responses for tests

### AccessLog

```go
import (
	"bytes"
	"net/http"

	"github.com/goapt/httpx"
	"github.com/goapt/logger"
)

func exampleAccessLog() {
	l := logger.New(&logger.Config{Mode: logger.ModeStd})
	client := httpx.NewClient(httpx.WithMiddleware(httpx.AccessLog(l)))

	req, _ := http.NewRequest(http.MethodPost, "https://httpbin.org/anything", bytes.NewReader([]byte(`{"k":"v"}`)))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, _ := client.Do(req)
	defer resp.Body.Close()
}
```


### Debug

```go
client := httpx.NewClient(httpx.WithMiddleware(httpx.Debug()))
resp, _ := client.Get("https://httpbin.org/json")
defer resp.Body.Close()
```

### Trace (OpenTelemetry)

```go
client := httpx.NewClient(httpx.WithMiddleware(httpx.Trace()))
resp, _ := client.Get("https://httpbin.org/get")
defer resp.Body.Close()
```

### HTTP Mock (for tests)

```go
import (
	"bytes"
	"errors"

	"github.com/goapt/httpx"
)

func exampleMock() {
	suites := []httpx.MockSuite{
		{URI: "/get", ResponseBody: "ok"},
		{URI: "/user/id/.*", ResponseBody: "user"},
		{URI: "/find\\?id=.*", ResponseBody: "find"},
		{URI: "/bodymatch", MatchBody: map[string]any{"user_id": 1}, ResponseBody: "body-ok"},
		{URI: "/query", MatchQuery: map[string]any{"name": "test"}, ResponseBody: "query-ok"},
		{URI: "/error", Error: errors.New("mock error")},
	}

	client := httpx.NewClient(httpx.WithMiddleware(httpx.Mock(suites)))
	body := bytes.NewBufferString(`{"user_id":1}`)
	resp, _ := client.Post("/bodymatch", "application/json; charset=utf-8", body)
	defer resp.Body.Close()
}
```

## Custom Middleware

```go
import (
	"net/http"

	"github.com/goapt/httpx"
)

func exampleCustomMW() {
	logMW := func(next http.RoundTripper) http.RoundTripper {
		return httpx.RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			resp, err := next.RoundTrip(req)
			return resp, err
		})
	}

	client := httpx.NewClient(httpx.WithMiddleware(logMW))
	_ = client
}
```

## Custom Transport

```go
import (
	"crypto/tls"
	"net/http"

	"github.com/goapt/httpx"
)

func exampleTransport() {
	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.TLSClientConfig = &tls.Config{
		CipherSuites: []uint16{tls.TLS_AES_128_GCM_SHA256, tls.TLS_AES_256_GCM_SHA384},
	}

	client := httpx.NewClient(httpx.WithTransport(tr))
	resp, _ := client.Get("https://httpbin.org/json")
	defer resp.Body.Close()
}
```
