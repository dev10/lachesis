package node

import (
	"math/rand"

	"github.com/andrecronje/lachesis/src/peers"
)

//XXX PeerSelector needs major refactoring
// PeerSelector provides an interface for the lachesis node to 
// update the last peer it gossiped with and select the next peer
// to gossip with 
type PeerSelector interface {
	Peers() *peers.PeerSet
	UpdateLast(peer int64)
	Next() *peers.Peer
}

//+++++++++++++++++++++++++++++++++++++++
//RANDOM

type RandomPeerSelector struct {
	peers  *peers.PeerSet
	selfID int64
	last   int64
}

func NewRandomPeerSelector(peerSet *peers.PeerSet, selfID int64) *RandomPeerSelector {
	return &RandomPeerSelector{
		selfID: selfID,
		peers:  peerSet,
	}
}

func (ps *RandomPeerSelector) Peers() *peers.PeerSet {
	return ps.peers
}

func (ps *RandomPeerSelector) UpdateLast(peer int64) {
	ps.last = peer
}

func (ps *RandomPeerSelector) Next() *peers.Peer {
	selectablePeers := ps.peers.Peers

	if len(selectablePeers) > 1 {
		_, selectablePeers = peers.ExcludePeer(selectablePeers, ps.selfID)

		if len(selectablePeers) > 1 {
			_, selectablePeers = peers.ExcludePeer(selectablePeers, ps.last)
		}
	}

	i := rand.Intn(len(selectablePeers))

	peer := selectablePeers[i]

	return peer
}
