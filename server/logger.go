package server

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strconv"
	"time"

	"github.com/blendle/zapdriver"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func loggerHTTPMiddlewareDefault(logRequestBody bool, disabledEndpoints []string) func(http.Handler) http.Handler {
	// Make a map lookup for disabled endpoints
	disabled := make(map[string]struct{})
	for _, d := range disabledEndpoints {
		disabled[d] = struct{}{}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// If Disabled
			if _, ok := disabled[r.RequestURI]; ok {
				next.ServeHTTP(w, r)
				return
			}

			start := time.Now()

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			var response *bytes.Buffer
			if logRequestBody {
				response = new(bytes.Buffer)
				ww.Tee(response)
			}

			next.ServeHTTP(ww, r)

			fields := []zapcore.Field{
				zap.Int("status", ww.Status()),
				zap.Duration("duration", time.Since(start)),
				zap.String("path", r.RequestURI),
				zap.String("method", r.Method),
				zap.String("package", "server.http"),
			}

			if reqID := r.Context().Value(middleware.RequestIDKey); reqID != nil {
				fields = append(fields, zap.String("request-id", reqID.(string)))
			}

			if logRequestBody {
				if req, err := httputil.DumpRequest(r, true); err == nil {
					fields = append(fields, zap.ByteString("request", req))
				}
				fields = append(fields, zap.ByteString("response", response.Bytes()))
			}

			// If we have an x-Forwarded-For header, use that for the remote
			if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
				fields = append(fields, zap.String("remote", forwardedFor))
			} else {
				fields = append(fields, zap.String("remote", r.RemoteAddr))
			}
			zap.L().Info("HTTP Request", fields...)
		})
	}
}

// Returns a middleware function for logging requests
func loggerHTTPMiddlewareStackdriver(logRequestBody bool, disabledEndpoints []string) func(http.Handler) http.Handler {
	// Make a map lookup for disabled endpoints
	disabled := make(map[string]struct{})
	for _, d := range disabledEndpoints {
		disabled[d] = struct{}{}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// If Disabled
			if _, ok := disabled[r.RequestURI]; ok {
				next.ServeHTTP(w, r)
				return
			}

			start := time.Now()

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			var response *bytes.Buffer
			if logRequestBody {
				response = new(bytes.Buffer)
				ww.Tee(response)
			}

			next.ServeHTTP(ww, r)

			// If the remote IP is being proxied, use the real IP
			remoteIP := r.Header.Get("X-Forwarded-For")
			if remoteIP == "" {
				remoteIP = r.RemoteAddr
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
				zap.String("package", "server.http"),
			}

			if reqID := r.Context().Value(middleware.RequestIDKey); reqID != nil {
				fields = append(fields, zap.String("request-id", reqID.(string)))
			}

			if logRequestBody {
				if req, err := httputil.DumpRequest(r, true); err == nil {
					fields = append(fields, zap.ByteString("request", req))
				}
				fields = append(fields, zap.ByteString("response", response.Bytes()))
			}

			zap.L().Info("HTTP Request", fields...)
		})
	}
}
