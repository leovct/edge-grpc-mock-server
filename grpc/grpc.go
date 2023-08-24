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
	// Log is the package-level variable used for logging messages and errors.
	log zerolog.Logger

	// Mode.
	mode         modes.Mode
	errWrongMode = fmt.Errorf("wrong mode")

	// Mock data.
	mockData MockData

	// Keep track of the number of `status` requests made to the mock server. The zero-prover constantly
	// sends those requests, in order to be aware of new blocks and to start proving as soon as possible.
	requestCounter int

	// Number of requests after which the server returns new data, block and trace (used in `dynamic` mode).
	updateDataThreshold int

	// Number of requests after which the server increments the block number (used in `random` mode).
	updateBlockNumberThreshold int
)

type ServerConfig struct {
	LogLevel zerolog.Level
	Port     int
	Mode     modes.Mode
	MockData MockData
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
func StartgRPCServer(config ServerConfig) error {
	// Set up the logger.
	lc := logger.LoggerConfig{
		Level:       config.LogLevel,
		CallerField: "grpc-server",
	}
	log = logger.NewLogger(lc)

	// Set up other parameters.
	mode = config.Mode
	mockData = config.MockData

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

	// Load block number from file or increment block number based on the request counter.
	var height int64
	switch mode {
	case modes.StaticMode:
		// Parse the block mock file and return the header number.
		var err error
		height, err = getBlockNumberFromBlockFile(mockData.BlockFile)
		if err != nil {
			return nil, err
		}

	case modes.DynamicMode:
		// List the block mock files under the block mock directory.
		files, err := filepath.Glob(filepath.Join(mockData.BlockDir, "*.json"))
		if err != nil {
			return nil, err
		}

		// Parse the block mock file at the current index and return the header number.
		fileIndex := requestCounter % updateDataThreshold
		file := files[fileIndex]
		fmt.Printf("\n[debug]\nfile_index:%d\nfile:%s\n[debug])", fileIndex, file) // TODO: remove.
		height, err = getBlockNumberFromBlockFile(file)
		if err != nil {
			return nil, err
		}

	case modes.RandomMode:
		// Increment the constant block number based on request counter.
		height = int64(constantBlockHeight + requestCounter%updateBlockNumberThreshold)
		fmt.Printf("\n[debug]\nnew_height:%d\n[debug]", height) // TODO: remove.

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
	switch mode {
	case modes.StaticMode:
		// Parse the block mock file and return the raw data.
		var mockBlock pb.BlockData
		if err := loadDataFromFile(mockData.BlockFile, &mockBlock); err != nil {
			return nil, err
		}
		rawData = mockBlock.Data

	case modes.DynamicMode:
		// List the block mock files under the block mock directory.
		files, err := filepath.Glob(filepath.Join(mockData.BlockDir, "*.json"))
		if err != nil {
			return nil, err
		}

		// Parse the block mock file at the current index and return the raw data.
		fileIndex := requestCounter % updateDataThreshold
		file := files[fileIndex]
		fmt.Printf("\n[debug]\nfile_index:%d\nfile:%s\n[debug])", fileIndex, file) // TODO: remove.
		var mockBlock pb.BlockData
		if err := loadDataFromFile(file, &mockBlock); err != nil {
			return nil, err
		}
		rawData = mockBlock.Data

	case modes.RandomMode:
		// Return a random block data.
		height := uint64(constantBlockHeight + requestCounter%updateBlockNumberThreshold)
		fmt.Printf("\n[debug]\nnew_height:%d\n[debug]", height) // TODO: remove.
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
	switch mode {
	case modes.StaticMode:
		// Parse the trace mock data file and return the raw trace.
		var mockTrace pb.Trace
		if err := loadDataFromFile(mockData.TraceFile, &mockTrace); err != nil {
			return nil, err
		}
		rawTrace = mockTrace.Trace

	case modes.DynamicMode:
		// List the block trace files under the trace mock directory.
		files, err := filepath.Glob(filepath.Join(mockData.TraceDir, "*.json"))
		if err != nil {
			return nil, err
		}

		// Parse the trace mock file at the current index and return the raw trace.
		fileIndex := requestCounter % updateDataThreshold
		file := files[fileIndex]
		fmt.Printf("\n[debug]\nfile_index:%d\nfile:%s\n[debug])", fileIndex, file) // TODO: remove.
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
			fmt.Println("BlockTrace encoding failed:", err)
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
				log.Debug().Msg("Decoding TraceResponce transactionTraces txn fields (RLP encoded)...")
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
