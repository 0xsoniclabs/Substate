package rlp

import (
	"math/big"

	"github.com/Fantom-foundation/Substate/types"
)

const (
	berlinBlock = 37_455_223
	londonBlock = 37_534_833
)

// IsLondonFork returns true if block is part of the london fork block range
func IsLondonFork(block uint64) bool {
	return block >= londonBlock
}

// IsBerlinFork returns true if block is part of the berlin fork block range
func IsBerlinFork(block uint64) bool {
	return block >= berlinBlock && block < londonBlock
}

// legacySubstateRLP represents legacy RLP structure between before Berlin fork thus before berlinBlock
type legacySubstateRLP struct {
	InputAlloc  WorldState
	OutputAlloc WorldState
	Env         *legacyEnv
	Message     *legacyMessage
	Result      *Result
}

func (r legacySubstateRLP) toLondon() *RLP {
	return &RLP{
		InputSubstate:  r.InputAlloc,
		OutputSubstate: r.OutputAlloc,
		Env:            r.Env.toLondon(),
		Message:        r.Message.toLondon(),
		Result:         r.Result,
	}
}

type legacyMessage struct {
	Nonce      uint64
	CheckNonce bool
	GasPrice   *big.Int
	Gas        uint64

	From  types.Address
	To    *types.Address `rlp:"nil"` // nil means contract creation
	Value *big.Int
	Data  []byte

	InitCodeHash *types.Hash `rlp:"nil"` // NOT nil for contract creation
}

func (m legacyMessage) toLondon() *Message {
	return &Message{
		Nonce:        m.Nonce,
		CheckNonce:   m.CheckNonce,
		GasPrice:     m.GasPrice,
		Gas:          m.Gas,
		From:         m.From,
		To:           m.To,
		Value:        new(big.Int).Set(m.Value),
		Data:         m.Data,
		InitCodeHash: m.InitCodeHash,
		AccessList:   nil, // access list was not present before berlin fork?

		// Same behavior as AccessListTx.gasFeeCap() and AccessListTx.gasTipCap()
		GasFeeCap: m.GasPrice,
		GasTipCap: m.GasPrice,
	}
}

type legacyEnv struct {
	Coinbase    types.Address
	Difficulty  *big.Int
	GasLimit    uint64
	Number      uint64
	Timestamp   uint64
	BlockHashes [][2]types.Hash
}

func (e legacyEnv) toLondon() *Env {
	return &Env{
		Coinbase:    e.Coinbase,
		Difficulty:  e.Difficulty,
		GasLimit:    e.GasLimit,
		Number:      e.Number,
		Timestamp:   e.Timestamp,
		BlockHashes: e.BlockHashes,
	}
}

// berlinRLP represents legacy RLP structure between Berlin and London fork starting at berlinBlock ending at londonBlock
type berlinRLP struct {
	InputAlloc  WorldState
	OutputAlloc WorldState
	Env         *legacyEnv
	Message     *berlinMessage
	Result      *Result
}

func (r berlinRLP) toLondon() *RLP {
	return &RLP{
		InputSubstate:  r.InputAlloc,
		OutputSubstate: r.OutputAlloc,
		Env:            r.Env.toLondon(),
		Message:        r.Message.toLondon(),
		Result:         r.Result,
	}

}

type berlinMessage struct {
	Nonce      uint64
	CheckNonce bool
	GasPrice   *big.Int
	Gas        uint64

	From  types.Address
	To    *types.Address `rlp:"nil"` // nil means contract creation
	Value *big.Int
	Data  []byte

	InitCodeHash *types.Hash `rlp:"nil"` // NOT nil for contract creation

	AccessList types.AccessList // missing in substate DB from Geth v1.9.x
}

func (m berlinMessage) toLondon() *Message {
	return &Message{
		Nonce:        m.Nonce,
		CheckNonce:   m.CheckNonce,
		GasPrice:     m.GasPrice,
		Gas:          m.Gas,
		From:         m.From,
		To:           m.To,
		Value:        new(big.Int).Set(m.Value),
		Data:         m.Data,
		InitCodeHash: m.InitCodeHash,
		AccessList:   m.AccessList,

		// Same behavior as AccessListTx.gasFeeCap() and AccessListTx.gasTipCap()
		GasFeeCap: m.GasPrice,
		GasTipCap: m.GasPrice,
	}
}