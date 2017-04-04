package blscgo

import "testing"
import "strconv"

var curve = CurveFp382_1
var unitN = 0

func TestPre(t *testing.T) {
	t.Log("init")
	Init(curve)
	unitN = GetOpUnitSize()
	var err error
	{
		var id ID
		id.Set([]uint64{6, 5, 4, 3, 2, 1}[0:unitN])

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
		sec.SetArray([]uint64{1, 2, 3, 4, 5, 6}[0:unitN])
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
	Init(curve)
	k := 3000
	var sec SecretKey
	sec.Init()

	// make master secret key
	msk := sec.GetMasterSecretKey(k)

	n := k
	secVec := make([]SecretKey, n)
	idVec := make([]ID, n)
	for i := 0; i < n; i++ {
		idVec[i].Set([]uint64{1, 2, 3, 4, 5, uint64(i), 5, 6}[0:unitN])
		secVec[i].Set(msk, &idVec[i])
	}
	// recover sec2 from secVec and idVec
	var sec2 SecretKey
	sec2.Recover(secVec, idVec)
	if sec.String() != sec2.String() {
		t.Errorf("Mismatch in recovered secret key:\n  %s\n  %s.", sec.String(), sec2.String())
	}
}

func TestSign(t *testing.T) {
	m := "testSign"
	t.Log(m)
	Init(curve)

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

	for i := 0; i < n; i++ {
		idVec[i].Set([]uint64{idTbl[i], 0, 0, 0, 0, 0}[0:unitN])
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
	Init(curve)
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
	Init(curve)
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

func BenchmarkPubkeyFromSeckey(b *testing.B) {
	b.StopTimer()
	Init(curve)
	var sec SecretKey
	for n := 0; n < b.N; n++ {
		sec.Init()
		b.StartTimer()
		sec.GetPublicKey()
		b.StopTimer()
	}
}

func BenchmarkSigning(b *testing.B) {
	b.StopTimer()
	Init(curve)
	var sec SecretKey
	for n := 0; n < b.N; n++ {
		sec.Init()
		b.StartTimer()
		sec.Sign(strconv.Itoa(n))
		b.StopTimer()
	}
}

func BenchmarkValidation(b *testing.B) {
	b.StopTimer()
	Init(curve)
	var sec SecretKey
	for n := 0; n < b.N; n++ {
		sec.Init()
		pub := sec.GetPublicKey()
		m := strconv.Itoa(n)
		sig := sec.Sign(m)
		b.StartTimer()
		sig.Verify(pub, m)
		b.StopTimer()
	}
}

func benchmarkDeriveSeckeyShare(k int, b *testing.B) {
	b.StopTimer()
	Init(curve)
	var sec SecretKey
	sec.Init()
	msk := sec.GetMasterSecretKey(k)
	var id ID
	for n := 0; n < b.N; n++ {
		id.Set([]uint64{1, 2, 3, 4, 5, uint64(n)})
		b.StartTimer()
		sec.Set(msk, &id)
		b.StopTimer()
	}
}

//func BenchmarkDeriveSeckeyShare100(b *testing.B)  { benchmarkDeriveSeckeyShare(100, b) }
//func BenchmarkDeriveSeckeyShare200(b *testing.B)  { benchmarkDeriveSeckeyShare(200, b) }
func BenchmarkDeriveSeckeyShare500(b *testing.B) { benchmarkDeriveSeckeyShare(500, b) }

//func BenchmarkDeriveSeckeyShare1000(b *testing.B) { benchmarkDeriveSeckeyShare(1000, b) }

func benchmarkRecoverSeckey(k int, b *testing.B) {
	b.StopTimer()
	Init(curve)
	var sec SecretKey
	sec.Init()
	msk := sec.GetMasterSecretKey(k)

	// derive n shares
	n := k
	secVec := make([]SecretKey, n)
	idVec := make([]ID, n)
	for i := 0; i < n; i++ {
		idVec[i].Set([]uint64{1, 2, 3, 4, 5, uint64(i)})
		secVec[i].Set(msk, &idVec[i])
	}

	// recover from secVec and idVec
	var sec2 SecretKey
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		sec2.Recover(secVec, idVec)
	}
}

func BenchmarkRecoverSeckey100(b *testing.B)  { benchmarkRecoverSeckey(100, b) }
func BenchmarkRecoverSeckey200(b *testing.B)  { benchmarkRecoverSeckey(200, b) }
func BenchmarkRecoverSeckey500(b *testing.B)  { benchmarkRecoverSeckey(500, b) }
func BenchmarkRecoverSeckey1000(b *testing.B) { benchmarkRecoverSeckey(1000, b) }

func benchmarkRecoverSignature(k int, b *testing.B) {
	b.StopTimer()
	Init(curve)
	var sec SecretKey
	sec.Init()
	msk := sec.GetMasterSecretKey(k)

	// derive n shares
	n := k
	idVec := make([]ID, n)
	secVec := make([]SecretKey, n)
	signVec := make([]Sign, n)
	for i := 0; i < n; i++ {
		idVec[i].Set([]uint64{1, 2, 3, 4, 5, uint64(i)})
		secVec[i].Set(msk, &idVec[i])
		signVec[i] = *secVec[i].Sign("test message")
	}

	// recover signature
	var sig Sign
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		sig.Recover(signVec, idVec)
	}
}

func BenchmarkRecoverSignature100(b *testing.B)  { benchmarkRecoverSignature(100, b) }
func BenchmarkRecoverSignature200(b *testing.B)  { benchmarkRecoverSignature(200, b) }
func BenchmarkRecoverSignature500(b *testing.B)  { benchmarkRecoverSignature(500, b) }
func BenchmarkRecoverSignature1000(b *testing.B) { benchmarkRecoverSignature(1000, b) }
