package httpx

import (
	"log/slog"
	"net/http"

	"github.com/goapt/logger/sloghttp"
)

func AccessLog(logger *slog.Logger) Middleware {
	return func(rt http.RoundTripper) http.RoundTripper {
		return sloghttp.NewRoundTripper(logger, rt, sloghttp.Config{
			Level:              slog.LevelInfo,
			WithUserAgent:      true,
			WithRequestBody:    true,
			WithRequestHeader:  true,
			WithResponseBody:   true,
			WithResponseHeader: true,
		})
	}
}
