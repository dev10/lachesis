package peers

import (
	"fmt"
	"math"

	"github.com/andrecronje/lachesis/src/crypto"
)

//XXX exclude peers should be in here

/* Constructors */

func NewEmptyPeerSet() *PeerSet {
	return &PeerSet{
		ByPubKey: make(map[string]*Peer),
		ById:     make(map[int64]*Peer),
	}
}

//NewPeerSet creates a new PeerSet from a list of Peers
func NewPeerSet(peers []*Peer) *PeerSet {
	peerSet := NewEmptyPeerSet()
	for _, peer := range peers {
		if peer.ID == 0 {
			peer.ComputeID()
		}

		peerSet.ByPubKey[peer.PubKeyHex] = peer
		peerSet.ById[peer.ID] = peer
	}

	peerSet.Peers = peers

	return peerSet
}

//WithNewPeer returns a new PeerSet with a list of peers including the new one.
func (peerSet *PeerSet) WithNewPeer(peer *Peer) *PeerSet {
	peers := append(peerSet.Peers, peer)
	newPeerSet := NewPeerSet(peers)
	return newPeerSet
}

//WithRemovedPeer returns a new PeerSet with a list of peers exluding the
//provided one
func (peerSet *PeerSet) WithRemovedPeer(peer *Peer) *PeerSet {
	peers := []*Peer{}
	for _, p := range peerSet.Peers {
		if p.PubKeyHex != peer.PubKeyHex {
			peers = append(peers, p)
		}
	}
	newPeerSet := NewPeerSet(peers)
	return newPeerSet
}

/* ToSlice Methods */

//PubKeys returns the PeerSet's slice of public keys
func (c *PeerSet) PubKeys() []string {
	res := []string{}

	for _, peer := range c.Peers {
		res = append(res, peer.PubKeyHex)
	}

	return res
}

//IDs returns the PeerSet's slice of IDs
func (c *PeerSet) IDs() []int64 {
	res := []int64{}

	for _, peer := range c.Peers {
		res = append(res, peer.ID)
	}

	return res
}

/* Utilities */

//Len returns the number of Peers in the PeerSet
func (c *PeerSet) Len() int {
	return len(c.ByPubKey)
}

//Hash uniquely identifies a PeerSet. It is computed by sorting the peers set
//by ID, and hashing (SHA256) their public keys together, one by one.
func (c *PeerSet) Hash() ([]byte, error) {
	if len(c.Hash_) == 0 {
		hash := []byte{}
		for _, p := range c.Peers {
			pk, _ := p.PubKeyBytes()
			hash = crypto.SimpleHashFromTwoHashes(hash, pk)
		}
		c.Hash_ = hash
	}
	return c.Hash_, nil
}

//Hex is the hexadecimal representation of Hash
func (c *PeerSet) Hex() string {
	if len(c.Hex_) == 0 {
		hash, _ := c.Hash()
		c.Hex_ = fmt.Sprintf("0x%X", hash)
	}
	return c.Hex_
}

//SuperMajority return the number of peers that forms a strong majortiy (+2/3)
//in the PeerSet
func (c *PeerSet) SuperMajority() int64 {
	if c.SuperMajority_ == 0 {
		val := int64(2*c.Len()/3 + 1)
		c.SuperMajority_ = val
	}
	return c.SuperMajority_
}

func (c *PeerSet) TrustCount() int64 {
	if c.TrustCount_ == 0 {
		val := int64(math.Ceil(float64(c.Len()) / float64(3)))
		c.TrustCount_ = val
	}
	return c.TrustCount_
}

func (c *PeerSet) clearCache() {
	c.Hash_ = []byte{}
	c.Hex_ = ""
	c.SuperMajority_ = 0
}
