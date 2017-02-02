package bls

import "fmt"

// PrintCtrs -- print counters counting various operations
func PrintCtrs() {
	fmt.Printf("--- Crypto-Ops\n  %s\n  %s\n  %s\n", SeckeyCtrs(), PubkeyCtrs(), SignatureCtrs())
}
