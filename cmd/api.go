package cmd

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	cli "github.com/spf13/cobra"
	"go.uber.org/zap/zapcore"

	"github.com/snowzach/gorestapi/embed"
	"github.com/snowzach/gorestapi/gorestapi/mainrpc"
	"github.com/snowzach/gorestapi/pkg/conf"
	"github.com/snowzach/gorestapi/pkg/log"
	"github.com/snowzach/gorestapi/pkg/server"
	"github.com/snowzach/gorestapi/pkg/server/logger"
	"github.com/snowzach/gorestapi/pkg/server/metrics"
	"github.com/snowzach/gorestapi/pkg/signal"
	"github.com/snowzach/gorestapi/pkg/version"
	"github.com/snowzach/gorestapi/store/postgres"
)

func init() {
	rootCmd.AddCommand(apiCmd)
}

var (
	apiCmd = &cli.Command{
		Use:   "api",
		Short: "Start API",
		Long:  `Start API`,
		Run: func(cmd *cli.Command, args []string) { // Initialize the databse

			var err error

			// Create the router and server config
			router, err := newRouter()
			if err != nil {
				log.Fatalf("router config error: %v", err)
			}

			// Create the database
			db, err := newDatabase()
			if err != nil {
				log.Fatalf("database config error: %v", err)
			}

			// Version endpoint
			router.Get("/version", version.GetVersion())

			// ThingRPC
			if err = mainrpc.Setup(router, db); err != nil {
				log.Fatalf("Could not setup mainrpc: %v", err)
			}

			// Serve api-docs and swagger-ui
			docsFileServer := http.FileServer(http.FS(embed.PublicHTMLFS()))
			router.Mount("/api-docs", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Vary", "Accept-Encoding")
				w.Header().Set("Cache-Control", "no-cache")
				docsFileServer.ServeHTTP(w, r)
			}))

			// Create a server
			s, err := newServer(router)
			if err != nil {
				log.Fatalf("could not create server error: %v", err)
			}

			// Start the listener and service connections.
			go func() {
				if err = s.ListenAndServe(); err != nil {
					log.Errorf("Server error: %v", err)
					signal.Stop.Stop()
				}
			}()
			log.Infof("API listening on %s", s.Addr)

			// Register signal handler and wait
			signal.Stop.OnSignal(signal.DefaultStopSignals...)
			<-signal.Stop.Chan() // Wait until Stop
			signal.Stop.Wait()   // Wait until everyone cleans up
			log.Flush()          // Flush the logger

		},
	}
)

func newRouter() (chi.Router, error) {

	router := chi.NewRouter()
	router.Use(
		middleware.Recoverer, // Recover from panics
		middleware.RequestID, // Inject request-id
	)

	// Request logger
	if conf.C.Bool("server.log.enabled") {
		var loggerConfig logger.Config
		if err := conf.C.Unmarshal(&loggerConfig, conf.UnmarshalConf{Path: "server.log"}); err != nil {
			return nil, fmt.Errorf("could not parser server.log config: %w", err)
		}
		switch conf.C.String("logger.encoding") {
		case "stackdriver":
			router.Use(logger.LoggerStackdriverMiddleware(log.Base.Named("server"), loggerConfig))
		default:
			router.Use(logger.LoggerStandardMiddleware(log.Base.Named("server"), loggerConfig))
		}
	}

	// CORS handler
	if conf.C.Bool("server.cors.enabled") {
		var corsOptions cors.Options
		if err := conf.C.Unmarshal(&corsOptions, conf.UnmarshalConf{
			Path: "server.cors",
			DecoderConfig: conf.DefaultDecoderConfig(
				conf.WithMatchName(conf.MatchSnakeCaseConfig),
			),
		}); err != nil {
			return nil, fmt.Errorf("could not parser server.cors config: %w", err)
		}
		router.Use(cors.New(corsOptions).Handler)
	}

	// If we have server metrics enabled, enable the middleware to collect them on the server.
	if conf.C.Bool("server.metrics.enabled") {
		var metricsConfig metrics.Config
		if err := conf.C.Unmarshal(&metricsConfig, conf.UnmarshalConf{
			Path:          "server.metrics",
			DecoderConfig: conf.DefaultDecoderConfig(),
		}); err != nil {
			return nil, fmt.Errorf("could not parser server.metrics config: %w", err)
		}
		router.Use(metrics.MetricsMiddleware(metricsConfig))
	}

	return router, nil

}

func newServer(handler http.Handler) (*server.Server, error) {

	// Parse the config
	var serverConfig = &server.Config{Handler: handler}
	if err := conf.C.Unmarshal(serverConfig, conf.UnmarshalConf{Path: "server"}); err != nil {
		fmt.Errorf("could not parse server config: %v", err)
	}

	// Create the server
	s, err := server.New(serverConfig)
	if err != nil {
		fmt.Errorf("could not create server: %v", err)
	}

	return s, nil

}

func newDatabase() (*postgres.Client, error) {

	var err error

	// Database config
	var postgresConfig = &postgres.Config{}
	if err := conf.C.Unmarshal(postgresConfig, conf.UnmarshalConf{Path: "database"}); err != nil {
		return nil, fmt.Errorf("could not parse database config: %v", err)
	}

	// Loggers
	postgresConfig.Logger = log.NewWrapper(log.Base.Named("store.postgres"), zapcore.InfoLevel)
	postgresConfig.QueryLogger = log.NewWrapper(log.Base.Named("store.postgres.query"), zapcore.DebugLevel)

	// if conf.C.Bool("store.postgres.log_queries") {
	// 	postgresConfig.QueryLogger = log.NewWrapper(log.Base.Named("store.postgres.query"), zapcore.DebugLevel)
	// }

	// Migrations
	postgresConfig.MigrationSource, err = embed.MigrationSource()
	if err != nil {
		return nil, fmt.Errorf("could not get database migrations error: %w", err)
	}

	// Create database
	db, err := postgres.New(postgresConfig)
	if err != nil {
		return nil, fmt.Errorf("could not create database client: %w", err)
	}

	return db, nil

}
