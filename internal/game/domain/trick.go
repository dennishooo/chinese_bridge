package domain

import (
	"fmt"
	"time"
)

// Trick represents a single trick in the game
type Trick struct {
	ID        string                     `json:"id"`
	Leader    PlayerPosition             `json:"leader"`
	Plays     map[PlayerPosition]*Formation `json:"plays"`
	Winner    string                     `json:"winner"`
	Points    int                        `json:"points"`
	TrumpSuit *Suit                     `json:"trump_suit,omitempty"`
	LedSuit   *Suit                     `json:"led_suit,omitempty"`
	IsComplete bool                      `json:"is_complete"`
	CreatedAt time.Time                  `json:"created_at"`
	CompletedAt *time.Time               `json:"completed_at,omitempty"`
}

// NewTrick creates a new trick
func NewTrick(id string, leader PlayerPosition) *Trick {
	return &Trick{
		ID:         id,
		Leader:     leader,
		Plays:      make(map[PlayerPosition]*Formation),
		Winner:     "",
		Points:     0,
		IsComplete: false,
		CreatedAt:  time.Now(),
	}
}

// AddPlay adds a player's formation to the trick
func (t *Trick) AddPlay(position PlayerPosition, formation *Formation, trumpSuit Suit) error {
	if t.IsComplete {
		return fmt.Errorf("trick is already complete")
	}

	if _, exists := t.Plays[position]; exists {
		return fmt.Errorf("player at position %s has already played", position.String())
	}

	// Validate the play follows suit rules
	if err := t.validatePlay(position, formation, trumpSuit); err != nil {
		return err
	}

	t.Plays[position] = formation
	t.TrumpSuit = &trumpSuit

	// Set led suit from the first play
	if len(t.Plays) == 1 {
		if formation.IsTrump(trumpSuit) {
			// If trump is led, led suit is trump
			t.LedSuit = &trumpSuit
		} else {
			t.LedSuit = &formation.Suit
		}
	}

	// Check if trick is complete (all 4 players have played)
	if len(t.Plays) == 4 {
		t.completeTrick(trumpSuit)
	}

	return nil
}

// validatePlay validates that a formation can be legally played
func (t *Trick) validatePlay(position PlayerPosition, formation *Formation, trumpSuit Suit) error {
	if len(t.Plays) == 0 {
		// First play (leader) can play anything
		return nil
	}

	// Get the led formation
	leaderFormation := t.Plays[t.Leader]
	if leaderFormation == nil {
		return fmt.Errorf("leader formation not found")
	}

	// Must match formation type
	if formation.Type != leaderFormation.Type {
		return fmt.Errorf("must match led formation type %s", leaderFormation.Type.String())
	}

	// If following suit, formation is valid
	if formation.Suit == leaderFormation.Suit {
		return nil
	}

	// If not following suit, must be void in led suit or playing trump
	// This validation would require access to player's hand, which should be done at game state level
	// For now, we assume the play is valid if it reaches this point
	return nil
}

// completeTrick determines the winner and calculates points
func (t *Trick) completeTrick(trumpSuit Suit) {
	if len(t.Plays) != 4 {
		return
	}

	leaderFormation := t.Plays[t.Leader]
	if leaderFormation == nil {
		return
	}

	winningPosition := t.Leader
	winningFormation := leaderFormation

	// Compare all plays to find the winner
	currentPos := t.Leader.GetNextPosition()
	for i := 0; i < 3; i++ {
		currentFormation := t.Plays[currentPos]
		if currentFormation != nil {
			// Compare formations
			comparison := currentFormation.Compare(winningFormation, trumpSuit, *t.LedSuit)
			if comparison > 0 {
				winningPosition = currentPos
				winningFormation = currentFormation
			}
		}
		currentPos = currentPos.GetNextPosition()
	}

	// Calculate total points in the trick
	totalPoints := 0
	for _, formation := range t.Plays {
		totalPoints += formation.GetPointValue()
	}

	t.Winner = winningPosition.String()
	t.Points = totalPoints
	t.IsComplete = true
	now := time.Now()
	t.CompletedAt = &now
}

// GetWinningFormation returns the formation that won the trick
func (t *Trick) GetWinningFormation() *Formation {
	if !t.IsComplete {
		return nil
	}

	for position, formation := range t.Plays {
		if position.String() == t.Winner {
			return formation
		}
	}
	return nil
}

