#!/bin/bash

# Send gRPC requests to the mock server in a loop.
for ((i = 1; i <= $1; i++)); do
    grpcurl -plaintext 127.0.0.1:8546 v1.System/GetStatus | jq
    grpcurl -plaintext  -d '{"number": 1}' 127.0.0.1:8546 v1.System/BlockByNumber | jq
    grpcurl -plaintext  -d '{"number": 1}' 127.0.0.1:8546 v1.System/GetTrace | jq
done
