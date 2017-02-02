package sim

import (
	"dfinity/beacon/bls"
	"dfinity/beacon/state"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
)

/// simulator section

// ProcessSimulator -- encode everything needed to simulate a single process
type ProcessSimulator struct {
	sec     bls.Seckey
	reginfo state.Node
	rseed   bls.Rand
	// rseed is the seed used for the internal randomness of the process, it did not seed the secret key
	sharesSource   map[common.Address]bls.SeckeyMap
	sharesCombined bls.SeckeyMap
}

// NewProcessSimulator -- create a new simulator given process data such as seed and private key
func NewProcessSimulator(sec bls.Seckey, seed bls.Rand) (p ProcessSimulator) {
	p.sec = sec
	p.reginfo = state.NodeFromSeckey(sec)
	p.rseed = seed
	p.sharesSource = make(map[common.Address]bls.SeckeyMap)
	p.sharesCombined = bls.SeckeyMap{}
	return
}

// NewProcessSimulatorDet -- create a new _deterministic_ simulator given only the private key
// The seed is derived from the private key.
// This makes the process simulator DETERMINISTIC by setting rseed to the process' own address.
// Since the address is public the process' behaviour becomes predictable from the outside.
// This will benefit testing.
func NewProcessSimulatorDet(sec bls.Seckey) ProcessSimulator {
	node := state.NodeFromSeckey(sec)
	// assign temporary variable to make the value addressable
	tmp := node.Address()
	return NewProcessSimulator(sec, bls.RandFromBytes(tmp[:]))
}

// Address -- return the address of the simulated process
func (p *ProcessSimulator) Address() common.Address {
	return p.reginfo.Address()
}

// SetGroupShare -- set the incoming shares from other group members
func (p *ProcessSimulator) SetGroupShare(addr common.Address, source common.Address, share bls.Seckey, vvec []bls.Pubkey) {
	//	fmt.Printf("Setting source share: (proc)%.4x (grp)%.2x (src)%.4x (sec)%.4s\n", p.Address(), addr, source, share.String())
	// verify share
	if Vvec {
		if bls.SharePubkey(vvec, p.reginfo.ID()).String() != bls.PubkeyFromSeckey(share).String() {
			fmt.Println("Error: Received secret share does not match committed verification vector")
		}
	}

	// if key source does not exist yet then make a bls.SeckeyMap
	_, exists := p.sharesSource[addr]
	if !exists {
		p.sharesSource[addr] = bls.SeckeyMap{}
	}
	// store source share
	p.sharesSource[addr][source] = share
	return
}

// AggregateGroupShares -- aggregate (sum up) all the shares that came in from members of the given group
func (p *ProcessSimulator) AggregateGroupShares(g state.Group) {
	addr := g.Address()
	vlist := make([]bls.Seckey, len(p.sharesSource[addr]))
	i := 0
	for _, sec := range p.sharesSource[addr] {
		vlist[i] = sec
		i++
	}
	p.sharesCombined[addr] = bls.AggregateSeckeys(vlist)
	return
}

// GetAggregatedGroupShare -- return the aggregated share received from the given group
func (p *ProcessSimulator) GetAggregatedGroupShare(g state.Group) bls.Seckey {
	return p.sharesCombined[g.Address()]
}

// GetSeckeyForGroup -- return the own secret provided for the group setup (function of internal seed and group address)
func (p *ProcessSimulator) GetSeckeyForGroup(g state.Group) (sec bls.Seckey) {
	addr := g.Address()
	gseed := p.rseed.DerivedRand(addr[:])
	sec = bls.SeckeyFromRand(gseed.Deri(0))
	//	fmt.Printf("sec for group: %s\n", sec.String())
	return
}

// GetSeckeySharesForGroup -- take own secret for the group setup (function of internal seed and group address) and split it up in shares for all group members
// from the process seed (rseed) and derive a per-group seed based on the group's address
func (p *ProcessSimulator) GetSeckeySharesForGroup(g state.Group) (bls.SeckeyMap, []bls.Pubkey) {
	addr := g.Address()
	gseed := p.rseed.DerivedRand(addr[:])
	// from the per-group seed derive a vector of k seckeys as the master seckey where k is the threshold
	// the master seckey defines a polynomial of degree k-1
	k := g.Threshold()
	msec := make([]bls.Seckey, k)
	vvec := make([]bls.Pubkey, k)
	for i := 0; i < k; i++ {
		msec[i] = bls.SeckeyFromRand(gseed.Deri(i))
		vvec[i] = bls.PubkeyFromSeckey(msec[i])
	}
	shares := bls.SeckeyMap{}
	for _, m := range g.Members() {
		shares[m] = bls.ShareSeckeyByAddr(msec, m)
	}
	return shares, vvec
}

// SignForGroup -- return the signature share for the given message and group
func (p *ProcessSimulator) SignForGroup(g state.Group, msg []byte) bls.Signature {
	sec := p.sharesCombined[g.Address()]
	//	fmt.Printf("sign for group: (grp)%.2x (sec)%x\n", g.Address(), sec.String())
	return bls.Sign(sec, msg)
}

// Sign -- return the own individual signature for the given message
func (p *ProcessSimulator) Sign(msg []byte) bls.Signature {
	return bls.Sign(p.sec, msg)
}

// Log -- print the state of the simulated process
func (p *ProcessSimulator) Log() {
	fmt.Printf("Process simulator: % x\n", p.reginfo.Address())
	fmt.Printf("  Seckey: % x\n", p.sec.Bytes())
	fmt.Printf("  rseed: % x\n", p.rseed)
	p.reginfo.Log()
}

// String -- return a very short summary of the state of the simulated process
func (p *ProcessSimulator) String() string {
	return fmt.Sprintf("Proc: (sec)%s (seed)%x %s", p.sec.String()[:4], p.rseed.String()[:2], p.reginfo.String())
}
