#!/bin/bash

# Send a few gRPC/HTTP requests to the mock server.
# Make sure the server is started!
# $ go run main.go
set -x

echo "Sending gRPC requests..."
grpcurl -plaintext  127.0.0.1:8546 v1.System/GetStatus | jq
grpcurl -plaintext  -d '{"number": 1}' 127.0.0.1:8546 v1.System/BlockByNumber | jq
grpcurl -plaintext  -d '{"number": 1}' 127.0.0.1:8546 v1.System/GetTrace | jq
grpcurl -plaintext  -d '{"number": 1}' 127.0.0.1:8546 v1.System/GetTrace | jq -r .trace | base64 -d | jq

echo "Sending HTTP requests..."
curl -s -X POST -H "Content-Type: application/json" -d '{"name": "salamander", "type": "fire"}' http://127.0.0.1:8080/save
cat out/1.json | jq
