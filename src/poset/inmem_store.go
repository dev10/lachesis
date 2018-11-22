package poset

import (
	_ "fmt"
	"strconv"

	cm "github.com/Fantom-foundation/go-lachesis/src/common"
	"github.com/Fantom-foundation/go-lachesis/src/peers"
)

type InmemStore struct {
	cacheSize              int
	eventCache             *cm.LRU          //hash => Event
	roundCache             *cm.LRU          //round number => Round
	blockCache             *cm.LRU          //index => Block
	frameCache             *cm.LRU          //round received => Frame
	consensusCache         *cm.RollingIndex //consensus index => hash
	totConsensusEvents     int64
	peerSetCache           *PeerSetCache //start round => PeerSet
	repertoireByPubKey     map[string]*peers.Peer
	repertoireByID         map[int64]*peers.Peer
	participantEventsCache *ParticipantEventsCache //pubkey => Events
	rootsByParticipant     map[string]*Root        //[participant] => Root
	rootsBySelfParent      map[string]*Root        //[Root.SelfParent.Hash] => Root
	lastRound              int64
	lastConsensusEvents    map[string]string //[participant] => hex() of last consensus event
	lastBlock              int64
}

func NewInmemStore(peerSet *peers.PeerSet, cacheSize int) *InmemStore {
	store := &InmemStore{
		cacheSize:              cacheSize,
		eventCache:             cm.NewLRU(cacheSize, nil),
		roundCache:             cm.NewLRU(cacheSize, nil),
		blockCache:             cm.NewLRU(cacheSize, nil),
		frameCache:             cm.NewLRU(cacheSize, nil),
		consensusCache:         cm.NewRollingIndex("ConsensusCache", cacheSize),
		peerSetCache:           NewPeerSetCache(),
		repertoireByPubKey:     make(map[string]*peers.Peer),
		repertoireByID:         make(map[int64]*peers.Peer),
		participantEventsCache: NewParticipantEventsCache(cacheSize, peerSet),
		rootsByParticipant:     make(map[string]*Root),
		rootsBySelfParent:      make(map[string]*Root),
		lastRound:              -1,
		lastBlock:              -1,
		lastConsensusEvents:    map[string]string{},
	}

	store.SetPeerSet(0, peerSet)

	return store
}

func NewInmemStoreFromFrame(frame *Frame, cacheSize int) (*InmemStore, error) {
	store := &InmemStore{
		cacheSize: cacheSize,
	}

	if err := store.Reset(frame); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *InmemStore) CacheSize() int {
	return s.cacheSize
}

func (s *InmemStore) GetPeerSet(round int64) (*peers.PeerSet, error) {
	return s.peerSetCache.Get(round)
}

func (s *InmemStore) GetLastPeerSet() (*peers.PeerSet, error) {
	return s.peerSetCache.GetLast()
}

func (s *InmemStore) SetPeerSet(round int64, peerSet *peers.PeerSet) error {
	//Update PeerSetCache
	err := s.peerSetCache.Set(round, peerSet)
	if err != nil {
		return err
	}

	//Extend PartipantEventsCache and Roots with new peers
	for id, p := range peerSet.ById {
		if _, ok := s.participantEventsCache.participants.ById[id]; !ok {
			if err := s.participantEventsCache.AddPeer(p); err != nil {
				return err
			}
		}

		if _, ok := s.rootsByParticipant[p.PubKeyHex]; !ok {
			root := NewBaseRoot(p.ID)
			s.rootsByParticipant[p.PubKeyHex] = root
			s.rootsBySelfParent[root.SelfParent.Hash] = root
		}

		s.repertoireByPubKey[p.PubKeyHex] = p
		s.repertoireByID[p.ID] = p
	}

	return nil
}

func (s *InmemStore) RepertoireByPubKey() map[string]*peers.Peer {
	return s.repertoireByPubKey
}

