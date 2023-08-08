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
}

func main() {
	var config Config
	var rootCmd = &cobra.Command{
		Use:   "mock",
		Short: "Edge gRPC mock server",
		Run: func(cmd *cobra.Command, args []string) {
			// Start the gRPC server.
			go func() {
				log.Fatal(grpc.StartgRPCServer(config.GRPCServerPort))
			}()

			// Start the HTTP server.
			log.Fatal(http.StartHTTPServer(config.HTTPServerPort, config.ProofsOutputDir))
		},
	}

	// Define flags for configuration
	rootCmd.Flags().IntVar(&config.GRPCServerPort, "grpc", 8546, "gRPC server port")
	rootCmd.Flags().IntVar(&config.HTTPServerPort, "http", 8080, "HTTP server port")
	rootCmd.Flags().StringVarP(&config.ProofsOutputDir, "output-dir", "o", "out", "Proofs output directory")
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
