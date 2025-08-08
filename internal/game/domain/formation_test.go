package domain

import (
	"testing"
)

func TestFormation_NewSingle(t *testing.T) {
	card := NewCard(Hearts, King, 1)
	formation := NewSingle(card)

	if formation.Type != Single {
		t.Errorf("Expected Single formation, got %v", formation.Type)
	}
	if len(formation.Cards) != 1 {
		t.Errorf("Expected 1 card, got %d", len(formation.Cards))
	}
	if !formation.Cards[0].IsEqual(card) {
		t.Error("Formation card doesn't match input card")
	}
	if formation.Suit != Hearts {
		t.Errorf("Expected Hearts suit, got %v", formation.Suit)
	}
}

func TestFormation_NewSingleJoker(t *testing.T) {
	joker := NewJoker(BigJoker, 1)
	formation := NewSingle(joker)

	if formation.Type != Single {
		t.Errorf("Expected Single formation, got %v", formation.Type)
	}
	if formation.Suit != Spades {
		t.Errorf("Expected default Spades suit for joker, got %v", formation.Suit)
	}
}

func TestFormation_NewPair(t *testing.T) {
	card1 := NewCard(Hearts, King, 1)
	card2 := NewCard(Hearts, King, 2)
	
	formation, err := NewPair(card1, card2)
	if err != nil {
		t.Errorf("NewPair failed: %v", err)
	}

	if formation.Type != Pair {
		t.Errorf("Expected Pair formation, got %v", formation.Type)
	}
	if len(formation.Cards) != 2 {
		t.Errorf("Expected 2 cards, got %d", len(formation.Cards))
	}
	if formation.Suit != Hearts {
		t.Errorf("Expected Hearts suit, got %v", formation.Suit)
	}
}

func TestFormation_NewPairInvalid(t *testing.T) {
	card1 := NewCard(Hearts, King, 1)
	card2 := NewCard(Hearts, Queen, 1) // Different rank
	
	_, err := NewPair(card1, card2)
	if err == nil {
		t.Error("Expected error for non-matching cards")
	}
}

func TestFormation_NewTractor(t *testing.T) {
	// Create consecutive pairs: King-King, Ace-Ace
	kingPair := []Card{
		NewCard(Hearts, King, 1),
		NewCard(Hearts, King, 2),
	}
	acePair := []Card{
		NewCard(Hearts, Ace, 1),
		NewCard(Hearts, Ace, 2),
	}
	
	pairs := [][]Card{kingPair, acePair}
	formation, err := NewTractor(pairs, Spades)
	
	if err != nil {
		t.Errorf("NewTractor failed: %v", err)
	}
	if formation.Type != Tractor {
		t.Errorf("Expected Tractor formation, got %v", formation.Type)
	}
	if len(formation.Cards) != 4 {
		t.Errorf("Expected 4 cards, got %d", len(formation.Cards))
	}
	if formation.Suit != Hearts {
		t.Errorf("Expected Hearts suit, got %v", formation.Suit)
	}
}