func (s *InmemStore) RepertoireByID() map[int64]*peers.Peer {
	return s.repertoireByID
}

func (s *InmemStore) RootsBySelfParent() map[string]*Root {
	return s.rootsBySelfParent
}

func (s *InmemStore) GetEvent(key string) (*Event, error) {
	res, ok := s.eventCache.Get(key)
	if !ok {
		return nil, cm.NewStoreErr("EventCache", cm.KeyNotFound, key)
	}

	return res.(*Event), nil
}

func (s *InmemStore) SetEvent(event *Event) error {
	key := event.Hex()
	_, err := s.GetEvent(key)
	if err != nil && !cm.Is(err, cm.KeyNotFound) {
		return err
	}
	if cm.Is(err, cm.KeyNotFound) {
		if err := s.addParticpantEvent(event.Creator(), key, event.Index()); err != nil {
			return err
		}
	}

	// fmt.Println("Adding event to cache", event.Hex())
	s.eventCache.Add(key, event)

	return nil
}

func (s *InmemStore) addParticpantEvent(participant string, hash string, index int64) error {
	return s.participantEventsCache.Set(participant, hash, index)
}

func (s *InmemStore) ParticipantEvents(participant string, skip int64) ([]string, error) {
	return s.participantEventsCache.Get(participant, skip)
}

func (s *InmemStore) ParticipantEvent(participant string, index int64) (string, error) {
	ev, err := s.participantEventsCache.GetItem(participant, index)
	if err != nil {
		root, ok := s.rootsByParticipant[participant]
		if !ok {
			return "", cm.NewStoreErr("InmemStore.Roots", cm.NoRoot, participant)
		}
		if root.SelfParent.Index == index {
			ev = root.SelfParent.Hash
			err = nil
		}
	}
	return ev, err
}

func (s *InmemStore) LastEventFrom(participant string) (last string, isRoot bool, err error) {
	//try to get the last event from this participant
	last, err = s.participantEventsCache.GetLast(participant)

	//if there is none, grab the root
	if err != nil && cm.Is(err, cm.Empty) {
		root, ok := s.rootsByParticipant[participant]
		if ok {
			last = root.SelfParent.Hash
			isRoot = true
			err = nil
		} else {
			err = cm.NewStoreErr("InmemStore.Roots", cm.NoRoot, participant)
		}
	}
	return
}

func (s *InmemStore) LastConsensusEventFrom(participant string) (last string, isRoot bool, err error) {
	//try to get the last consensus event from this participant
	last, ok := s.lastConsensusEvents[participant]
	//if there is none, grab the root
	if !ok {
		root, ok := s.rootsByParticipant[participant]
		if ok {
			last = root.SelfParent.Hash
			isRoot = true
		} else {
			err = cm.NewStoreErr("InmemStore.Roots", cm.NoRoot, participant)
		}
	}
	return
}

func (s *InmemStore) KnownEvents() map[int64]int64 {
	known := s.participantEventsCache.Known()
	lastPeerSet, _ := s.GetLastPeerSet()
	if lastPeerSet != nil {
		for p, pid := range lastPeerSet.ByPubKey {
			if known[pid.ID] == -1 {
				root, ok := s.rootsByParticipant[p]
				if ok {
					known[pid.ID] = root.SelfParent.Index
				}
			}
		}
	}
	return known
}

func (s *InmemStore) ConsensusEvents() []string {
	lastWindow, _ := s.consensusCache.GetLastWindow()
	res := make([]string, len(lastWindow))
	for i, item := range lastWindow {
		res[i] = item.(string)
	}
	return res
}

func (s *InmemStore) ConsensusEventsCount() int64 {
	return s.totConsensusEvents
}

func (s *InmemStore) AddConsensusEvent(event *Event) error {
	s.consensusCache.Set(event.Hex(), s.totConsensusEvents)
	s.totConsensusEvents++
	s.lastConsensusEvents[event.Creator()] = event.Hex()
	return nil
}

