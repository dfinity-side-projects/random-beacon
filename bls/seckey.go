package bls

import (
	"dfinity/beacon/blscgo"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"log"
	"math/big"
)

// Logging counters
var secAggCalls, secAggLen, secShareCalls, secShareLen, secRecoverCalls, secRecoverLen int

// SeckeyCtrs -- one-line summary of counter values related to secret key operations
func SeckeyCtrs() string {
	return fmt.Sprintf("(sec:agg,shr)     %d/%d,%d/%d", secAggCalls, secAggLen, secShareCalls, secShareLen)
}

// Constants

// R --
var R = Decimal2Big("16798108731015832284940804142231733909759579603404752749028378864165570215949")

// types

// Seckey -- represented by a big.Int modulo R
type Seckey struct {
	secret *big.Int
}

// SeckeyMap -- a map from addresses to Seckey
type SeckeyMap map[common.Address]Seckey

// Getters

// Bytes --
func (sec Seckey) Bytes() []byte {
	// big endian
	return sec.secret.Bytes()
}

// String --
func (sec Seckey) String() string {
	// big endian
	return sec.secret.String()
}

// BigInt --
func (sec Seckey) BigInt() *big.Int {
	return sec.secret
}

// Hex --
func (sec Seckey) Hex() string {
	return fmt.Sprintf("0x%x", sec.secret)
}

// SecretKey -- convert the Seckey to blscgo.SecretKey
func (sec Seckey) SecretKey() (sk *blscgo.SecretKey) {
	sk = new(blscgo.SecretKey)
	err := sk.SetStr(sec.String())
	if err != nil {
		log.Fatalln("Error in SecretKey conversion to blscgo.")
	}
	return
}

// Constructors

// SeckeyFromBytes --
func SeckeyFromBytes(b []byte) (sec Seckey) {
	// the secret has to be cut off at 31 bytes to make it smaller than the constant R
	// R has 254 bits
	// TODO mask only the two highest bits with zeros
	if len(b) > 31 {
		b = b[:31]
	}
	i := Bytes2Big(b)
	sec.secret = &i
	return
}

// SeckeyFromRand --
func SeckeyFromRand(seed Rand) Seckey {
	return SeckeyFromBytes(seed.Bytes())
}

// SeckeyFromBigInt --
func SeckeyFromBigInt(b *big.Int) (sec Seckey) {
	sec.secret = b
	return
}

// SeckeyFromInt --
func SeckeyFromInt(i int64) (sec Seckey) {
	sec.secret = big.NewInt(i)
	return
}

// AggregateSeckeys -- Aggregate multiple seckeys into one by summing up
func AggregateSeckeys(secs []Seckey) (sec Seckey) {
	secAggCalls++
	secAggLen += len(secs)
	sec.secret = big.NewInt(0)
	for _, s := range secs {
		sec.secret.Add(sec.secret, s.secret)
	}
	sec.secret.Mod(sec.secret, &R)
	return
}

// ShareSeckey -- Derive shares from master through polynomial substitution
func ShareSeckey(msec []Seckey, id ID) (sec Seckey) {
	secShareCalls++
	secShareLen += len(msec)
	sec.secret = big.NewInt(0)
	// degree of polynomial, need k >= 1, i.e. len(msec) >= 2
	k := len(msec) - 1
	// msec = c_0, c_1, ..., c_k
	// evaluate polynomial f(x) with coefficients c0, ..., ck
	sec.secret.Set(msec[k].secret)
	for j := k - 1; j >= 0; j-- {
		sec.secret.Mul(sec.secret, &id.value)
		//sec.secret.Mod(&sec.secret, &R)
		sec.secret.Add(sec.secret, msec[j].secret)
		sec.secret.Mod(sec.secret, &R)
	}
	return
}

// ShareSeckeyByAddr -- wrapper around sharing by ID
func ShareSeckeyByAddr(msec []Seckey, addr common.Address) (sec Seckey) {
	return ShareSeckey(msec, IDFromAddress(addr))
}

// RecoverSeckey -- Recover master from shares through Lagrange interpolation
func RecoverSeckey(secs []Seckey, ids []ID) (sec Seckey) {
	secRecoverCalls++
	secRecoverLen += len(secs)
	sec.secret = big.NewInt(0)
	k := len(secs)
	// need len(ids) = k > 0
	for i := 0; i < k; i++ {
		// compute delta_i depending on ids only
		var delta, num, den, diff *big.Int = big.NewInt(1), big.NewInt(1), big.NewInt(1), big.NewInt(0)
		for j := 0; j < k; j++ {
			if j != i {
				num.Mul(num, &ids[j].value)
				num.Mod(num, &R)
				diff.Sub(&ids[j].value, &ids[i].value)
				den.Mul(den, diff)
				den.Mod(den, &R)
			}
		}
		// delta = num / den
		den.ModInverse(den, &R)
		delta.Mul(num, den)
		delta.Mod(delta, &R)
		// apply delta to secs[i]
		delta.Mul(delta, secs[i].secret)
		// skip reducing delta modulo R here
		sec.secret.Add(sec.secret, delta)
		sec.secret.Mod(sec.secret, &R)
	}
	return
}

// RecoverSeckeyByMap --
func RecoverSeckeyByMap(m SeckeyMap, k int) (sec Seckey) {
	ids := make([]ID, k)
	secs := make([]Seckey, k)
	i := 0
	for a, s := range m {
		ids[i] = IDFromAddress(a)
		secs[i] = s
		i++
		if i >= k {
			break
		}
	}
	return RecoverSeckey(secs, ids)
}
