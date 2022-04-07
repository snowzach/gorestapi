package cmd

import "net/http"

// defaults loads the default config for the app
func defaults() map[string]interface{} {
	return map[string]interface{}{
		// Logger Defaults
		"logger.level":              "info",
		"logger.encoding":           "console",
		"logger.color":              true,
		"logger.dev_mode":           true,
		"logger.disable_caller":     false,
		"logger.disable_stacktrace": true,

		// Metrics, profiler, pidfile
		"metrics.enabled":  true,
		"metrics.host":     "",
		"metrics.port":     "6060",
		"profiler.enabled": true,
		"pidfile":          "",

		// Server Configuration
		"server.host":     "",
		"server.port":     "8080",
		"server.tls":      false,
		"server.devcert":  false,
		"server.certfile": "server.crt",
		"server.keyfile":  "server.key",
		// Server Log
		"server.log.enabled":       true,
		"server.log.level":         "info",
		"server.log.request_body":  false,
		"server.log.response_body": false,
		"server.log.ignore_paths":  []string{"/version"},
		// Server CORS
		"server.cors.enabled":           true,
		"server.cors.allowed_origins":   []string{"*"},
		"server.cors.allowed_methods":   []string{http.MethodHead, http.MethodOptions, http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch},
		"server.cors.allowed_headers":   []string{"*"},
		"server.cors.allow_credentials": false,
		"server.cors.max_age":           300,
		// Server Metrics
		"server.metrics.enabled":      true,
		"server.metrics.ignore_paths": []string{"/version"},

		// Database Settings
		"database.username":              "postgres",
		"database.password":              "password",
		"database.host":                  "postgres",
		"database.port":                  5432,
		"database.database":              "gorestapi",
		"database.auto_create":           true,
		"database.search_path":           "",
		"database.sslmode":               "disable",
		"database.sslcert":               "",
		"database.sslkey":                "",
		"database.sslrootcert":           "",
		"database.retries":               5,
		"database.sleep_between_retries": "7s",
		"database.max_connections":       40,
		"database.log_queries":           false,
		"database.wipe_confirm":          false,
	}
}
