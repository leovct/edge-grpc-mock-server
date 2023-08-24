// Package grpc provides functionalities to start and handle a gRPC server.
package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
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
	blockMockFilePath  string
	tracesMockFilePath string

	// Increase the block height on every GetStatus request made.
	counter     int
	counterStep int = 50
)

type ServerConfig struct {
	LogLevel zerolog.Level
	Port     int
	Mode     modes.Mode
	MockData Mock
}

// Mock data config.
type Mock struct {
	Dir       string
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
	blockMockFilePath = fmt.Sprintf("%s/%s", config.MockData.Dir, config.MockData.BlockFile)
	tracesMockFilePath = fmt.Sprintf("%s/%s", config.MockData.Dir, config.MockData.TraceFile)

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

	var height int64
	switch mode {
	case modes.StaticMode:
		// Parse the block mock data file and return the header number.
		var mockBlock pb.BlockData
		if err := loadDataFromFile(blockMockFilePath, &mockBlock); err != nil {
			return nil, err
		}
		block, err := parseAndPrintRawBlockData(mockBlock.Data)
		if err != nil {
			return nil, err
		}
		height = int64(block.Header.Number)
	case modes.RandomMode, modes.DynamicMode:
		height = int64(constantBlockHeight + counter)
		counter += counterStep
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

	var rawData []byte
	switch mode {
	case modes.StaticMode:
		// Parse the block mock data file and return the raw data.
		var mockBlock pb.BlockData
		if err := loadDataFromFile(blockMockFilePath, &mockBlock); err != nil {
			return nil, err
		}
		rawData = mockBlock.Data
	case modes.RandomMode, modes.DynamicMode:
		height := constantBlockHeight + counter
		block := edge.GenerateRandomEdgeBlock(uint64(height), uint64(10))
		rawData = block.MarshalRLP()
	default:
		return nil, errWrongMode
	}

	if _, err := parseAndPrintRawBlockData(rawData); err != nil {
		return nil, err
	}
	log.Debug().Msgf("BlockResponse encoded data: %v", rawData)

	return &pb.BlockData{
		Data: rawData,
	}, nil
}

func (s *server) GetTrace(context.Context, *pb.BlockNumber) (*pb.Trace, error) {
	log.Info().Msg("gRPC /GetTrace request received")

	var rawTrace []byte
	switch mode {
	case modes.StaticMode:
		// Parse the trace mock data file and return the raw trace.
		var mockTrace pb.Trace
		if err := loadDataFromFile(tracesMockFilePath, &mockTrace); err != nil {
			return nil, err
		}
		rawTrace = mockTrace.Trace
	case modes.RandomMode, modes.DynamicMode:
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

	if err := parseAndPrintRawTrace(rawTrace); err != nil {
		return nil, err
	}
	log.Debug().Msgf("TraceResponse encoded trace: %v", rawTrace)

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
			log.Debug().Msgf("BlockResponse decoded data: %v", string(data))
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
			log.Debug().Msgf("TraceResponce decoded trace: %v", string(data))

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
						log.Debug().Msgf("Transaction #%d decoded: %v", i+1, string(data))
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
