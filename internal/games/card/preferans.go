package card

type CardStatus int

const (
	NOT_DEALT = CardStatus(iota)
	IN_TALON  = CardStatus(iota)
	IN_HAND   = CardStatus(iota)
	IN_TRICK  = CardStatus(iota)
	PLAYED    = CardStatus(iota)
	DISCARDED = CardStatus(iota)
)

type PreferansState struct {
	LeadSuit      Suit
	CurrentPlayer int
	Owner         map[Card]int
	Status        map[Card]CardStatus
}
