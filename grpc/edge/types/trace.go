package types

import (
	"encoding/hex"
	"math/big"
	"strings"
)

type ArgBytes []byte

func (b ArgBytes) MarshalText() ([]byte, error) {
	return encodeToHex(b), nil
}

func (b *ArgBytes) UnmarshalText(input []byte) error {
	hh, err := decodeToHex(input)
	if err != nil {
		return nil
	}

	aux := make([]byte, len(hh))
	copy(aux[:], hh[:])
	*b = aux

	return nil
}

func decodeToHex(b []byte) ([]byte, error) {
	str := string(b)
	str = strings.TrimPrefix(str, "0x")

	if len(str)%2 != 0 {
		str = "0" + str
	}

	return hex.DecodeString(str)
}

func encodeToHex(b []byte) []byte {
	str := hex.EncodeToString(b)
	if len(str)%2 != 0 {
		str = "0" + str
	}

	return []byte("0x" + str)
}

type Trace struct {
	// AccountTrie is the partial trie for the account merkle trie touched during the block
	AccountTrie map[string]string `json:"accountTrie"`

	// StorageTrie is the partial trie for the storage tries touched during the block
	StorageTrie map[string]string `json:"storageTrie"`

	// ParentStateRoot is the parent state root for this block
	ParentStateRoot Hash `json:"parentStateRoot"`

	// TxnTraces is the list of traces per transaction in the block
	TxnTraces []*TxnTrace `json:"transactionTraces"`
}

type TxnTrace struct {
	// Transaction is the RLP encoding of the transaction
	Transaction ArgBytes `json:"txn"`

	// Delta is the list of updates per account during this transaction
	Delta map[Address]*JournalEntry `json:"delta"`
}

type JournalEntry struct {
	// Addr is the address of the account affected by the
	// journal change
	Addr Address `json:"address"`

	// Balance tracks changes in the account Balance
	Balance *big.Int `json:"balance,omitempty"`

	// Nonce tracks changes in the account Nonce
	Nonce *uint64 `json:"nonce,omitempty"`

	// Storage track changes in the storage
	Storage map[Hash]Hash `json:"storage,omitempty"`

	// Code tracks the initialization of the contract Code
	Code []byte `json:"code,omitempty"`

	// Suicide tracks whether the contract has been self destructed
	Suicide *bool `json:"suicide,omitempty"`

	// Touched tracks whether the account has been touched/created
	Touched *bool `json:"touched,omitempty"`
}
