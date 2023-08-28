package main

import (
	"fmt"
	"log"
	"path/filepath"
	"zero-provers/server/grpc"
	"zero-provers/server/http"
	"zero-provers/server/logger"
	"zero-provers/server/modes"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

type Config struct {
	//// Server configuration.
	// Port of the gRPC server.
	GRPCServerPort int
	// Port of the HTTP server.
	HTTPServerPort int
	// URL path of the HTTP server save endpoint.
	HTTPServerSaveEndpoint string

	// Mode of the mock server, either static or dynamic.
	// - static: the server always return the same mock block data.
	// - dynamic: the server returns new mock block data every x requests.
	// - random: the server returns random block data every requests.
	Mode string

	//// Static mode configuration.
	// Mock data files loaded in static mode.
	MockBlockDir string
	MockTraceDir string

	//// Dynamic mode configuration.
	// Mock data directories (and underlying files) used in dynamic mode.
	MockBlockFile string
	MockTraceFile string
	// Number of requests after which the server returns new data, block and trace (used in `dynamic` mode).
	UpdateDataThreshold int

	//// Random mode configuration.
	// Number of requests after which the server increments the block number (used in `random` mode).
	UpdateBlockNumberThreshold int

	//// Other parameters.
	// Directory in which proofs are stored.
	ProofsOutputDir string
	// Verbosity of the logs.
	Verbosity int8
}

func main() {
	var config Config
	var rootCmd = &cobra.Command{
		Use:   "edge-grpc-mock-server",
		Short: "Edge gRPC mock server",
		Run: func(cmd *cobra.Command, args []string) {
			// Set up the logger.
			logLevel := zerolog.Level(config.Verbosity)
			lc := logger.LoggerConfig{
				Level:       logLevel,
				CallerField: "root",
			}
			customLog := logger.NewLogger(lc)

			// Check the mode.
			switch modes.Mode(config.Mode) {
			case modes.StaticMode, modes.RandomMode:
				// Valid modes, no action needed.
			case modes.DynamicMode:
				// Check that the number of block and trace files are the same.
				blockFiles, err := filepath.Glob(filepath.Join(config.MockBlockDir, "*.json"))
				if err != nil {
					return
				}
				traceFiles, err := filepath.Glob(filepath.Join(config.MockTraceDir, "*.json"))
				if err != nil {
					return
				}
				if len(blockFiles) != len(traceFiles) {
					customLog.Fatal().Msg("When running the mock server in dynamic mode, you need the same number of block and trace files")
					return
				}
			default:
				customLog.Fatal().Msgf("Mode '%s' is not supported... Please either use '%s', '%s' or '%s'.",
					config.Mode, modes.StaticMode, modes.DynamicMode, modes.RandomMode)
				return
			}

			// Start the gRPC server.
			go func() {
				log.Fatal(grpc.StartgRPCServer(grpc.ServerConfig{
					LogLevel:                   logLevel,
					Port:                       config.GRPCServerPort,
					Mode:                       modes.Mode(config.Mode),
					UpdateDataThreshold:        config.UpdateDataThreshold,
					UpdateBlockNumberThreshold: config.UpdateBlockNumberThreshold,
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

	// Static mode configuration.
	rootCmd.PersistentFlags().StringVar(&config.MockBlockFile, "mock-data-block-file", "data/blocks/block-57.json", "The mock data block file path (used in static mode)")
	rootCmd.PersistentFlags().StringVar(&config.MockTraceFile, "mock-data-trace-file", "data/traces/trace-57.json", "The mock data trace file path (used in static mode)")
	rootCmd.PersistentFlags().IntVar(&config.UpdateDataThreshold, "update-data-threshold", 30, "The number of requests after which the server returns new data, block and trace (used in dynamic mode).")

	// Dynamic mode configuration.
	rootCmd.PersistentFlags().StringVar(&config.MockBlockDir, "mock-data-block-dir", "data/blocks", "The mock data block directory (used in dynamic mode)")
	rootCmd.PersistentFlags().StringVar(&config.MockTraceDir, "mock-data-trace-dir", "data/traces", "The mock data trace directory (used in dynamic mode)")

	// Random mode configuration.
	rootCmd.PersistentFlags().IntVar(&config.UpdateBlockNumberThreshold, "update-block-number-threshold", 30, "The number of requests after which the server increments the block number (used in random mode)")

	// Other parameters.
	rootCmd.PersistentFlags().StringVarP(&config.ProofsOutputDir, "output-dir", "o", "out", "The proofs output directory")
	rootCmd.PersistentFlags().Int8VarP(&config.Verbosity, "verbosity", "v", int8(zerolog.InfoLevel),
		fmt.Sprintf("Verbosity level from %d (%s) to %d (%s)",
			int8(zerolog.PanicLevel), zerolog.PanicLevel, int8(zerolog.TraceLevel), zerolog.TraceLevel))

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
