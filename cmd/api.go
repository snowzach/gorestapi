package cmd

import (
	"net/http"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	cli "github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/snowzach/gorestapi/conf"
	"github.com/snowzach/gorestapi/embed"
	"github.com/snowzach/gorestapi/gorestapi/thingrpc"
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

			pg, err := postgres.New()
			if err != nil {
				logger.Fatalw("Database Error", "error", err)
			}

			// Create the server
			s, err := server.New()
			if err != nil {
				logger.Fatalw("Could not create server", "error", err)
			}

			// Main RPC
			if err = thingrpc.Setup(s.Router(), pg); err != nil {
				logger.Fatalw("Could not setup bcrpc", "error", err)
			}

			// Serve api-docs and swagger-ui
			docsFileServer := http.FileServer(&assetfs.AssetFS{Asset: embed.Asset, AssetDir: embed.AssetDir, AssetInfo: embed.AssetInfo, Prefix: "public"})
			s.Router().Get("/api/api-docs/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Vary", "Accept-Encoding")
				w.Header().Set("Cache-Control", "no-cache")
				http.StripPrefix("/api", docsFileServer).ServeHTTP(w, r)
			}))

			err = s.ListenAndServe()
			if err != nil {
				logger.Fatalw("Could not start server",
					"error", err,
				)
			}

			conf.Stop.SetupInterrupt()
			<-conf.Stop.Chan() // Wait until StopChan

			conf.Stop.Wait() // Wait until everyone cleans up
			zap.L().Sync()   // Flush the logger

		},
	}
)
