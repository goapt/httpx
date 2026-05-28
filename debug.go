package httpx

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/goapt/logger"
	"github.com/goapt/logger/sloghttp"
)

func Debug() Middleware {
	return func(rt http.RoundTripper) http.RoundTripper {
		logger := logger.New(logger.NewJSONHandler(os.Stdout, logger.WithLevel(slog.LevelDebug)))

		return sloghttp.NewRoundTripper(logger, rt, sloghttp.Config{
			Level:              slog.LevelDebug,
			WithUserAgent:      true,
			WithRequestBody:    true,
			WithRequestHeader:  true,
			WithResponseBody:   true,
			WithResponseHeader: true,
		})
	}
}
