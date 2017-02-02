package bls

import "math/big"

// Generators

// Decimal2Big --
func Decimal2Big(s string) (x big.Int) {
	x.SetString(s, 10)
	return
}

// Hex2Big --
func Hex2Big(s string) (x big.Int) {
	x.SetString(s, 16)
	return
}

// Bytes2Big --
func Bytes2Big(b []byte) (x big.Int) {
	// big endian
	x.SetBytes(b)
	return
}
