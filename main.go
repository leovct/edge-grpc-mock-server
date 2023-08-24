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
	MockBlockDir string
	MockTraceDir string

	// Mock data file paths.
	MockBlockFile string
	MockTraceFile string

	// Mode of the mock server, either static or dynamic.
	// - static: the server always return the same mock block data.
	// - dynamic: the server returns new mock block data every x requests.
	// - random: the server returns random block data every requests.
	Mode string

	// Verbosity of the logs.
	Verbosity int8
}

func main() {
	var config Config
	var rootCmd = &cobra.Command{
		Use:   "edge-grpc-mock-server",
		Short: "Edge gRPC mock server",
		Run: func(cmd *cobra.Command, args []string) {
			logLevel := zerolog.Level(config.Verbosity)

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
					LogLevel: logLevel,
					Port:     config.GRPCServerPort,
					Mode:     modes.Mode(config.Mode),
					MockData: grpc.MockData{
						BlockDir:  config.MockBlockDir,
						TraceDir:  config.MockTraceDir,
						BlockFile: config.MockBlockFile,
						TraceFile: config.MockTraceFile,
					},
				}))
			}()

			// Start the HTTP server.
			log.Fatal(http.StartHTTPServer(http.ServerConfig{
				LogLevel:        logLevel,
				Port:            config.HTTPServerPort,
				SaveEndpoint:    config.HTTPServerSaveEndpoint,
				ProofsOutputDir: config.ProofsOutputDir,
			}))
		},
	}

	// Server configuration.
	rootCmd.PersistentFlags().IntVarP(&config.GRPCServerPort, "grpc-port", "g", 8546, "gRPC server port")
	rootCmd.PersistentFlags().IntVarP(&config.HTTPServerPort, "http-port", "p", 8080, "HTTP server port")
	rootCmd.PersistentFlags().StringVarP(&config.HTTPServerSaveEndpoint, "http-save-endpoint", "e", "/save", "HTTP server save endpoint")

	// Server mode.
	rootCmd.PersistentFlags().StringVarP(&config.Mode, "mode", "m", string(modes.StaticMode),
		`Mode of the mock server.
- static: the server always return the same mock block data.
- dynamic: the server returns new mock block data every {n} requests.
- random: the server returns random block data every requests.
`)

	// Mock data files loaded in static mode.
	rootCmd.PersistentFlags().StringVar(&config.MockBlockFile, "mock-data-block-file", "data/blocks/block.json", "Mock data block file path")
	rootCmd.PersistentFlags().StringVar(&config.MockTraceFile, "mock-data-trace-file", "data/traces/encoded/trace3.json", "Mock data trace file path")

	// Mock data directories (and underlying files) used in dynamic mode.
	rootCmd.PersistentFlags().StringVar(&config.MockBlockDir, "mock-data-block-dir", "data/blocks", "Mock data block directory")
	rootCmd.PersistentFlags().StringVar(&config.MockTraceDir, "mock-data-trace-dir", "data/traces/encoded", "Mock data trace directory")

	// Other parameters.
	rootCmd.PersistentFlags().StringVarP(&config.ProofsOutputDir, "output-dir", "o", "out", "Proofs output directory")
	rootCmd.PersistentFlags().Int8VarP(&config.Verbosity, "verbosity", "v", int8(zerolog.InfoLevel),
		fmt.Sprintf("Verbosity level from %d (%s) to %d (%s)",
			int8(zerolog.PanicLevel), zerolog.PanicLevel, int8(zerolog.TraceLevel), zerolog.TraceLevel))

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
