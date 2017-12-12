package network

import (
	"sort"
	"strings"
)

// Stuff for implementing the Stellar Consensus Protocol. See:
// https://www.stellar.org/papers/stellar-consensus-protocol.pdf
// When there are frustrating single-letter variable names, it's because we are
// making the names line up with the protocol paper.

// For now each block just has a list of comments.
// This isn't supposed to be useful, it's just for testing.
type SlotValue struct {
	Comments []string
}

func MakeSlotValue(comment string) SlotValue {
	return SlotValue{Comments: []string{comment}}
}

func Combine(a SlotValue, b SlotValue) SlotValue {
	joined := append(a.Comments, b.Comments...)
	sort.Strings(joined)
	answer := []string{}
	for _, item := range joined {
		if len(answer) == 0 || answer[len(answer)-1] != item {
			answer = append(answer, item)
		}
	}
	return SlotValue{Comments: answer}
}

func HasSlotValue(list []SlotValue, v SlotValue) bool {
	k := strings.Join(v.Comments, ",")
	for _, s := range list {
		if strings.Join(s.Comments, ",") == k {
			return true
		}
	}
	return false
}

type QuorumSlice struct {
	// Members is a list of public keys for nodes that occur in the quorum slice.
	// Typically includes ourselves.
	Members []string

	// The number of members we require for consensus, including ourselves.
	// The protocol can support other sorts of slices, like weighted or any wacky
	// thing, but for now we only do this simple "any k out of these n" voting.
	Threshold int
}

type NominateMessage struct {
	// What slot we are nominating values for
	I int

	// The values we have voted to nominate
	X []SlotValue

	// The values we have accepted as nominated
	Y []SlotValue
	
	D QuorumSlice
}

func (m *NominateMessage) MessageType() string {
	return "Nominate"
}

// See page 21 of the protocol paper for more detail here.
type NominationState struct {
	// The values we have voted to nominate
	X []SlotValue

	// The values we have accepted as nominated
	Y []SlotValue

	// The values that we consider to be candidates 
	Z []SlotValue

	// The last NominateMessage received from each node
	N map[string]*NominateMessage
}

func NewNominationState() *NominationState {
	return &NominationState{
		X: make([]SlotValue, 0),
		Y: make([]SlotValue, 0),
		Z: make([]SlotValue, 0),
		N: make(map[string]*NominateMessage),
	}
}

// HasNomination tells you whether this nomination state can currently send out
// a nominate message.
// If we have never received a nomination from a peer, and haven't had SetDefault
// called ourselves, then we won't have a nomination.
func (s *NominationState) HasNomination() bool {
	return len(s.X) > 0
}

func (s *NominationState) SetDefault(v SlotValue) {
	if s.HasNomination() {
		// We already have something to nominate
		return
	}
	s.X = []SlotValue{v}
}

// Handles an incoming nomination message from a peer
func (s *NominationState) Handle(node string, m *NominateMessage) {
	// TODO
}

// Ballot phases
type Phase int
const (
	Prepare Phase = iota
	Confirm
	Externalize
)

type Ballot struct {
	// An increasing counter, n >= 1, to ensure we can always have more ballots
	n int

	// The value this ballot proposes
	x SlotValue
}

// See page 23 of the protocol paper for more detail here.
type BallotState struct {
	// The current ballot we are trying to prepare and commit.
	b *Ballot

	// The highest two ballots that are accepted as prepared.
	// p is the highest, pPrime the next.
	// It's nil if there is no such ballot.
	p *Ballot
	pPrime *Ballot

	// In the Prepare phase, c is the lowest and h is the highest ballot
	// for which we have voted to commit but not accepted abort.
	// In the Confirm phase, c is the lowest and h is the highest ballot
	// for which we accepted commit.
	// In the Externalize phase, c is the lowest and h is the highest ballot
	// for which we confirmed commit.
	// If c is not nil, then c <= h <= b.
	c *Ballot
	h *Ballot

	// The value to use in the next ballot
	z SlotValue
	
	// The latest PrepareMessage, ConfirmMessage, or ExternalizeMessage from each peer
	M map[string]Message
}

// PrepareMessage is the first phase of the three-phase ballot protocol
type PrepareMessage struct {
	// What slot we are preparing ballots for
	I int

	// The current ballot we are trying to prepare
	Bn int
	Bx SlotValue

	// The contents of state.p
	Pn int
	Px SlotValue

	// The contents of state.pPrime
	Ppn int
	Ppx SlotValue

	// Ballot numbers for c and h
	Cn int
	Hn int

	D QuorumSlice
}

func (m *PrepareMessage) MessageType() string {
	return "Prepare"
}

// ConfirmMessage is the second phase of the three-phase ballot protocol
type ConfirmMessage struct {
	// What slot we are confirming ballots for
	I int

	// The current ballot we are trying to confirm
	Bn int
	Bx SlotValue

	// state.p.n
	Pn int

	// state.c.n
	Cn int

	// state.h.n
	Hn int

	D QuorumSlice
}

func (m *ConfirmMessage) MessageType() string {
	return "Confirm"
}

// ExternalizeMessage is the third phase of the three-phase ballot protocol
// Sent after we have confirmed a commit
type ExternalizeMessage struct {
	// What slot we are externalizing
	I int

	// The value at this slot
	X SlotValue

	// state.c.n
	Cn int

	// state.h.n
	Hn int

	D QuorumSlice
}

func (m *ExternalizeMessage) MessageType() string {
	return "Externalize"
}

type StateBuilder struct {
	// Which slot is actively being built
	slot int

	// Values for past slots that have already achieved consensus
	values map[int]SlotValue

	nState NominationState
	bState BallotState
}

func NewStateBuilder() *StateBuilder {
	return &StateBuilder{
		slot: 1,
		values: make(map[int]SlotValue),
		nState: NewNominationState(),
	}
}

