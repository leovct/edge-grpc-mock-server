package config

import "github.com/rs/zerolog"

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
	MockDataDir        string
	MockDataBlockFile  string
	MockDataStatusFile string
	MockDataTraceFile  string

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
