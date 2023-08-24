#!/bin/bash

# Send gRPC requests to the mock server in a loop.
n=100
for ((i=1; i<=$n; i++)); do
    grpcurl -plaintext 127.0.0.1:8546 v1.System/GetStatus | jq
    grpcurl -plaintext 127.0.0.1:8546 v1.System/GetStatus | jq
    grpcurl -plaintext 127.0.0.1:8546 v1.System/GetStatus | jq
    grpcurl -plaintext  -d '{"number": 1}' 127.0.0.1:8546 v1.System/BlockByNumber | jq
    grpcurl -plaintext  -d '{"number": 1}' 127.0.0.1:8546 v1.System/GetTrace | jq
done
