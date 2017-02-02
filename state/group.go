package state

import (
	"dfinity/beacon/bls"
	dfn "dfinity/beacon/common"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"log"
)

// Group -- encodes all data of a group as recorded on the blockchain
type Group struct {
	members []common.Address
	// group pubkey
	pub       bls.Pubkey
	threshold uint16
}

// NewGroup -- create a new Group struct with list of members and empty pubkey
func NewGroup(addresses []common.Address, k uint16) Group {
	return Group{addresses, bls.Pubkey{}, k}
}

// SetPubkey -- set the group's pubkey and threshold
func (g *Group) SetPubkey(pub bls.Pubkey, k uint16) {
	g.pub = pub
	g.threshold = k
}

// Getters

// Address - the group address
func (g Group) Address() (a common.Address) {
	// hash of all member addresses
	d := sha3.NewKeccak256()
	addresses := g.members
	dfn.SortAddresses(addresses)
	var err error
	for _, addr := range addresses {
		_, err = d.Write(addr[:])
		if err != nil {
			log.Fatalln("Error when calling Keccak256")
		}
	}
	var h common.Hash
	d.Sum(h[:0])
	return common.BytesToAddress(h[:])
}

// Pubkey -- the group pubkey
func (g Group) Pubkey() bls.Pubkey {
	return g.pub
}

// Members -- the list of members
func (g Group) Members() []common.Address {
	return g.members
}

// Threshold -- the threshold used in the setup
func (g Group) Threshold() int {
	return int(g.threshold)
}

// Size -- the number of members
func (g Group) Size() int {
	return len(g.members)
}

// Log -- print multi-line group state
func (g Group) Log() {
	fmt.Println("    members: ", len(g.members))
	for _, m := range g.members {
		fmt.Printf("      address: % x\n", m)
	}
	fmt.Printf("    addr: % x\n", g.pub.Address())
	fmt.Println("    threshold: ", g.threshold)
}

// String -- one-line summary representation
func (g Group) String() string {
	a := g.Address()
	mem := "["
	for i, m := range g.members {
		if i > 0 {
			mem += ","
		}
		mem += fmt.Sprintf("%x", string(m[:2]))
	}
	mem += "]"
	return fmt.Sprintf("GrpR: (addr)%x (pub)%.8s (n)%d (k)%d (mem)%s", a[:2], g.pub.String(), len(g.members), g.threshold, mem)
}

// isValid --
/* TODO: check if group pubkey is individually signed by enough group members */
func (g Group) isValid() bool {
	return true
}
