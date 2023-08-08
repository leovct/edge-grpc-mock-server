#!/bin/bash

# Send a few gRPC/HTTP requests to the mock server.
# Make sure the server is started!
# go run main.go --http-port 8080 --http-save-endpoint /save --grpc-port 8546 --output-dir out

echo "Sending gRPC requests..."
grpcurl -plaintext  127.0.0.1:8546 v1.System/GetStatus	
grpcurl -plaintext  -d '{"number": 1}' 127.0.0.1:8546 v1.System/BlockByNumber
grpcurl -plaintext  -d '{"number": 1}' 127.0.0.1:8546 v1.System/GetTrace

echo "Sending HTTP requests..."
curl -X POST -H "Content-Type: application/json" -d '{"name": "salamander", "type": "fire"}' http://127.0.0.1:8080/save
cat out/1.json