// GetPlayOrder returns the order in which players should play
func (t *Trick) GetPlayOrder() []PlayerPosition {
	order := make([]PlayerPosition, 4)
	current := t.Leader
	for i := 0; i < 4; i++ {
		order[i] = current
		current = current.GetNextPosition()
	}
	return order
}

// GetNextToPlay returns the next player who needs to play
func (t *Trick) GetNextToPlay() *PlayerPosition {
	if t.IsComplete {
		return nil
	}

	playOrder := t.GetPlayOrder()
	for _, position := range playOrder {
		if _, hasPlayed := t.Plays[position]; !hasPlayed {
			return &position
		}
	}
	return nil
}

// HasPlayerPlayed checks if a player at the given position has played
func (t *Trick) HasPlayerPlayed(position PlayerPosition) bool {
	_, exists := t.Plays[position]
	return exists
}

// GetPlayerFormation returns the formation played by a player
func (t *Trick) GetPlayerFormation(position PlayerPosition) *Formation {
	return t.Plays[position]
}

// CanPlayerPlay checks if a player can play in this trick
func (t *Trick) CanPlayerPlay(position PlayerPosition) bool {
	if t.IsComplete {
		return false
	}

	// Check if it's the player's turn
	nextToPlay := t.GetNextToPlay()
	return nextToPlay != nil && *nextToPlay == position
}

// GetTrickSummary returns a summary of the trick
func (t *Trick) GetTrickSummary() map[string]interface{} {
	summary := map[string]interface{}{
		"id":           t.ID,
		"leader":       t.Leader.String(),
		"is_complete":  t.IsComplete,
		"points":       t.Points,
		"created_at":   t.CreatedAt,
		"plays_count":  len(t.Plays),
	}

	if t.Winner != "" {
		summary["winner"] = t.Winner
	}

	if t.LedSuit != nil {
		summary["led_suit"] = t.LedSuit.String()
	}

	if t.TrumpSuit != nil {
		summary["trump_suit"] = t.TrumpSuit.String()
	}

	if t.CompletedAt != nil {
		summary["completed_at"] = *t.CompletedAt
	}

	// Add play details
	plays := make(map[string]interface{})
	for position, formation := range t.Plays {
		plays[position.String()] = map[string]interface{}{
			"type":        formation.Type.String(),
			"suit":        formation.Suit.String(),
			"cards_count": len(formation.Cards),
			"points":      formation.GetPointValue(),
		}
	}
	summary["plays"] = plays

	return summary
}

// String returns a string representation of the trick
func (t *Trick) String() string {
	status := "In Progress"
	if t.IsComplete {
		status = fmt.Sprintf("Won by %s", t.Winner)
	}

	return fmt.Sprintf("Trick %s: Leader=%s, Plays=%d/4, Points=%d, Status=%s",
		t.ID, t.Leader.String(), len(t.Plays), t.Points, status)
}

// ValidateFormationAgainstTrick validates if a formation can be played in this trick
func (t *Trick) ValidateFormationAgainstTrick(position PlayerPosition, formation *Formation, playerHand []Card, trumpSuit Suit) error {
	if t.IsComplete {
		return fmt.Errorf("trick is already complete")
	}

	if t.HasPlayerPlayed(position) {
		return fmt.Errorf("player has already played in this trick")
	}

	if !t.CanPlayerPlay(position) {
		return fmt.Errorf("not player's turn to play")
	}

	// Validate formation itself
	if err := formation.IsValid(); err != nil {
		return fmt.Errorf("invalid formation: %w", err)
	}

	// Validate player has all cards in the formation
	for _, card := range formation.Cards {
		hasCard := false
		for _, handCard := range playerHand {
			if handCard.IsEqual(card) {
				hasCard = true
				break
			}
		}
		if !hasCard {
			return fmt.Errorf("player does not have card: %s", card.String())
		}
	}

	// If this is the first play, any valid formation is allowed
	if len(t.Plays) == 0 {
		return nil
	}

	// Validate against suit-following rules
	return t.validatePlay(position, formation, trumpSuit)
}

// GetRemainingPositions returns positions that haven't played yet
func (t *Trick) GetRemainingPositions() []PlayerPosition {
	remaining := make([]PlayerPosition, 0, 4)
	playOrder := t.GetPlayOrder()
	
	for _, position := range playOrder {
		if !t.HasPlayerPlayed(position) {
			remaining = append(remaining, position)
		}
	}
	
	return remaining
}