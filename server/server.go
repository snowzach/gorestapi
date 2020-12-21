package server

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/knadh/koanf"
	"github.com/snowzach/certtools"
	"github.com/snowzach/certtools/autocert"
	"go.uber.org/zap"
)

type Server struct {
	logger *zap.SugaredLogger
	router chi.Router
	server *http.Server
}

// New will setup the API listener
func New(config *koanf.Koanf) (*Server, error) {

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	// Log Requests - Use appropriate format depending on the encoding
	if config.Bool("server.log_requests") {
		switch config.String("logger.encoding") {
		case "stackdriver":
			r.Use(loggerHTTPMiddlewareStackdriver(config.Bool("server.log_requests_body"), config.Strings("server.log_disabled_http")))
		default:
			r.Use(loggerHTTPMiddlewareDefault(config.Bool("server.log_requests_body"), config.Strings("server.log_disabled_http")))
		}
	}

	// CORS Config
	r.Use(cors.New(cors.Options{
		AllowedOrigins:   config.Strings("server.cors.allowed_origins"),
		AllowedMethods:   config.Strings("server.cors.allowed_methods"),
		AllowedHeaders:   config.Strings("server.cors.allowed_headers"),
		AllowCredentials: config.Bool("server.cors.allowed_credentials"),
		MaxAge:           config.Int("server.cors.max_age"),
	}).Handler)

	s := &Server{
		logger: zap.S().With("package", "server"),
		router: r,
	}

	return s, nil

}

// ListenAndServe will listen for requests
func (s *Server) ListenAndServe(config *koanf.Koanf) error {

	s.server = &http.Server{
		Addr:    net.JoinHostPort(config.String("server.host"), config.String("server.port")),
		Handler: s.router,
	}

	// Listen
	listener, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		return fmt.Errorf("Could not listen on %s: %v", s.server.Addr, err)
	}

	// Enable TLS?
	if config.Bool("server.tls") {
		var cert tls.Certificate
		if config.Bool("server.devcert") {
			s.logger.Warn("WARNING: This server is using an insecure development tls certificate. This is for development only!!!")
			cert, err = autocert.New(autocert.InsecureStringReader("localhost"))
			if err != nil {
				return fmt.Errorf("Could not autocert generate server certificate: %v", err)
			}
		} else {
			// Load keys from file
			cert, err = tls.LoadX509KeyPair(config.String("server.certfile"), config.String("server.keyfile"))
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
	s.logger.Infow("API Listening", "address", s.server.Addr, "tls", config.Bool("server.tls"))

	// Enable profiler
	if config.Bool("server.profiler_enabled") && config.String("server.profiler_path") != "" {
		zap.S().Debugw("Profiler enabled on API", "path", config.String("server.profiler_path"))
		s.router.Mount(config.String("server.profiler_path"), middleware.Profiler())
	}

	return nil

}

// Router returns the router
func (s *Server) Router() chi.Router {
	return s.router
}
