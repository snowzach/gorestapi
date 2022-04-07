package logger

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"time"

	"github.com/blendle/zapdriver"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Level        zapcore.Level `conf:"level"`
	RequestBody  bool          `conf:"request_body"`
	ResponseBody bool          `conf:"response_body"`
	IgnorePaths  []string      `conf:"ignore_paths"`
}

func LoggerStandardMiddleware(logger *zap.Logger, config Config) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Check if the prefix should be ignored
			for _, prefix := range config.IgnorePaths {
				if strings.HasPrefix(r.RequestURI, prefix) {
					next.ServeHTTP(w, r)
					return
				}
			}

			start := time.Now()

			// See if we're already using a wrapped response writer and if not make one.
			ww, ok := w.(middleware.WrapResponseWriter)
			if !ok {
				ww = middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			}

			// If we should log bodies, setup buffer to capture
			var responseBody *bytes.Buffer
			if config.ResponseBody {
				responseBody = new(bytes.Buffer)
				ww.Tee(responseBody)
			}

			next.ServeHTTP(ww, r)

			// If the remote IP is being proxied, use the real IP
			remoteIP := r.Header.Get("x-real-ip")
			if remoteIP == "" {
				remoteIP = r.Header.Get("x-forwarded-for")
				if remoteIP == "" {
					remoteIP = r.RemoteAddr
				}
			}

			fields := []zapcore.Field{
				zap.Int("status", ww.Status()),
				zap.Duration("duration", time.Since(start)),
				zap.String("path", r.RequestURI),
				zap.String("method", r.Method),
				zap.String("protocol", r.Proto),
				zap.String("agent", r.UserAgent()),
				zap.String("remote", remoteIP),
			}

			if reqID := middleware.GetReqID(r.Context()); reqID != "" {
				fields = append(fields, zap.String("request-id", reqID))
			}

			if config.RequestBody {
				if req, err := httputil.DumpRequest(r, true); err == nil {
					fields = append(fields, zap.ByteString("request", req))
				}
			}
			if config.ResponseBody {
				fields = append(fields, zap.ByteString("response", responseBody.Bytes()))
			}

			// Write the log entry assuming we're logging at that level.
			if entry := logger.Check(config.Level, "HTTP Request"); entry != nil {
				entry.Write(fields...)
			}
		})
	}
}

// Returns a middleware function for logging requests
func LoggerStackdriverMiddleware(logger *zap.Logger, config Config) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Check if the prefix should be ignored
			for _, prefix := range config.IgnorePaths {
				if strings.HasPrefix(r.RequestURI, prefix) {
					next.ServeHTTP(w, r)
					return
				}
			}

			start := time.Now()

			// See if we're already using a wrapped response writer and if not make one.
			ww, ok := w.(middleware.WrapResponseWriter)
			if !ok {
				ww = middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			}

			// If we should log bodies, setup buffer to capture
			var responseBody *bytes.Buffer
			if config.ResponseBody {
				responseBody = new(bytes.Buffer)
				ww.Tee(responseBody)
			}

			next.ServeHTTP(ww, r)

			// If the remote IP is being proxied, use the real IP
			remoteIP := r.Header.Get("x-real-ip")
			if remoteIP == "" {
				remoteIP = r.Header.Get("x-forwarded-for")
				if remoteIP == "" {
					remoteIP = r.RemoteAddr
				}
			}

			fields := []zapcore.Field{
				zapdriver.HTTP(&zapdriver.HTTPPayload{
					RequestMethod: r.Method,
					RequestURL:    r.RequestURI,
					RequestSize:   strconv.FormatInt(r.ContentLength, 10),
					Status:        ww.Status(),
					ResponseSize:  strconv.Itoa(ww.BytesWritten()),
					UserAgent:     r.UserAgent(),
					RemoteIP:      remoteIP,
					Referer:       r.Referer(),
					Latency:       fmt.Sprintf("%fs", time.Since(start).Seconds()),
					Protocol:      r.Proto,
				}),
			}

			if reqID := r.Context().Value(middleware.RequestIDKey); reqID != nil {
				fields = append(fields, zap.String("request-id", reqID.(string)))
			}

			if config.RequestBody {
				if req, err := httputil.DumpRequest(r, true); err == nil {
					fields = append(fields, zap.ByteString("request", req))
				}
			}
			if config.ResponseBody {
				fields = append(fields, zap.ByteString("response", responseBody.Bytes()))
			}

			// Write the log entry assuming we're logging at that level.
			if entry := logger.Check(config.Level, "HTTP Request"); entry != nil {
				entry.Write(fields...)
			}

		})
	}
}
