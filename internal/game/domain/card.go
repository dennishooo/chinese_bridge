package domain

import (
	"fmt"
)

// Suit represents the four card suits plus trump indicators
type Suit int

const (
	Spades Suit = iota
	Hearts
	Clubs
	Diamonds
)

func (s Suit) String() string {
	switch s {
	case Spades:
		return "Spades"
	case Hearts:
		return "Hearts"
	case Clubs:
		return "Clubs"
	case Diamonds:
		return "Diamonds"
	default:
		return "Unknown"
	}
}

// Rank represents card ranks from 2 to Ace
type Rank int

const (
	Two Rank = iota + 2
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace
)

func (r Rank) String() string {
	switch r {
	case Two:
		return "2"
	case Three:
		return "3"
	case Four:
		return "4"
	case Five:
		return "5"
	case Six:
		return "6"
	case Seven:
		return "7"
	case Eight:
		return "8"
	case Nine:
		return "9"
	case Ten:
		return "10"
	case Jack:
		return "J"
	case Queen:
		return "Q"
	case King:
		return "K"
	case Ace:
		return "A"
	default:
		return "Unknown"
	}
}

// GetPointValue returns the point value for scoring cards
func (r Rank) GetPointValue() int {
	switch r {
	case Five:
		return 5
	case Ten:
		return 10
	case King:
		return 10
	default:
		return 0
	}
}

// JokerType represents the two types of jokers
type JokerType int

const (
	BigJoker JokerType = iota
	SmallJoker
)

func (j JokerType) String() string {
	switch j {
	case BigJoker:
		return "Big Joker"
	case SmallJoker:
		return "Small Joker"
	default:
		return "Unknown"
	}
}

// Card represents a single playing card with suit, rank, deck ID, and joker type
type Card struct {
	Suit      Suit      `json:"suit"`
	Rank      Rank      `json:"rank"`
	DeckID    int       `json:"deck_id"`    // 1 or 2 for duplicate cards
	IsJoker   bool      `json:"is_joker"`
	JokerType JokerType `json:"joker_type,omitempty"`
}

// NewCard creates a new regular card
func NewCard(suit Suit, rank Rank, deckID int) Card {
	return Card{
		Suit:    suit,
		Rank:    rank,
		DeckID:  deckID,
		IsJoker: false,
	}
}

// NewJoker creates a new joker card
func NewJoker(jokerType JokerType, deckID int) Card {
	return Card{
		DeckID:    deckID,
		IsJoker:   true,
		JokerType: jokerType,
	}
}

// String returns a string representation of the card
func (c Card) String() string {
	if c.IsJoker {
		return fmt.Sprintf("%s (Deck %d)", c.JokerType.String(), c.DeckID)
	}
	return fmt.Sprintf("%s of %s (Deck %d)", c.Rank.String(), c.Suit.String(), c.DeckID)
}

// GetPointValue returns the point value of the card for scoring
func (c Card) GetPointValue() int {
	if c.IsJoker {
		return 0
	}
	return c.Rank.GetPointValue()
}

// IsEqual checks if two cards are identical (same suit, rank, and deck)
func (c Card) IsEqual(other Card) bool {
	if c.IsJoker && other.IsJoker {
		return c.JokerType == other.JokerType && c.DeckID == other.DeckID
	}
	if c.IsJoker || other.IsJoker {
		return false
	}
	return c.Suit == other.Suit && c.Rank == other.Rank && c.DeckID == other.DeckID
}

// IsSameFace checks if two cards have the same face value (ignoring deck ID)
func (c Card) IsSameFace(other Card) bool {
	if c.IsJoker && other.IsJoker {
		return c.JokerType == other.JokerType
	}
	if c.IsJoker || other.IsJoker {
		return false
	}
	return c.Suit == other.Suit && c.Rank == other.Rank
}

