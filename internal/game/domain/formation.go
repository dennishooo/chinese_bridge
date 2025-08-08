package domain

import (
	"fmt"
	"sort"
)

// FormationType represents the different types of card formations
type FormationType int

const (
	Single FormationType = iota
	Pair
	Tractor
)

func (f FormationType) String() string {
	switch f {
	case Single:
		return "Single"
	case Pair:
		return "Pair"
	case Tractor:
		return "Tractor"
	default:
		return "Unknown"
	}
}

// Formation represents a valid combination of cards that can be played
type Formation struct {
	Type  FormationType `json:"type"`
	Cards []Card        `json:"cards"`
	Suit  Suit          `json:"suit"`
}

// NewSingle creates a single card formation
func NewSingle(card Card) *Formation {
	suit := card.Suit
	if card.IsJoker {
		// Jokers are considered trump suit for formation purposes
		suit = Spades // Default, will be overridden by trump suit context
	}
	
	return &Formation{
		Type:  Single,
		Cards: []Card{card},
		Suit:  suit,
	}
}

// NewPair creates a pair formation from two identical cards
func NewPair(card1, card2 Card) (*Formation, error) {
	if !card1.IsSameFace(card2) {
		return nil, fmt.Errorf("cards must have the same face value to form a pair")
	}

	suit := card1.Suit
	if card1.IsJoker {
		suit = Spades // Default, will be overridden by trump suit context
	}

	return &Formation{
		Type:  Pair,
		Cards: []Card{card1, card2},
		Suit:  suit,
	}, nil
}

// NewTractor creates a tractor formation from consecutive pairs
func NewTractor(pairs [][]Card, trumpSuit Suit) (*Formation, error) {
	if len(pairs) < 2 {
		return nil, fmt.Errorf("tractor must have at least 2 pairs")
	}

	// Validate each pair
	allCards := make([]Card, 0, len(pairs)*2)
	pairRanks := make([]Rank, 0, len(pairs))

	for _, pair := range pairs {
		if len(pair) != 2 {
			return nil, fmt.Errorf("each element must be a pair of 2 cards")
		}

		if !pair[0].IsSameFace(pair[1]) {
			return nil, fmt.Errorf("cards in pair must have the same face value")
		}

		// Jokers and 2s cannot be part of tractors
		if pair[0].IsJoker || pair[0].Rank == Two {
			return nil, fmt.Errorf("jokers and 2s cannot be part of tractors")
		}

		allCards = append(allCards, pair...)
		pairRanks = append(pairRanks, pair[0].Rank)
	}

	// Check all pairs are from the same suit
	firstSuit := pairs[0][0].Suit
	for _, pair := range pairs {
		if pair[0].Suit != firstSuit {
			return nil, fmt.Errorf("all pairs in tractor must be from the same suit")
		}
	}

	// Sort ranks and check they are consecutive
	sort.Slice(pairRanks, func(i, j int) bool {
		return pairRanks[i] < pairRanks[j]
	})

	for i := 1; i < len(pairRanks); i++ {
		if pairRanks[i] != pairRanks[i-1]+1 {
			return nil, fmt.Errorf("tractor pairs must be consecutive ranks")
		}
	}

	return &Formation{
		Type:  Tractor,
		Cards: allCards,
		Suit:  firstSuit,
	}, nil
}

// IsValid checks if the formation is valid according to Chinese Bridge rules
func (f *Formation) IsValid() error {
	switch f.Type {
	case Single:
		if len(f.Cards) != 1 {
			return fmt.Errorf("single formation must have exactly 1 card")
		}
	case Pair:
		if len(f.Cards) != 2 {
			return fmt.Errorf("pair formation must have exactly 2 cards")
		}
		if !f.Cards[0].IsSameFace(f.Cards[1]) {
			return fmt.Errorf("pair formation cards must have the same face value")
		}
	case Tractor:
		if len(f.Cards) < 4 || len(f.Cards)%2 != 0 {
			return fmt.Errorf("tractor formation must have at least 4 cards in pairs")
		}
		// Additional tractor validation would be implemented here
	default:
		return fmt.Errorf("unknown formation type")
	}
	return nil
}

// GetHighestCard returns the highest ranking card in the formation
func (f *Formation) GetHighestCard(trumpSuit Suit) Card {
	if len(f.Cards) == 0 {
		return Card{} // Empty card
	}

	highest := f.Cards[0]
	for _, card := range f.Cards[1:] {
		if card.GetTrumpHierarchy(trumpSuit) > highest.GetTrumpHierarchy(trumpSuit) {
			highest = card
		} else if card.GetTrumpHierarchy(trumpSuit) == highest.GetTrumpHierarchy(trumpSuit) {
			// Same trump hierarchy, compare suit hierarchy
			if card.GetSuitHierarchy() > highest.GetSuitHierarchy() {
				highest = card
			}
		}
	}
	return highest
}

