package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/GhanshyamJha05/fifa-cli/internal/server"
	"github.com/spf13/cobra"
)

func serveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Start the REST API server",
		Long:  "Starts the HTTP API on the configured host and port (default 0.0.0.0:8080).",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := svc.Config()
			srv := server.New(cfg, svc, logger)
			fmt.Fprintf(os.Stderr, "API listening on http://%s:%d\n", cfg.ServerHost, cfg.ServerPort)

			go func() {
				if err := srv.Start(); err != nil {
					logger.Error("server stopped", "error", err)
				}
			}()

			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
			<-quit

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			return srv.Shutdown(ctx)
		},
	}
}
