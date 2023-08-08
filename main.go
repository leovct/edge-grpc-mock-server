package main

import (
	"log"
	"zero-provers/server/grpc"
	"zero-provers/server/http"

	"github.com/spf13/cobra"
)

type Config struct {
	GRPCServerPort  int
	HTTPServerPort  int
	ProofsOutputDir string
	MockDataDir     string
}

func main() {
	var config Config
	var rootCmd = &cobra.Command{
		Use:   "mock",
		Short: "Edge gRPC mock server",
		Run: func(cmd *cobra.Command, args []string) {
			// Start the gRPC server.
			go func() {
				log.Fatal(grpc.StartgRPCServer(config.GRPCServerPort, config.MockDataDir))
			}()

			// Start the HTTP server.
			log.Fatal(http.StartHTTPServer(config.HTTPServerPort, config.ProofsOutputDir))
		},
	}

	// Define flags for configuration
	rootCmd.Flags().IntVar(&config.GRPCServerPort, "grpc", 8546, "gRPC server port")
	rootCmd.Flags().IntVar(&config.HTTPServerPort, "http", 8080, "HTTP server port")
	rootCmd.Flags().StringVarP(&config.ProofsOutputDir, "output-dir", "o", "out", "Proofs output directory")
	rootCmd.Flags().StringVarP(&config.MockDataDir, "mock-data-dir", "m", "data", "Mock data directory containing mock status (status.json), block (block.json) and trace (trace.json) files")
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