// GetPointValue returns the total point value of all cards in the formation
func (f *Formation) GetPointValue() int {
	total := 0
	for _, card := range f.Cards {
		total += card.GetPointValue()
	}
	return total
}

// CanFollow checks if this formation can follow the led formation
func (f *Formation) CanFollow(led *Formation, trumpSuit Suit) bool {
	// Must match formation type
	if f.Type != led.Type {
		return false
	}

	// Must match suit if possible
	if f.Suit == led.Suit {
		return true
	}

	// If void in led suit, can ruff with trump or sluff with any suit
	return true
}

// IsTrump checks if the formation contains trump cards
func (f *Formation) IsTrump(trumpSuit Suit) bool {
	for _, card := range f.Cards {
		if card.GetTrumpHierarchy(trumpSuit) > 0 {
			return true
		}
	}
	return false
}

// Compare compares two formations to determine which wins
// Returns positive if f wins, negative if other wins, 0 if equal
func (f *Formation) Compare(other *Formation, trumpSuit Suit, ledSuit Suit) int {
	// Different formation types cannot be compared directly
	if f.Type != other.Type {
		return 0
	}

	fIsTrump := f.IsTrump(trumpSuit)
	otherIsTrump := other.IsTrump(trumpSuit)

	// Trump formations beat non-trump formations
	if fIsTrump && !otherIsTrump {
		return 1
	}
	if !fIsTrump && otherIsTrump {
		return -1
	}

	// Both trump or both non-trump, compare highest cards
	fHighest := f.GetHighestCard(trumpSuit)
	otherHighest := other.GetHighestCard(trumpSuit)

	fHierarchy := fHighest.GetTrumpHierarchy(trumpSuit)
	otherHierarchy := otherHighest.GetTrumpHierarchy(trumpSuit)

	if fHierarchy != otherHierarchy {
		if fHierarchy > otherHierarchy {
			return 1
		}
		return -1
	}

	// Same trump hierarchy, compare suit hierarchy
	fSuitHierarchy := fHighest.GetSuitHierarchy()
	otherSuitHierarchy := otherHighest.GetSuitHierarchy()

	if fSuitHierarchy > otherSuitHierarchy {
		return 1
	}
	if fSuitHierarchy < otherSuitHierarchy {
		return -1
	}

	return 0
}

// String returns a string representation of the formation
func (f *Formation) String() string {
	cardStrs := make([]string, len(f.Cards))
	for i, card := range f.Cards {
		cardStrs[i] = card.String()
	}
	return fmt.Sprintf("%s: [%s]", f.Type.String(), fmt.Sprintf("%v", cardStrs))
}

// ValidateFormation validates a set of cards can form the specified formation type
func ValidateFormation(cards []Card, formationType FormationType, trumpSuit Suit) error {
	switch formationType {
	case Single:
		if len(cards) != 1 {
			return fmt.Errorf("single formation requires exactly 1 card")
		}
	case Pair:
		if len(cards) != 2 {
			return fmt.Errorf("pair formation requires exactly 2 cards")
		}
		if !cards[0].IsSameFace(cards[1]) {
			return fmt.Errorf("pair formation requires two cards with the same face value")
		}
	case Tractor:
		if len(cards) < 4 || len(cards)%2 != 0 {
			return fmt.Errorf("tractor formation requires at least 4 cards in pairs")
		}
		
		// Group cards into pairs
		pairs := make([][]Card, 0, len(cards)/2)
		cardMap := make(map[string][]Card)
		
		// Group cards by face value
		for _, card := range cards {
			key := fmt.Sprintf("%s_%s", card.Suit.String(), card.Rank.String())
			if card.IsJoker {
				key = fmt.Sprintf("joker_%s", card.JokerType.String())
			}
			cardMap[key] = append(cardMap[key], card)
		}
		
		// Validate each group has exactly 2 cards
		for _, group := range cardMap {
			if len(group) != 2 {
				return fmt.Errorf("tractor formation requires each rank to appear exactly twice")
			}
			pairs = append(pairs, group)
		}
		
		// Validate tractor formation
		_, err := NewTractor(pairs, trumpSuit)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown formation type")
	}
	
	return nil
}