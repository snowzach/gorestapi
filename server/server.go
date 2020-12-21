package server

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/snowzach/certtools"
	"github.com/snowzach/certtools/autocert"
	"go.uber.org/zap"

	"github.com/snowzach/gorestapi/conf"
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

	// Log Requests - Use appropriate format depending on the encoding
	if conf.C.Bool("server.log_requests") {
		switch conf.C.String("logger.encoding") {
		case "stackdriver":
			r.Use(loggerHTTPMiddlewareStackdriver(conf.C.Bool("server.log_requests_body"), conf.C.Strings("server.log_disabled_http")))
		default:
			r.Use(loggerHTTPMiddlewareDefault(conf.C.Bool("server.log_requests_body"), conf.C.Strings("server.log_disabled_http")))
		}
	}

	// CORS Config
	r.Use(cors.New(cors.Options{
		AllowedOrigins:   conf.C.Strings("server.cors.allowed_origins"),
		AllowedMethods:   conf.C.Strings("server.cors.allowed_methods"),
		AllowedHeaders:   conf.C.Strings("server.cors.allowed_headers"),
		AllowCredentials: conf.C.Bool("server.cors.allowed_credentials"),
		MaxAge:           conf.C.Int("server.cors.max_age"),
	}).Handler)

	s := &Server{
		logger: zap.S().With("package", "server"),
		router: r,
	}

	return s, nil

}

// ListenAndServe will listen for requests
func (s *Server) ListenAndServe() error {

	s.server = &http.Server{
		Addr:    net.JoinHostPort(conf.C.String("server.host"), conf.C.String("server.port")),
		Handler: s.router,
	}

	// Listen
	listener, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		return fmt.Errorf("Could not listen on %s: %v", s.server.Addr, err)
	}

	// Enable TLS?
	if conf.C.Bool("server.tls") {
		var cert tls.Certificate
		if conf.C.Bool("server.devcert") {
			s.logger.Warn("WARNING: This server is using an insecure development tls certificate. This is for development only!!!")
			cert, err = autocert.New(autocert.InsecureStringReader("localhost"))
			if err != nil {
				return fmt.Errorf("Could not autocert generate server certificate: %v", err)
			}
		} else {
			// Load keys from file
			cert, err = tls.LoadX509KeyPair(conf.C.String("server.certfile"), conf.C.String("server.keyfile"))
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
	s.logger.Infow("API Listening", "address", s.server.Addr, "tls", conf.C.Bool("server.tls"))

	// Enable profiler
	if conf.C.Bool("server.profiler_enabled") && conf.C.String("server.profiler_path") != "" {
		zap.S().Debugw("Profiler enabled on API", "path", conf.C.String("server.profiler_path"))
		s.router.Mount(conf.C.String("server.profiler_path"), middleware.Profiler())
	}

	return nil

}

// Router returns the router
func (s *Server) Router() chi.Router {
	return s.router
}
