package election

import (
	"fmt"
	"math/big"

	"github.com/Fantom-foundation/go-lachesis/src/hash"
)

// calculate SfWitness votes only for the new root.
// If this root sees that the current election is decided, then @return decided SfWitness
func (el *Election) ProcessRoot(newRoot hash.Event, newRootSlot RootSlot) (*ElectionRes, error) {
	if len(el.decidedRoots) == len(el.nodes) {
		// current election is already decided
		return el.chooseSfWitness()
	}

	if newRootSlot.Frame <= el.frameToDecide {
		// too old root, out of interest for current election
		return nil, nil
	}
	round := newRootSlot.Frame - el.frameToDecide

	notDecidedRoots := el.notDecidedRoots()
	for _, nodeIdSubject := range notDecidedRoots {
		slotSubject := RootSlot{
			Frame:  el.frameToDecide,
			Nodeid: nodeIdSubject,
		}
		vote := voteValue{}

		if round == 1 {
			// in initial round, vote "yes" if subject is strongly seen
			seenRoot := el.stronglySee(newRoot, slotSubject)
			vote.yes = seenRoot != nil
			vote.decided = false
			if seenRoot != nil {
				vote.seenRoot = *seenRoot
			}
		} else if round > 1 {
			seenRoots := el.stronglySeenRoots(newRoot, newRootSlot.Frame-1)

			yesVotes := new(big.Int)
			noVotes := new(big.Int)

			// calc number of "yes" and "no", weighted by node's stake
			var subjectHash *hash.Event
			for _, seenRoot := range seenRoots {
				vid := voteId{
					forNodeid: nodeIdSubject,
					fromRoot:  seenRoot.root,
				}

				if vote, ok := el.votes[vid]; ok {
					if vote.yes && subjectHash != nil && *subjectHash != vote.seenRoot {
						msg := "2 fork roots are strongly seen => more than 1/3n are Byzantine (%s != %s, election frame=%d, nodeid=%s)"
						return nil, fmt.Errorf(msg, subjectHash.String(), vote.seenRoot.String(), el.frameToDecide, nodeIdSubject.String())
					}

					if vote.yes {
						subjectHash = &vote.seenRoot
						yesVotes = yesVotes.Add(yesVotes, seenRoot.stakeAmount)
					} else {
						noVotes = noVotes.Add(noVotes, seenRoot.stakeAmount)
					}
				} else {
					el.Fatal("Every root must vote for every not decided subject. Possibly roots are processed out of order, root=", newRoot.String())
				}
			}
			// sanity checks
			if new(big.Int).Add(yesVotes, noVotes).Cmp(el.superMajority) < 0 {
				el.Fatal("Root must see at least 2/3n of prev roots. Possibly roots are processed out of order, root=", newRoot.String())
			}
			if new(big.Int).Add(yesVotes, noVotes).Cmp(el.totalStake) > 0 {
				el.Fatal("Root cannot see more than 100% of prev roots, root=", newRoot.String())
			}

			// vote as majority of votes
			vote.yes = yesVotes.Cmp(noVotes) >= 0
			if vote.yes && subjectHash != nil {
				vote.seenRoot = *subjectHash
			}

			// If supermajority is seen, then the final decision may be made.
			// It's guaranteed to be final and consistent unless more than 1/3n are Byzantine.
			vote.decided = yesVotes.Cmp(el.superMajority) >= 0 || noVotes.Cmp(el.superMajority) >= 0
			if vote.decided {
				el.decidedRoots[nodeIdSubject] = vote
			}
		}
		// save vote for next rounds
		vid := voteId{
			fromRoot:  newRoot,
			forNodeid: slotSubject.Nodeid,
		}
		el.votes[vid] = vote
	}

	frameDecided := len(el.decidedRoots) == len(el.nodes)
	if frameDecided {
		return el.chooseSfWitness()
	}
	return nil, nil
}
