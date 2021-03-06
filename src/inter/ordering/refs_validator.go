package ordering

import (
	"fmt"

	"github.com/Fantom-foundation/go-lachesis/src/hash"
	"github.com/Fantom-foundation/go-lachesis/src/inter"
)

type refsValidator struct {
	event    *inter.Event
	creators map[hash.Peer]struct{}
}

func newRefsValidator(e *inter.Event) *refsValidator {
	return &refsValidator{
		event:    e,
		creators: make(map[hash.Peer]struct{}, len(e.Parents)),
	}
}

func (v *refsValidator) AddUniqueParent(node hash.Peer) error {
	if _, ok := v.creators[node]; ok {
		return fmt.Errorf("event %s has double refer to node %s",
			v.event.Hash().String(),
			node.String())
	}
	v.creators[node] = struct{}{}
	return nil

}

func (v *refsValidator) CheckSelfParent() error {
	if _, ok := v.creators[v.event.Creator]; !ok {
		return fmt.Errorf("event %s has no refer to self-node %s",
			v.event.Hash().String(),
			v.event.Creator.String())
	}
	return nil
}
