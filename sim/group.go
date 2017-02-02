package sim

import (
	"dfinity/beacon/bls"
	"dfinity/beacon/state"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"time"
)

// GroupSimulator -- encodes everything needed to simulate the behavior of a group
type GroupSimulator struct {
	// the group secret key (sec) is here only to enable optional double-checks
	sec      bls.Seckey
	reginfo  state.Group
	proclist []*ProcessSimulator
	procmap  map[common.Address]*ProcessSimulator
}

// ExchangeSeckeyShares -- make all group members exchange secret shares with each other
func ExchangeSeckeyShares(g state.Group, members []*ProcessSimulator) {
	for _, p := range members {
		// get secret shares for all other processes
		shares, vvec := p.GetSeckeySharesForGroup(g)
		// send shares out to all other individual processes
		for _, q := range members {
			q.SetGroupShare(g.Address(), p.Address(), shares[q.Address()], vvec)
		}
		// optional double-check of the group secret
		if DoubleCheck {
			sec := p.GetSeckeyForGroup(g)
			recovered := bls.RecoverSeckeyByMap(shares, g.Threshold())
			if sec.String() != recovered.String() {
				fmt.Println("Error: recovered seckey share (ByMap) does not match.")
			}
		}
	}
}

// NewGroupSimulator -- create a new group simulator, given simulators of its members
func NewGroupSimulator(members []*ProcessSimulator, k uint16) GroupSimulator {
	m := len(members)
	// collect all members' addresses in a Group struct with empty Pubkey
	addresses := make([]common.Address, m)
	pmap := make(map[common.Address]*ProcessSimulator)
	for i, p := range members {
		addresses[i] = p.Address()
		pmap[p.Address()] = p
	}
	g := state.NewGroup(addresses, k)

	// get all members' contribution to the group secret
	ExchangeSeckeyShares(g, members)

	// build group pubkey
	pubs := make([]bls.Pubkey, m)
	for i, p := range members {
		pubs[i] = bls.PubkeyFromSeckey(p.GetSeckeyForGroup(g))
	}
	pub := bls.AggregatePubkeys(pubs)

	// set group pubkey in Group struct
	g.SetPubkey(pub, k)

	// tell each process to aggregate their shares
	// processes need their aggregated shares for signing later
	for _, q := range members {
		q.AggregateGroupShares(g)
	}

	var sec bls.Seckey
	if DoubleCheck {
		// fetch the combined shares from each process into a SeckeyMap
		// (every process does this so we wouldn't need to in the group simulator)
		aggShares := bls.SeckeyMap{}
		for _, p := range members {
			aggShares[p.Address()] = p.GetAggregatedGroupShare(g)
		}

		// recover the combined group secret from combined shares
		// choose k random shares, combine and compare
		sec = bls.RecoverSeckeyByMap(aggShares, int(k))
		pubDup := bls.PubkeyFromSeckey(sec)

		// optional double-check: aggregate all contributions into the group secret and compare
		secs := make([]bls.Seckey, m)
		for i, p := range members {
			secs[i] = p.GetSeckeyForGroup(g)
		}
		secDup := bls.AggregateSeckeys(secs)
		if sec.String() != secDup.String() {
			fmt.Println("Error: recovered aggregated seckey does not match.")
		}

		if pub.String() != pubDup.String() {
			fmt.Println("Error: recovered aggregated pubkey does not match.")
		}
	}

	return GroupSimulator{sec, g, members, pmap}
}

// Sign -- make the group members jointly create a group signature
func (g GroupSimulator) Sign(msg []byte) bls.Signature {
	sigmap := make(map[common.Address]bls.Signature)
	// get signature share from each process
	t0 := time.Now()
	for _, p := range g.proclist {
		sigmap[p.Address()] = p.SignForGroup(g.reginfo, msg)
	}
	delta1 := time.Since(t0)
	t1 := time.Now()
	sig1 := bls.RecoverSignatureByMap(sigmap, g.reginfo.Threshold())
	delta2 := time.Since(t1)
	if Timing {
		fmt.Printf("Time for group signatures with %d shares: %v (%vus / share) + %v (recovery).\n", len(g.proclist), delta1, (delta1.Nanoseconds()/1000)/int64(len(g.proclist)), delta2)
	}

	// optional verification
	if DoubleCheck {
		sig2 := bls.Sign(g.sec, msg)
		if sig1.String() != sig2.String() {
			fmt.Println("Error in Group sign: Recovered signature does not match.")
		}
	}

	return sig1
}

// Address -- return the address under which the simulated group is registered
func (g GroupSimulator) Address() common.Address {
	return g.reginfo.Address()
}

// Log -- print the current state of the simulated group
func (g GroupSimulator) Log() {
	fmt.Printf("Group simulator: % x\n", g.reginfo.Address())
	fmt.Printf("  Seckey: % x\n", g.sec.Bytes())
	g.reginfo.Log()
}

// String -- return a very short summary of the state of the simulated group
func (g *GroupSimulator) String() string {
	return fmt.Sprintf("GrpP: (sec)%s %s", g.sec.String()[:4], g.reginfo.String())
}
