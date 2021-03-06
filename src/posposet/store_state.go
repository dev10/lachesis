package posposet

import (
	"github.com/Fantom-foundation/go-lachesis/src/hash"
	"github.com/Fantom-foundation/go-lachesis/src/posposet/wire"
	"github.com/Fantom-foundation/go-lachesis/src/state"
)

// StateDB returns state database.
func (s *Store) StateDB(from hash.Hash) *state.DB {
	db, err := state.New(from, s.table.Balances)
	if err != nil {
		s.Fatal(err)
	}
	return db
}

// SetState stores state.
// State is seldom read; so no cache.
func (s *Store) SetState(st *State) {
	const key = "current"
	s.set(s.table.States, []byte(key), st.ToWire())

}

// GetState returns stored state.
// State is seldom read; so no cache.
func (s *Store) GetState() *State {
	const key = "current"
	w, _ := s.get(s.table.States, []byte(key), &wire.State{}).(*wire.State)
	return WireToState(w)
}
