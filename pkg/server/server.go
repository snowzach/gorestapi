package server

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"

	"github.com/snowzach/certtools"
	"github.com/snowzach/certtools/autocert"
)

type Config struct {
	Host     string `conf:"host"`
	Port     string `conf:"port"`
	TLS      bool   `conf:"tls"`
	DevCert  bool   `conf:"devcert"`
	CertFile string `conf:"certfile"`
	KeyFile  string `conf:"keyfile"`
	Handler  http.Handler
}

type Server struct {
	config *Config
	*http.Server
}

// New will setup the API listener
func New(config *Config) (*Server, error) {
	return &Server{
		config: config,
		Server: &http.Server{
			Handler: config.Handler,
			Addr:    net.JoinHostPort(config.Host, config.Port),
		},
	}, nil
}

// ListenAndServe will listen for requests
func (s *Server) ListenAndServe() error {

	// Listen
	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return fmt.Errorf("could not listen on %s: %w", s.Addr, err)
	}

	// Enable TLS?
	if s.config.TLS {
		var cert tls.Certificate
		if s.config.DevCert {
			cert, err = autocert.New(autocert.InsecureStringReader("localhost"))
			if err != nil {
				return fmt.Errorf("could not generate autocert server certificate: %w", err)
			}
		} else {
			// Load keys from file
			cert, err = tls.LoadX509KeyPair(s.config.CertFile, s.config.KeyFile)
			if err != nil {
				return fmt.Errorf("could not load server certificate: %w", err)
			}
		}

		// Sane/Safe defaults
		s.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   certtools.SecureTLSMinVersion(),
			CipherSuites: certtools.SecureTLSCipherSuites(),
		}
		// Wrap the listener in a TLS Listener
		listener = tls.NewListener(listener, s.TLSConfig)
	}

	return s.Serve(listener)

}
