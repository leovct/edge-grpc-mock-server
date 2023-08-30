package edge

import (
	"encoding/hex"
	"math/big"
	"strconv"
	"strings"

	"github.com/0xPolygon/polygon-edge/types"
)

// BlockRPC represents a block returned by the edge RPC.
type BlockRPC struct {
	ParentHash      types.Hash       `json:"parentHash"`
	Sha3Uncles      types.Hash       `json:"sha3Uncles"`
	Miner           argBytes         `json:"miner"`
	StateRoot       types.Hash       `json:"stateRoot"`
	TxRoot          types.Hash       `json:"transactionsRoot"`
	ReceiptsRoot    types.Hash       `json:"receiptsRoot"`
	LogsBloom       types.Bloom      `json:"logsBloom"`
	Difficulty      argUint64        `json:"difficulty"`
	TotalDifficulty argUint64        `json:"totalDifficulty"`
	Size            argUint64        `json:"size"`
	Number          argUint64        `json:"number"`
	GasLimit        argUint64        `json:"gasLimit"`
	GasUsed         argUint64        `json:"gasUsed"`
	Timestamp       argUint64        `json:"timestamp"`
	ExtraData       argBytes         `json:"extraData"`
	MixHash         types.Hash       `json:"mixHash"`
	Nonce           types.Nonce      `json:"nonce"`
	Hash            types.Hash       `json:"hash"`
	Transactions    []TransactionRPC `json:"transactions"`
	Uncles          []types.Hash     `json:"uncles"`
	BaseFee         argUint64        `json:"baseFeePerGas,omitempty"`
}

func (b *BlockRPC) ToBlockGrpc() *types.Block {
	header := types.Header{
		ParentHash:   b.ParentHash,
		Sha3Uncles:   b.Sha3Uncles,
		Miner:        b.Miner,
		StateRoot:    b.StateRoot,
		TxRoot:       b.TxRoot,
		ReceiptsRoot: b.ReceiptsRoot,
		LogsBloom:    b.LogsBloom,
		Difficulty:   uint64(b.Difficulty),
		Number:       uint64(b.Number),
		GasLimit:     uint64(b.GasLimit),
		GasUsed:      uint64(b.GasUsed),
		Timestamp:    uint64(b.Timestamp),
		ExtraData:    b.ExtraData,
		MixHash:      b.MixHash,
		Nonce:        b.Nonce,
		Hash:         b.Hash,
		BaseFee:      uint64(b.BaseFee),
	}

	transactions := make([]*types.Transaction, len(b.Transactions))
	for _, txGrpc := range b.Transactions {
		txRPC := txGrpc.toTransactionGrpc()
		if txRPC != nil {
			transactions = append(transactions, txRPC)
		}
	}

	// Note: we don't parse uncles for the moment.
	var uncles []*types.Header

	return &types.Block{
		Header:       &header,
		Transactions: transactions,
		Uncles:       uncles,
	}
}

func (b *BlockRPC) Copy() *BlockRPC {
	bb := new(BlockRPC)
	*bb = *b

	bb.Miner = make([]byte, len(b.Miner))
	copy(bb.Miner[:], b.Miner[:])

	bb.ExtraData = make([]byte, len(b.ExtraData))
	copy(bb.ExtraData[:], b.ExtraData[:])

	return bb
}

// TransactionRPC represents a transaction returned by the edge RPC.
type TransactionRPC struct {
	Nonce     argUint64      `json:"nonce"`
	GasPrice  *argBig        `json:"gasPrice,omitempty"`
	GasTipCap *argBig        `json:"maxPriorityFeePerGas,omitempty"`
	GasFeeCap *argBig        `json:"maxFeePerGas,omitempty"`
	Gas       argUint64      `json:"gas"`
	To        *types.Address `json:"to"`
	Value     argBig         `json:"value"`
	Input     argBytes       `json:"input"`
	V         argBig         `json:"v"`
	R         argBig         `json:"r"`
	S         argBig         `json:"s"`
	Hash      types.Hash     `json:"hash"`
	From      types.Address  `json:"from"`

	// Additional fields.
	BlockHash   *types.Hash `json:"blockHash"`
	BlockNumber *argUint64  `json:"blockNumber"`
	TxIndex     *argUint64  `json:"transactionIndex"`
	ChainID     *argBig     `json:"chainId,omitempty"`
	Type        argUint64   `json:"type"`
}

func (tx *TransactionRPC) toTransactionGrpc() *types.Transaction {
	return &types.Transaction{
		Nonce:     uint64(tx.Nonce),
		GasPrice:  (*big.Int)(tx.GasPrice),
		GasTipCap: (*big.Int)(tx.GasTipCap),
		GasFeeCap: (*big.Int)(tx.GasFeeCap),
		Gas:       uint64(tx.Gas),
		To:        (*types.Address)(tx.To),
		Value:     (*big.Int)(&tx.Value),
		Input:     tx.Input,
		V:         (*big.Int)(&tx.V),
		R:         (*big.Int)(&tx.R),
		S:         (*big.Int)(&tx.S),
		Hash:      types.Hash(tx.Hash),
		From:      types.Address(tx.From),
		Type:      types.TxType(tx.Type),
		ChainID:   (*big.Int)(tx.ChainID),
	}
}

type argUint64 uint64

func (u argUint64) MarshalText() ([]byte, error) {
	buf := make([]byte, 2, 10)
	copy(buf, `0x`)
	buf = strconv.AppendUint(buf, uint64(u), 16)

	return buf, nil
}

func (u *argUint64) UnmarshalText(input []byte) error {
	str := strings.TrimPrefix(string(input), "0x")
	num, err := strconv.ParseUint(str, 16, 64)

	if err != nil {
		return err
	}

	*u = argUint64(num)

	return nil
}

type argBytes []byte

func (b argBytes) MarshalText() ([]byte, error) {
	return encodeToHex(b), nil
}

func (b *argBytes) UnmarshalText(input []byte) error {
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

type argBig big.Int

func (a *argBig) UnmarshalText(input []byte) error {
	buf, err := decodeToHex(input)
	if err != nil {
		return err
	}

	b := new(big.Int)
	b.SetBytes(buf)
	*a = argBig(*b)

	return nil
}

func (a argBig) MarshalText() ([]byte, error) {
	b := (*big.Int)(&a)

	return []byte("0x" + b.Text(16)), nil
}
