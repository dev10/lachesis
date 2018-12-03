package node

import (
	"fmt"

	"github.com/Fantom-foundation/go-lachesis/src/poset"
)

type Infos struct {
	ParticipantEvents map[string]map[string]*poset.Event
	Rounds            []*poset.RoundInfo
  	Blocks            []*poset.Block
}

type Graph struct {
	*Node
}

func (g *Graph) GetBlocks() []*poset.Block {
	res := []*poset.Block{}
 	blockIdx := int64(0)
	store := g.Node.core.poset.Store

 	for blockIdx <= store.LastBlockIndex() {
		r, err := store.GetBlock(blockIdx)
 		if err != nil {
			break
		}
 		res = append(res, r)
 		blockIdx++
	}
 	return res
}

func (g *Graph) GetParticipantEvents() (map[string]map[string]*poset.Event, error) {
	res := make(map[string]map[string]*poset.Event)

	store := g.Node.core.poset.Store
	peers := g.Node.core.poset.Peers

	for _, p := range peers.ByPubKey {
		root, err := store.GetRoot(p.PubKeyHex)

		if err != nil {
			return res, err
		}

		evs, err := store.ParticipantEvents(p.PubKeyHex, root.SelfParent.Index)

		if err != nil {
			return res, err
		}

		res[p.PubKeyHex] = make(map[string]*poset.Event)

		selfParent := fmt.Sprintf("Root%d", p.ID)

		flagTable := make(map[string]int64)
		flagTable[selfParent] = 1

		// Create and save the first Event
		initialEvent := poset.NewEvent([][]byte{},
			[]poset.InternalTransaction{},
			[]poset.BlockSignature{},
			[]string{}, []byte{}, 0, flagTable)

		res[p.PubKeyHex][root.SelfParent.Hash] = initialEvent

		for _, e := range evs {
			event, err := store.GetEvent(e)
			if err != nil {
				return res, err
			}

			hash := event.Hex()

			res[p.PubKeyHex][hash] = event
		}
	}

	return res, nil
}

func (g *Graph) GetRounds() []*poset.RoundInfo {
	res := []*poset.RoundInfo{}

	round := int64(0)

	store := g.Node.core.poset.Store

	for round <= store.LastRound() {
		r, err := store.GetRound(round)

		if err != nil || !r.IsQueued() {
			break
		}

		res = append(res, r)

		round++
	}

	return res
}

func (g *Graph) GetInfos() (Infos, error) {
	participantEvents, err := g.GetParticipantEvents()
	if err != nil {
		return Infos{}, err
	}

	return Infos{
		ParticipantEvents: participantEvents,
		Rounds:            g.GetRounds(),
    	Blocks:            g.GetBlocks(),
	}, nil
}

func NewGraph(n *Node) *Graph {
	return &Graph{
		Node: n,
	}
}
