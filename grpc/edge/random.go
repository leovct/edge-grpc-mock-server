package edge

import (
	"crypto/rand"
	"math/big"
	"time"
	"zero-provers/server/grpc/edge/types"
)

// GenerateRandomEdgeTrace generates a random `Trace` with random data.
func GenerateRandomEdgeTrace(accountTriesAmount, storageTriesAmount, storageEntriesAmount, txnTracesAmount int) *types.Trace {
	trace := &types.Trace{
		AccountTrie:     make(map[string]string),
		StorageTrie:     make(map[string]string),
		ParentStateRoot: generateRandomHash(),
		TxnTraces:       []*types.TxnTrace{},
	}

	// Add some random accountTrie entries.
	for i := 0; i < accountTriesAmount; i++ {
		key := generateRandomHash()
		value := generateRandomHash()
		trace.AccountTrie[key.String()] = value.String()
	}

	// Add some random storageTrie entries.
	for i := 0; i < storageTriesAmount; i++ {
		key := generateRandomHash()
		value := generateRandomHash()
		trace.StorageTrie[key.String()] = value.String()
	}

	// Add some random TxnTraces.
	generateRandomBool := func() *bool {
		b, _ := rand.Int(rand.Reader, big.NewInt(2))
		res := b.Int64() == 1
		return &res
	}

	generateRandomNonce := func() *uint64 {
		n, _ := rand.Int(rand.Reader, big.NewInt(100))
		nonce := uint64(n.Uint64())
		return &nonce
	}

	generateRandomJournalEntry := func() *types.JournalEntry {
		entry := &types.JournalEntry{
			Addr:    generateRandomAddress(),
			Balance: generateRandomBigInt(),
			Nonce:   generateRandomNonce(),
			Storage: make(map[types.Hash]types.Hash),
			Code:    generateRandomBytes(64),
			Suicide: generateRandomBool(),
			Touched: generateRandomBool(),
		}

		// Add some random storage entries.
		for i := 0; i < storageEntriesAmount; i++ {
			key := generateRandomHash()
			value := generateRandomHash()
			entry.Storage[key] = value
		}

		return entry
	}

	generateRandomTxnTrace := func(nonce uint64) *types.TxnTrace {
		txn := generateRandomTx(nonce)
		return &types.TxnTrace{
			Transaction: txn.MarshalRLP(),
			Delta: map[types.Address]*types.JournalEntry{
				generateRandomAddress(): generateRandomJournalEntry(),
			},
		}
	}

	for i := 0; i < txnTracesAmount; i++ {
		trace.TxnTraces = append(trace.TxnTraces, generateRandomTxnTrace(uint64(i)))
	}

	return trace
}

// GenerateRandomEdgeBlock generates a random `Block` with random data.
func GenerateRandomEdgeBlock(number, txnTracesAmount uint64) *types.BlockGrpc {
	// Generate a random EdgeBlock.
	header := &types.Header{
		ParentHash:   generateRandomHash(),
		Sha3Uncles:   generateRandomHash(),
		Miner:        []byte{1, 2, 3},
		StateRoot:    generateRandomHash(),
		TxRoot:       generateRandomHash(),
		ReceiptsRoot: generateRandomHash(),
		LogsBloom:    types.Bloom{},
		Difficulty:   12345,
		Number:       number,
		GasLimit:     21000000,
		GasUsed:      200000,
		Timestamp:    uint64(time.Now().Unix()),
		ExtraData:    []byte{4, 5, 6},
		MixHash:      generateRandomHash(),
		Nonce:        types.Nonce{7, 8, 9, 10, 11, 12, 13, 14},
		Hash:         generateRandomHash(),
		BaseFee:      5,
	}

	// Generate a list of random transactions.
	var transactions []*types.TransactionGrpc
	var i uint64
	for i = 0; i < txnTracesAmount; i++ {
		transactions = append(transactions, generateRandomTx(i))
	}

	// Generate a list of random uncles.
	var uncles []*types.Header
	for i := 0; i < 2; i++ {
		uncles = append(uncles, &types.Header{
			ParentHash:   generateRandomHash(),
			Sha3Uncles:   generateRandomHash(),
			Miner:        []byte{1, 2, 3},
			StateRoot:    generateRandomHash(),
			TxRoot:       generateRandomHash(),
			ReceiptsRoot: generateRandomHash(),
			LogsBloom:    types.Bloom{},
			Difficulty:   12345,
			Number:       67890,
			GasLimit:     21000000,
			GasUsed:      200000,
			Timestamp:    uint64(time.Now().Unix()),
			ExtraData:    []byte{4, 5, 6},
			MixHash:      generateRandomHash(),
			Nonce:        types.Nonce{7, 8, 9, 10, 11, 12, 13, 14},
			Hash:         generateRandomHash(),
			BaseFee:      5,
		})
	}

	return &types.BlockGrpc{
		Header:       header,
		Transactions: transactions,
		Uncles:       uncles,
	}
}

func generateRandomTx(nonce uint64) *types.TransactionGrpc {
	randomAddress := generateRandomAddress()
	return &types.TransactionGrpc{
		Nonce:     nonce,
		GasPrice:  generateRandomBigInt(),
		GasTipCap: generateRandomBigInt(),
		GasFeeCap: generateRandomBigInt(),
		Gas:       21000,
		To:        &randomAddress,
		Value:     generateRandomBigInt(),
		Input:     []byte{1, 2, 3},
		V:         generateRandomBigInt(),
		R:         generateRandomBigInt(),
		S:         generateRandomBigInt(),
		Hash:      generateRandomHash(),
		From:      generateRandomAddress(),
		Type:      types.TxType(0),
		ChainID:   generateRandomBigInt(),
	}
}

func generateRandomHash() types.Hash {
	bytes := generateRandomBytes(types.HashLength)
	var hash types.Hash
	copy(hash[:], bytes)
	return hash
}

func generateRandomAddress() types.Address {
	bytes := generateRandomBytes(types.AddressLength)
	var address types.Address
	copy(address[:], bytes)
	return address
}

func generateRandomBigInt() *big.Int {
	n, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	return n
}

// generateRandomBytes generates a slice of random bytes with the given length.
func generateRandomBytes(length int) []byte {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return nil
	}
	return b
}
