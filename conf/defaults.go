package conf

import (
	"net/http"

	config "github.com/spf13/viper"
)

func init() {

	// Logger Defaults
	config.SetDefault("logger.level", "info")
	config.SetDefault("logger.encoding", "console")
	config.SetDefault("logger.color", true)
	config.SetDefault("logger.dev_mode", true)
	config.SetDefault("logger.disable_caller", false)
	config.SetDefault("logger.disable_stacktrace", true)

	// Pidfile
	config.SetDefault("pidfile", "")

	// Profiler config
	config.SetDefault("profiler.enabled", false)
	config.SetDefault("profiler.host", "")
	config.SetDefault("profiler.port", "6060")

	// Server Configuration
	config.SetDefault("server.host", "")
	config.SetDefault("server.port", "8900")
	config.SetDefault("server.tls", false)
	config.SetDefault("server.devcert", false)
	config.SetDefault("server.certfile", "server.crt")
	config.SetDefault("server.keyfile", "server.key")
	config.SetDefault("server.log_requests", true)
	config.SetDefault("server.log_requests_body", false)
	config.SetDefault("server.log_disabled_http", []string{"/version"})
	config.SetDefault("server.profiler_enabled", false)
	config.SetDefault("server.profiler_path", "/debug")
	config.SetDefault("server.cors.allowed_origins", []string{"*"})
	config.SetDefault("server.cors.allowed_methods", []string{http.MethodHead, http.MethodOptions, http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch})
	config.SetDefault("server.cors.allowed_headers", []string{"*"})
	config.SetDefault("server.cors.allowed_credentials", false)
	config.SetDefault("server.cors.max_age", 300)

	// Database Settings
	config.SetDefault("storage.type", "postgres")
	config.SetDefault("storage.username", "postgres")
	config.SetDefault("storage.password", "password")
	config.SetDefault("storage.host", "postgres")
	config.SetDefault("storage.port", 5432)
	config.SetDefault("storage.database", "gorestapi")
	config.SetDefault("storage.sslmode", "disable")
	config.SetDefault("storage.retries", 5)
	config.SetDefault("storage.sleep_between_retries", "7s")
	config.SetDefault("storage.max_connections", 80)
	config.SetDefault("storage.wipe_confirm", false)

}