// GetTrumpHierarchy returns the trump hierarchy value for card comparison
// Higher values beat lower values in trump hierarchy
func (c Card) GetTrumpHierarchy(trumpSuit Suit) int {
	if c.IsJoker {
		switch c.JokerType {
		case BigJoker:
			return 1000 // Highest trump
		case SmallJoker:
			return 999 // Second highest trump
		}
	}

	// Trump 2s are permanent trumps
	if c.Rank == Two {
		if c.Suit == trumpSuit {
			return 998 // Trump suit 2s
		}
		return 997 // Off-suit 2s
	}

	// Trump suit cards (excluding 2s)
	if c.Suit == trumpSuit {
		return 900 + int(c.Rank) // Trump suit cards ranked by rank
	}

	// Off-suit cards have no trump hierarchy
	return 0
}

// GetSuitHierarchy returns the hierarchy value within a suit
func (c Card) GetSuitHierarchy() int {
	if c.IsJoker {
		return 0 // Jokers don't have suit hierarchy
	}
	return int(c.Rank)
}

// Deck represents a complete deck of 108 cards for Chinese Bridge
type Deck struct {
	Cards []Card `json:"cards"`
}

// NewDeck creates a new deck with 2 standard 52-card decks plus 4 jokers
func NewDeck() *Deck {
	deck := &Deck{
		Cards: make([]Card, 0, 108),
	}

	// Add two standard 52-card decks
	for deckID := 1; deckID <= 2; deckID++ {
		// Add regular cards for each suit
		for suit := Spades; suit <= Diamonds; suit++ {
			for rank := Two; rank <= Ace; rank++ {
				deck.Cards = append(deck.Cards, NewCard(suit, rank, deckID))
			}
		}
	}

	// Add 4 jokers (2 big, 2 small)
	for deckID := 1; deckID <= 2; deckID++ {
		deck.Cards = append(deck.Cards, NewJoker(BigJoker, deckID))
		deck.Cards = append(deck.Cards, NewJoker(SmallJoker, deckID))
	}

	return deck
}

// Shuffle randomizes the order of cards in the deck
func (d *Deck) Shuffle() {
	// Implementation would use crypto/rand for secure shuffling
	// This is a placeholder for the actual shuffle algorithm
}

// Deal removes and returns the specified number of cards from the top of the deck
func (d *Deck) Deal(count int) ([]Card, error) {
	if count > len(d.Cards) {
		return nil, fmt.Errorf("cannot deal %d cards, only %d cards remaining", count, len(d.Cards))
	}

	dealt := make([]Card, count)
	copy(dealt, d.Cards[:count])
	d.Cards = d.Cards[count:]

	return dealt, nil
}

// Remaining returns the number of cards left in the deck
func (d *Deck) Remaining() int {
	return len(d.Cards)
}

// ValidateDeckComposition ensures the deck has the correct composition
func (d *Deck) ValidateDeckComposition() error {
	if len(d.Cards) != 108 {
		return fmt.Errorf("deck must have exactly 108 cards, found %d", len(d.Cards))
	}

	// Count cards by type
	suitCounts := make(map[Suit]map[Rank]int)
	jokerCounts := make(map[JokerType]int)

	for suit := Spades; suit <= Diamonds; suit++ {
		suitCounts[suit] = make(map[Rank]int)
	}

	for _, card := range d.Cards {
		if card.IsJoker {
			jokerCounts[card.JokerType]++
		} else {
			suitCounts[card.Suit][card.Rank]++
		}
	}

	// Validate each rank appears exactly twice in each suit
	for suit := Spades; suit <= Diamonds; suit++ {
		for rank := Two; rank <= Ace; rank++ {
			if suitCounts[suit][rank] != 2 {
				return fmt.Errorf("rank %s of %s must appear exactly twice, found %d", 
					rank.String(), suit.String(), suitCounts[suit][rank])
			}
		}
	}

	// Validate jokers
	if jokerCounts[BigJoker] != 2 {
		return fmt.Errorf("must have exactly 2 big jokers, found %d", jokerCounts[BigJoker])
	}
	if jokerCounts[SmallJoker] != 2 {
		return fmt.Errorf("must have exactly 2 small jokers, found %d", jokerCounts[SmallJoker])
	}

	return nil
}

// GetTotalPoints calculates the total point value of all cards in the deck
func (d *Deck) GetTotalPoints() int {
	total := 0
	for _, card := range d.Cards {
		total += card.GetPointValue()
	}
	return total
}