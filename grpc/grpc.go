// Package grpc provides functionalities to start and handle a gRPC server.
package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"zero-provers/server/grpc/edge"
	edgetypes "zero-provers/server/grpc/edge/types"
	pb "zero-provers/server/grpc/pb"
	"zero-provers/server/logger"
	"zero-provers/server/modes"

	empty "google.golang.org/protobuf/types/known/emptypb"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Constant dummy block height returned by the `/GetStatus` endpoint.
const constantBlockHeight = 100_000_000_000_000_000

var (
	config ServerConfig

	// Log is the package-level variable used for logging messages and errors.
	log zerolog.Logger

	errWrongMode = fmt.Errorf("wrong mode")

	// Keep track of the number of `status` requests made to the mock server. The zero-prover constantly
	// sends those requests, in order to be aware of new blocks and to start proving as soon as possible.
	requestCounter int
)

type ServerConfig struct {
	LogLevel                   zerolog.Level
	Port                       int
	Mode                       modes.Mode
	UpdateDataThreshold        int
	UpdateBlockNumberThreshold int
	MockData                   MockData
}

// Mock data config.
type MockData struct {
	BlockDir  string
	TraceDir  string
	BlockFile string
	TraceFile string
}

// server is an internal implementation of the gRPC server.
type server struct {
	pb.UnimplementedSystemServer
}

// StartgRPCServer starts a gRPC server on the specified port.
// It listens for incoming TCP connections and handles gRPC requests using the internal server
// implementation. The server continues to run until it is manually stopped or an error occurs.
func StartgRPCServer(_config ServerConfig) error {
	config = _config

	// Set up the logger.
	lc := logger.LoggerConfig{
		Level:       config.LogLevel,
		CallerField: "grpc-server",
	}
	log = logger.NewLogger(lc)

	// Create a listener on the specified port.
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Port))
	if err != nil {
		return err
	}

	// Create a new gRPC server instance with reflection and system services.
	s := grpc.NewServer()
	reflection.Register(s)
	pb.RegisterSystemServer(s, &server{})

	// Start serving incoming gRPC requests on the listener.
	log.Info().Msgf("gRPC server is starting on port %d", config.Port)
	if err := s.Serve(listener); err != nil {
		log.Error().Err(err).Msg("Unable to start gRPC server")
		return err
	}
	return nil
}

// GetStatus is the implementation of the `GetStatus` RPC method.
func (s *server) GetStatus(context.Context, *empty.Empty) (*pb.ChainStatus, error) {
	log.Info().Msg("gRPC /GetStatus request received")

	requestCounter++
	log.Debug().Msgf("Request counter: %d", requestCounter)

	// Load block number from file or increment block number based on the request counter.
	var height int64
	switch config.Mode {
	case modes.StaticMode:
		// Parse the block mock file and return the header number.
		var err error
		height, err = getBlockNumberFromBlockFile(config.MockData.BlockFile)
		if err != nil {
			return nil, err
		}

	case modes.DynamicMode:
		// List the block mock files under the block mock directory.
		files, err := filepath.Glob(filepath.Join(config.MockData.BlockDir, "*.json"))
		if err != nil {
			return nil, err
		}

		// Parse the block mock file at the current index and return the header number.
		fileIndex := computeIndex(requestCounter, config.UpdateDataThreshold, len(files))
		file := files[fileIndex]
		height, err = getBlockNumberFromBlockFile(file)
		if err != nil {
			return nil, err
		}

	case modes.RandomMode:
		// Increment the constant block number based on request counter.
		height = int64(constantBlockHeight + requestCounter%config.UpdateBlockNumberThreshold)

	default:
		return nil, errWrongMode
	}

	log.Debug().Msgf("StatusResponse number: %v", height)
	return &pb.ChainStatus{
		Current: &pb.ChainStatus_Block{
			Number: height,
		},
	}, nil

}

// BlockByNumber is the implementation of the `BlockByNumber` RPC method.
func (s *server) BlockByNumber(context.Context, *pb.BlockNumber) (*pb.BlockData, error) {
	log.Info().Msg("gRPC /BlockByNumber request received")

	// Load block data from file or generate random data.
	var rawData []byte
	switch config.Mode {
	case modes.StaticMode:
		// Parse the block mock file and return the raw data.
		var mockBlock pb.BlockData
		if err := loadDataFromFile(config.MockData.BlockFile, &mockBlock); err != nil {
			return nil, err
		}
		rawData = mockBlock.Data

	case modes.DynamicMode:
		// List the block mock files under the block mock directory.
		files, err := filepath.Glob(filepath.Join(config.MockData.BlockDir, "*.json"))
		if err != nil {
			return nil, err
		}

		// Parse the block mock file at the current index and return the raw data.
		fileIndex := computeIndex(requestCounter, config.UpdateDataThreshold, len(files))
		file := files[fileIndex]
		var mockBlock pb.BlockData
		if err := loadDataFromFile(file, &mockBlock); err != nil {
			return nil, err
		}
		rawData = mockBlock.Data

	case modes.RandomMode:
		// Return a random block data.
		height := uint64(constantBlockHeight + requestCounter%config.UpdateBlockNumberThreshold)
		txnTracesAmount := uint64(10)
		block := edge.GenerateRandomEdgeBlock(height, txnTracesAmount)
		rawData = block.MarshalRLP()

	default:
		return nil, errWrongMode
	}

	// Parse the data.
	if _, err := parseAndPrintRawBlockData(rawData); err != nil {
		return nil, err
	}
	log.Trace().Msgf("BlockResponse encoded data: %v", rawData)

	return &pb.BlockData{
		Data: rawData,
	}, nil
}

