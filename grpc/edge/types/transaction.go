package types

import (
	"fmt"
	"math/big"

	"github.com/umbracle/fastrlp"
)

const (
	LegacyTx     TxType = 0x0
	StateTx      TxType = 0x7f
	DynamicFeeTx TxType = 0x02
)

type (
	// TransactionRPC represents a transaction returned by the edge RPC.
	TransactionRPC struct {
		Nonce       uint64   `json:"nonce"`
		GasPrice    *ArgBig  `json:"gasPrice,omitempty"`
		GasTipCap   *ArgBig  `json:"maxPriorityFeePerGas,omitempty"`
		GasFeeCap   *ArgBig  `json:"maxFeePerGas,omitempty"`
		Gas         uint64   `json:"gas"`
		To          *Address `json:"to"`
		Value       ArgBig   `json:"value"`
		Input       []byte   `json:"input"`
		V           ArgBig   `json:"v"`
		R           ArgBig   `json:"r"`
		S           ArgBig   `json:"s"`
		Hash        Hash     `json:"hash"`
		From        Address  `json:"from"`
		BlockHash   *Hash    `json:"blockHash"`
		BlockNumber *uint64  `json:"blockNumber"`
		TxIndex     *uint64  `json:"transactionIndex"`
		ChainID     *ArgBig  `json:"chainId,omitempty"`
		Type        uint64   `json:"type"`
	}

	// TransactionGrpc represents a transaction returned by edge gRPC server.
	TransactionGrpc struct {
		Hash      Hash
		Nonce     uint64
		From      Address
		To        *Address
		Value     *big.Int
		GasPrice  *big.Int
		Gas       uint64
		GasTipCap *big.Int
		GasFeeCap *big.Int
		Input     []byte
		V, R, S   *big.Int
		Type      TxType
		ChainID   *big.Int
	}

	ArgBig big.Int

	TxType byte
)

func (tx *TransactionRPC) toTransactionGrpc() TransactionGrpc {
	return TransactionGrpc{
		Hash:      Hash(tx.Hash),
		Nonce:     tx.Nonce,
		From:      Address(tx.From),
		To:        (*Address)(tx.To),
		Value:     (*big.Int)(&tx.Value),
		GasPrice:  (*big.Int)(tx.GasPrice),
		Gas:       tx.Gas,
		GasTipCap: (*big.Int)(tx.GasTipCap),
		GasFeeCap: (*big.Int)(tx.GasFeeCap),
		Input:     tx.Input,
		V:         (*big.Int)(&tx.V),
		R:         (*big.Int)(&tx.R),
		S:         (*big.Int)(&tx.S),
		Type:      TxType(tx.Type),
		ChainID:   (*big.Int)(tx.ChainID),
	}
}

func (a *ArgBig) UnmarshalText(input []byte) error {
	buf, err := decodeToHex(input)
	if err != nil {
		return err
	}

	b := new(big.Int)
	b.SetBytes(buf)
	*a = ArgBig(*b)

	return nil
}

func (a ArgBig) MarshalText() ([]byte, error) {
	b := (*big.Int)(&a)

	return []byte("0x" + b.Text(16)), nil
}

func (t *TransactionGrpc) MarshalRLP() []byte {
	return t.MarshalRLPTo(nil)
}

func (t *TransactionGrpc) MarshalRLPTo(dst []byte) []byte {
	return MarshalRLPTo(t.MarshalRLPWith, dst)
}

func (t *TransactionGrpc) MarshalRLPWith(arena *fastrlp.Arena) *fastrlp.Value {
	vv := arena.NewArray()

	// Check Transaction1559Payload there https://eips.ethereum.org/EIPS/eip-1559#specification
	if t.Type == DynamicFeeTx {
		vv.Set(arena.NewBigInt(t.ChainID))
	}

	vv.Set(arena.NewUint(t.Nonce))

	if t.Type == DynamicFeeTx {
		// Add EIP-1559 related fields.
		// For non-dynamic-fee-tx gas price is used.
		vv.Set(arena.NewBigInt(t.GasTipCap))
		vv.Set(arena.NewBigInt(t.GasFeeCap))
	} else {
		vv.Set(arena.NewBigInt(t.GasPrice))
	}

	vv.Set(arena.NewUint(t.Gas))

	// Address may be empty.
	if t.To != nil {
		vv.Set(arena.NewCopyBytes(t.To.Bytes()))
	} else {
		vv.Set(arena.NewNull())
	}

	vv.Set(arena.NewBigInt(t.Value))
	vv.Set(arena.NewCopyBytes(t.Input))

	// Specify access list as per spec.
	// This is needed to have the same format as other EVM chains do.
	// There is no access list feature here, so it is always empty just to be compatible.
	// Check Transaction1559Payload there https://eips.ethereum.org/EIPS/eip-1559#specification
	if t.Type == DynamicFeeTx {
		vv.Set(arena.NewArray())
	}

	// Signature values.
	vv.Set(arena.NewBigInt(t.V))
	vv.Set(arena.NewBigInt(t.R))
	vv.Set(arena.NewBigInt(t.S))

	if t.Type == StateTx {
		vv.Set(arena.NewCopyBytes(t.From.Bytes()))
	}

	return vv
}

func (t *TransactionGrpc) UnmarshalRLP(input []byte) error {
	return UnmarshalRlp(t.UnmarshalRLPFrom, input)
}

// UnmarshalRLPFrom unmarshals a Transaction in RLP format.
func (t *TransactionGrpc) UnmarshalRLPFrom(p *fastrlp.Parser, v *fastrlp.Value) error {
	elems, err := v.GetElems()
	if err != nil {
		return err
	}

	if len(elems) < 9 {
		return fmt.Errorf("incorrect number of elements to decode transaction, expected 9 but found %d", len(elems))
	}

	p.Hash(t.Hash[:0], v)

	// Nonce.
	if t.Nonce, err = elems[0].GetUint64(); err != nil {
		return err
	}
	// GasPrice.
	t.GasPrice = new(big.Int)
	if err = elems[1].GetBigInt(t.GasPrice); err != nil {
		return err
	}
	// Gas.
	if t.Gas, err = elems[2].GetUint64(); err != nil {
		return err
	}
	// To.
	if vv, _ := v.Get(3).Bytes(); len(vv) == 20 {
		// Address.
		addr := BytesToAddress(vv)
		t.To = &addr
	} else {
		// Reset To.
		t.To = nil
	}
	// Value.
	t.Value = new(big.Int)
	if err = elems[4].GetBigInt(t.Value); err != nil {
		return err
	}
	// Input.
	if t.Input, err = elems[5].GetBytes(t.Input[:0]); err != nil {
		return err
	}

	// V.
	t.V = new(big.Int)
	if err = elems[6].GetBigInt(t.V); err != nil {
		return err
	}

	// R.
	t.R = new(big.Int)
	if err = elems[7].GetBigInt(t.R); err != nil {
		return err
	}
	// S.
	t.S = new(big.Int)
	if err = elems[8].GetBigInt(t.S); err != nil {
		return err
	}

	return nil
}

func BytesToAddress(b []byte) Address {
	var a Address

	size := len(b)
	min := min(size, AddressLength)

	copy(a[AddressLength-min:], b[len(b)-min:])

	return a
}

func min(i, j int) int {
	if i < j {
		return i
	}

	return j
}
