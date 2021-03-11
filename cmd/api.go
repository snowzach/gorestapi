package cmd

import (
	"net/http"

	cli "github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/snowzach/gorestapi/conf"
	"github.com/snowzach/gorestapi/embed"
	"github.com/snowzach/gorestapi/gorestapi/mainrpc"
	"github.com/snowzach/gorestapi/server"
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

			migrationSource, err := embed.MigrationSource()
			if err != nil {
				logger.Fatalw("Could not get database migrations", "error", err)
			}

			// Database
			pg, err := postgres.New(conf.C, migrationSource)
			if err != nil {
				logger.Fatalw("Database error", "error", err)
			}

			// Create the server
			s, err := server.New(conf.C)
			if err != nil {
				logger.Fatalw("Could not create server", "error", err)
			}

			s.Router().Get("/version", conf.GetVersion())

			// ThingRPC
			if err = mainrpc.Setup(s.Router(), pg); err != nil {
				logger.Fatalw("Could not setup thingrpc", "error", err)
			}

			// Serve api-docs and swagger-ui
			docsFileServer := http.FileServer(http.FS(embed.PublicHTMLFS()))
			s.Router().Mount("/api-docs", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Vary", "Accept-Encoding")
				w.Header().Set("Cache-Control", "no-cache")
				docsFileServer.ServeHTTP(w, r)
			}))

			if err = s.ListenAndServe(conf.C); err != nil {
				logger.Fatalw("Could not start server", "error", err)
			}

			conf.Stop.InitInterrupt()
			<-conf.Stop.Chan() // Wait until Stop
			conf.Stop.Wait()   // Wait until everyone cleans up
			_ = zap.L().Sync() // Flush the logger

		},
	}
)