func TestFormation_NewTractorInvalid(t *testing.T) {
	tests := []struct {
		name  string
		pairs [][]Card
	}{
		{
			name: "Only one pair",
			pairs: [][]Card{
				{NewCard(Hearts, King, 1), NewCard(Hearts, King, 2)},
			},
		},
		{
			name: "Non-consecutive ranks",
			pairs: [][]Card{
				{NewCard(Hearts, King, 1), NewCard(Hearts, King, 2)},
				{NewCard(Hearts, Nine, 1), NewCard(Hearts, Nine, 2)}, // Skip Queen, Jack, Ten
			},
		},
		{
			name: "Different suits",
			pairs: [][]Card{
				{NewCard(Hearts, King, 1), NewCard(Hearts, King, 2)},
				{NewCard(Spades, Ace, 1), NewCard(Spades, Ace, 2)},
			},
		},
		{
			name: "Contains jokers",
			pairs: [][]Card{
				{NewJoker(BigJoker, 1), NewJoker(BigJoker, 2)},
				{NewCard(Hearts, Ace, 1), NewCard(Hearts, Ace, 2)},
			},
		},
		{
			name: "Contains 2s",
			pairs: [][]Card{
				{NewCard(Hearts, Two, 1), NewCard(Hearts, Two, 2)},
				{NewCard(Hearts, Three, 1), NewCard(Hearts, Three, 2)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewTractor(tt.pairs, Spades)
			if err == nil {
				t.Errorf("Expected error for %s", tt.name)
			}
		})
	}
}

func TestFormation_IsValid(t *testing.T) {
	tests := []struct {
		name      string
		formation *Formation
		wantError bool
	}{
		{
			name:      "Valid single",
			formation: NewSingle(NewCard(Hearts, King, 1)),
			wantError: false,
		},
		{
			name: "Valid pair",
			formation: &Formation{
				Type:  Pair,
				Cards: []Card{NewCard(Hearts, King, 1), NewCard(Hearts, King, 2)},
				Suit:  Hearts,
			},
			wantError: false,
		},
		{
			name: "Invalid single - too many cards",
			formation: &Formation{
				Type:  Single,
				Cards: []Card{NewCard(Hearts, King, 1), NewCard(Hearts, Queen, 1)},
				Suit:  Hearts,
			},
			wantError: true,
		},
		{
			name: "Invalid pair - wrong number of cards",
			formation: &Formation{
				Type:  Pair,
				Cards: []Card{NewCard(Hearts, King, 1)},
				Suit:  Hearts,
			},
			wantError: true,
		},
		{
			name: "Invalid pair - different faces",
			formation: &Formation{
				Type:  Pair,
				Cards: []Card{NewCard(Hearts, King, 1), NewCard(Hearts, Queen, 1)},
				Suit:  Hearts,
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.formation.IsValid()
			if (err != nil) != tt.wantError {
				t.Errorf("IsValid() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestFormation_GetHighestCard(t *testing.T) {
	trumpSuit := Hearts
	
	// Test with trump cards
	trumpKing := NewCard(Hearts, King, 1)
	trumpAce := NewCard(Hearts, Ace, 1)
	formation := &Formation{
		Type:  Pair,
		Cards: []Card{trumpKing, trumpAce},
		Suit:  Hearts,
	}
	
	highest := formation.GetHighestCard(trumpSuit)
	if !highest.IsEqual(trumpAce) {
		t.Errorf("Expected Ace to be highest, got %s", highest.String())
	}
}

func TestFormation_GetPointValue(t *testing.T) {
	// Formation with King (10 points) and Five (5 points)
	formation := &Formation{
		Type:  Pair,
		Cards: []Card{NewCard(Hearts, King, 1), NewCard(Hearts, Five, 1)},
		Suit:  Hearts,
	}
	
	points := formation.GetPointValue()
	expected := 15 // 10 + 5
	if points != expected {
		t.Errorf("Expected %d points, got %d", expected, points)
	}
}

func TestFormation_IsTrump(t *testing.T) {
	trumpSuit := Hearts
	
	tests := []struct {
		name      string
		formation *Formation
		expected  bool
	}{
		{
			name: "Trump suit formation",
			formation: &Formation{
				Cards: []Card{NewCard(Hearts, King, 1)},
				Suit:  Hearts,
			},
			expected: true,
		},
		{
			name: "Non-trump suit formation",
			formation: &Formation{
				Cards: []Card{NewCard(Spades, King, 1)},
				Suit:  Spades,
			},
			expected: false,
		},
		{
			name: "Joker formation",
			formation: &Formation{
				Cards: []Card{NewJoker(BigJoker, 1)},
				Suit:  Spades,
			},
			expected: true,
		},
		{
			name: "Two formation (permanent trump)",
			formation: &Formation{
				Cards: []Card{NewCard(Spades, Two, 1)},
				Suit:  Spades,
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.formation.IsTrump(trumpSuit); got != tt.expected {
				t.Errorf("IsTrump() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFormation_Compare(t *testing.T) {
	trumpSuit := Hearts
	ledSuit := Spades
	
	// Create formations for testing
	trumpFormation := &Formation{
		Type:  Single,
		Cards: []Card{NewCard(Hearts, King, 1)},
		Suit:  Hearts,
	}
	
	nonTrumpFormation := &Formation{
		Type:  Single,
		Cards: []Card{NewCard(Spades, Ace, 1)},
		Suit:  Spades,
	}
	
	higherTrumpFormation := &Formation{
		Type:  Single,
		Cards: []Card{NewCard(Hearts, Ace, 1)},
		Suit:  Hearts,
	}

	tests := []struct {
		name      string
		formation *Formation
		other     *Formation
		expected  int // positive if formation wins, negative if other wins, 0 if equal
	}{
		{
			name:      "Trump beats non-trump",
			formation: trumpFormation,
			other:     nonTrumpFormation,
			expected:  1,
		},
		{
			name:      "Non-trump loses to trump",
			formation: nonTrumpFormation,
			other:     trumpFormation,
			expected:  -1,
		},
		{
			name:      "Higher trump beats lower trump",
			formation: higherTrumpFormation,
			other:     trumpFormation,
			expected:  1,
		},
		{
			name:      "Different formation types return 0",
			formation: &Formation{Type: Single, Cards: []Card{NewCard(Hearts, King, 1)}},
			other:     &Formation{Type: Pair, Cards: []Card{NewCard(Hearts, King, 1), NewCard(Hearts, King, 2)}},
			expected:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.formation.Compare(tt.other, trumpSuit, ledSuit)
			if (result > 0 && tt.expected <= 0) || (result < 0 && tt.expected >= 0) || (result == 0 && tt.expected != 0) {
				t.Errorf("Compare() = %d, want %d", result, tt.expected)
			}
		})
	}
}

func TestValidateFormation(t *testing.T) {
	trumpSuit := Hearts
	
	tests := []struct {
		name          string
		cards         []Card
		formationType FormationType
		wantError     bool
	}{
		{
			name:          "Valid single",
			cards:         []Card{NewCard(Hearts, King, 1)},
			formationType: Single,
			wantError:     false,
		},
		{
			name:          "Valid pair",
			cards:         []Card{NewCard(Hearts, King, 1), NewCard(Hearts, King, 2)},
			formationType: Pair,
			wantError:     false,
		},
		{
			name:          "Invalid single - too many cards",
			cards:         []Card{NewCard(Hearts, King, 1), NewCard(Hearts, Queen, 1)},
			formationType: Single,
			wantError:     true,
		},
		{
			name:          "Invalid pair - different ranks",
			cards:         []Card{NewCard(Hearts, King, 1), NewCard(Hearts, Queen, 1)},
			formationType: Pair,
			wantError:     true,
		},
		{
			name: "Valid tractor",
			cards: []Card{
				NewCard(Hearts, King, 1), NewCard(Hearts, King, 2),
				NewCard(Hearts, Ace, 1), NewCard(Hearts, Ace, 2),
			},
			formationType: Tractor,
			wantError:     false,
		},
		{
			name:          "Invalid tractor - odd number of cards",
			cards:         []Card{NewCard(Hearts, King, 1), NewCard(Hearts, King, 2), NewCard(Hearts, Ace, 1)},
			formationType: Tractor,
			wantError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFormation(tt.cards, tt.formationType, trumpSuit)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateFormation() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}