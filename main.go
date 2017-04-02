package main

import (
	"dfinity/beacon/bls"
	"dfinity/beacon/blscgo"
	"dfinity/beacon/sim"
	"flag"
	"fmt"
)

func main() {
	var l, n, k, N, m uint
	var seedstr string
	var bist, vvec, timing bool
	var curve string
	flag.UintVar(&l, "l", 20, "Length of chain (number of blocks to create)")
	flag.UintVar(&n, "n", 3, "Group size")
	flag.UintVar(&k, "k", 2, "Threshold")
	flag.UintVar(&N, "N", 8, "Number of processes")
	flag.UintVar(&m, "m", 5, "Number of groups")
	flag.StringVar(&seedstr, "seed", "DFINITY", "Random seed")
	flag.BoolVar(&bist, "bist", false, "Enable Built-in self test")
	flag.BoolVar(&vvec, "vvec", false, "Enable validation against verification vector")
	flag.BoolVar(&timing, "timing", false, "Enable output of timing information")
	flag.StringVar(&curve, "curve", "bn382_1", "Pairing type")
	flag.Parse()

	// init Cgo
	if curve == "bn254" {
		fmt.Println("bn254")
		blscgo.Init(blscgo.CurveFp254BNb)
	} else if curve == "bn382_1" {
		fmt.Println("bn382_1")
		blscgo.Init(blscgo.CurveFp382_1)
	} else if curve == "bn382_2" {
		fmt.Println("bn382_2")
		blscgo.Init(blscgo.CurveFp382_2)
	} else {
		fmt.Printf("not supported curve %s\n", curve)
		return
	}

	seed := bls.RandFromBytes([]byte(seedstr))
	sim.DoubleCheck = bist
	sim.Vvec = vvec
	sim.Timing = timing
	// seed, groupSize, threshold, nProcesses, nGroups
	mysim := sim.NewBlockchainSimulator(seed, uint16(n), uint16(k), N, uint16(m))
	fmt.Println("--- Genesis block ")
	fmt.Printf("%d: %s", mysim.Length(), mysim.Tip().String(true))
	fmt.Printf("--- Blockchain states: (l)%d\n", l)
	for i := uint(0); i < l; i++ {
		mysim.Advance(1, false)
		fmt.Printf("%3d: %s\n", mysim.Length(), mysim.Tip().String(false))
	}

	if timing {
		bls.PrintCtrs()
		fmt.Println("--- Info")
		fmt.Println("Expected Crypto-Ops:")
		fmt.Println("  Seckey calls:    m*n/m*n^2, m*n^2/m*n^2*k")
		fmt.Println("  Pubkey calls:    N+m*n+m*n^2, m*n^2/m*n^2*k, m/m*n   (if --vvec enabled)")
		// pubkey generation: N is process generation, m*n is vvec generation, m*n^2 is rhs of vvec verification
		// pubkey sharing: m*n^2/m*n^2*k is lhs of vvec verification
		// pubkey aggregation: m/m*n is generation of group pubkey from member shares
		fmt.Println("  Pubkey calls:    N, 0/0, m/m*n                       (if --vvec disabled)")
		fmt.Println("  Signature calls: N+l*n, N, l/l*k")
	}
}
