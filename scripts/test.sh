#!/bin/bash

# Send a few gRPC/HTTP requests to the mock server.
# Make sure the server is started!
# go run main.go --http-port 8080 --http-save-endpoint /save --grpc-port 8546 --output-dir out

wait_for_service() {
  port=$1
  name=$2
  {
    while ! echo -n > /dev/tcp/127.0.0.1/$port; do
      echo "$name port is not open yet. Waiting for 5 seconds"
      sleep 5
    done
  } 2>/dev/null
  echo "$name port is now open."
}

echo "Waiting for services..."
wait_for_service 8546 "gRPC"
wait_for_service 8080 "HTTP"

echo "Sending gRPC requests..."
grpcurl -plaintext  127.0.0.1:8546 v1.System/GetStatus	
grpcurl -plaintext  -d '{"number": 1}' 127.0.0.1:8546 v1.System/BlockByNumber
grpcurl -plaintext  -d '{"number": 1}' 127.0.0.1:8546 v1.System/GetTrace

echo "Sending HTTP requests..."
curl -X POST -H "Content-Type: application/json" -d '{"name": "salamander", "type": "fire"}' http://127.0.0.1:8080/save
cat out/1.json
