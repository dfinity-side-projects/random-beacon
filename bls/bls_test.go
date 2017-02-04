package bls

import "testing"
import "dfinity/beacon/blscgo"

func TestComparison(t *testing.T) {
	t.Log("testComparison")
	blscgo.Init()
	b := Decimal2Big("16798108731015832284940804142231733909759579603404752749028378864165570215948")
	sec := SeckeyFromBigInt(&b)
	t.Log("sec.Hex: ", sec.Hex())
	t.Log("sec.String: ", sec.String())

	// Add Seckeys
	sum := AggregateSeckeys([]Seckey{sec, sec})
	t.Log("sum: ", sum.Hex())

	sk := sec.SecretKey()
	t.Log("sk = sec.SecretKey(): ", sk.String())

	// Pubkey
	pk := sk.GetPublicKey()
	t.Log("pk: ", pk.String())
	pub := PubkeyFromSeckey(sec)
	t.Log("pub: ", pub.String())
	//pub2 := PublicKeyFromSeckey(sec)
	//t.Log("pub2: ", pub2.String())

	// Add SecretKeys
	sk.Add(sk)
	t.Log("sksum: ", sk.String())

	if sk.String() != sum.Hex() {
		t.Error("Mismatch in secret key addition")
	}

	// Sig
	sig := Sign(sec, []byte("hi"))
	asig := AggregateSigs([]Signature{sig, sig})
	if !VerifyAggregateSig([]Pubkey{pub, pub}, []byte("hi"), asig) {
		t.Error("Aggregated signature does not verify")
	}
}
