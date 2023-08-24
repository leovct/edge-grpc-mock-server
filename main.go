package main

import (
	"fmt"
	"log"
	"zero-provers/server/grpc"
	"zero-provers/server/http"
	"zero-provers/server/modes"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

type Config struct {
	// Port of the gRPC server.
	GRPCServerPort int
	// Port of the HTTP server.
	HTTPServerPort int
	// URL path of the HTTP server save endpoint.
	HTTPServerSaveEndpoint string

	// Directory in which proofs are stored.
	ProofsOutputDir string
	// Directory in which mock data is provided.
	MockDataDir       string
	MockDataBlockFile string
	MockDataTraceFile string

	// Mode of the mock server, either static or dynamic.
	// - static: the server always return the same mock block data.
	// - dynamic: the server returns new mock block data every x requests.
	// - random: the server returns random block data every requests.
	Mode string

	// Set to true if debug mode is enabled.
	Debug bool
	// Verbosity of the logs.
	LogLevel zerolog.Level
}

func main() {
	var config Config
	var rootCmd = &cobra.Command{
		Use:   "edge-grpc-mock-server",
		Short: "Edge gRPC mock server",
		Run: func(cmd *cobra.Command, args []string) {
			// Determine log level based on debug flag.
			if config.Debug {
				config.LogLevel = zerolog.DebugLevel
			} else {
				config.LogLevel = zerolog.InfoLevel
			}

			// Check the mode.
			switch modes.Mode(config.Mode) {
			case modes.StaticMode, modes.DynamicMode, modes.RandomMode:
				// Valid modes, no action needed.
			default:
				fmt.Printf("Mode '%s' is not supported... Please either use '%s', '%s' or '%s'.",
					config.Mode, modes.StaticMode, modes.DynamicMode, modes.RandomMode)
				return
			}
			fmt.Println(config.Mode)

			// Start the gRPC server.
			go func() {
				log.Fatal(grpc.StartgRPCServer(grpc.ServerConfig{
					LogLevel: config.LogLevel,
					Port:     config.GRPCServerPort,
					Mode:     modes.Mode(config.Mode),
					MockData: grpc.Mock{
						Dir:       config.MockDataDir,
						BlockFile: config.MockDataBlockFile,
						TraceFile: config.MockDataTraceFile,
					},
				}))
			}()

			// Start the HTTP server.
			log.Fatal(http.StartHTTPServer(http.ServerConfig{
				LogLevel:        config.LogLevel,
				Port:            config.HTTPServerPort,
				SaveEndpoint:    config.HTTPServerSaveEndpoint,
				ProofsOutputDir: config.ProofsOutputDir,
			}))
		},
	}

	// Define flags for configuration.
	rootCmd.PersistentFlags().IntVarP(&config.GRPCServerPort, "grpc-port", "g", 8546, "gRPC server port")
	rootCmd.PersistentFlags().IntVarP(&config.HTTPServerPort, "http-port", "p", 8080, "HTTP server port")
	rootCmd.PersistentFlags().StringVarP(&config.HTTPServerSaveEndpoint, "http-save-endpoint", "e", "/save", "HTTP server save endpoint")

	rootCmd.PersistentFlags().StringVarP(&config.ProofsOutputDir, "output-dir", "o", "out", "Proofs output directory")
	rootCmd.PersistentFlags().StringVar(&config.MockDataDir, "mock-data-dir", "data", "Mock data directory")
	rootCmd.PersistentFlags().StringVar(&config.MockDataBlockFile, "mock-data-block-file", "block.json", "Mock data block file (in the mock data dir)")
	rootCmd.PersistentFlags().StringVar(&config.MockDataTraceFile, "mock-data-trace-file", "trace3.json", "Mock data trace file (in the mock data dir)")

	rootCmd.PersistentFlags().StringVarP(&config.Mode, "mode", "m", string(modes.StaticMode),
		`Mode of the mock server.
- static: the server always return the same mock block data.
- dynamic: the server returns new mock block data every {n} requests.
- random: the server returns random block data every requests.
`)

	rootCmd.PersistentFlags().BoolVarP(&config.Debug, "debug", "d", false, "Enable verbose mode")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
