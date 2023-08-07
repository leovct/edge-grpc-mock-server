package main

import (
	"log"
	"zero-provers/server/grpc"
	"zero-provers/server/http"
)

const (
	gRPCServerPort  = 8546
	HTTPServerPort  = 8080
	ProofsOutputDir = "out"
)

func main() {
	// Start the gRPC server.
	go func() {
		log.Fatal(grpc.StartgRPCServer(gRPCServerPort), nil)
	}()

	// Start the HTTP server.
	log.Fatal(http.StartHTTPServer(HTTPServerPort, ProofsOutputDir), nil)
}
