package main

import (
	"log"
	"zero-provers/server/grpc"
	"zero-provers/server/http"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

type Config struct {
	// Port of the gRPC server.
	gRPCServerPort int
	// Port of the HTTP server.
	hTTPServerPort int
	// URL path of the HTTP server save endpoint.
	hTTPServerSaveEndpoint string
	// Directory in which proofs are stored.
	proofsOutputDir string
	// Directory in which mock data is provided.
	mockDataDir string
	// Set to true if debug mode is enabled.
	debug bool
	// Verbosity of the logs.
	logLevel zerolog.Level
}

func main() {
	var config Config
	var rootCmd = &cobra.Command{
		Use:   "mock",
		Short: "Edge gRPC mock server",
		Run: func(cmd *cobra.Command, args []string) {
			// Start the gRPC server.
			go func() {
				log.Fatal(grpc.StartgRPCServer(config.logLevel, config.gRPCServerPort, config.mockDataDir))
			}()

			// Start the HTTP server.
			log.Fatal(http.StartHTTPServer(config.logLevel, config.hTTPServerPort, config.hTTPServerSaveEndpoint, config.proofsOutputDir))
		},
	}

	// Define flags for configuration
	rootCmd.PersistentFlags().IntVarP(&config.gRPCServerPort, "grpc-port", "g", 8546, "gRPC server port")
	rootCmd.PersistentFlags().IntVarP(&config.hTTPServerPort, "http-port", "p", 8080, "HTTP server port")
	rootCmd.PersistentFlags().StringVarP(&config.hTTPServerSaveEndpoint, "http-save-endpoint", "e", "/save", "HTTP server save endpoint")
	rootCmd.PersistentFlags().StringVarP(&config.proofsOutputDir, "output-dir", "o", "out", "Proofs output directory")
	rootCmd.PersistentFlags().StringVarP(&config.mockDataDir, "mock-data-dir", "m", "data", "Mock data directory containing mock status (status.json), block (block.json) and trace (trace.json) files")
	rootCmd.PersistentFlags().BoolVarP(&config.debug, "debug", "d", false, "Enable verbose mode")
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}

	// Set log level.
	config.logLevel = zerolog.InfoLevel
	if config.debug {
		config.logLevel = zerolog.DebugLevel
	}
}
