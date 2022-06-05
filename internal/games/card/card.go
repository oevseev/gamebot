package card

type Rank int

const (
	RANK_NONE = Rank(iota)
	ACE       = Rank(iota)
	TWO       = Rank(iota)
	THREE     = Rank(iota)
	FOUR      = Rank(iota)
	FIVE      = Rank(iota)
	SIX       = Rank(iota)
	SEVEN     = Rank(iota)
	EIGHT     = Rank(iota)
	NINE      = Rank(iota)
	TEN       = Rank(iota)
	JACK      = Rank(iota)
	QUEEN     = Rank(iota)
	KING      = Rank(iota)
)

type Suit int

const (
	SUIT_NONE = Suit(iota)
	CLUBS     = Suit(iota)
	DIAMONDS  = Suit(iota)
	HEARTS    = Suit(iota)
	SPADES    = Suit(iota)
)

type Card struct {
	Rank Rank
	Suit Suit
}
