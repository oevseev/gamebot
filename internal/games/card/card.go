package card

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

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

func (c Card) MarshalText() (text []byte, err error) {
	return []byte(fmt.Sprintf("%d,%d", c.Rank, c.Suit)), nil
}

func (c *Card) UnmarshalText(text []byte) error {
	if !utf8.Valid(text) {
		return errors.New("text is not valid UTF-8")
	}
	s := string(text)

	substrings := strings.SplitN(s, ",", 2)
	rank, err := strconv.ParseInt(substrings[0], 10, 32)
	if err != nil {
		return err
	}
	suit, err := strconv.ParseInt(substrings[0], 10, 32)
	if err != nil {
		return err
	}

	c.Rank = Rank(rank)
	c.Suit = Suit(suit)

	return nil
}
