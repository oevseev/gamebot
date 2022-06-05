package preferans

import (
	"crypto/rand"
	"math/big"

	"github.com/oevseev/gamebot/internal/games/card"
)

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
	LeadSuit       card.Suit                `json:"leadSuit,omitempty"`
	CurrentPlayer  int                      `json:"currentPlayer,omitempty"`
	CardsRemaining map[int]int              `json:"cardsRemaining"`
	Owner          map[card.Card]int        `json:"owner"`
	Status         map[card.Card]CardStatus `json:"status"`
}

type Bid struct {
	Level  int       `json:"level"`
	Strain card.Suit `json:"strain"`
}

var Pass = Bid{
	Level:  0,
	Strain: card.SUIT_NONE,
}

type Bidding struct {
	Bids  []Bid         `json:"bids"`
	Whist map[int]Whist `json:"whist"`
}

type Score struct {
	Heap  map[int]int         `json:"heap"`
	Pool  map[int]int         `json:"pool"`
	Whist map[int]map[int]int `json:"whist"`
}

type Deal struct {
	Dealer  int
	Bidding Bidding
	Table   Table
	Score   Score
}

var preferansRanks = []card.Rank{
	card.SEVEN, card.EIGHT, card.NINE, card.TEN,
	card.JACK, card.QUEEN, card.KING, card.ACE,
}

func NewDeal(dealer int) Deal {
	cards := make([]card.Card, 0, 32)
	for _, rank := range preferansRanks {
		for suit := card.CLUBS; suit <= card.SPADES; suit++ {
			cards = append(cards, card.Card{
				Rank: rank,
				Suit: suit,
			})
		}
	}
	for i := 31; i > 0; i-- {
		jRaw, _ := rand.Int(rand.Reader, big.NewInt(int64(i)+1))
		j := int32(jRaw.Int64())
		cards[i], cards[j] = cards[j], cards[i]
	}

	owner := make(map[card.Card]int)
	status := make(map[card.Card]CardStatus)
	for i, card := range cards {
		if i < 30 {
			owner[card] = i / 10
			status[card] = IN_HAND
		} else if i == 30 {
			status[card] = IN_TALON_FIRST
		} else if i == 31 {
			status[card] = IN_TALON_SECOND
		}
	}

	return Deal{
		Dealer: dealer,
		Bidding: Bidding{
			Bids:  make([]Bid, 0),
			Whist: make(map[int]Whist),
		},
		Table: Table{
			Owner:  owner,
			Status: status,
		},
		Score: Score{
			Heap:  make(map[int]int),
			Pool:  make(map[int]int),
			Whist: make(map[int]map[int]int),
		},
	}
}

type PublicState struct {
	Dealer  int         `json:"dealer"`
	Bidding Bidding     `json:"bidding"`
	Table   PublicTable `json:"table"`
	Pool    []Score     `json:"pool"`
}

type Config struct {
	PoolLimit   int           `json:"poolLimit"`
	PlayerOrder map[int64]int `json:"playerOrder"`
}

type Game struct {
	Config Config
	deals  []Deal
}

func NewGame(creator int64) *Game {
	dealerIndexRaw, _ := rand.Int(rand.Reader, big.NewInt(3))
	dealerIndex := int(dealerIndexRaw.Int64()) + 1

	return &Game{
		Config: Config{
			PoolLimit: 10,
			PlayerOrder: map[int64]int{
				creator: 1,
			},
		},
		deals: []Deal{
			NewDeal(dealerIndex),
		},
	}
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
