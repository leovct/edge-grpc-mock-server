package types

import (
	"fmt"

	"github.com/umbracle/fastrlp"
)

type (
	// BlockRPC represents a block returned by the edge RPC.
	BlockRPC struct {
		ParentHash      Hash             `json:"parentHash"`
		Sha3Uncles      Hash             `json:"sha3Uncles"`
		Miner           []byte           `json:"miner"`
		StateRoot       Hash             `json:"stateRoot"`
		TxRoot          Hash             `json:"transactionsRoot"`
		ReceiptsRoot    Hash             `json:"receiptsRoot"`
		LogsBloom       Bloom            `json:"logsBloom"`
		Difficulty      uint64           `json:"difficulty"`
		TotalDifficulty uint64           `json:"totalDifficulty"`
		Size            uint64           `json:"size"`
		Number          uint64           `json:"number"`
		GasLimit        uint64           `json:"gasLimit"`
		GasUsed         uint64           `json:"gasUsed"`
		Timestamp       uint64           `json:"timestamp"`
		ExtraData       []byte           `json:"extraData"`
		MixHash         Hash             `json:"mixHash"`
		Nonce           Nonce            `json:"nonce"`
		Hash            Hash             `json:"hash"`
		Transactions    []TransactionRPC `json:"transactions"`
		Uncles          []Hash           `json:"uncles"`
		BaseFee         uint64           `json:"baseFeePerGas,omitempty"`
	}

	// BlockGrpc represents a block returned by edge gRPC server.
	BlockGrpc struct {
		Header       *Header
		Transactions []*TransactionGrpc
		Uncles       []*Header
	}
)

func (b *BlockRPC) ToBlockGrpc() *BlockGrpc {
	header := Header{
		ParentHash:      b.ParentHash,
		Sha3Uncles:      b.Sha3Uncles,
		Miner:           b.Miner,
		StateRoot:       b.StateRoot,
		TxRoot:          b.TxRoot,
		ReceiptsRoot:    b.ReceiptsRoot,
		LogsBloom:       b.LogsBloom,
		Difficulty:      b.Difficulty,
		TotalDifficulty: b.TotalDifficulty,
		Size:            b.Size,
		Number:          b.Number,
		GasLimit:        b.GasLimit,
		GasUsed:         b.GasUsed,
		Timestamp:       b.Timestamp,
		ExtraData:       b.ExtraData,
		MixHash:         b.MixHash,
		Nonce:           b.Nonce,
		Hash:            b.Hash,
		BaseFee:         b.BaseFee,
	}

	transactions := make([]*TransactionGrpc, len(b.Transactions))
	for _, txGrpc := range b.Transactions {
		txRPC := txGrpc.toTransactionGrpc()
		transactions = append(transactions, &txRPC)
	}

	// Note: we don't parse uncles for the moment.
	var uncles []*Header

	return &BlockGrpc{
		Header:       &header,
		Transactions: transactions,
		Uncles:       uncles,
	}
}

func (b *BlockGrpc) MarshalRLP() []byte {
	return b.MarshalRLPTo(nil)
}

func (b *BlockGrpc) MarshalRLPTo(dst []byte) []byte {
	return MarshalRLPTo(b.MarshalRLPWith, dst)
}

func (b *BlockGrpc) MarshalRLPWith(ar *fastrlp.Arena) *fastrlp.Value {
	vv := ar.NewArray()
	vv.Set(b.Header.MarshalRLPWith(ar))

	if len(b.Transactions) == 0 {
		vv.Set(ar.NewNullArray())
	} else {
		v0 := ar.NewArray()
		for _, tx := range b.Transactions {
			if tx.Type != LegacyTx {
				v0.Set(ar.NewCopyBytes([]byte{byte(tx.Type)}))
			}

			v0.Set(tx.MarshalRLPWith(ar))
		}
		vv.Set(v0)
	}

	if len(b.Uncles) == 0 {
		vv.Set(ar.NewNullArray())
	} else {
		v1 := ar.NewArray()
		for _, uncle := range b.Uncles {
			v1.Set(uncle.MarshalRLPWith(ar))
		}
		vv.Set(v1)
	}

	return vv
}

func (b *BlockGrpc) UnmarshalRLP(input []byte) error {
	return UnmarshalRlp(b.UnmarshalRLPFrom, input)
}

func (b *BlockGrpc) UnmarshalRLPFrom(p *fastrlp.Parser, v *fastrlp.Value) error {
	elems, err := v.GetElems()
	if err != nil {
		return err
	}

	if len(elems) < 3 {
		return fmt.Errorf("incorrect number of elements to decode block, expected 3 but found %d", len(elems))
	}

	// Header.
	b.Header = &Header{}
	if err = b.Header.UnmarshalRLPFrom(p, elems[0]); err != nil {
		return err
	}

	// Transactions.
	txns, err := elems[1].GetElems()
	if err != nil {
		return err
	}

	for _, txn := range txns {
		bTxn := &TransactionGrpc{}
		if err = bTxn.UnmarshalRLPFrom(p, txn); err != nil {
			return err
		}

		b.Transactions = append(b.Transactions, bTxn)
	}

	// Uncles.
	uncles, err := elems[2].GetElems()
	if err != nil {
		return err
	}

	for _, uncle := range uncles {
		bUncle := &Header{}
		if err := bUncle.UnmarshalRLPFrom(p, uncle); err != nil {
			return err
		}

		b.Uncles = append(b.Uncles, bUncle)
	}

	return nil
}
