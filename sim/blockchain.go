package sim

import (
	"dfinity/beacon/bls"
	"dfinity/beacon/state"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
)

// BlockchainSimulator -- Encodes the state of all processes, groups and the blockchain
type BlockchainSimulator struct {
	groupSize uint16
	threshold uint16
	seed      bls.Rand
	proc      []ProcessSimulator
	group     []GroupSimulator
	grpmap    map[common.Address]*GroupSimulator
	chain     []state.State
}

// DoubleCheck -- enable optional double-checks for verification
var DoubleCheck = true

// Vvec -- enable checks involving the verification vectors
var Vvec = true

// Timing -- enable output of timing information
var Timing = false

// InitProcs -- initialize the individual processes for the genesis block
func (sim *BlockchainSimulator) InitProcs(n uint) {
	sim.proc = make([]ProcessSimulator, n)
	rsec := sim.seed.Ders("InitProcs_sec")
	rseed := sim.seed.Ders("InitProcs_seed")
	for i := 0; i < int(n); i++ {
		sim.proc[i] = NewProcessSimulator(bls.SeckeyFromRand(rsec.Deri(i)), rseed.Deri(i))
		fmt.Println(sim.proc[i].String())
	}
}

// InitGroups -- initialize the groups for the genesis block
func (sim *BlockchainSimulator) InitGroups(n uint16) {
	sim.group = make([]GroupSimulator, n)
	sim.grpmap = make(map[common.Address]*GroupSimulator)
	r := sim.seed.Ders("InitGroups")
	// build a temporary state datastructure from processes
	/* s := state.NewState()
	for _, p := range sim.proc {
		s.AddNode(p.reginfo)
	} */
	// create n groups
	for i := 0; i < int(n); i++ {
		// choose members based on r
		/* groupinfo := s.NewRandomGroup(r.Deri(i), sim.groupSize)
		   groupinfo.Log() */
		// LATER: replace the following using groupinfo
		indices := r.Deri(i).RandomPerm(len(sim.proc), int(sim.groupSize))
		members := make([]*ProcessSimulator, sim.groupSize)
		for j, idx := range indices {
			members[j] = &(sim.proc[idx])
		}
		sim.group[i] = NewGroupSimulator(members, sim.threshold)
		sim.grpmap[sim.group[i].Address()] = &sim.group[i]
		fmt.Println(sim.group[i].String())
	}
}

// NewBlockchainSimulator -- create a new blockchain simulation
// set the seed and define parameters like group size, threshold, number of processes etc.
func NewBlockchainSimulator(seed bls.Rand, groupSize uint16, threshold uint16, nProcesses uint, nGroups uint16) BlockchainSimulator {
	sim := BlockchainSimulator{seed: seed, groupSize: groupSize, threshold: threshold}
	sim.Log()

	// Start the processes first
	fmt.Printf("--- Process setup: (N)%d\n", nProcesses)
	sim.InitProcs(nProcesses)

	// Start the groups
	fmt.Printf("--- Group setup: (m)%d\n", nGroups)
	sim.InitGroups(nGroups)

	// Build the genesis block
	genesis := state.NewState()
	for _, p := range sim.proc {
		genesis.AddNode(p.reginfo)
		// this includes verification of proof-of-possession
	}
	for _, g := range sim.group {
		genesis.AddGroup(g.reginfo)
	}
	// the sig field remains empty because the genesis block is not signed

	// print op counts
	if Timing {
		bls.PrintCtrs()
	}

	// Build the chain with 1 block
	sim.chain = append(sim.chain, genesis)

	return sim
}

// Advance -- carry out the simulation for the given number of steps (blocks)
func (sim *BlockchainSimulator) Advance(n uint, verbose bool) {
	if n == 0 {
		return
	}
	// choose tip
	tip := sim.Tip()
	// select pre-determined random group from tip
	a := tip.SelectedGroupAddress()
	g := sim.grpmap[a]
	// get new group signature
	sig := g.Sign(tip.Rand().Bytes())
	if DoubleCheck {
		if !bls.VerifySig(tip.GroupPubkey(a), tip.Rand().Bytes(), sig) {
			fmt.Println("Error: group signature not valid.")
		}
	}

	// the new state is identical to the curren tip, except that we overwrite the signature
	newstate := tip

	// sign new state by group
	newstate.SetSignature(sig)

	// append new state
	sim.chain = append(sim.chain, newstate)

	// recurse
	sim.Advance(n-1, verbose)
	return
}

// Log -- print out a short form of the current state of the random beacon
func (sim *BlockchainSimulator) Log() {
	seed := sim.seed.Bytes()
	fmt.Printf("BlkCh: (n)%d (k)%d (seed)%x\n", sim.groupSize, sim.threshold, seed[:8])
	/*
		fmt.Println("  groups: ", len(sim.group))
		fmt.Println("  processes: ", len(sim.proc))
		fmt.Println("  chain height: ", len(sim.chain))
		sim.chain[len(sim.chain)-1].Log()
		for _, p := range sim.proc {
			p.Log()
		}
		for _, g := range sim.group {
			g.Log()
		}
	*/
}

// Length -- return the current block height
func (sim *BlockchainSimulator) Length() int {
	return len(sim.chain)
}

// Tip -- return the current state at the tip of the chain
func (sim *BlockchainSimulator) Tip() state.State {
	return sim.chain[len(sim.chain)-1]
}
