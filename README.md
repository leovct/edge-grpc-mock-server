# Edge gRPC Mock Server

## Table of Contents

- [Introduction](#introduction)
- [Usage](#usage)
  - [Static mode (default)](#static-mode-default)
  - [Dynamic mode](#dynamic-mode)
  - [Random mode](#random-mode)
- [Use Case](#use-case)
  - [1. Start the mock server](#1-start-the-mock-server)
  - [2. Start the zero-prover setup](#2-start-the-zero-prover-setup)
  - [3. Benchmark proof generation time](#3-benchmark-proof-generation-time)
- [Datasets](#datasets)
- [Contributing](#contributing)

## Introduction

Simple mock of a [polygon-edge](https://github.com/0xPolygon/polygon-edge) gRPC server node meant to be used along a [zero-prover](https://github.com/mir-protocol/zero-provers) leader/worker setup. This component makes it easy to send specific blocks and traces to the zero-prover, to see how it behaves, without having to deploy an entire blockchain network such as edge.

It consists of two servers:

1. A gRPC server that mocks the functioning of an edge node. It only implements a subset of all the [methods](https://github.com/0xPolygon/polygon-edge/blob/feat/zero/server/proto/system.proto#L10) such as `GetStatus`, `BlockByNumber` and `GetTrace`. You can get the list of available methods using `make list` (make sure you started the server!). By default, the server returns mock data (see `data/` folder) but it can also be randomly generated using the `random` flag.

2. An HTTP server that either saves HTTP POST request data to the filesystem.

## Usage

```sh
$ go run main.go --help
Edge gRPC mock server

Usage:
  edge-grpc-mock-server [flags]

Flags:
  -g, --grpc-port int                       gRPC server port (default 8546)
  -h, --help                                help for edge-grpc-mock-server
  -p, --http-port int                       HTTP server port (default 8080)
  -e, --http-save-endpoint string           HTTP server save endpoint (default "/save")
      --mock-data-block-dir string          The mock data block directory (used in dynamic mode) (default "data/blocks")
      --mock-data-block-file string         The mock data block file path (used in static mode) (default "data/blocks/block-57.json")
      --mock-data-trace-dir string          The mock data trace directory (used in dynamic mode) (default "data/traces")
      --mock-data-trace-file string         The mock data trace file path (used in static mode) (default "data/traces/trace-57.json")
  -m, --mode string                         Mode of the mock server.
                                            - static: the server always return the same mock block data.
                                            - dynamic: the server returns new mock block data every {n} requests.
                                            - random: the server returns random block data every requests.
                                             (default "static")
  -o, --output-dir string                   The proofs output directory (default "out")
      --update-block-number-threshold int   The number of requests after which the server increments the block number (used in random mode) (default 30)
      --update-data-threshold int           The number of requests after which the server returns new data, block and trace (used in dynamic mode). (default 30)
  -v, --verbosity int8                      Verbosity level from 5 (panic) to -1 (trace) (default 1)
```

### Static mode (default)

In `static` mode, the server will always return the same mock data.

By default, it returns `data/blocks/block-57.json` and `data/traces/trace-57.json`.

```sh
go run main.go \
  --grpc-port 8546 \
  --http-port 8080 \
  --http-save-endpoint /save \
  --mock-data-block-file data/blocks/block-57.json \
  --mock-data-trace-file data/traces/trace-57.json \
  --mode static \
  --output-dir out \
  --verbosity 0
```

### Dynamic mode

In `dynamic` mode, the server returns dynamic mock data meaning after a certain number of requests, it will update the data it returns.

By default, the `--update-data-threshold` flag is set to 30 which means that the mock data will be updated each time the server receives 30 `/GetStatus` requests. Those requests are made by the zero-prover leader to check for new blocks.

The command also accepts directory flags as input, `--mock-data-block-dir` and `--mock-data-trace-dir`. In these folders, you should place all your mock block and trace files. The server will arrange the files in these directories in alphabetical order and will begin by providing the contents of the first files on the list. When the specified threshold for updating the data is reached, the server will increase the file index. It will continue to supply new block and trace files until no new files are available. After that point, it will consistently provide the last block and trace files in the list.

```sh
go run main.go \
  --grpc-port 8546 \
  --http-port 8080 \
  --http-save-endpoint /save \
  --mock-data-block-dir data/blocks \
  --mock-data-trace-dir data/traces \
  --mode dynamic \
  --update-data-threshold 30 \
  --output-dir out \
  --verbosity 0
```

### Random mode

In `random` mode, the server will generate and return random blocks and traces.

The server will accept an `-update-block-number-threshold` flag which represents the number of requests after which the server increments the block number. By default, it is set to 30.

```sh
go run main.go \
  --grpc-port 8546 \
  --http-port 8080 \
  --http-save-endpoint /save \
  --mode random \
  --update-block-number-threshold 30 \
  --output-dir out \
  --verbosity 0
```

## Use Case

### 1. Start the mock server

We use the `dynamic` mode of the mock server to be able to return dynamic block and trace mock data. You can use `go run main.go --help` to see the other options and the default values.

Here, the mock data will be updated every 30 `/GetStatus` requests received by the mock server. At the beginning, the mock server will return the first block and trace mock files of the directories. Then, after `n` requests, it will return the files at index `n`. Once the server has iterated over all the files, it will simply return the last block and trace mock files.

```sh
$ go run main.go \
  --grpc-port 8546 \
  --http-port 8080 \
  --http-save-endpoint /save \
  --mock-data-block-dir edge-grpc-mock-server/data/mock-uniswap-snowball-2/blocks \
  --mock-data-trace-dir edge-grpc-mock-server/data/mock-uniswap-snowball-2/traces \
  --mode dynamic \
  --update-data-threshold 30 \
  --output-dir out \
  --verbosity 0
Thu Sep  7 11:14:50 UTC 2023 INF http/http.go:63 > HTTP server is listening on port 8080
Thu Sep  7 11:14:50 UTC 2023 DBG http/http.go:64 > Config: {LogLevel:debug Port:8080 SaveEndpoint:/save ProofsOutputDir:out}
Thu Sep  7 11:14:50 UTC 2023 INF grpc/grpc.go:92 > gRPC server is listening on port 8546
Thu Sep  7 11:14:50 UTC 2023 DBG grpc/grpc.go:93 > Config: {LogLevel:debug Port:8546 Mode:dynamic UpdateDataThreshold:30 UpdateBlockNumberThreshold:30 MockData:{BlockDir:edge-grpc-mock-server/data/mock-uniswap-snowball-2/blocks TraceDir:edge-grpc-mock-server/data/mock-uniswap-snowball-2/traces BlockFile:data/blocks/block-57.json TraceFile:data/traces/trace-57.json}}
```

### 2. Start the zero-prover setup

First, start the zero-prover worker.

```sh
zero_prover_worker \
  --leader-notif-min-delay 1sec \
  -a http://127.0.0.1:9002 \
  -i 127.0.0.1 \
  -p 9002 \
  http://127.0.0.1:9001
```

Then, start the zero-prover leader.

```sh
$ zero_prover_leader \
    --secret-key-path prover.key \
    --contract-address 0x0000000000000000000000000000000000000000 \
    --rpc-url http://change_me.com \
    --full-node-endpoint http://127.0.0.1:8546 \
    --proof-complete-endpoint http://127.0.0.1:8080/save \
    --commit-height-delta-before-generating-proofs 0 \
    -i 127.0.0.1 \
    -p 9001
 INFO  zero_prover_leader::full_node_adapter::edge::edge_node_adapter > Received new block for height 57
Received payload for 57!
Starting proof for height 57...
Processing txn 0...
Partial s_root check for 0xdaa3…c8a1 == 0xdaa3…c8a1
BlockProofInitProofPayload { block_metadata: BlockMetadata { block_beneficiary: 0x719e1a8231e68b1cccdf5af1a1f496b3f317367f, block_timestamp: 1692971648, block_number: 57, block_difficulty: 1, block_gaslimit: 30000000, block_chain_id: 2001, block_base_fee: 0 }, skip_previous_block_proof: true, num_txns_in_block: 1 }
```

Soon, you will see that the leader sends gRPC `GetStatus` requests to the mock server.

```sh
Thu Sep  7 11:15:09 UTC 2023 INF edge-grpc-mock-server/grpc/grpc.go:103 > gRPC /GetStatus request received
Thu Sep  7 11:15:09 UTC 2023 DBG edge-grpc-mock-server/grpc/grpc.go:107 > Request counter: 1
Thu Sep  7 11:15:09 UTC 2023 DBG edge-grpc-mock-server/grpc/grpc.go:280 > Fetching mock data from edge-grpc-mock-server/data/mock-uniswap-snowball-2/blocks/block_1.json
Thu Sep  7 11:15:09 UTC 2023 DBG edge-grpc-mock-server/grpc/grpc.go:290 > Mock data loaded from edge-grpc-mock-server/data/mock-uniswap-snowball-2/blocks/block_1.json
Thu Sep  7 11:15:09 UTC 2023 DBG edge-grpc-mock-server/grpc/grpc.go:148 > StatusResponse number: 1
```

After some time, the leader will ask for the block metadata and the trace.

```sh
Thu Sep  7 11:15:09 UTC 2023 INF edge-grpc-mock-server/grpc/grpc.go:159 > gRPC /BlockByNumber request received
Thu Sep  7 11:15:09 UTC 2023 DBG edge-grpc-mock-server/grpc/grpc.go:280 > Fetching mock data from edge-grpc-mock-server/data/mock-uniswap-snowball-2/blocks/block_1.json
Thu Sep  7 11:15:09 UTC 2023 DBG edge-grpc-mock-server/grpc/grpc.go:280 > Fetching mock data from edge-grpc-mock-server/data/mock-uniswap-snowball-2/blocks/block_1.json
Thu Sep  7 11:15:09 UTC 2023 DBG edge-grpc-mock-server/grpc/grpc.go:290 > Mock data loaded from edge-grpc-mock-server/data/mock-uniswap-snowball-2/blocks/block_1.json
Thu Sep  7 11:15:09 UTC 2023 DBG edge-grpc-mock-server/grpc/grpc.go:290 > Mock data loaded from edge-grpc-mock-server/data/mock-uniswap-snowball-2/blocks/block_1.json
Thu Sep  7 11:15:09 UTC 2023 DBG edge-grpc-mock-server/grpc/grpc.go:202 > Decoded block header: {ParentHash:0xf25702fbe429a1f548c9ae8c346e641498eccdc4a8e2e1ad57c00b4490d3e79b Sha3Uncles:0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347 Miner:[70 197 150 32 144 185 61 96 80 181 234 193 35 65 185 205 175 6 235 35] StateRoot:0x0775babfe2f3484bc5152051b91d11367fe4eeafbd71bfd32798e8f2172bd134 TxRoot:0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421 ReceiptsRoot:0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421 LogsBloom:0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000 Difficulty:1 Number:1 GasLimit:30000000 GasUsed:0 Timestamp:1693509241 ExtraData:[0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 248 108 192 192 194 128 128 248 101 128 1 160 2 169 42 81 26 206 16 248 215 57 28 160 129 117 88 51 130 185 87 97 86 37 218 227 6 231 80 250 104 29 234 150 160 2 169 42 81 26 206 16 248 215 57 28 160 129 117 88 51 130 185 87 97 86 37 218 227 6 231 80 250 104 29 234 150 160 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0] MixHash:0xadce6e5230abe012342a44e4e9b6d05997d6f015387ae0e59be924afc7ec70c1 Nonce:0x0000000000000000 Hash:0xe7c6886ee88392ecced3ae3bce52ea0165da0287d75bf21e30bc4ab71a1c217e BaseFee:878822934}
Thu Sep  7 11:15:09 UTC 2023 DBG edge-grpc-mock-server/grpc/grpc.go:203 > Number of transactions: 0
Thu Sep  7 11:15:09 UTC 2023 DBG edge-grpc-mock-server/grpc/grpc.go:207 > Number of uncles: 0
```

```sh
Thu Sep  7 11:15:09 UTC 2023 INF edge-grpc-mock-server/grpc/grpc.go:221 > gRPC /GetTrace request received
Thu Sep  7 11:15:09 UTC 2023 DBG edge-grpc-mock-server/grpc/grpc.go:280 > Fetching mock data from edge-grpc-mock-server/data/mock-uniswap-snowball-2/traces/trace_1.json
Thu Sep  7 11:15:09 UTC 2023 DBG edge-grpc-mock-server/grpc/grpc.go:290 > Mock data loaded from edge-grpc-mock-server/data/mock-uniswap-snowball-2/traces/trace_1.json
```

Finally, you'll see this kind of message once the zero prover has finished generating the proof.

```sh
Thu Sep  7 11:15:09 UTC 2023 INF http/http.go:70 > POST request received on /save endpoint
Thu Sep  7 11:15:09 UTC 2023 INF http/http.go:94 > Proof saved to disk
```

Here is how it looks on the leader side after a few minutes.

```sh
Got a completed message for ProofKey { intern: ProofKeyIntern { b_height: 1, underlying_txns: ProofUnderlyingTxns { txn_idxs: 0..1 } }, p_type: Txn }
debug2: channel 0: window 999415 sent adjust 49161
Got a completed message for ProofKey { intern: ProofKeyIntern { b_height: 1, underlying_txns: ProofUnderlyingTxns { txn_idxs: 1..2 } }, p_type: Txn }
Got a completed message for ProofKey { intern: ProofKeyIntern { b_height: 1, underlying_txns: ProofUnderlyingTxns { txn_idxs: 0..2 } }, p_type: Agg }
Got a completed message for ProofKey { intern: ProofKeyIntern { b_height: 1, underlying_txns: ProofUnderlyingTxns { txn_idxs: 0..2 } }, p_type: Block }
Got a completed message for ProofKey { intern: ProofKeyIntern { b_height: 1, underlying_txns: ProofUnderlyingTxns { txn_idxs: 0..2 } }, p_type: CompressedBlock }
```

You can check the content of the proofs folder (by default, proofs are stored under `out/`). Note that the number of proofs may vary depending on how long you let the mock server run. To decode a proof simply use the following command: `cat <proof_file> | jq -r .trace | base64 -d | jq`.

### 3. Benchmark proof generation time

To assess the time required for the leader/worker configuration to produce proof for a specific trace, you can monitor logs.

When you observe the log entry `gRPC /GetTrace request received`, it signifies that the leader has initiated a request for the block trace. This happens after the leader has requested other details such as block metadata and has decided that it should generate a proof for a block at a given height. In this process, distinct tasks are assigned to the workers, which involve the generation of diverse types of proofs like transaction, aggregation, block, or compressed block proofs.

```sh
Thu Sep  7 11:15:09 UTC 2023 INF edge-grpc-mock-server/grpc/grpc.go:221 > gRPC /GetTrace request received
```

After the proof generation phase is concluded, you'll encounter the log entry `POST request received on /save endpoint`. At this point, the leader forwards the compressed block proof to the designated HTTP server.

```sh
Thu Sep  7 11:15:09 UTC 2023 INF http/http.go:77 > POST request received on /save endpoint
```

Given these logs, we can estimate the proof took approximately one minute to generate.

## Datasets

We provide different edge block and trace datasets to be used along a zero-prover setup under `data/archives/`. They have been manually generated using a real edge blockchain network and some load-testing tools like [`polycli loadtest`](https://github.com/maticnetwork/polygon-cli/blob/main/doc/polycli_loadtest.md). Some only include ERC721 mints while other include [Snowball](https://github.com/maticnetwork/jhilliard/blob/main/snowball/src/Snowball.sol) and Uniswap calls.

```sh
$ tree data/archives
data/archives
├── mock-erc721-mints.tar.bz2
├── mock-mix-and-uniswap.tar.bz2
├── mock-sstore-and-sha3.tar.bz2
└── mock-uniswap-snowball.tar.bz2

1 directory, 4 files
```

To extract those files and start using them, you can use the following command.

```sh
# For this example, we'll imagine you want to use the `mock-uniswap-snowball` dataset.
# We'll extract the content of the archive under `data/mock-uniswap-snowball`.
$ tar -xf data/archives/mock-uniswap-snowball.tar.bz2 -C data

$ tree -L 1 data
data
├── archives
├── blocks
├── mock-uniswap-snowball
└── traces

5 directories, 0 files
```

## Contributing

First, clone the repository.

```sh
git clone https://github.com/leovct/edge-grpc-mock-server && cd edge-grpc-mock-server
```

Install the [protobuf compiler](https://grpc.io/docs/protoc-installation/).

```sh
# On Ubuntu using `apt`.
apt install protobuf-compiler
protoc --version

# On MacOS using `homebrew`.
brew install protobuf
protoc --version
```

When making any change to the gRPC `System` service (see `grpc/pb/server.proto`), make sure to compile protocol buffers and generate the code with `make gen`. This step is very important to make sure the gRPC service is always up to date and aligned with the `proto` definition.

You can then run the server and experiment with it.

Use `go run main.go --help` to list all the different flags available.

Unit tests have not been implemented yet but you can run some HTTP/gRPC requests using [curl](https://curl.se/) and [grpcurl](https://github.com/fullstorydev/grpcurl) to test the behavior of the mock server. We provided a handy script called `scripts/test.sh` that you can execute using `make test` for this purpose.

To integrate last changes from `polygon-edge@feat/zero` [branch](https://github.com/0xPolygon/polygon-edge/tree/feat/zero), copy the last commit you want to use (i.e. [9071047](https://github.com/0xPolygon/polygon-edge/commit/907104765c64fae5cf4f2a40a8561c7ff6184058)). Then modify the `replace` statement at the end of `go.mod`. It should look like the following.

```diff
diff --git a/go.mod b/go.mod
index 183d3e7..c3731a8 100644
--- a/go.mod
+++ b/go.mod
@@ -34,4 +34,4 @@ require (

 // Use polygon-edge@feat/zero last commit.
 // https://github.com/0xPolygon/polygon-edge/tree/feat/zero
-replace github.com/0xPolygon/polygon-edge => github.com/0xPolygon/polygon-edge v1.1.1-0.20230929152933-907104765c64
+replace github.com/0xPolygon/polygon-edge => github.com/0xPolygon/polygon-edge 9071047
```

Finally, simply run `go mod tidy`. It will automatically update dependencies and reformat `go.mod`.

```go
replace github.com/0xPolygon/polygon-edge => github.com/0xPolygon/polygon-edge v1.1.1-0.20230929152933-907104765c64
```
