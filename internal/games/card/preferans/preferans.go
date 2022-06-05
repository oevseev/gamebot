package preferans

import "github.com/oevseev/gamebot/internal/games/card"

type CardStatus int

const (
	NOT_DEALT       = CardStatus(iota)
	IN_TALON_FIRST  = CardStatus(iota)
	IN_TALON_SECOND = CardStatus(iota)
	IN_HAND         = CardStatus(iota)
	IN_TRICK        = CardStatus(iota)
	PLAYED          = CardStatus(iota)
	DISCARDED       = CardStatus(iota)
)

type Whist int

const (
	WHIST_PASS     = CardStatus(iota)
	WHIST          = CardStatus(iota)
	HALF_WHIST     = CardStatus(iota)
	WHIST_RETURNED = CardStatus(iota)
)

type State struct {
	LeadSuit      card.Suit
	CurrentPlayer int
	Owner         map[card.Card]int
	Status        map[card.Card]CardStatus
}

type PublicState struct {
	LeadSuit       card.Suit
	CurrentPlayer  int
	CardsRemaining map[int]int
	Owner          map[card.Card]int
	Status         map[card.Card]CardStatus
}

type Bid struct {
	Level  int
	Strain card.Suit
}

var Pass = Bid{
	Level:  0,
	Strain: card.SUIT_NONE,
}

type Bidding struct {
	Bids  []Bid
	Whist map[int]Whist
}

type Score struct {
	Heap  map[int]int
	Pool  map[int]int
	Whist map[int]map[int]int
}

type PreferansDeal struct {
	Dealer  int
	Bidding Bidding
	State   State
	Score   Score
}

type PreferansConfig struct {
	PoolLimit int
}

type PreferansGame struct {
	Config PreferansConfig
	Deals  []PreferansDeal
}
