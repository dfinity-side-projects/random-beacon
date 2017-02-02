package bls

import (
	"dfinity/beacon/blscgo"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
)

/// Crypto
// Debugging counters
var pubGenCalls, pubAggCalls, pubAggLen, pubShareCalls, pubShareLen int

// PubkeyCtrs --
func PubkeyCtrs() string {
	return fmt.Sprintf("(pub:gen,shr,agg) %d,%d/%d,%d/%d", pubGenCalls, pubShareCalls, pubShareLen, pubAggCalls, pubAggLen)
}

// types

// Pubkey -
type Pubkey struct {
	value []byte
}

// PubkeyMap --
type PubkeyMap map[common.Address]Pubkey

// Getters

// Hash -- hash & id
func (pub Pubkey) Hash() common.Hash {
	return crypto.Keccak256Hash(pub.value)
}

// Address --
func (pub Pubkey) Address() common.Address {
	h := pub.Hash()

	return common.BytesToAddress(h[:])
	//	return Trace2Address(pub.trace)
	//	pubBytes := []byte("pubkey")
}

// String --
func (pub Pubkey) String() string {
	return string(pub.value)
	//	a := pub.Address()
	//	return fmt.Sprintf("%x", pub.Address())
}

// PublicKey --
func (pub Pubkey) PublicKey() (pk *blscgo.PublicKey) {
	pk = new(blscgo.PublicKey)
	err := pk.SetStr(pub.String())
	if err != nil {
		log.Fatalln("Error in PublicKey conversion to blscgo.")
	}
	return
}

// Generation

// PubkeyFromSeckey -- derive the pubkey from seckey
func PubkeyFromSeckey(sec Seckey) (pub Pubkey) {
	//	pubkey_ctr++
	pubGenCalls++
	// Convert via blscgo
	pk := sec.SecretKey().GetPublicKey()
	pub.value = []byte(pk.String())
	return
}

// AggregatePubkeys -- aggregate multiple into one by summing up
func AggregatePubkeys(pubs []Pubkey) (pub Pubkey) {
	pubAggCalls++
	pubAggLen += len(pubs)
	// initialize sum to zero
	var s blscgo.SecretKey
	err := s.SetStr("0")
	if err != nil {
		log.Fatalln("Error in PublicKey conversion to blscgo.")
	}
	sum := s.GetPublicKey()
	// sum it up
	for _, p := range pubs {
		sum.Add(p.PublicKey())
	}
	// convert back from blscgo
	pub.value = []byte(sum.String())
	return
}

// SharePubkey -- Derive shares from master through polynomial substitution
func SharePubkey(mpub []Pubkey, id ID) (pub Pubkey) {
	pubShareCalls++
	pubShareLen += len(mpub)

	// convert to blscgo master
	mpk := make([]blscgo.PublicKey, len(mpub))
	for i, p := range mpub {
		mpk[i] = *(p.PublicKey())
	}

	// derive gshare
	var pk blscgo.PublicKey
	cgoid := id.CgoID()
	pk.Set(mpk, &cgoid)

	// convert back from blscgo
	pub.value = []byte(pk.String())
	return
}
