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

	empty "google.golang.org/protobuf/types/known/emptypb"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	// Mock data file paths.
	StatusFile = "status.json"
	BlocksFile = "block.json"
	TraceFile  = "trace3.json"

	// Constant dummy block height returned by the `/GetStatus` endpoint.
	constantBlockHeight = 100_000_000_000_000_000
)

var (
	// log is the package-level variable used for logging messages and errors.
	log zerolog.Logger

	// Increase the block height on every GetStatus request made.
	counter     int
	counterStep int = 50

	// Mock data provided by the user.
	mockDir        string
	mockStatusData *pb.StatusResponse
	mockBlockData  *pb.BlockResponse
	mockTraceData  *pb.TraceResponse
)

// server is an internal implementation of the gRPC server.
type server struct {
	pb.UnimplementedSystemServer
}

// StartgRPCServer starts a gRPC server on the specified port.
// It listens for incoming TCP connections and handles gRPC requests using the internal server
// implementation. The server continues to run until it is manually stopped or an error occurs.
func StartgRPCServer(logLevel zerolog.Level, port int, setRandomMode bool, mockDataDir string) error {
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
		log.Info().Msgf("Fetching mock data from `%s` directory", mockDataDir)
		mockDir = mockDataDir
		mockStatusData, mockBlockData, mockTraceData, err = loadMockData()
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

// Load mock data if provided by the user.
func loadMockData() (*pb.StatusResponse, *pb.BlockResponse, *pb.TraceResponse, error) {
	// Load status mock data.
	var mockStatus pb.StatusResponse
	statusMockFilePath := fmt.Sprintf("%s/%s", mockDir, StatusFile)
	if _, err := os.Stat(statusMockFilePath); err == nil {
		data, err := os.ReadFile(statusMockFilePath)
		if err != nil {
			fmt.Println("Error reading mock status file:", err)
			return nil, nil, nil, err
		}

		if err := json.Unmarshal(data, &mockStatus); err != nil {
			fmt.Println("Error unmarshaling mock status JSON:", err)
			return nil, nil, nil, err
		}
		log.Info().Msg("Mock status data loaded")
	}

	// Load block mock data.
	var mockBlock pb.BlockResponse
	blocksMockFilePath := fmt.Sprintf("%s/%s", mockDir, BlocksFile)
	if _, err := os.Stat(blocksMockFilePath); err == nil {
		data, err := os.ReadFile(blocksMockFilePath)
		if err != nil {
			fmt.Println("Error reading mock blocks file:", err)
			return nil, nil, nil, err
		}

		if err := json.Unmarshal(data, &mockBlock); err != nil {
			fmt.Println("Error unmarshaling mock blocks JSON:", err)
			return nil, nil, nil, err
		}
		log.Info().Msg("Mock blocks data loaded")
	}

	// Load trace mock data.
	var mockTrace pb.TraceResponse
	tracesMockFilePath := fmt.Sprintf("%s/%s", mockDir, TraceFile)
	if _, err := os.Stat(tracesMockFilePath); err == nil {
		data, err := os.ReadFile(tracesMockFilePath)
		if err != nil {
			fmt.Println("Error reading mock traces file:", err)
			return nil, nil, nil, err
		}

		if err := json.Unmarshal(data, &mockTrace); err != nil {
			fmt.Println("Error unmarshaling mock traces JSON:", err)
			return nil, nil, nil, err
		}
		log.Info().Msg("Mock traces data loaded")
	}
	return &mockStatus, &mockBlock, &mockTrace, nil
}

// GetStatus is the implementation of the `GetStatus` RPC method.
// It returns a constant `ServerStatus` response.
func (s *server) GetStatus(context.Context, *empty.Empty) (*pb.StatusResponse, error) {
	log.Info().Msg("gRPC /GetStatus request received")

	// Return mock data if provided.
	if mockStatusData != nil {
		log.Info().Msgf("Mock StatusResponse number: %v", mockStatusData.Current.Number)
		return mockStatusData, nil
	}

	// Else, return dummy data.
	height := int64(constantBlockHeight + counter)
	counter += counterStep
	log.Info().Msgf("StatusResponse number: %v", height)
	return &pb.StatusResponse{
		Current: &pb.StatusResponse_Block{
			Number: height,
		},
	}, nil
}

// BlockByNumber is the implementation of the `BlockByNumber` RPC method.
// It returns a constant `BlockResponse` containing a single byte.
func (s *server) BlockByNumber(context.Context, *pb.BlockNumberRequest) (*pb.BlockResponse, error) {
	log.Info().Msg("gRPC /BlockByNumber request received")

	var rawData []byte
	if mockBlockData != nil {
		// Return mock data if provided.
		rawData = mockBlockData.Data
		log.Info().Msgf("Mock BlockResponse encoded data: %v", mockBlockData.Data)
	} else {
		// Else, return dummy data.
		height := constantBlockHeight + counter
		block := edge.GenerateDummyEdgeBlock(uint64(height))
		rawData = block.MarshalRLP()
		log.Info().Msgf("BlockResponse encoded data: %v", rawData)
	}

	// TODO: remove after debug session
	decodedBlock := edgetypes.Block{}
	if err := decodedBlock.UnmarshalRLP(rawData); err != nil {
		log.Error().Err(err).Msg("BlockResponse decoding failed")
		//return nil, err
	} else {
		data, err := json.MarshalIndent(decodedBlock, "", "  ")
		if err != nil {
			log.Error().Err(err).Msg("Unable to format JSON struct")
			//return nil, err
		} else {
			log.Info().Msg("BlockResponse decoded data")
			fmt.Println(string(data))
		}
	}

	return &pb.BlockResponse{
		Data: rawData,
	}, nil
}

func (s *server) GetTrace(context.Context, *pb.BlockNumberRequest) (*pb.TraceResponse, error) {
	log.Info().Msg("gRPC /GetTrace request received")

	var rawTrace []byte
	if mockTraceData != nil {
		// Return mock data if provided.
		log.Info().Msgf("Mock TraceResponse encoded data: %v", mockTraceData.Trace)
		rawTrace = mockTraceData.Trace
	} else {
		// Else, return dummy data.
		trace := *edge.GenerateDummyEdgeTrace()
		var err error
		rawTrace, err = json.Marshal(trace)
		if err != nil {
			fmt.Println("BlockTrace encoding failed:", err)
			return nil, err
		}
		log.Info().Msgf("TraceResponse encoded trace: %v", rawTrace)
	}

	// TODO: remove after debug session
	var decodedTrace *edgetypes.Trace
	if err := json.Unmarshal(rawTrace, &decodedTrace); err != nil {
		log.Error().Err(err).Msg("BlockTrace decoding failed")
		//return nil, err
	} else {
		data, err := json.MarshalIndent(decodedTrace, "", "  ")
		if err != nil {
			log.Error().Err(err).Msg("Unable to format JSON struct")
			//return nil, err
		} else {
			log.Info().Msg("TraceResponce decoded trace")
			fmt.Println(string(data))

			traces := decodedTrace.TxnTraces
			if len(traces) > 0 {
				log.Info().Msg("Decoding TraceResponce transactionTraces txn fields (RLP encoded)...")
			}
			for i, trace := range traces {
				decodedTxn := edgetypes.Transaction{}
				txnBytes := []byte(trace.Transaction)
				if err := decodedTxn.UnmarshalRLP(txnBytes); err != nil {
					log.Error().Err(err).Msgf("Transaction #%d decoding failed", i)
					return nil, err
				} else {
					data, err := json.MarshalIndent(decodedTxn, "", "  ")
					if err != nil {
						log.Error().Err(err).Msg("Unable to format JSON struct")
						//return nil, err
					} else {
						log.Info().Msgf("Transaction #%d decoded", i)
						fmt.Println(string(data))
					}
				}
			}
		}
	}

	return &pb.TraceResponse{
		Trace: rawTrace,
	}, nil
}
