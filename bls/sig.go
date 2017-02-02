package bls

import (
	"dfinity/beacon/blscgo"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"log"
)

// Debugging counters
var sigGenCalls, sigVerifyCalls, sigAggCalls, sigAggLen, sigRecoverCalls, sigRecoverLen int

// SignatureCtrs --
func SignatureCtrs() string {
	return fmt.Sprintf("(sig:gen,ver,rec) %d,%d,%d/%d", sigGenCalls, sigVerifyCalls, sigRecoverCalls, sigRecoverLen)
}

// types

// Signature --
type Signature struct {
	value []byte
}

// SignatureMap --
type SignatureMap map[common.Address]Signature

// Conversion

// Rand --
func (sig Signature) Rand() Rand {
	return RandFromBytes(sig.value)
}

// String --
func (sig Signature) String() string {
	return string(sig.value)
}

// Signing

// Sig -- convert Signature to blscgo Sign
func (sig Signature) Sig() (sign *blscgo.Sign) {
	sign = new(blscgo.Sign)
	err := sign.SetStr(sig.String())
	if err != nil {
		log.Fatalln("Error in Signature conversion to blscgo.")
	}
	return
}

// Sign -- sign a message with secret key
func Sign(sec Seckey, msg []byte) (sig Signature) {
	sigGenCalls++
	// convert Seckey to blscgo.SecretKey
	sk := sec.SecretKey()
	// sign
	sign := sk.Sign(string(msg))
	// convert back from blscgo
	sig.value = []byte(sign.String())
	return
}

// Verifying

// VerifySig -- verify message and signature against public key
func VerifySig(pub Pubkey, msg []byte, sig Signature) bool {
	sigVerifyCalls++
	// convert to blscgo and verify
	return sig.Sig().Verify(pub.PublicKey(), string(msg))
}

// VerifyAggregateSig --
func VerifyAggregateSig(pubs []Pubkey, msg []byte, asig Signature) bool {
	return VerifySig(AggregatePubkeys(pubs), msg, asig)
}

// BatchVerify --
func BatchVerify(pubs []Pubkey, msg []byte, sigs []Signature) bool {
	return VerifyAggregateSig(pubs, msg, AggregateSigs(sigs))
}

// Aggregation and Recovery

// AggregateSigs -- aggregate multiple into one by summing up
func AggregateSigs(sigs []Signature) (sig Signature) {
	sigAggCalls++
	sigAggLen += len(sigs)
	// convert to blscgo
	sum := sigs[0].Sig()
	// sum it up
	for _, s := range sigs[1:] {
		fmt.Println("agg sigs")
		sum.Add(s.Sig())
	}
	// convert back from blscgo
	sig.value = []byte(sum.String())
	return
}

// RecoverSignature -- Recover master from shares through Lagrange interpolation
func RecoverSignature(sigs []Signature, ids []ID) (sig Signature) {
	sigRecoverCalls++
	sigRecoverLen += len(sigs)

	// convert sigs to blscgo
	signVec := make([]blscgo.Sign, len(sigs))
	for i, s := range sigs {
		signVec[i] = *(s.Sig())
	}
	// convert ids to blscgo
	idVec := make([]blscgo.ID, len(ids))
	for i, id := range ids {
		idVec[i] = id.CgoID()
	}

	var sign blscgo.Sign
	sign.Recover(signVec, idVec)

	sig.value = []byte(sign.String())
	return
}

// RecoverSignatureByMap --
func RecoverSignatureByMap(m SignatureMap, k int) (sec Signature) {
	ids := make([]ID, k)
	sigs := make([]Signature, k)
	i := 0
	for a, s := range m {
		ids[i] = IDFromAddress(a)
		sigs[i] = s
		i++
		if i >= k {
			break
		}
	}
	return RecoverSignature(sigs, ids)
}
