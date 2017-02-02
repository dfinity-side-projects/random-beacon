package state

import (
	"dfinity/beacon/bls"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
)

// Node -- encodes the information with which a process (=node) is registered on the blockchain
type Node struct {
	pub bls.Pubkey
	pop bls.Pop
}

// Constructors

// NodeFromSeckey --
func NodeFromSeckey(sec bls.Seckey) Node {
	pub := bls.PubkeyFromSeckey(sec)
	return Node{pub, bls.GeneratePop(sec, pub)}
}

// Getters

// Address --
func (n Node) Address() common.Address {
	return n.pub.Address()
}

// ID --
func (n Node) ID() bls.ID {
	return bls.IDFromBig(n.Address().Big())
}

// hasPop --
func (n Node) hasPop() bool {
	return bls.VerifyPop(n.pub, n.pop)
}

// Log --
func (n Node) Log() {
	fmt.Printf("    pub: % x\n", n.pub.Address())
	//	fmt.Printf("  Seckey: % x\n", p.sec.Bytes())
	fmt.Println("    pop: ", n.pop)
}

// String --
func (n Node) String() string {
	a := n.pub.Address()
	return fmt.Sprintf("Node: (addr)%x (pub)%s", string(a[:2]), n.pub.String()[:8])
}
