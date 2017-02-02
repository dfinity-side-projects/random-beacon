package bls

// types

// Pop --
type Pop Signature

// Proof-of-Possesion

// GeneratePop --
func GeneratePop(sec Seckey, pub Pubkey) Pop {
	tmp := &pub
	return Pop(Sign(sec, []byte(tmp.String())))
}

// Verification

// VerifyPop --
func VerifyPop(pub Pubkey, pop Pop) bool {
	return VerifySig(pub, []byte(pub.String()), Signature(pop))
}