func (s *InmemStore) GetRound(r int64) (*RoundInfo, error) {
	res, ok := s.roundCache.Get(r)
	if !ok {
		return nil, cm.NewStoreErr("RoundCache", cm.KeyNotFound, strconv.FormatInt(r, 10))
	}
	return res.(*RoundInfo), nil
}

func (s *InmemStore) SetRound(r int64, round *RoundInfo) error {
	s.roundCache.Add(r, round)
	if r > s.lastRound {
		s.lastRound = r
	}
	return nil
}

func (s *InmemStore) LastRound() int64 {
	return s.lastRound
}

func (s *InmemStore) RoundWitnesses(r int64) []string {
	round, err := s.GetRound(r)
	if err != nil {
		return []string{}
	}
	return round.Witnesses()
}

func (s *InmemStore) RoundEvents(r int64) int {
	round, err := s.GetRound(r)
	if err != nil {
		return 0
	}
	return len(round.Message.CreatedEvents)
}

func (s *InmemStore) GetRoot(participant string) (*Root, error) {
	res, ok := s.rootsByParticipant[participant]
	if !ok {
		return nil, cm.NewStoreErr("RootCache", cm.KeyNotFound, participant)
	}
	return res, nil
}

func (s *InmemStore) GetBlock(index int64) (*Block, error) {
	res, ok := s.blockCache.Get(index)
	if !ok {
		return nil, cm.NewStoreErr("BlockCache", cm.KeyNotFound, strconv.FormatInt(index, 10))
	}
	return res.(*Block), nil
}

func (s *InmemStore) SetBlock(block *Block) error {
	index := block.Index()
	_, err := s.GetBlock(index)
	if err != nil && !cm.Is(err, cm.KeyNotFound) {
		return err
	}
	s.blockCache.Add(index, block)
	if index > s.lastBlock {
		s.lastBlock = index
	}
	return nil
}

func (s *InmemStore) LastBlockIndex() int64 {
	return s.lastBlock
}

func (s *InmemStore) GetFrame(index int64) (*Frame, error) {
	res, ok := s.frameCache.Get(index)
	if !ok {
		return nil, cm.NewStoreErr("FrameCache", cm.KeyNotFound, strconv.FormatInt(index, 10))
	}
	return res.(*Frame), nil
}

func (s *InmemStore) SetFrame(frame *Frame) error {
	index := frame.Round
	_, err := s.GetFrame(index)
	if err != nil && !cm.Is(err, cm.KeyNotFound) {
		return err
	}
	s.frameCache.Add(index, frame)
	return nil
}

func (s *InmemStore) Reset(frame *Frame) error {
	//Reset Root caches
	s.rootsByParticipant = frame.Roots
	rootsBySelfParent := make(map[string]*Root, len(frame.Roots))
	for _, r := range s.rootsByParticipant {
		rootsBySelfParent[r.SelfParent.Hash] = r
	}
	s.rootsBySelfParent = rootsBySelfParent

	//Reset Peer caches
	peerSet := peers.NewPeerSet(frame.Peers)
	s.peerSetCache.Set(frame.Round, peerSet)
	s.participantEventsCache = NewParticipantEventsCache(s.cacheSize, peerSet)

	//Reset Events and Rounds caches
	s.eventCache = cm.NewLRU(s.cacheSize, nil)
	s.roundCache = cm.NewLRU(s.cacheSize, nil)
	s.frameCache = cm.NewLRU(s.cacheSize, nil)
	s.consensusCache = cm.NewRollingIndex("ConsensusCache", s.cacheSize)

	s.lastRound = -1
	s.lastBlock = -1

	//Set Frame
	return s.SetFrame(frame)
}

func (s *InmemStore) Close() error {
	return nil
}

func (s *InmemStore) NeedBoostrap() bool {
	return false
}

func (s *InmemStore) StorePath() string {
	return ""
}
