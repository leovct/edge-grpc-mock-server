package types

import (
	"encoding/binary"
	"fmt"

	"github.com/0xPolygon/polygon-edge/helper/keccak"
	"github.com/umbracle/fastrlp"
)

const BloomByteLength = 256

type Header struct {
	ParentHash   Hash
	Sha3Uncles   Hash
	Miner        []byte
	StateRoot    Hash
	TxRoot       Hash
	ReceiptsRoot Hash
	LogsBloom    Bloom
	Difficulty   uint64
	Number       uint64
	GasLimit     uint64
	GasUsed      uint64
	Timestamp    uint64
	ExtraData    []byte
	MixHash      Hash
	Nonce        Nonce
	Hash         Hash
	// BaseFee was added by EIP-1559 and is ignored in legacy headers.
	BaseFee uint64 `json:"baseFeePerGas"`
}

type Bloom [BloomByteLength]byte

type Nonce [8]byte

func (h *Header) SetNonce(i uint64) {
	binary.BigEndian.PutUint64(h.Nonce[:], i)
}

var HeaderHash = defHeaderHash
var marshalArenaPool fastrlp.ArenaPool

func (h *Header) ComputeHash() *Header {
	h.Hash = HeaderHash(h)
	return h
}

func defHeaderHash(h *Header) (hash Hash) {
	// default header hashing
	ar := marshalArenaPool.Get()
	hasher := keccak.DefaultKeccakPool.Get()

	v := h.MarshalRLPWith(ar)
	hasher.WriteRlp(hash[:0], v)

	marshalArenaPool.Put(ar)
	keccak.DefaultKeccakPool.Put(hasher)

	return
}

func (h *Header) MarshalRLPWith(arena *fastrlp.Arena) *fastrlp.Value {
	vv := arena.NewArray()

	vv.Set(arena.NewCopyBytes(h.ParentHash.Bytes()))
	vv.Set(arena.NewCopyBytes(h.Sha3Uncles.Bytes()))
	vv.Set(arena.NewCopyBytes(h.Miner[:]))
	vv.Set(arena.NewCopyBytes(h.StateRoot.Bytes()))
	vv.Set(arena.NewCopyBytes(h.TxRoot.Bytes()))
	vv.Set(arena.NewCopyBytes(h.ReceiptsRoot.Bytes()))
	vv.Set(arena.NewCopyBytes(h.LogsBloom[:]))

	vv.Set(arena.NewUint(h.Difficulty))
	vv.Set(arena.NewUint(h.Number))
	vv.Set(arena.NewUint(h.GasLimit))
	vv.Set(arena.NewUint(h.GasUsed))
	vv.Set(arena.NewUint(h.Timestamp))

	vv.Set(arena.NewCopyBytes(h.ExtraData))
	vv.Set(arena.NewCopyBytes(h.MixHash.Bytes()))
	vv.Set(arena.NewCopyBytes(h.Nonce[:]))

	vv.Set(arena.NewUint(h.BaseFee))

	return vv
}

func (h *Header) UnmarshalRLP(input []byte) error {
	return UnmarshalRlp(h.UnmarshalRLPFrom, input)
}

func (h *Header) UnmarshalRLPFrom(p *fastrlp.Parser, v *fastrlp.Value) error {
	elems, err := v.GetElems()
	if err != nil {
		return err
	}

	if len(elems) < 15 {
		return fmt.Errorf("incorrect number of elements to decode header, expected 15 but found %d", len(elems))
	}

	// parentHash
	if err = elems[0].GetHash(h.ParentHash[:]); err != nil {
		return err
	}
	// sha3uncles
	if err = elems[1].GetHash(h.Sha3Uncles[:]); err != nil {
		return err
	}
	// miner
	if h.Miner, err = elems[2].GetBytes(h.Miner[:]); err != nil {
		return err
	}
	// stateroot
	if err = elems[3].GetHash(h.StateRoot[:]); err != nil {
		return err
	}
	// txroot
	if err = elems[4].GetHash(h.TxRoot[:]); err != nil {
		return err
	}
	// receiptroot
	if err = elems[5].GetHash(h.ReceiptsRoot[:]); err != nil {
		return err
	}
	// logsBloom
	if _, err = elems[6].GetBytes(h.LogsBloom[:0], 256); err != nil {
		return err
	}
	// difficulty
	if h.Difficulty, err = elems[7].GetUint64(); err != nil {
		return err
	}
	// number
	if h.Number, err = elems[8].GetUint64(); err != nil {
		return err
	}
	// gasLimit
	if h.GasLimit, err = elems[9].GetUint64(); err != nil {
		return err
	}
	// gasused
	if h.GasUsed, err = elems[10].GetUint64(); err != nil {
		return err
	}
	// timestamp
	if h.Timestamp, err = elems[11].GetUint64(); err != nil {
		return err
	}
	// extraData
	if h.ExtraData, err = elems[12].GetBytes(h.ExtraData[:0]); err != nil {
		return err
	}
	// mixHash
	if err = elems[13].GetHash(h.MixHash[:0]); err != nil {
		return err
	}
	// nonce
	nonce, err := elems[14].GetUint64()
	if err != nil {
		return err
	}

	h.SetNonce(nonce)

	// compute the hash after the decoding
	h.ComputeHash()

	return err
}
