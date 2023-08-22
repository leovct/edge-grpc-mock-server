# Edge gRPC Mock Server

## Table of Contents

- [Introduction](#introduction)
- [Usage](#usage)
- [Use Case](#use-case)
- [Contributing](#contributing)

## Introduction

Simple mock of an [edge](https://github.com/0xPolygon/polygon-edge) gRPC server node meant to be used along a [zero-prover](https://github.com/mir-protocol/zero-provers) leader/worker setup.

It consists of two servers:

1. A gRPC server that mocks the functioning of an Edge node. It only implements a subset of all the [methods](https://github.com/0xPolygon/polygon-edge/blob/zero-trace/server/proto/system.proto) such as `GetStatus`, `GetTrace` and `BlockByNumber`. By default, the data is mocked (see `data/`) but it can also be randomly generated using the `random` flag.

2. An HTTP server that either displays HTTP POST request data or saves it to the filesystem.

## Usage

```sh
$ go run main.go --help
Edge gRPC mock server

Usage:
  edge-grpc-mock-server [flags]

Flags:
  -d, --debug                          Enable verbose mode
  -g, --grpc-port int                  gRPC server port (default 8546)
  -h, --help                           help for edge-grpc-mock-server
  -p, --http-port int                  HTTP server port (default 8080)
  -e, --http-save-endpoint string      HTTP server save endpoint (default "/save")
      --mock-data-block-file string    Mock data block file (in the mock data dir) (default "block.json")
  -m, --mock-data-dir string           Mock data directory (default "data")
      --mock-data-status-file string   Mock data status file (in the mock data dir) (default "status.json")
      --mock-data-trace-file string    Mock data trace file (in the mock data dir) (default "trace3.json")
  -o, --output-dir string              Proofs output directory (default "out")
  -r, --random                         Generate random trace data instead of relying on mocks (default false)
```

## Use Case

1. Start the mock server.

By default, the mock server will return mock data for status, block and trace. Use `go run main.go --help` to see the files loaded by default and check the `data/` folder to inspect the content of the mocks.

```sh
$ go run main.go \
  --debug \
  --grpc-port 8546 \
  --http-port 8080 \
  --http-save-endpoint /save \
  --mock-data-dir data \
  --output-dir out
Mon Aug 21 12:40:48 CEST 2023 INF http/http.go:56 > HTTP server save endpoint: /save ready
Mon Aug 21 12:40:48 CEST 2023 INF http/http.go:57 > HTTP server is starting on port 8080
Mon Aug 21 12:40:48 CEST 2023 DBG grpc/grpc.go:75 > Fetching mock data from `data` directory
Mon Aug 21 12:40:48 CEST 2023 DBG grpc/grpc.go:105 > Fetching mock status file from data/status.json
Mon Aug 21 12:40:48 CEST 2023 DBG grpc/grpc.go:108 > Mock status data loaded
Mon Aug 21 12:40:48 CEST 2023 DBG grpc/grpc.go:105 > Fetching mock block file from data/block.json
Mon Aug 21 12:40:48 CEST 2023 INF grpc/grpc.go:125 > Mock blocks data loaded
Mon Aug 21 12:40:48 CEST 2023 DBG grpc/grpc.go:105 > Fetching mock trace file from data/trace3.json
Mon Aug 21 12:40:48 CEST 2023 INF grpc/grpc.go:142 > Mock traces data loaded
Mon Aug 21 12:40:48 CEST 2023 INF grpc/grpc.go:84 > gRPC server is starting on port 8546
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
Mon Aug 21 12:40:48 CEST 2023 INF grpc/grpc.go:161 > gRPC /GetStatus request received
Mon Aug 21 12:40:48 CEST 2023 DBG grpc/grpc.go:165 > Mock StatusResponse number: 206672
...
Mon Aug 21 12:41:49 CEST 2023 INF grpc/grpc.go:183 > gRPC /BlockByNumber request received
Mon Aug 21 12:41:49 CEST 2023 DBG grpc/grpc.go:189 > Mock BlockResponse encoded data: [249 2 211 249 2 206 160 249 194 9 216 192 190 43 207 165 141 200 215 120 190 45 105 16 88 122 69 38 143 63 193 131 21 160 13 206 131 108 37 160 29 204 77 232 222 199 93 122 171 133 181 103 182 204 212 26 211 18 69 27 148 138 116 19 240 161 66 253 64 212 147 71 148 145 216 93 68 100 122 75 7 75 231 153 166 122 83 71 28 77 94 48 62 160 8 208 221 208 125 10 188 154 174 139 88 147 5 43 229 181 113 89 157 206 28 247 11 74 247 152 46 212 25 170 40 160 160 86 232 31 23 27 204 85 166 255 131 69 230 146 192 248 110 91 72 224 27 153 108 173 192 1 98 47 181 227 99 180 33 160 86 232 31 23 27 204 85 166 255 131 69 230 146 192 248 110 91 72 224 27 153 108 173 192 1 98 47 181 227 99 180 33 185 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 1 1 132 1 201 195 128 128 132 100 195 229 196 184 211 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 248 177 195 192 192 128 192 248 67 184 64 38 81 4 150 224 152 123 162 86 243 120 184 195 60 48 210 14 59 233 137 120 169 146 151 221 98 64 136 147 176 162 54 6 109 16 215 28 173 224 153 89 191 83 103 211 96 196 93 1 64 163 129 85 120 240 74 85 8 200 33 187 12 123 122 13 248 101 128 1 160 132 175 2 234 136 144 62 209 170 215 232 170 13 81 70 48 32 212 109 179 43 48 145 78 250 121 61 250 186 74 65 246 160 132 175 2 234 136 144 62 209 170 215 232 170 13 81 70 48 32 212 109 179 43 48 145 78 250 121 61 250 186 74 65 246 160 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 160 173 206 110 82 48 171 224 18 52 42 68 228 233 182 208 89 151 214 240 21 56 122 224 229 155 233 36 175 199 236 112 193 136 0 0 0 0 0 0 0 0 132 52 97 198 22 192 192]
Mon Aug 21 12:41:49 CEST 2023 DBG grpc/grpc.go:209 > BlockResponse decoded data: {
  "Header": {
    "ParentHash": "0xf9c209d8c0be2bcfa58dc8d778be2d6910587a45268f3fc18315a00dce836c25",
    "Sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
    "Miner": "kdhdRGR6SwdL55mmelNHHE1eMD4=",
    "StateRoot": "0x08d0ddd07d0abc9aae8b5893052be5b571599dce1cf70b4af7982ed419aa28a0",
    "TxRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
    "ReceiptsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
    "LogsBloom": [
      0,
      ...,
      0
    ],
    "Difficulty": 1,
    "Number": 1,
    "GasLimit": 30000000,
    "GasUsed": 0,
    "Timestamp": 1690559940,
    "ExtraData": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAD4scPAwIDA+EO4QCZRBJbgmHuiVvN4uMM8MNIOO+mJeKmSl91iQIiTsKI2Bm0Q1xyt4JlZv1Nn02DEXQFAo4FVePBKVQjIIbsMe3oN+GWAAaCErwLqiJA+0arX6KoNUUYwINRtsyswkU76eT36ukpB9qCErwLqiJA+0arX6KoNUUYwINRtsyswkU76eT36ukpB9qAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA==",
    "MixHash": "0xadce6e5230abe012342a44e4e9b6d05997d6f015387ae0e59be924afc7ec70c1",
    "Nonce": [
      0,
      ...,
      0
    ],
    "Hash": "0x9e6fa83b9754e8ccbc35ea6b7516c2df3e6d9224991ae03c2627d507863b2a9f",
    "baseFeePerGas": 0
  },
  "Transactions": null,
  "Uncles": null
}
...
Mon Aug 21 12:39:21 CEST 2023 INF grpc/grpc.go:220 > gRPC /GetTrace request received
Mon Aug 21 12:39:21 CEST 2023 DBG grpc/grpc.go:225 > Mock TraceResponse encoded data: [123 10 32 32 ... 93 10 125 10]
Mon Aug 21 12:39:21 CEST 2023 DBG grpc/grpc.go:250 > TraceResponce decoded trace: {
  "accountTrie": {
    "29fc6b8d0b979fe92b13053075fb24b19cb30a68cb043af8995c5ee7440f7aba": "f851808080a0b9b7e24499d6a857e9eb2e5ea5c8b30e068376c3519503472e41bbe5d947f8d98080808080a07f56f62b4ccf2e070be4d7c2e2c527367d6f5092f0fcbc504c96b77bcf24420b80808080808080",
    "66924a3ec75079b9da77770472bae89720cbed69f477de0793cbb27d48ecc57b": "f873a02080c7b7ae81a58eb98d9c78de4a1fd7fd9535fc953ed2be602daaa41767312ab850f84e808ad3c21bcecceda1000000a056e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421a0c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470",
    "8b19ac819d2a7610bce525fa551308dec3268e8ae20276b5b4414d24f21b109b": "f90191a03f84f317ceea01d67e9da286f143aea01d93cc275c3d9fad2722ce1aeb8a39bf8080a0da6faea54d8a227cc3d7a7e901ca20c63a333ccf60df131f36a6b02239983985a0c9c8c5324bd2eef39d061b5ad2409d077f4f9e03a6f2b89da8ffc03f77f8cbc5a0ecfb2faf207149229c871e123d5a6e90a27325f495dc7a679ef626a24f264c81a069c217c002a710db060939d0c0e132862fe78accf4d9bd0e4386b7a36834f195a09cb3c4dacefe3caac1388a8a39378db830fb2bf3983ccc5507d2d9728f67e676a029fc6b8d0b979fe92b13053075fb24b19cb30a68cb043af8995c5ee7440f7abaa0e6e82a109ac4e5daed8f76043d8393d0298472e0518ff96a48b8492e479881fba0b1f14f8b14818bd3b02ead0ecb2437bb07c5c3abcfba5be21ff33a572054b9c9a0aa83d885b05246cb62598c50ea169ab8e2495493cee532a4ee6a7571e3ee93d3a0bca602c40edd041089ea7c49cf65d6320427d889b547c88ca764e2d3ea77fb7f80a02c7fadb7de8ceb143581c40f241f304aabb0baca2368258b5c11f6f1d95053078080",
    "9b4ceef74bba9462847d8f4e4ad70505f36f9afe3a70b9580d0359894cfccff2": "f85180a0fd0838ecfbd0807d3b2cd7329e96f47c98dd9285e814ece852b12e9797b9ef6080808080808080808080a0b352e4c72d853473be5eac69cf10adfdc61a3747ec26a4b424f0a29292fc1d3780808080",
    "b1f14f8b14818bd3b02ead0ecb2437bb07c5c3abcfba5be21ff33a572054b9c9": "f85180a031d1eb10fdf37ef90417d463e9e3d1dad5bab6a3e8b38ccc4eaa744250317457808080a09b4ceef74bba9462847d8f4e4ad70505f36f9afe3a70b9580d0359894cfccff28080808080808080808080",
    "b9b7e24499d6a857e9eb2e5ea5c8b30e068376c3519503472e41bbe5d947f8d9": "f873a020e8b6a6f1f350490419424e9677d77d7356920536f752f0ed245592df04905db850f84e808ad3c21bcecceda1000000a056e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421a0c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470",
    "bca602c40edd041089ea7c49cf65d6320427d889b547c88ca764e2d3ea77fb7f": "f85180808080a06b795828f9b394ba531a3e97e87dce2b3e1153e14339019980089334951c7e28808080a02b552522c790b0433d466408318c5c61bfa7e2dbb9c2e13bfec3d1d51290addd8080808080808080",
    "ecfb2faf207149229c871e123d5a6e90a27325f495dc7a679ef626a24f264c81": "f871808080a066924a3ec75079b9da77770472bae89720cbed69f477de0793cbb27d48ecc57b80a03469ef1dfb42bfdc6b32c6352ae4e463fa021a1b03d000d33818e9a4ceb46b23808080a08bf437b987501fdf4f79bc77f3ea409f43bcae24574eecd7a8b99318d32ee18d80808080808080",
    "fd0838ecfbd0807d3b2cd7329e96f47c98dd9285e814ece852b12e9797b9ef60": "f8749f3eec2b84f0ba344fd4b4d2f022469febe7a772c4789acfc119eb558ab1da3db852f850808c033b2e3c9fd0803ce8000000a056e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421a0c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"
  },
  "storageTrie": null,
  "parentStateRoot": "0x8b19ac819d2a7610bce525fa551308dec3268e8ae20276b5b4414d24f21b109b",
  "transactionTraces": [
    {
      "txn": "0xf92f6...765b",
      "delta": {
        "0x0000000000000000000000000000000000000000": {
          "address": "0x0000000000000000000000000000000000000000",
          "read": true
        },
        "0x6FdA56C57B0Acadb96Ed5624aC500C0429d59429": {
          "address": "0x6FdA56C57B0Acadb96Ed5624aC500C0429d59429",
          "nonce": 1,
          "code": "YIBg...ADM=",
          "read": true
        },
        "0x84eb9227FCD22c94ED6e53Baf27C070018802D47": {
          "address": "0x84eb9227FCD22c94ED6e53Baf27C070018802D47",
          "read": true
        },
        "0x85dA99c8a7C2C95964c8EfD687E95E632Fc533D6": {
          "address": "0x85dA99c8a7C2C95964c8EfD687E95E632Fc533D6",
          "nonce": 1,
          "read": true
        }
      }
    }
  ]
}
Mon Aug 21 12:39:21 CEST 2023 INF grpc/grpc.go:255 > Decoding TraceResponce transactionTraces txn fields (RLP encoded)...
Mon Aug 21 12:39:21 CEST 2023 INF grpc/grpc.go:269 > Transaction #1 decoded
{
  "Nonce": 0,
  "GasPrice": 0,
  "GasTipCap": null,
  "GasFeeCap": null,
  "Gas": 2644387,
  "To": null,
  "Value": 0,
  "Input": "YIBg...Mw==",
  "V": 4038,
  "R": 107519757195806997439305138420673972387394027891232797455314865585635523889381,
  "S": 38851201517106677587590729829229679179700367655464385643439700858251583911515,
  "Hash": "0x2d6f97b42e8744513cfff5ba0f7ebbabb0644b41080a19f9c1c4b25cec82016f",
  "From": "0x0000000000000000000000000000000000000000",
  "Type": 0,
  "ChainID": null
}
...
Mon Aug 21 12:43:35 CEST 2023 INF http/http.go:70 > POST request received on /save endpoint
Mon Aug 21 12:43:35 CEST 2023 INF http/http.go:94 > Proof saved to disk
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

You can check the content of the proofs folder. Note that the number of proofs may vary depending on how long you let the mock server run.

```sh
# Check the content of the proofs folder.
$ tree out
out
├── 1.json
├── 2.json
├── 3.json
├── 4.json
└── 5.json

1 directory, 5 files

# Check the content of the first proof.
$ cat out/1.json
{
  "trace": "ewog...Cn0K"
}

# Decode the first proof.
# Note that this proof is very simple, if you want to see a more complex proof, check `data/trace2.json` and `data/decoded_trace2.json`.
$ cat out/1.json | jq -r .trace | base64 -d | jq
{
  "accountTrie": {
    "29fc6b8d0b979fe92b13053075fb24b19cb30a68cb043af8995c5ee7440f7aba": "f851808080a0b9b7e24499d6a857e9eb2e5ea5c8b30e068376c3519503472e41bbe5d947f8d98080808080a07f56f62b4ccf2e070be4d7c2e2c527367d6f5092f0fcbc504c96b77bcf24420b80808080808080",
    "66924a3ec75079b9da77770472bae89720cbed69f477de0793cbb27d48ecc57b": "f873a02080c7b7ae81a58eb98d9c78de4a1fd7fd9535fc953ed2be602daaa41767312ab850f84e808ad3c21bcecceda1000000a056e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421a0c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470",
    "8b19ac819d2a7610bce525fa551308dec3268e8ae20276b5b4414d24f21b109b": "f90191a03f84f317ceea01d67e9da286f143aea01d93cc275c3d9fad2722ce1aeb8a39bf8080a0da6faea54d8a227cc3d7a7e901ca20c63a333ccf60df131f36a6b02239983985a0c9c8c5324bd2eef39d061b5ad2409d077f4f9e03a6f2b89da8ffc03f77f8cbc5a0ecfb2faf207149229c871e123d5a6e90a27325f495dc7a679ef626a24f264c81a069c217c002a710db060939d0c0e132862fe78accf4d9bd0e4386b7a36834f195a09cb3c4dacefe3caac1388a8a39378db830fb2bf3983ccc5507d2d9728f67e676a029fc6b8d0b979fe92b13053075fb24b19cb30a68cb043af8995c5ee7440f7abaa0e6e82a109ac4e5daed8f76043d8393d0298472e0518ff96a48b8492e479881fba0b1f14f8b14818bd3b02ead0ecb2437bb07c5c3abcfba5be21ff33a572054b9c9a0aa83d885b05246cb62598c50ea169ab8e2495493cee532a4ee6a7571e3ee93d3a0bca602c40edd041089ea7c49cf65d6320427d889b547c88ca764e2d3ea77fb7f80a02c7fadb7de8ceb143581c40f241f304aabb0baca2368258b5c11f6f1d95053078080",
    "9b4ceef74bba9462847d8f4e4ad70505f36f9afe3a70b9580d0359894cfccff2": "f85180a0fd0838ecfbd0807d3b2cd7329e96f47c98dd9285e814ece852b12e9797b9ef6080808080808080808080a0b352e4c72d853473be5eac69cf10adfdc61a3747ec26a4b424f0a29292fc1d3780808080",
    "b1f14f8b14818bd3b02ead0ecb2437bb07c5c3abcfba5be21ff33a572054b9c9": "f85180a031d1eb10fdf37ef90417d463e9e3d1dad5bab6a3e8b38ccc4eaa744250317457808080a09b4ceef74bba9462847d8f4e4ad70505f36f9afe3a70b9580d0359894cfccff28080808080808080808080",
    "b9b7e24499d6a857e9eb2e5ea5c8b30e068376c3519503472e41bbe5d947f8d9": "f873a020e8b6a6f1f350490419424e9677d77d7356920536f752f0ed245592df04905db850f84e808ad3c21bcecceda1000000a056e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421a0c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470",
    "bca602c40edd041089ea7c49cf65d6320427d889b547c88ca764e2d3ea77fb7f": "f85180808080a06b795828f9b394ba531a3e97e87dce2b3e1153e14339019980089334951c7e28808080a02b552522c790b0433d466408318c5c61bfa7e2dbb9c2e13bfec3d1d51290addd8080808080808080",
    "ecfb2faf207149229c871e123d5a6e90a27325f495dc7a679ef626a24f264c81": "f871808080a066924a3ec75079b9da77770472bae89720cbed69f477de0793cbb27d48ecc57b80a03469ef1dfb42bfdc6b32c6352ae4e463fa021a1b03d000d33818e9a4ceb46b23808080a08bf437b987501fdf4f79bc77f3ea409f43bcae24574eecd7a8b99318d32ee18d80808080808080",
    "fd0838ecfbd0807d3b2cd7329e96f47c98dd9285e814ece852b12e9797b9ef60": "f8749f3eec2b84f0ba344fd4b4d2f022469febe7a772c4789acfc119eb558ab1da3db852f850808c033b2e3c9fd0803ce8000000a056e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421a0c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"
  },
  "storageTrie": null,
  "parentStateRoot": "0x8b19ac819d2a7610bce525fa551308dec3268e8ae20276b5b4414d24f21b109b",
  "transactionTraces": [
    {
      "txn": "0xf92f...765b",
      "delta": {
        "0x0000000000000000000000000000000000000000": {
          "address": "0x0000000000000000000000000000000000000000",
          "read": true
        },
        "0x6FdA56C57B0Acadb96Ed5624aC500C0429d59429": {
          "address": "0x6FdA56C57B0Acadb96Ed5624aC500C0429d59429",
          "nonce": 1,
          "code": "YIBg...ADM=",
          "read": true
        },
        "0x84eb9227FCD22c94ED6e53Baf27C070018802D47": {
          "address": "0x84eb9227FCD22c94ED6e53Baf27C070018802D47",
          "read": true
        },
        "0x85dA99c8a7C2C95964c8EfD687E95E632Fc533D6": {
          "address": "0x85dA99c8a7C2C95964c8EfD687E95E632Fc533D6",
          "nonce": 1,
          "read": true
        }
      }
    }
  ]
}
```

4. Finally, to assess the time required for the leader/worker configuration to produce a proof for a specific trace, you can monitor specific log entries:

When you observe the log entry `gRPC /GetTrace request received`, it signifies that the leader has initiated a request for the block trace. This happens after the leader has requested other details such as block metadata and has decided that it should generate a proof for a block at a given height. In this process, distinct tasks are assigned to the workers, which involve the generation of diverse types of proofs like transaction, aggregation, block, or compressed block proofs.

```sh
Fri Aug 18 14:08:15 CEST 2023 INF grpc/grpc.go:209 > gRPC /GetTrace request received
```

After the proof generation phase is concluded, you'll encounter the log entry `POST request received on /save endpoint`. At this point, the leader forwards the compressed block proof to the designated HTTP server.

```sh
Fri Aug 18 15:08:15 CEST 2023 INF http/http.go:70 > POST request received on /save endpoint
```

Given these logs, we can estimate the proof took approximately one minute to generate.

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
