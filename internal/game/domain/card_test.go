package domain

import (
	"testing"
)

func TestCard_NewCard(t *testing.T) {
	card := NewCard(Hearts, King, 1)
	
	if card.Suit != Hearts {
		t.Errorf("Expected suit Hearts, got %v", card.Suit)
	}
	if card.Rank != King {
		t.Errorf("Expected rank King, got %v", card.Rank)
	}
	if card.DeckID != 1 {
		t.Errorf("Expected deck ID 1, got %d", card.DeckID)
	}
	if card.IsJoker {
		t.Error("Expected non-joker card")
	}
}

func TestCard_NewJoker(t *testing.T) {
	joker := NewJoker(BigJoker, 2)
	
	if !joker.IsJoker {
		t.Error("Expected joker card")
	}
	if joker.JokerType != BigJoker {
		t.Errorf("Expected BigJoker, got %v", joker.JokerType)
	}
	if joker.DeckID != 2 {
		t.Errorf("Expected deck ID 2, got %d", joker.DeckID)
	}
}

func TestCard_GetPointValue(t *testing.T) {
	tests := []struct {
		name     string
		card     Card
		expected int
	}{
		{"King has 10 points", NewCard(Spades, King, 1), 10},
		{"Ten has 10 points", NewCard(Hearts, Ten, 1), 10},
		{"Five has 5 points", NewCard(Clubs, Five, 1), 5},
		{"Ace has 0 points", NewCard(Diamonds, Ace, 1), 0},
		{"Two has 0 points", NewCard(Spades, Two, 1), 0},
		{"Big Joker has 0 points", NewJoker(BigJoker, 1), 0},
		{"Small Joker has 0 points", NewJoker(SmallJoker, 1), 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.card.GetPointValue(); got != tt.expected {
				t.Errorf("GetPointValue() = %d, want %d", got, tt.expected)
			}
		})
	}
}

