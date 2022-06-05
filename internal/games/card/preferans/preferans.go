package preferans

import "github.com/oevseev/gamebot/internal/games/card"

type CardStatus int

const (
	STATUS_UNKNOWN  = CardStatus(iota)
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

type Table struct {
	LeadSuit      card.Suit
	CurrentPlayer int
	Owner         map[card.Card]int
	Status        map[card.Card]CardStatus
}

type PublicTable struct {
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

type Deal struct {
	Dealer  int
	Bidding Bidding
	Table   Table
	Score   Score
}

type PublicState struct {
	Dealer  int
	Bidding Bidding
	Table   PublicTable
	Pool    []Score
}

type Config struct {
	PoolLimit   int
	PlayerOrder map[int64]int
}

type Game struct {
	Config Config
	deals  []Deal
}

func (g *Game) GetPublicState(playerId int64) PublicState {
	currentDeal := g.deals[len(g.deals)-1]
	selectedPlayerIdx := g.Config.PlayerOrder[playerId]

	pool := make([]Score, 0, len(g.deals))
	for _, deal := range g.deals {
		pool = append(pool, deal.Score)
	}

	cardsRemaining := make(map[int]int)
	for _, ownerIdx := range currentDeal.Table.Owner {
		cardsRemaining[ownerIdx] += 1
	}

	owner := make(map[card.Card]int)
	status := make(map[card.Card]CardStatus)
	for card, ownerIdx := range currentDeal.Table.Owner {
		cardStatus := currentDeal.Table.Status[card]
		if ownerIdx != selectedPlayerIdx || cardStatus != IN_HAND {
			continue
		}
		owner[card] = ownerIdx
		status[card] = cardStatus
	}

	table := PublicTable{
		LeadSuit:       currentDeal.Table.LeadSuit,
		CurrentPlayer:  currentDeal.Table.CurrentPlayer,
		CardsRemaining: cardsRemaining,
		Owner:          owner,
		Status:         status,
	}

	return PublicState{
		Dealer:  currentDeal.Dealer,
		Bidding: currentDeal.Bidding,
		Table:   table,
		Pool:    pool,
	}
}