func (s *server) GetTrace(context.Context, *pb.BlockNumber) (*pb.Trace, error) {
	log.Info().Msg("gRPC /GetTrace request received")

	// Load trace data from file or generate random data.
	var rawTrace []byte
	switch config.Mode {
	case modes.StaticMode:
		// Parse the trace mock data file and return the raw trace.
		var mockTrace pb.Trace
		if err := loadDataFromFile(config.MockData.TraceFile, &mockTrace); err != nil {
			return nil, err
		}
		rawTrace = mockTrace.Trace

	case modes.DynamicMode:
		// List the block trace files under the trace mock directory.
		files, err := filepath.Glob(filepath.Join(config.MockData.TraceDir, "*.json"))
		if err != nil {
			return nil, err
		}

		// Parse the trace mock file at the current index and return the raw trace.
		fileIndex := computeIndex(requestCounter, config.UpdateDataThreshold, len(files))
		file := files[fileIndex]
		var mockTrace pb.Trace
		if err := loadDataFromFile(file, &mockTrace); err != nil {
			return nil, err
		}
		rawTrace = mockTrace.Trace

	case modes.RandomMode:
		trace := *edge.GenerateRandomEdgeTrace(10, 10, 10, 10)
		var err error
		rawTrace, err = json.Marshal(trace)
		if err != nil {
			log.Error().Err(err).Msg("BlockTrace encoding failed")
			return nil, err
		}

	default:
		return nil, errWrongMode
	}

	// Parse the raw trace.
	if err := parseAndPrintRawTrace(rawTrace); err != nil {
		return nil, err
	}
	log.Trace().Msgf("TraceResponse encoded trace: %v", rawTrace)

	return &pb.Trace{
		Trace: rawTrace,
	}, nil
}

// Parse a raw block data and display its content.
func parseAndPrintRawBlockData(rawBlockData []byte) (edgetypes.Block, error) {
	decodedBlock := edgetypes.Block{}
	if err := decodedBlock.UnmarshalRLP(rawBlockData); err != nil {
		log.Error().Err(err).Msg("BlockData decoding failed")
		return edgetypes.Block{}, err
	} else {
		data, err := json.MarshalIndent(decodedBlock, "", "  ")
		if err != nil {
			log.Error().Err(err).Msg("Unable to format JSON struct")
			return edgetypes.Block{}, err
		} else {
			log.Trace().Msgf("BlockResponse decoded data: %v", string(data))
		}
	}
	return decodedBlock, nil
}

// Parse a raw trace and display its content.
func parseAndPrintRawTrace(rawTrace []byte) error {
	// Decode the raw trace.
	var decodedTrace *edgetypes.Trace
	if err := json.Unmarshal(rawTrace, &decodedTrace); err != nil {
		log.Error().Err(err).Msg("Raw trace decoding failed")
		return err
	} else {
		// Marshal the decoded trace to JSON.
		data, err := json.MarshalIndent(decodedTrace, "", "  ")
		if err != nil {
			log.Error().Err(err).Msg("Unable to format JSON struct")
			return err
		} else {
			log.Trace().Msgf("TraceResponce decoded trace: %v", string(data))

			// Decode each transaction of the trace.
			traces := decodedTrace.TxnTraces
			if len(traces) > 0 {
				log.Debug().Msgf("Decoding %d transaction trace(s)...", len(traces))
			}
			for i, trace := range traces {
				decodedTxn := edgetypes.Transaction{}
				txnBytes := []byte(trace.Transaction)
				if err := decodedTxn.UnmarshalRLP(txnBytes); err != nil {
					log.Error().Err(err).Msgf("Transaction #%d decoding failed", i+1)
					return err
				} else {
					data, err := json.MarshalIndent(decodedTxn, "", "  ")
					if err != nil {
						log.Error().Err(err).Msg("Unable to format JSON struct")
						return err
					} else {
						log.Trace().Msgf("Transaction #%d decoded: %v", i+1, string(data))
					}
				}
			}
		}
	}
	return nil
}

// Compute the file index in the case of dynamic mode.
// Iterate over all the indexes and if the index is greater than the number of files, return
// the index of the last file.
func computeIndex(requestCounter, updateThreshold int, numberOfFiles int) int {
	index := (requestCounter - 1) / updateThreshold
	if index > numberOfFiles-1 {
		return numberOfFiles - 1
	}
	return index
}

// Load data from file.
func loadDataFromFile(filePath string, target interface{}) error {
	log.Debug().Msgf("Fetching mock data from %s", filePath)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading mock file: %w", err)
	}

	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("error unmarshaling mock JSON: %w", err)
	}

	log.Debug().Msgf("Mock data loaded from %s", filePath)
	return nil
}

// Load block number from block file.
func getBlockNumberFromBlockFile(filePath string) (int64, error) {
	var rawBlock pb.BlockData
	if err := loadDataFromFile(filePath, &rawBlock); err != nil {
		return 0, err
	}
	block, err := parseAndPrintRawBlockData(rawBlock.Data)
	if err != nil {
		return 0, err
	}
	return int64(block.Header.Number), nil
}
