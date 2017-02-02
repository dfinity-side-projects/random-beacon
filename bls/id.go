package bls

import (
	"dfinity/beacon/blscgo"
	"github.com/ethereum/go-ethereum/common"
	"log"
	"math/big"
)

// ID -- id for secret sharing, represented by big.Int
type ID struct {
	value big.Int
}

// Setters

// SetBig --
func (id ID) SetBig(b *big.Int) {
	id.value = *b
}

// Getters

// CgoID --
func (id ID) CgoID() (cgoid blscgo.ID) {
	err := cgoid.SetStr(id.value.String())
	if err != nil {
		log.Fatalln("Error in ID conversion to blscgo.")
	}
	return
}

// Constructors

// IDFromBig --
func IDFromBig(b *big.Int) (id ID) {
	id.value = *b
	return
}

// IDFromInt64 --
func IDFromInt64(i int64) (id ID) {
	id.value = *big.NewInt(i)
	return
}

// IDFromAddress --
func IDFromAddress(addr common.Address) ID {
	return IDFromBig(addr.Big())
}
