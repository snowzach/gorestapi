package conf

import (
	"net/http"
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
)

// C is the global configuration with "." for delimeter
var C = koanf.New(".")

func init() {

	// Set Defaults
	_ = C.Load(confmap.Provider(map[string]interface{}{
		// Logger Defaults
		"logger.level":              "info",
		"logger.encoding":           "console",
		"logger.color":              true,
		"logger.dev_mode":           true,
		"logger.disable_caller":     false,
		"logger.disable_stacktrace": true,

		// Pidfile
		"pidfile": "",

		// Profiler config
		"profiler.enabled": false,
		"profiler.host":    "",
		"profiler.port":    "6060",

		// Server Configuration
		"server.host":                     "",
		"server.port":                     "8080",
		"server.tls":                      false,
		"server.devcert":                  false,
		"server.certfile":                 "server.crt",
		"server.keyfile":                  "server.key",
		"server.log_requests":             true,
		"server.log_requests_body":        false,
		"server.log_disabled_http":        []string{"/version"},
		"server.profiler_enabled":         false,
		"server.profiler_path":            "/debug",
		"server.cors.allowed_origins":     []string{"*"},
		"server.cors.allowed_methods":     []string{http.MethodHead, http.MethodOptions, http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch},
		"server.cors.allowed_headers":     []string{"*"},
		"server.cors.allowed_credentials": false,
		"server.cors.max_age":             300,

		// Database Settings
		"database.type":                  "postgres",
		"database.username":              "postgres",
		"database.password":              "password",
		"database.host":                  "postgres",
		"database.port":                  5432,
		"database.database":              "gorestapi",
		"database.sslmode":               "disable",
		"database.retries":               5,
		"database.sleep_between_retries": "7s",
		"database.max_connections":       40,
		"database.wipe_confirm":          false,
		"database.log_queries":           false,
	}, "."), nil)

	// All underscores in environment variables to dots
	envReplacer := strings.NewReplacer("_", ".")
	// Build a map of existing config items with all underscores replaced with dots so `thing.that_value` can
	// be replaced by environment variable THING_THAT_VALUE instead of it trying to replace `thing.that.value`
	envLookup := make(map[string]string) //
	for _, key := range C.Keys() {
		envLookup[envReplacer.Replace(key)] = key
	}
	// Load the environment variables, compare to our lookup of existing values and set override value
	_ = C.Load(env.Provider("", ".", func(s string) string {
		// Convert environemnt variable to lower case and change underscore to dot
		key := envReplacer.Replace(strings.ToLower(s))
		if replacement, found := envLookup[key]; found {
			return replacement
		}
		return "" // No existing variable, skip it
	}), nil)

}
