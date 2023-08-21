package main

import (
	"fmt"
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
	mockDataDir        string
	mockDataBlockFile  string
	mockDataStatusFile string
	mockDataTraceFile  string
	// Generate random trace data instead of relying on mocks.
	setRandomMode bool
	// Set to true if debug mode is enabled.
	debug bool
	// Verbosity of the logs.
	logLevel zerolog.Level
}

func main() {
	var config Config
	var rootCmd = &cobra.Command{
		Use:   "edge-grpc-mock-server",
		Short: "Edge gRPC mock server",
		Run: func(cmd *cobra.Command, args []string) {
			// Determine log level based on debug flag.
			if config.debug {
				config.logLevel = zerolog.DebugLevel
			} else {
				config.logLevel = zerolog.InfoLevel
			}
			fmt.Printf("Log level set to %s\n", config.logLevel)

			// Start the gRPC server.
			go func() {
				mock := grpc.Mock{
					Dir:        config.mockDataDir,
					StatusFile: config.mockDataStatusFile,
					BlockFile:  config.mockDataBlockFile,
					TraceFile:  config.mockDataTraceFile,
				}
				log.Fatal(grpc.StartgRPCServer(config.logLevel, config.gRPCServerPort, config.setRandomMode, mock))
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
	rootCmd.PersistentFlags().StringVarP(&config.mockDataDir, "mock-data-dir", "m", "data", "Mock data directory")
	rootCmd.PersistentFlags().StringVar(&config.mockDataStatusFile, "mock-data-status-file", "status.json", "Mock data status file (in the mock data dir)")
	rootCmd.PersistentFlags().StringVar(&config.mockDataBlockFile, "mock-data-block-file", "block.json", "Mock data block file (in the mock data dir)")
	rootCmd.PersistentFlags().StringVar(&config.mockDataTraceFile, "mock-data-trace-file", "trace3.json", "Mock data trace file (in the mock data dir)")
	rootCmd.PersistentFlags().BoolVarP(&config.setRandomMode, "random", "r", false, "Generate random trace data instead of relying on mocks (default false)")
	rootCmd.PersistentFlags().BoolVarP(&config.debug, "debug", "d", false, "Enable verbose mode")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
