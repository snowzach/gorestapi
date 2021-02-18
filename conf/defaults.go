package conf

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
)

// C is the global configuration with "." for delimeter
var C = koanf.New(".")

// Defaults loads the default config for the app
func Defaults(c *koanf.Koanf) error {
	return c.Load(confmap.Provider(map[string]interface{}{
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
	}, "."), nil)
}

// File loads configuration from a file
func File(c *koanf.Koanf, configFile string) error {
	ext := filepath.Ext(configFile)
	switch ext {
	case ".yaml", ".yml":
		return c.Load(file.Provider(configFile), yaml.Parser())
	case ".json":
		return c.Load(file.Provider(configFile), json.Parser())
	case ".toml":
		return c.Load(file.Provider(configFile), toml.Parser())
	}
	return fmt.Errorf("unknown config extension %s", ext)
}

// Env environment configuration overrides
func Env(c *koanf.Koanf) error {
	// All underscores in environment variables to dots
	envReplacer := strings.NewReplacer("_", ".")
	// Build a map of existing config items with all underscores replaced with dots so `thing.that_value` can
	// be replaced by environment variable THING_THAT_VALUE instead of it trying to replace `thing.that.value`
	envLookup := make(map[string]string) //
	for _, key := range C.Keys() {
		envLookup[envReplacer.Replace(key)] = key
	}
	// Load the environment variables, compare to our lookup of existing values and set override value
	return c.Load(env.ProviderWithValue("", ".", func(key string, value string) (string, interface{}) {
		// Convert environemnt variable to lower case and change underscore to dot
		key = envReplacer.Replace(strings.ToLower(key))
		if replacement, found := envLookup[key]; found {
			// Check the existing type of the variable, and allow modifying
			switch C.Get(replacement).(type) {
			case []interface{}, []string: // If existing value is string slice, split on space
				return replacement, strings.Split(value, " ")
			}
			return replacement, value
		}
		return "", nil // No existing variable, skip it
	}), nil)
}