func TestCard_IsEqual(t *testing.T) {
	card1 := NewCard(Hearts, King, 1)
	card2 := NewCard(Hearts, King, 1)
	card3 := NewCard(Hearts, King, 2) // Different deck
	card4 := NewCard(Spades, King, 1) // Different suit
	joker1 := NewJoker(BigJoker, 1)
	joker2 := NewJoker(BigJoker, 1)
	joker3 := NewJoker(SmallJoker, 1)

	tests := []struct {
		name     string
		card1    Card
		card2    Card
		expected bool
	}{
		{"Same cards are equal", card1, card2, true},
		{"Different deck IDs are not equal", card1, card3, false},
		{"Different suits are not equal", card1, card4, false},
		{"Same jokers are equal", joker1, joker2, true},
		{"Different joker types are not equal", joker1, joker3, false},
		{"Joker and regular card are not equal", joker1, card1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.card1.IsEqual(tt.card2); got != tt.expected {
				t.Errorf("IsEqual() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCard_IsSameFace(t *testing.T) {
	card1 := NewCard(Hearts, King, 1)
	card2 := NewCard(Hearts, King, 2) // Same face, different deck
	card3 := NewCard(Spades, King, 1) // Different suit
	joker1 := NewJoker(BigJoker, 1)
	joker2 := NewJoker(BigJoker, 2)

	tests := []struct {
		name     string
		card1    Card
		card2    Card
		expected bool
	}{
		{"Same face different deck", card1, card2, true},
		{"Different suits same rank", card1, card3, false},
		{"Same joker type different deck", joker1, joker2, true},
		{"Joker and regular card", joker1, card1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.card1.IsSameFace(tt.card2); got != tt.expected {
				t.Errorf("IsSameFace() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCard_GetTrumpHierarchy(t *testing.T) {
	trumpSuit := Hearts
	bigJoker := NewJoker(BigJoker, 1)
	smallJoker := NewJoker(SmallJoker, 1)
	trumpTwo := NewCard(Hearts, Two, 1)
	offSuitTwo := NewCard(Spades, Two, 1)
	trumpKing := NewCard(Hearts, King, 1)
	offSuitKing := NewCard(Spades, King, 1)

	tests := []struct {
		name     string
		card     Card
		expected int
	}{
		{"Big Joker highest", bigJoker, 1000},
		{"Small Joker second", smallJoker, 999},
		{"Trump Two third", trumpTwo, 998},
		{"Off-suit Two fourth", offSuitTwo, 997},
		{"Trump King has trump hierarchy", trumpKing, 900 + int(King)},
		{"Off-suit King has no trump hierarchy", offSuitKing, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.card.GetTrumpHierarchy(trumpSuit); got != tt.expected {
				t.Errorf("GetTrumpHierarchy() = %d, want %d", got, tt.expected)
			}
		})
	}
}

func TestDeck_NewDeck(t *testing.T) {
	deck := NewDeck()
	
	if len(deck.Cards) != 108 {
		t.Errorf("Expected 108 cards, got %d", len(deck.Cards))
	}

	// Validate deck composition
	if err := deck.ValidateDeckComposition(); err != nil {
		t.Errorf("Deck composition validation failed: %v", err)
	}

	// Check total points
	totalPoints := deck.GetTotalPoints()
	expectedPoints := 200 // 8 Kings * 10 + 8 Tens * 10 + 8 Fives * 5
	if totalPoints != expectedPoints {
		t.Errorf("Expected total points %d, got %d", expectedPoints, totalPoints)
	}
}

func TestDeck_Deal(t *testing.T) {
	deck := NewDeck()
	initialCount := len(deck.Cards)

	// Deal 5 cards
	dealt, err := deck.Deal(5)
	if err != nil {
		t.Errorf("Deal failed: %v", err)
	}

	if len(dealt) != 5 {
		t.Errorf("Expected 5 dealt cards, got %d", len(dealt))
	}

	if len(deck.Cards) != initialCount-5 {
		t.Errorf("Expected %d remaining cards, got %d", initialCount-5, len(deck.Cards))
	}

	// Try to deal more cards than available
	_, err = deck.Deal(200)
	if err == nil {
		t.Error("Expected error when dealing more cards than available")
	}
}

func TestDeck_ValidateDeckComposition(t *testing.T) {
	// Test valid deck
	validDeck := NewDeck()
	if err := validDeck.ValidateDeckComposition(); err != nil {
		t.Errorf("Valid deck failed validation: %v", err)
	}

	// Test invalid deck - wrong number of cards
	invalidDeck := &Deck{Cards: make([]Card, 100)}
	if err := invalidDeck.ValidateDeckComposition(); err == nil {
		t.Error("Expected error for deck with wrong number of cards")
	}

	// Test deck with wrong card counts
	wrongCountDeck := &Deck{Cards: make([]Card, 108)}
	// Fill with all Aces of Spades
	for i := 0; i < 108; i++ {
		wrongCountDeck.Cards[i] = NewCard(Spades, Ace, 1)
	}
	if err := wrongCountDeck.ValidateDeckComposition(); err == nil {
		t.Error("Expected error for deck with wrong card distribution")
	}
}

func TestRank_String(t *testing.T) {
	tests := []struct {
		rank     Rank
		expected string
	}{
		{Two, "2"},
		{Three, "3"},
		{Ten, "10"},
		{Jack, "J"},
		{Queen, "Q"},
		{King, "K"},
		{Ace, "A"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.rank.String(); got != tt.expected {
				t.Errorf("Rank.String() = %s, want %s", got, tt.expected)
			}
		})
	}
}

func TestSuit_String(t *testing.T) {
	tests := []struct {
		suit     Suit
		expected string
	}{
		{Spades, "Spades"},
		{Hearts, "Hearts"},
		{Clubs, "Clubs"},
		{Diamonds, "Diamonds"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.suit.String(); got != tt.expected {
				t.Errorf("Suit.String() = %s, want %s", got, tt.expected)
			}
		})
	}
}

func TestJokerType_String(t *testing.T) {
	tests := []struct {
		jokerType JokerType
		expected  string
	}{
		{BigJoker, "Big Joker"},
		{SmallJoker, "Small Joker"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.jokerType.String(); got != tt.expected {
				t.Errorf("JokerType.String() = %s, want %s", got, tt.expected)
			}
		})
	}
}