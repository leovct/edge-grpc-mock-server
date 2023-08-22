// Package grpc provides functionalities to start and handle a gRPC server.
package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"
	"zero-provers/server/grpc/edge"
	edgetypes "zero-provers/server/grpc/edge/types"
	pb "zero-provers/server/grpc/pb"
	"zero-provers/server/logger"

	empty "google.golang.org/protobuf/types/known/emptypb"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Constant dummy block height returned by the `/GetStatus` endpoint.
const constantBlockHeight = 100_000_000_000_000_000

var (
	// log is the package-level variable used for logging messages and errors.
	log zerolog.Logger

	// Increase the block height on every GetStatus request made.
	counter     int
	counterStep int = 50

	// Mock data provided by the user.
	mockStatusData *pb.ChainStatus
	mockBlockData  *pb.BlockData
	mockTraceData  *pb.Trace
	traceMutex     sync.Mutex
)

// server is an internal implementation of the gRPC server.
type server struct {
	pb.UnimplementedSystemServer
}

// Mock data config.
type Mock struct {
	Dir        string
	StatusFile string
	BlockFile  string
	TraceFile  string
}

// StartgRPCServer starts a gRPC server on the specified port.
// It listens for incoming TCP connections and handles gRPC requests using the internal server
// implementation. The server continues to run until it is manually stopped or an error occurs.
func StartgRPCServer(logLevel zerolog.Level, port int, setRandomMode bool, mockData Mock) error {
	// Set up the logger.
	lc := logger.LoggerConfig{
		Level:       logLevel,
		CallerField: "grpc-server",
	}
	log = logger.NewLogger(lc)

	// Create a listener on the specified port.
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	// Create a new gRPC server instance with reflection and system services.
	s := grpc.NewServer()
	reflection.Register(s)
	pb.RegisterSystemServer(s, &server{})

	// Load mock data if provided.
	if !setRandomMode {
		log.Debug().Msgf("Fetching mock data from `%s` directory", mockData.Dir)
		mockStatusData, mockBlockData, mockTraceData, err = loadMockData(mockData)
		if err != nil {
			log.Error().Err(err).Msg("Unable to load mock data")
			return err
		}
	}

	// Start serving incoming gRPC requests on the listener.
	log.Info().Msgf("gRPC server is starting on port %d", port)
	if err := s.Serve(listener); err != nil {
		log.Error().Err(err).Msg("Unable to start gRPC server")
		return err
	}
	return nil
}

// GetStatus is the implementation of the `GetStatus` RPC method.
// It returns a constant `ServerStatus` response.
func (s *server) GetStatus(context.Context, *empty.Empty) (*pb.ChainStatus, error) {
	log.Info().Msg("gRPC /GetStatus request received")

	// Return mock data if provided.
	if mockStatusData != nil {
		log.Debug().Msgf("Mock StatusResponse number: %v", mockStatusData.Current.Number)
		return mockStatusData, nil
	}

	// Else, return dummy data.
	height := int64(constantBlockHeight + counter)
	counter += counterStep
	log.Debug().Msgf("StatusResponse number: %v", height)
	return &pb.ChainStatus{
		Current: &pb.ChainStatus_Block{
			Number: height,
		},
	}, nil
}

// BlockByNumber is the implementation of the `BlockByNumber` RPC method.
// It returns a constant `BlockResponse` containing a single byte.
func (s *server) BlockByNumber(context.Context, *pb.BlockNumber) (*pb.BlockData, error) {
	log.Info().Msg("gRPC /BlockByNumber request received")

	var rawData []byte
	if mockBlockData != nil {
		// Return mock data if provided.
		rawData = mockBlockData.Data
		log.Debug().Msgf("Mock BlockResponse encoded data: %v", mockBlockData.Data)
	} else {
		// Else, return random data.
		height := constantBlockHeight + counter
		block := edge.GenerateRandomEdgeBlock(uint64(height), uint64(10))
		rawData = block.MarshalRLP()
		log.Debug().Msgf("BlockResponse encoded data: %v", rawData)
	}
	if err := parseAndPrintRawBlockData(rawData); err != nil {
		return nil, err
	}

	return &pb.BlockData{
		Data: rawData,
	}, nil
}

