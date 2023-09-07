// Package grpc provides functionalities to start and handle a gRPC server.
package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"zero-provers/server/grpc/edge"
	pb "zero-provers/server/grpc/pb"
	"zero-provers/server/logger"
	"zero-provers/server/modes"

	empty "google.golang.org/protobuf/types/known/emptypb"

	"github.com/0xPolygon/polygon-edge/types"
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
	lock           sync.RWMutex
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
	log.Debug().Msgf("gRPC server config: %+v", config)
	log.Info().Msgf("gRPC server is listening on port %d", config.Port)
	if err := s.Serve(listener); err != nil {
		log.Error().Err(err).Msg("Unable to start gRPC server")
		return err
	}
	return nil
}

// GetStatus is the implementation of the `GetStatus` RPC method.
func (s *server) GetStatus(context.Context, *empty.Empty) (*pb.ChainStatus, error) {
	log.Info().Msg("gRPC /GetStatus request received")

	lock.Lock()
	requestCounter++
	log.Debug().Msgf("Request counter: %d", requestCounter)
	lock.Unlock()

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
		// List and sort the block mock files under the block mock directory.
		files, err := getFilesInDir(config.MockData.BlockDir)
		if err != nil {
			return nil, err
		}

		// Parse the block mock file at the current index and return the header number.
		lock.RLock()
		fileIndex := computeIndex(requestCounter, config.UpdateDataThreshold, len(files))
		lock.RUnlock()
		file := files[fileIndex]
		height, err = getBlockNumberFromBlockFile(file)
		if err != nil {
			return nil, err
		}

	case modes.RandomMode:
		// Increment the constant block number based on request counter.
		lock.RLock()
		height = int64(constantBlockHeight + requestCounter%config.UpdateBlockNumberThreshold)
		lock.RUnlock()

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
	var block *types.Block
	switch config.Mode {
	case modes.StaticMode:
		// Parse the block mock file in the edge RPC format and convert it to the GRPC format.
		var mockBlockRPC edge.BlockRPC
		if err := loadDataFromFile(config.MockData.BlockFile, &mockBlockRPC); err != nil {
			return nil, err
		}
		block = mockBlockRPC.ToBlockGrpc()

	case modes.DynamicMode:
		// List and sort the block mock files under the block mock directory.
		files, err := getFilesInDir(config.MockData.BlockDir)
		if err != nil {
			return nil, err
		}

		// Parse the block mock file at the current index and convert it to the GRPC format.
		lock.RLock()
		fileIndex := computeIndex(requestCounter, config.UpdateDataThreshold, len(files))
		lock.RUnlock()
		file := files[fileIndex]
		var mockBlockRPC edge.BlockRPC
		if err := loadDataFromFile(file, &mockBlockRPC); err != nil {
			return nil, err
		}
		block = mockBlockRPC.ToBlockGrpc()

	case modes.RandomMode:
		// Return a random block data.
		lock.RLock()
		height := uint64(constantBlockHeight + requestCounter%config.UpdateBlockNumberThreshold)
		lock.RUnlock()
		txnTracesAmount := uint64(10)
		block = edge.GenerateRandomEdgeBlock(height, txnTracesAmount)

	default:
		return nil, errWrongMode
	}
	if block != nil {
		log.Debug().Msgf("Decoded block header: %+v", *block.Header)
		log.Debug().Msgf("Number of transactions: %d", len(block.Transactions))
		for i, tx := range block.Transactions {
			log.Debug().Msgf("Tx #%d: %+v", i, tx)
		}
		log.Debug().Msgf("Number of uncles: %d", len(block.Uncles))
		for i, uncle := range block.Uncles {
			log.Debug().Msgf("Uncle #%d: %+v", i, uncle)
		}
	}

	// Encode the block using RLP.
	encodedBlock := block.MarshalRLP()
	return &pb.BlockData{
		Data: encodedBlock,
	}, nil
}

func (s *server) GetTrace(context.Context, *pb.BlockNumber) (*pb.Trace, error) {
	log.Info().Msg("gRPC /GetTrace request received")

	// Load trace data from file or generate random data.
	var trace types.Trace
	switch config.Mode {
	case modes.StaticMode:
		// Parse the decoded trace mock file.
		if err := loadDataFromFile(config.MockData.TraceFile, &trace); err != nil {
			return nil, err
		}

	case modes.DynamicMode:
		// List the block trace files under the trace mock directory.
		files, err := getFilesInDir(config.MockData.TraceDir)
		if err != nil {
			return nil, err
		}

		// Parse the decoded trace mock file at the current index.
		lock.RLock()
		fileIndex := computeIndex(requestCounter, config.UpdateDataThreshold, len(files))
		lock.RUnlock()
		file := files[fileIndex]
		if err := loadDataFromFile(file, &trace); err != nil {
			return nil, err
		}

	case modes.RandomMode:
		trace = *edge.GenerateRandomEdgeTrace(10, 10, 10, 10)

	default:
		return nil, errWrongMode
	}
	log.Trace().Msgf("Decoded trace: %+v", trace)

	// Encode the trace using base64.
	encodedTrace, err := json.Marshal(trace)
	if err != nil {
		log.Error().Err(err).Msg("Trace encoding failed")
		return nil, err
	}
	return &pb.Trace{
		Trace: encodedTrace,
	}, nil
}

// Compute the file index in the case of dynamic mode.
// Iterate over all the indexes and if the index is greater than the number of files, return
// the index of the last file.
func computeIndex(counter, updateThreshold int, numberOfFiles int) int {
	index := (counter - 1) / updateThreshold
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

// Return file names in directory and sorted in natural order.
func getFilesInDir(dirPath string) ([]string, error) {
	files, err := filepath.Glob(filepath.Join(dirPath, "*.json"))
	if err != nil {
		return nil, err
	}

	sort.Slice(files, func(i, j int) bool {
		return naturalSort(files[i], files[j])
	})

	return files, nil
}

// Load block number from block file.
func getBlockNumberFromBlockFile(filePath string) (int64, error) {
	var mockBlockRPC edge.BlockRPC
	if err := loadDataFromFile(filePath, &mockBlockRPC); err != nil {
		return 0, err
	}
	return int64(mockBlockRPC.Number), nil
}

// Sort files in natural order.
func naturalSort(s1, s2 string) bool {
	parts1 := strings.FieldsFunc(s1, func(r rune) bool {
		return !((r >= '0' && r <= '9') || (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z'))
	})
	parts2 := strings.FieldsFunc(s2, func(r rune) bool {
		return !((r >= '0' && r <= '9') || (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z'))
	})

	for i := 0; i < len(parts1) && i < len(parts2); i++ {
		part1, part2 := parts1[i], parts2[i]
		if part1 != part2 {
			isDigit1 := part1[0] >= '0' && part1[0] <= '9'
			isDigit2 := part2[0] >= '0' && part2[0] <= '9'

			if isDigit1 && isDigit2 {
				num1, _ := strconv.Atoi(part1)
				num2, _ := strconv.Atoi(part2)
				return num1 < num2
			}

			return part1 < part2
		}
	}

	return len(parts1) < len(parts2)
}
