package state

import (
	"dfinity/beacon/bls"
	dfn "dfinity/beacon/common"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
)

// State -- encodes the state of the chain (state of 1 block)
type State struct {
	nodes  map[common.Address]Node
	groups map[common.Address]Group
	sig    bls.Signature
}

// NewState --
func NewState() State {
	s := State{}
	s.nodes = make(map[common.Address]Node)
	s.groups = make(map[common.Address]Group)
	return s
}

// AddNode --
func (s *State) AddNode(n Node) (valid bool) {
	valid = n.hasPop()
	if valid {
		s.nodes[n.Address()] = n
	}
	return
}

// AddGroup --
func (s *State) AddGroup(g Group) (valid bool) {
	valid = g.isValid()
	if valid {
		s.groups[g.Address()] = g
	}
	return
}

// SetSignature --
func (s *State) SetSignature(sig bls.Signature) {
	s.sig = sig
}

// Rand --
func (s State) Rand() bls.Rand {
	return s.sig.Rand()
}

// NodeAddressList --
func (s State) NodeAddressList() []common.Address {
	// return a sorted list of addresses of all nodes
	addresses := make([]common.Address, len(s.nodes))
	var i int
	for a := range s.nodes {
		addresses[i] = a
		i++
	}
	dfn.SortAddresses(addresses)
	return addresses
}

// NewRandomGroup --
func (s State) NewRandomGroup(r bls.Rand, n uint16) Group {
	N := len(s.nodes) // need n <= N
	fmt.Println(N, n)
	// get sorted list of nodes
	nodes := s.NodeAddressList()
	// choose members based on r
	indices := r.RandomPerm(N, int(n))
	members := make([]common.Address, int(n))
	for j, idx := range indices {
		members[j] = nodes[idx]
	}
	return Group{members, bls.Pubkey{}, 0}
}

// GroupAddressList --
func (s State) GroupAddressList() []common.Address {
	// return a sorted list of addresses of all groups
	addresses := make([]common.Address, len(s.groups))
	var i int
	for k := range s.groups {
		addresses[i] = k
		i++
	}
	dfn.SortAddresses(addresses)
	return addresses
}

// SelectedGroupAddress --
func (s State) SelectedGroupAddress() common.Address {
	i := s.Rand().Modulo(len(s.groups))
	return s.GroupAddressList()[i]
}

// GroupPubkey --
func (s State) GroupPubkey(a common.Address) bls.Pubkey {
	return s.groups[a].pub
}

// SelectedGroupPubkey --
func (s State) SelectedGroupPubkey() bls.Pubkey {
	return s.GroupPubkey(s.SelectedGroupAddress())
}

// Log --
func (s State) Log() {
	fmt.Println("State: ")
	fmt.Println("  sig:  ", s.sig)
	fmt.Printf("  rand: % x\n", s.Rand())
	fmt.Println("  nodes: ", len(s.nodes))
	for i, a := range s.NodeAddressList() {
		fmt.Printf("    %d. % x\n", i, a)
	}
	fmt.Println("  groups: ", len(s.groups))
	for i, a := range s.groups {
		fmt.Printf("    %d. % x\n", i, a)
	}
	fmt.Printf("  selected group(s):\n")
	fmt.Printf("    %d. % x\n", 1, s.SelectedGroupAddress())
}

// String --
func (s State) String(long bool) string {
	rnd := s.Rand().Bytes()
	str := fmt.Sprintf("Stat: (sig)%.8s (rnd)%.2x (N)%d (m)%d (grp)%.2x", s.sig.String(), rnd, len(s.nodes), len(s.groups), s.SelectedGroupAddress())
	if long {
		str += "\n"
		for i, a := range s.NodeAddressList() {
			str += fmt.Sprintf("  %3d. % s\n", i, s.nodes[a].String())
		}
		for i, a := range s.GroupAddressList() {
			str += fmt.Sprintf("  %3d. % s\n", i, s.groups[a].String())
		}
	}
	return str
}
