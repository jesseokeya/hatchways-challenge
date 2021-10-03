package api

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"
)

type StructuredLogger struct {
	*zerolog.Logger
	Debug bool
}

func (z *StructuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	logger := z.With()
	logger = logger.Timestamp()

	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		logger = logger.Str("req_id", reqID)
	}

	scheme := "http"
	host := r.Host

	if r.TLS != nil {
		scheme = "https"
	}
	if val := r.Header.Get("X-Forwarded-Host"); val != "" {
		host = val
	}

	fields := map[string]interface{}{
		"remote_ip":  r.RemoteAddr,
		"host":       r.Host,
		"uri":        fmt.Sprintf("%s://%s%s", scheme, host, r.RequestURI),
		"proto":      r.Proto,
		"method":     r.Method,
		"user_agent": r.Header.Get("User-Agent"),
		"bytes_in":   r.Header.Get("Content-Length"),
	}
	logger = logger.Fields(fields)

	sublogger := logger.Logger()
	sublogger.Info().Fields(fields).Msg("request started")

	return &StructuredLoggerEntry{Logger: logger}
}

type StructuredLoggerEntry struct {
	Logger zerolog.Context
}

func (entry *StructuredLoggerEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	entry.Logger = entry.Logger.Int("status", status)
	entry.Logger = entry.Logger.Int("bytes_out", bytes)
	entry.Logger = entry.Logger.Float64("resp_elapsed_ms", float64(elapsed.Nanoseconds())/1000000.0)

	logger := entry.Logger.Logger()
	if status > 399 && status < 500 {
		logger.Warn().Msg("invalid request")
	} else if status >= 500 {
		logger.Error().Msg("internal error")
	} else {
		logger.Info().Msg("request complete")
	}
}

func (entry *StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	entry.Logger = entry.Logger.Fields(map[string]interface{}{
		"stack": string(stack),
		"panic": fmt.Sprintf("%+v", v),
	})
}

// Run secondary functions based on severity.
type SeverityHook struct {
	AlertFn func(level zerolog.Level, msg string)
}

func (h SeverityHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if level >= zerolog.WarnLevel {
		if h.AlertFn != nil {
			buf := &bytes.Buffer{}
			logger := zerolog.New(buf)
			logger.Error().Dict("log_line", e).Msg(msg)
			h.AlertFn(level, buf.String())
		}
	}
}
