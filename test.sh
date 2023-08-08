#!/bin/bash

# Send a few gRPC/HTTP requests to the mock server.
# Make sure the server is started!
# go run main.go --http-port 8080 --http-save-endpoint /save --grpc-port 8546 --output-dir out

while ! echo -n > 127.0.0.1:8546; do
  echo "Mock server gRPC port is not open yet. Waiting for 5 seconds."
  sleep 5
done
rm 127.0.0.1:8546
echo "Mock server gRPC port is now open."

while ! echo -n > 127.0.0.1:8080; do
  echo "Mock server HTTP port is not open yet. Waiting for 5 seconds."
  sleep 5
done
rm 127.0.0.1:8080
echo "Mock server HTTP port is now open."

echo "Sending gRPC requests..."
grpcurl -plaintext  127.0.0.1:8546 v1.System/GetStatus	
grpcurl -plaintext  -d '{"number": 1}' 127.0.0.1:8546 v1.System/BlockByNumber
grpcurl -plaintext  -d '{"number": 1}' 127.0.0.1:8546 v1.System/GetTrace

echo "Sending HTTP requests..."
curl -X POST -H "Content-Type: application/json" -d '{"name": "salamander", "type": "fire"}' http://127.0.0.1:8080/save
cat out/1.json
