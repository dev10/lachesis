package peers

import (
	"encoding/hex"

	"github.com/andrecronje/lachesis/src/common"
)

const (
	jsonPeerPath = "peers.json"
)

func NewPeer(pubKeyHex, netAddr string) *Peer {
	peer := &Peer{
		PubKeyHex: pubKeyHex,
		NetAddr:   netAddr,
	}

	peer.ComputeID()

	return peer
}

func (this *Peer) Equals(that *Peer) bool {
	return this.ID == that.ID &&
		this.NetAddr == that.NetAddr &&
		this.PubKeyHex == that.PubKeyHex
}

func (p *Peer) PubKeyBytes() ([]byte, error) {
	return hex.DecodeString(p.PubKeyHex[2:])
}

func (p *Peer) ComputeID() error {
	// TODO: Use the decoded bytes from hex
	pubKey, err := p.PubKeyBytes()

	if err != nil {
		return err
	}

	p.ID = int64(common.Hash32(pubKey))

	return nil
}

// PeerStore provides an interface for persistent storage and
// retrieval of peers.
type PeerStore interface {
	// Peers returns the list of known peers.
	Peers() (*PeerSet, error)

	// SetPeers sets the list of known peers. This is invoked when a peer is
	// added or removed.
	SetPeers([]*Peer) error
}

// ExcludePeer is used to exclude a single peer from a list of peers.
func ExcludePeer(peers []*Peer, peer int64) (int64, []*Peer) {
	index := int64(-1)
	otherPeers := make([]*Peer, 0, len(peers))
	for i, p := range peers {
		if p.ID != peer {
			otherPeers = append(otherPeers, p)
		} else {
			index = int64(i)
		}
	}
	return index, otherPeers
}
