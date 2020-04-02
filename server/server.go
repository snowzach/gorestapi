package server

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/snowzach/certtools"
	"github.com/snowzach/certtools/autocert"
	config "github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Server struct {
	logger *zap.SugaredLogger
	router chi.Router
	server *http.Server
}

// New will setup the API listener
func New() (*Server, error) {

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// Log Requests - Use appropriate format depending on the encoding
	if config.GetBool("server.log_requests") {
		switch config.GetString("logger.encoding") {
		case "stackdriver":
			r.Use(loggerHTTPMiddlewareStackdriver(config.GetBool("server.log_requests_body"), config.GetStringSlice("server.log_disabled_http")))
		default:
			r.Use(loggerHTTPMiddlewareDefault(config.GetBool("server.log_requests_body"), config.GetStringSlice("server.log_disabled_http")))
		}
	}

	// Log Requests
	if config.GetBool("server.log_requests") {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				start := time.Now()
				var requestID string
				if reqID := r.Context().Value(middleware.RequestIDKey); reqID != nil {
					requestID = reqID.(string)
				}
				ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
				next.ServeHTTP(ww, r)

				latency := time.Since(start)

				fields := []zapcore.Field{
					zap.Int("status", ww.Status()),
					zap.Duration("took", latency),
					zap.String("remote", r.RemoteAddr),
					zap.String("request", r.RequestURI),
					zap.String("method", r.Method),
					zap.String("package", "server.request"),
				}
				if requestID != "" {
					fields = append(fields, zap.String("request-id", requestID))
				}
				zap.L().Info("API Request", fields...)
			})
		})
	}

	// CORS Config
	r.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{http.MethodHead, http.MethodOptions, http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler)

	s := &Server{
		logger: zap.S().With("package", "server"),
		router: r,
	}

	// Register our routes
	s.router.Get("/version", GetVersion())

	return s, nil

}

// ListenAndServe will listen for requests
func (s *Server) ListenAndServe() error {

	s.server = &http.Server{
		Addr:    net.JoinHostPort(config.GetString("server.host"), config.GetString("server.port")),
		Handler: s.router,
	}

	// Listen
	listener, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		return fmt.Errorf("Could not listen on %s: %v", s.server.Addr, err)
	}

	// Enable TLS?
	if config.GetBool("server.tls") {
		var cert tls.Certificate
		if config.GetBool("server.devcert") {
			s.logger.Warn("WARNING: This server is using an insecure development tls certificate. This is for development only!!!")
			cert, err = autocert.New(autocert.InsecureStringReader("localhost"))
			if err != nil {
				return fmt.Errorf("Could not autocert generate server certificate: %v", err)
			}
		} else {
			// Load keys from file
			cert, err = tls.LoadX509KeyPair(config.GetString("server.certfile"), config.GetString("server.keyfile"))
			if err != nil {
				return fmt.Errorf("Could not load server certificate: %v", err)
			}
		}

		// Enabed Certs - TODO Add/Get a cert
		s.server.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   certtools.SecureTLSMinVersion(),
			CipherSuites: certtools.SecureTLSCipherSuites(),
		}
		// Wrap the listener in a TLS Listener
		listener = tls.NewListener(listener, s.server.TLSConfig)
	}

	go func() {
		if err = s.server.Serve(listener); err != nil {
			s.logger.Fatalw("API Listen error", "error", err, "address", s.server.Addr)
		}
	}()
	s.logger.Infow("API Listening", "address", s.server.Addr, "tls", config.GetBool("server.tls"))

	// Enable profiler
	if config.GetBool("server.profiler_enabled") && config.GetString("server.profiler_path") != "" {
		zap.S().Debugw("Profiler enabled on API", "path", config.GetString("server.profiler_path"))
		s.router.Mount(config.GetString("server.profiler_path"), middleware.Profiler())
	}

	return nil

}

// Router returns the router
func (s *Server) Router() chi.Router {
	return s.router
}

// RenderOrErrInternal will render whatever you pass it (assuming it has Renderer) or prints an internal error
func RenderOrErrInternal(w http.ResponseWriter, r *http.Request, d render.Renderer) {
	if err := render.Render(w, r, d); err != nil {
		render.Render(w, r, ErrInternal(err))
		return
	}
}