func (s *server) GetTrace(context.Context, *pb.BlockNumber) (*pb.Trace, error) {
	log.Info().Msg("gRPC /GetTrace request received")

	traceMutex.Lock()
	defer traceMutex.Unlock()

	var rawTrace []byte
	if mockTraceData != nil {
		// Return mock data if provided.
		log.Debug().Msgf("Mock TraceResponse encoded data: %v", mockTraceData.Trace)
		rawTrace = mockTraceData.Trace
	} else {
		// Else, return random data.
		trace := *edge.GenerateRandomEdgeTrace(10, 10, 10, 10)
		var err error
		rawTrace, err = json.Marshal(trace)
		if err != nil {
			fmt.Println("BlockTrace encoding failed:", err)
			return nil, err
		}
		log.Debug().Msgf("TraceResponse encoded trace: %v", rawTrace)
	}
	if err := parseAndPrintRawTrace(rawTrace); err != nil {
		return nil, err
	}

	return &pb.Trace{
		Trace: rawTrace,
	}, nil
}

func (s *server) UpdateTrace(ctx context.Context, req *pb.Trace) (*pb.OperationStatus, error) {
	log.Info().Msg("gRPC /UpdateTrace request received")

	traceMutex.Lock()
	defer traceMutex.Unlock()

	// Extract the new trace data from the request.
	newRawTrace := []byte(req.Trace)
	if err := parseAndPrintRawTrace(newRawTrace); err != nil {
		return nil, err
	}

	// Update trace data.
	mockTraceData = &pb.Trace{
		Trace: newRawTrace,
	}
	log.Debug().Msg("Trace data updated successfully")
	return &pb.OperationStatus{
		Success: true,
	}, nil
}

// Parse a raw block data and display its content.
func parseAndPrintRawBlockData(rawBlockData []byte) error {
	decodedBlock := edgetypes.Block{}
	if err := decodedBlock.UnmarshalRLP(rawBlockData); err != nil {
		log.Error().Err(err).Msg("BlockData decoding failed")
		return err
	} else {
		data, err := json.MarshalIndent(decodedBlock, "", "  ")
		if err != nil {
			log.Error().Err(err).Msg("Unable to format JSON struct")
			return err
		} else {
			log.Debug().Msgf("BlockResponse decoded data: %v", string(data))
		}
	}
	return nil
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

// Load mock data if provided by the user.
func loadMockData(mockData Mock) (*pb.ChainStatus, *pb.BlockData, *pb.Trace, error) {
	// Status mock data.
	statusMockFilePath := fmt.Sprintf("%s/%s", mockData.Dir, mockData.StatusFile)
	var mockStatus pb.ChainStatus
	if err := loadDataFromFile(statusMockFilePath, &mockStatus); err != nil {
		return nil, nil, nil, err
	}

	// Block mock data.
	blocksMockFilePath := fmt.Sprintf("%s/%s", mockData.Dir, mockData.BlockFile)
	var mockBlock pb.BlockData
	if err := loadDataFromFile(blocksMockFilePath, &mockBlock); err != nil {
		return nil, nil, nil, err
	}

	// Load trace mock data.
	traceMutex.Lock()
	defer traceMutex.Unlock()

	tracesMockFilePath := fmt.Sprintf("%s/%s", mockData.Dir, mockData.TraceFile)
	var mockTrace pb.Trace
	if err := loadDataFromFile(tracesMockFilePath, &mockTrace); err != nil {
		return nil, nil, nil, err
	}

	return &mockStatus, &mockBlock, &mockTrace, nil
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
