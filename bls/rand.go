package bls

import (
	"strconv"
	"math/big"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

/// Rand

// RandLength --
const RandLength = 32 

// Rand --
type Rand [RandLength]byte

// Constructors

// RandFromBytes --
func RandFromBytes(b []byte) (r Rand) {
	h := crypto.Keccak256Hash(b)
        copy(r[:RandLength], h[:])
	return
}

// Getters

// Bytes --
func (r Rand) Bytes() []byte {
	return r[:]
}

// String --
func (r Rand) String() string {
	return string(r[:])
}

// DerivedRand -- Derived Randomness hierarchically
func (r Rand) DerivedRand(idx []byte) Rand {
	// Keccak is not susceptible to length-extension-attacks, so we can use it as-is to implement an HMAC
	return RandFromBytes(crypto.Keccak256(r.Bytes(), idx))
}

// Shortcuts to the derivation function

// Ders --
// ... by string
func (r Rand) Ders(s ...string) Rand {
	ri := r
	for _, si := range s {
		ri = ri.DerivedRand([]byte(si))
	}
	return ri
}

// Deri --
// ... by int
func (r Rand) Deri(i int) Rand {
	return r.Ders(strconv.Itoa(i))
}

// Modulo --
// Convert to a random integer from the interval [0,n-1]. 
func (r Rand) Modulo(n int) int {
	// modulo len(groups) with big.Ints (Mod method works on pointers)
	var b big.Int
	b.Mod(common.Bytes2Big(r.Bytes()), big.NewInt(int64(n)))
        return int(b.Int64())
}

// RandomPerm --
// Convert to a random permutation
func (r Rand) RandomPerm(n int, k int) []int {
	// modulo len(groups) with big.Ints (Mod method works on pointers)
	l := make([]int, n)
	for i := range l {
		l[i] = i
	}
	for i := 0; i < k; i++ {
		j := r.Deri(i).Modulo(n-i) + i
		l[i], l[j] = l[j], l[i]
	}
	return l[:k]	
}
