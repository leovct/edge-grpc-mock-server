# Edge gRPC Mock Server

## Table of Contents

- [Introduction](#introduction)
- [Usage](#usage)
- [Use Case](#use-case)
- [Contributing](#contributing)

## Introduction

Simple mock of an [edge](https://github.com/0xPolygon/polygon-edge) gRPC server node meant to be used along a [zero-prover](https://github.com/mir-protocol/zero-provers) leader/worker setup.

It consists of two servers:

1. A gRPC server that mocks the functioning of an Edge node. It only implements a subset of all the [methods](https://github.com/0xPolygon/polygon-edge/blob/zero-trace/server/proto/system.proto) such as `GetStatus`, `GetTrace` and `BlockByNumber`. Most of the data returned is randomly generated.

2. An HTTP server that either displays HTTP POST request data or saves it to the filesystem.

## Usage

```sh
$ go run main.go --help
Edge gRPC mock server

Usage:
  mock [flags]

Flags:
  -d, --debug                       Enable verbose mode
  -g, --grpc-port int               gRPC server port (default 8546)
  -h, --help                        help for mock
  -p, --http-port int               HTTP server port (default 8080)
  -e, --http-save-endpoint string   HTTP server save endpoint (default "/save")
  -m, --mock-data-dir string        Mock data directory containing mock status (status.json), block (block.json) and trace (trace.json) files (default "data")
  -o, --output-dir string           Proofs output directory (default "out")
```

## Use Case

1. Start the mock server.

```sh
$ go run main.go \
  --http-port 8080 \
  --http-save-endpoint /save \
  --grpc-port 8546 \
  --output-dir out \
  --mock-data-dir data
Tue Aug  8 13:35:00 CEST 2023 INF http/http.go:56 > HTTP server save endpoint: /save ready
Tue Aug  8 13:35:00 CEST 2023 INF http/http.go:57 > HTTP server is starting on port 8080
Tue Aug  8 13:35:00 CEST 2023 INF grpc/grpc.go:75 > Fetching mock data from `data` directory
Tue Aug  8 13:35:00 CEST 2023 INF grpc/grpc.go:108 > Mock status data loaded
Tue Aug  8 13:35:00 CEST 2023 INF grpc/grpc.go:125 > Mock blocks data loaded
Tue Aug  8 13:35:00 CEST 2023 INF grpc/grpc.go:142 > Mock traces data loaded
Tue Aug  8 13:35:00 CEST 2023 INF grpc/grpc.go:84 > gRPC server is starting on port 8546
```

2. Start a zero worker.

```sh
zero_prover_worker \
  --leader-notif-min-delay 1sec \
  -a http://127.0.0.1:9002 \
  -i 127.0.0.1 \
  -p 9002 \
  http://127.0.0.1:9001
```

3. Start the zero leader.

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
Received payload for 206672!
Starting proof for height 206672...
BlockProofInitProofPayload { block_metadata: BlockMetadata { block_beneficiary: 0x91d85d44647a4b074be799a67a53471c4d5e303e, block_timestamp: 1690559940, block_number: 1, block_difficulty: 1, block_gaslimit: 30000000, block_chain_id: 2001, block_base_fee: 878822934 }, skip_previous_block_proof: true, num_txns_in_block: 0 }
```

4. Soon, you will see that the leader sends gRPC requests to the mock server.

```sh
...
Tue Aug  8 13:38:00 CEST 2023 INF grpc/grpc.go:150 > gRPC /GetStatus request received
Tue Aug  8 13:38:00 CEST 2023 INF grpc/grpc.go:154 > Mock StatusResponse number: 206672
...
Tue Aug  8 13:39:00 CEST 2023 INF grpc/grpc.go:209 > gRPC /GetTrace request received
Tue Aug  8 13:39:00 CEST 2023 INF grpc/grpc.go:214 > Mock TraceResponse encoded data: [123 34 97 99 99 111 117 110 116 84 114 105 101 34 58 110 117 108 108 44 34 115 116 111 114 97 103 101 84 114 105 101 34 58 110 117 108 108 44 34 112 97 114 101 110 116 83 116 97 116 101 82 111 111 116 34 58 34 48 120 48 56 100 48 100 100 100 48 55 100 48 97 98 99 57 97 97 101 56 98 53 56 57 51 48 53 50 98 101 53 98 53 55 49 53 57 57 100 99 101 49 99 102 55 48 98 52 97 102 55 57 56 50 101 100 52 49 57 97 97 50 56 97 48 34 44 34 116 114 97 110 115 97 99 116 105 111 110 84 114 97 99 101 115 34 58 91 93 125]
Tue Aug  8 13:39:00 CEST 2023 INF grpc/grpc.go:239 > TraceResponce decoded trace
{
  "accountTrie": null,
  "storageTrie": null,
  "parentStateRoot": "0x08d0ddd07d0abc9aae8b5893052be5b571599dce1cf70b4af7982ed419aa28a0",
  "transactionTraces": []
}
...
Tue Aug  8 13:40:00 CEST 2023 INF grpc/grpc.go:172 > gRPC /BlockByNumber request received
Tue Aug  8 13:40:00 CEST 2023 INF grpc/grpc.go:178 > Mock BlockResponse encoded data: [249 2 211 249 2 206 160 249 194 9 216 192 190 43 207 165 141 200 215 120 190 45 105 16 88 122 69 38 143 63 193 131 21 160 13 206 131 108 37 160 29 204 77 232 222 199 93 122 171 133 181 103 182 204 212 26 211 18 69 27 148 138 116 19 240 161 66 253 64 212 147 71 148 145 216 93 68 100 122 75 7 75 231 153 166 122 83 71 28 77 94 48 62 160 8 208 221 208 125 10 188 154 174 139 88 147 5 43 229 181 113 89 157 206 28 247 11 74 247 152 46 212 25 170 40 160 160 86 232 31 23 27 204 85 166 255 131 69 230 146 192 248 110 91 72 224 27 153 108 173 192 1 98 47 181 227 99 180 33 160 86 232 31 23 27 204 85 166 255 131 69 230 146 192 248 110 91 72 224 27 153 108 173 192 1 98 47 181 227 99 180 33 185 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 1 1 132 1 201 195 128 128 132 100 195 229 196 184 211 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 248 177 195 192 192 128 192 248 67 184 64 38 81 4 150 224 152 123 162 86 243 120 184 195 60 48 210 14 59 233 137 120 169 146 151 221 98 64 136 147 176 162 54 6 109 16 215 28 173 224 153 89 191 83 103 211 96 196 93 1 64 163 129 85 120 240 74 85 8 200 33 187 12 123 122 13 248 101 128 1 160 132 175 2 234 136 144 62 209 170 215 232 170 13 81 70 48 32 212 109 179 43 48 145 78 250 121 61 250 186 74 65 246 160 132 175 2 234 136 144 62 209 170 215 232 170 13 81 70 48 32 212 109 179 43 48 145 78 250 121 61 250 186 74 65 246 160 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 160 173 206 110 82 48 171 224 18 52 42 68 228 233 182 208 89 151 214 240 21 56 122 224 229 155 233 36 175 199 236 112 193 136 0 0 0 0 0 0 0 0 132 52 97 198 22 192 192]
Tue Aug  8 13:40:00 CEST 2023 INF grpc/grpc.go:198 > BlockResponse decoded data
{
  "Header": {
    "ParentHash": "0xf9c209d8c0be2bcfa58dc8d778be2d6910587a45268f3fc18315a00dce836c25",
    "Sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
    "Miner": "kdhdRGR6SwdL55mmelNHHE1eMD4=",
    "StateRoot": "0x08d0ddd07d0abc9aae8b5893052be5b571599dce1cf70b4af7982ed419aa28a0",
    "TxRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
    "ReceiptsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
    "LogsBloom": [
      ...
    ],
    "Difficulty": 1,
    "Number": 1,
    "GasLimit": 30000000,
    "GasUsed": 0,
    "Timestamp": 1690559940,
    "ExtraData": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAD4scPAwIDA+EO4QCZRBJbgmHuiVvN4uMM8MNIOO+mJeKmSl91iQIiTsKI2Bm0Q1xyt4JlZv1Nn02DEXQFAo4FVePBKVQjIIbsMe3oN+GWAAaCErwLqiJA+0arX6KoNUUYwINRtsyswkU76eT36ukpB9qCErwLqiJA+0arX6KoNUUYwINRtsyswkU76eT36ukpB9qAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA==",
    "MixHash": "0xadce6e5230abe012342a44e4e9b6d05997d6f015387ae0e59be924afc7ec70c1",
    "Nonce": [
      ...
    ],
    "Hash": "0x9e6fa83b9754e8ccbc35ea6b7516c2df3e6d9224991ae03c2627d507863b2a9f",
    "baseFeePerGas": 0
  },
  "Transactions": null,
  "Uncles": null
}
...
Tue Aug  8 13:39:00 CEST 2023 INF http/http.go:70 > POST request received on /save endpoint
Tue Aug  8 13:39:00 CEST 2023 INF http/http.go:94 > Proof saved to disk
```

Here is how it looks on the leader side after a few minutes.

```sh
Got a got completed message for ProofKey { intern: ProofKeyIntern { b_height: 1, underlying_txns: ProofUnderlyingTxns { txn_idxs: 0..1 } }, p_type: Txn }
debug2: channel 0: window 999415 sent adjust 49161
Got a got completed message for ProofKey { intern: ProofKeyIntern { b_height: 1, underlying_txns: ProofUnderlyingTxns { txn_idxs: 1..2 } }, p_type: Txn }
Got a got completed message for ProofKey { intern: ProofKeyIntern { b_height: 1, underlying_txns: ProofUnderlyingTxns { txn_idxs: 0..2 } }, p_type: Agg }
Got a got completed message for ProofKey { intern: ProofKeyIntern { b_height: 1, underlying_txns: ProofUnderlyingTxns { txn_idxs: 0..2 } }, p_type: Block }
Got a got completed message for ProofKey { intern: ProofKeyIntern { b_height: 1, underlying_txns: ProofUnderlyingTxns { txn_idxs: 0..2 } }, p_type: CompressedBlock }
```

You can check the content of the proofs folder. Note that the number of proofs may vary depending on how long you let the mock server run.

```sh
$ tree out
out
├── 1.json
├── 2.json
├── 3.json
├── 4.json
└── 5.json

1 directory, 5 files

$ cat out/1.json
{
  "trace": "eyJhY2NvdW50VHJpZSI6bnVsbCwic3RvcmFnZVRyaWUiOm51bGwsInBhcmVudFN0YXRlUm9vdCI6IjB4MDhkMGRkZDA3ZDBhYmM5YWFlOGI1ODkzMDUyYmU1YjU3MTU5OWRjZTFjZjcwYjRhZjc5ODJlZDQxOWFhMjhhMCIsInRyYW5zYWN0aW9uVHJhY2VzIjpbXX0="
}
```

## Contributing

First, clone the repository.

```sh
git clone https://github.com/leovct/edge-grpc-mock-server && cd edge-grpc-mock-server
```

Install the [protobuf compiler](https://grpc.io/docs/protoc-installation/).

```sh
# On Ubuntu using `apt`.
apt-get install protobuf-compiler
protoc --version

# On MacOS using `homebrew`.
brew install protobuf
protoc --version
```

When making any change to the gRPC `System` service (see `grpc/pb/server.proto`), make sure to compile protocol buffers and generate the code with `make gen`. This step is very important to make sure the gRPC service is always up to date and aligned with the `proto` definition.

You can then run the server and experiment with it.

Use `go run main.go --help` to list all the different flags available.

Unit tests have not been implemented yet but you can run some HTTP/gRPC requests using [curl](https://curl.se/) and [grpcurl](https://github.com/fullstorydev/grpcurl) to test the behavior of the mock server. We provided a handy script called `test.sh` that you can execute using `make test` for this purpose.

