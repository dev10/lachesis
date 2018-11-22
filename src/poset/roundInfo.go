package poset

import (
	"github.com/andrecronje/lachesis/src/peers"
	"github.com/golang/protobuf/proto"
)

type pendingRound struct {
	Index   int64
	Decided bool
}

type RoundInfo struct {
	Message RoundInfoMessage
	queued bool
}

func NewRoundInfo(peers *peers.PeerSet) *RoundInfo {
	return &RoundInfo{
		Message: RoundInfoMessage{
			CreatedEvents: make(map[string]*RoundEvent),
			ReceivedEvents: []string{},
			PeerSet: peers,
		},
	}
}

func (r *RoundInfo) AddCreatedEvent(x string, witness bool) {
	_, ok := r.Message.CreatedEvents[x]
	if !ok {
		r.Message.CreatedEvents[x] = &RoundEvent{
			Witness: witness,
		}
	}
}

func (r *RoundInfo) AddReceivedEvent(x string) {
	r.Message.ReceivedEvents = append(r.Message.ReceivedEvents, x)
}

func (r *RoundInfo) SetFame(x string, f bool) {
	e, ok := r.Message.CreatedEvents[x]
	if !ok {
		e = &RoundEvent{
			Witness: true,
		}
	}
	if f {
		e.Famous = Trilean_TRUE
	} else {
		e.Famous = Trilean_FALSE
	}
	r.Message.CreatedEvents[x] = e
}

//return true if no witnesses' fame is left undefined
func (r *RoundInfo) WitnessesDecided() bool {
	c := int64(0)
	for _, e := range r.Message.CreatedEvents {
		if e.Witness && e.Famous != Trilean_UNDEFINED {
			c++
		}
	}
	return c >= r.Message.PeerSet.SuperMajority()
}

//return witnesses
func (r *RoundInfo) Witnesses() []string {
	var res []string
	for x, e := range r.Message.CreatedEvents {
		if e.Witness {
			res = append(res, x)
		}
	}
	return res
}

func (r *RoundInfo) RoundEvents() []string {
	var res []string
	for x, e := range r.Message.CreatedEvents {
		if !e.Consensus {
			res = append(res, x)
		}
	}
	return res
}

//return consensus events
func (r *RoundInfo) ConsensusEvents() []string {
	var res []string
	for x, e := range r.Message.CreatedEvents {
		if e.Consensus {
			res = append(res, x)
		}
	}
	return res
}

//return famous witnesses
func (r *RoundInfo) FamousWitnesses() []string {
	var res []string
	for x, e := range r.Message.CreatedEvents {
		if e.Witness && e.Famous == Trilean_TRUE {
			res = append(res, x)
		}
	}
	return res
}

func (r *RoundInfo) IsDecided(witness string) bool {
	w, ok := r.Message.CreatedEvents[witness]
	return ok && w.Witness && w.Famous != Trilean_UNDEFINED
}

func (r *RoundInfo) ProtoMarshal() ([]byte, error) {
	var bf proto.Buffer
	bf.SetDeterministic(true)
	if err := bf.Marshal(&r.Message); err != nil {
		return nil, err
	}
	return bf.Bytes(), nil
}

func (r *RoundInfo) ProtoUnmarshal(data []byte) error {
	return proto.Unmarshal(data, &r.Message)
}

func (r *RoundInfo) IsQueued() bool {
	return r.queued
}

func (this *RoundEvent) Equals(that *RoundEvent) bool {
	return this.Consensus == that.Consensus &&
		this.Witness == that.Witness &&
		this.Famous == that.Famous
}

func EqualsMapStringRoundEvent(this map[string]*RoundEvent, that map[string]*RoundEvent) bool {
	if len(this) != len(that) {
		return false
	}
	for k, v := range this {
		v2, ok := that[k]
		if !ok || !v2.Equals(v) {
			return false
		}
	}
	return true
}

func (this *RoundInfo) Equals(that *RoundInfo) bool {
	return this.queued == that.queued &&
		EqualsMapStringRoundEvent(this.Message.CreatedEvents, that.Message.CreatedEvents)
}
