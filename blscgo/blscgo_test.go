package blscgo

import "testing"

func TestPre(t *testing.T) {
	t.Log("init")
	Init()
	var err error
	{
		var id ID
		err = id.Set([]uint64{4, 3, 2, 1})
		if err != nil {
			t.Fatal(err)
		}

		t.Log("id :", id)
		var id2 ID
		err = id2.SetStr(id.String())
		if err != nil {
			t.Fatal(err)
		}
		t.Log("id2:", id2)
	}
	{
		var sec SecretKey
		err = sec.SetArray([]uint64{1, 2, 3, 4})
		if err != nil {
			t.Fatal(err)
		}
		t.Log("sec=", sec)
	}

	t.Log("create secret key")
	m := "this is a bls sample for go"
	var sec SecretKey
	sec.Init()
	t.Log("sec:", sec)
	t.Log("create public key")
	pub := sec.GetPublicKey()
	t.Log("pub:", pub)
	sign := sec.Sign(m)
	t.Log("sign:", sign)
	if !sign.Verify(pub, m) {
		t.Error("Signature does not verify")
	}

	// How to make array of SecretKey
	{
		sec := make([]SecretKey, 3)
		for i := 0; i < len(sec); i++ {
			sec[i].Init()
			t.Log("sec=", sec[i].String())
		}
	}
}

func TestRecoverSecretKey(t *testing.T) {
	t.Log("testRecoverSecretKey")
	Init()
	k := 3000
	var sec SecretKey
	sec.Init()

	// make master secret key
	msk := sec.GetMasterSecretKey(k)

	n := k
	secVec := make([]SecretKey, n)
	idVec := make([]ID, n)
	var err error
	for i := 0; i < n; i++ {
		err = idVec[i].Set([]uint64{1, 2, 3, uint64(i)})
		if err != nil {
			t.Fatal(err)
		}
		secVec[i].Set(msk, &idVec[i])
	}
	// recover sec2 from secVec and idVec
	var sec2 SecretKey
	sec2.Recover(secVec, idVec)
	if sec != sec2 {
		t.Errorf("Mismatch in recovered secret key:\n  %s\n  %s.", sec.String(), sec2.String())
	}
}

func TestSign(t *testing.T) {
	m := "testSign"
	t.Log(m)
	Init()

	var sec0 SecretKey
	sec0.Init()
	pub0 := sec0.GetPublicKey()
	s0 := sec0.Sign(m)
	if !s0.Verify(pub0, m) {
		t.Error("Signature does not verify")
	}

	k := 3
	msk := sec0.GetMasterSecretKey(k)
	mpk := GetMasterPublicKey(msk)

	idTbl := []uint64{3, 5, 193, 22, 15}
	n := len(idTbl)

	secVec := make([]SecretKey, n)
	pubVec := make([]PublicKey, n)
	signVec := make([]Sign, n)
	idVec := make([]ID, n)

	var err error
	for i := 0; i < n; i++ {
		err = idVec[i].Set([]uint64{idTbl[i], 0, 0, 0})
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("idVec[%d]=%s\n", i, idVec[i].String())

		secVec[i].Set(msk, &idVec[i])

		pubVec[i].Set(mpk, &idVec[i])
		t.Logf("pubVec[%d]=%s\n", i, pubVec[i].String())

		if pubVec[i].String() != secVec[i].GetPublicKey().String() {
			t.Error("Pubkey derivation does not match")
		}

		signVec[i] = *secVec[i].Sign(m)
		if !signVec[i].Verify(&pubVec[i], m) {
			t.Error("Pubkey derivation does not match")
		}
	}
	var sec1 SecretKey
	sec1.Recover(secVec, idVec)
	if sec0.String() != sec1.String() {
		t.Error("Mismatch in recovered seckey.")
	}
	var pub1 PublicKey
	pub1.Recover(pubVec, idVec)
	if pub0.String() != pub1.String() {
		t.Error("Mismatch in recovered pubkey.")
	}
	var s1 Sign
	s1.Recover(signVec, idVec)
	if s0.String() != s1.String() {
		t.Error("Mismatch in recovered signature.")
	}
}

func TestAdd(t *testing.T) {
	t.Log("testAdd")
	Init()
	var sec1 SecretKey
	var sec2 SecretKey
	sec1.Init()
	sec2.Init()

	pub1 := sec1.GetPublicKey()
	pub2 := sec2.GetPublicKey()

	m := "test test"
	sign1 := sec1.Sign(m)
	sign2 := sec2.Sign(m)

	t.Log("sign1    :", sign1)
	sign1.Add(sign2)
	t.Log("sign1 add:", sign1)
	pub1.Add(pub2)
	if !sign1.Verify(pub1, m) {
		t.Fail()
	}
}

func TestPop(t *testing.T) {
	t.Log("testPop")
	Init()
	var sec SecretKey
	sec.Init()
	pop := sec.GetPop()
	if !pop.VerifyPop(sec.GetPublicKey()) {
		t.Errorf("Valid Pop does not verify")
	}
	sec.Init()
	if pop.VerifyPop(sec.GetPublicKey()) {
		t.Errorf("Invalid Pop verifies")
	}
}
