#!/bin/bash
set -x

# Send a few gRPC/HTTP requests to the mock server.
# Make sure the server is started!
# go run main.go --http-port 8080 --http-save-endpoint /save --grpc-port 8546 --output-dir out
echo "Sending gRPC requests..."
grpcurl -plaintext  127.0.0.1:8546 v1.System/GetStatus | jq
grpcurl -plaintext  -d '{"number": 1}' 127.0.0.1:8546 v1.System/BlockByNumber | jq

# Get trace
grpcurl -plaintext  -d '{"number": 1}' 127.0.0.1:8546 v1.System/GetTrace | jq
grpcurl -plaintext  -d '{"number": 1}' 127.0.0.1:8546 v1.System/GetTrace | jq -r .trace | base64 -d | jq

# Update trace
grpcurl -plaintext -d "{\"trace\": $(cat data/trace1.json | jq .trace)}" 127.0.0.1:8546 v1.System/UpdateTrace | jq
grpcurl -plaintext  -d '{"number": 1}' 127.0.0.1:8546 v1.System/GetTrace | jq
grpcurl -plaintext  -d '{"number": 1}' 127.0.0.1:8546 v1.System/GetTrace | jq -r .trace | base64 -d | jq

echo "Sending HTTP requests..."
curl -s -X POST -H "Content-Type: application/json" -d '{"name": "salamander", "type": "fire"}' http://127.0.0.1:8080/save
cat out/1.json | jq
