package domain

import (
	"fmt"
	"time"
)

// GamePhase represents the current phase of the game
type GamePhase int

const (
	PhaseWaiting GamePhase = iota
	PhaseDealing
	PhaseBidding
	PhaseTrumpDeclaration
	PhaseKittyExchange
	PhasePlaying
	PhaseEnded
)

func (p GamePhase) String() string {
	switch p {
	case PhaseWaiting:
		return "Waiting"
	case PhaseDealing:
		return "Dealing"
	case PhaseBidding:
		return "Bidding"
	case PhaseTrumpDeclaration:
		return "Trump Declaration"
	case PhaseKittyExchange:
		return "Kitty Exchange"
	case PhasePlaying:
		return "Playing"
	case PhaseEnded:
		return "Ended"
	default:
		return "Unknown"
	}
}

// PlayerPosition represents the position of a player at the table
type PlayerPosition int

const (
	North PlayerPosition = iota
	East
	South
	West
)

func (p PlayerPosition) String() string {
	switch p {
	case North:
		return "North"
	case East:
		return "East"
	case South:
		return "South"
	case West:
		return "West"
	default:
		return "Unknown"
	}
}

// GetNextPosition returns the next position clockwise
func (p PlayerPosition) GetNextPosition() PlayerPosition {
	return PlayerPosition((int(p) + 1) % 4)
}

// GetPartnerPosition returns the partner's position
func (p PlayerPosition) GetPartnerPosition() PlayerPosition {
	return PlayerPosition((int(p) + 2) % 4)
}

// Player represents a player in the game
type Player struct {
	ID       string         `json:"id"`
	Name     string         `json:"name"`
	Position PlayerPosition `json:"position"`
	Hand     []Card         `json:"hand"`
	HasPassed bool          `json:"has_passed"` // For bidding phase
}

// NewPlayer creates a new player
func NewPlayer(id, name string, position PlayerPosition) *Player {
	return &Player{
		ID:       id,
		Name:     name,
		Position: position,
		Hand:     make([]Card, 0, 25),
		HasPassed: false,
	}
}

// AddCard adds a card to the player's hand
func (p *Player) AddCard(card Card) {
	p.Hand = append(p.Hand, card)
}

// AddCards adds multiple cards to the player's hand
func (p *Player) AddCards(cards []Card) {
	p.Hand = append(p.Hand, cards...)
}

// RemoveCard removes a card from the player's hand
func (p *Player) RemoveCard(card Card) error {
	for i, handCard := range p.Hand {
		if handCard.IsEqual(card) {
			p.Hand = append(p.Hand[:i], p.Hand[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("card not found in player's hand")
}

// RemoveCards removes multiple cards from the player's hand
func (p *Player) RemoveCards(cards []Card) error {
	for _, card := range cards {
		if err := p.RemoveCard(card); err != nil {
			return err
		}
	}
	return nil
}

// HasCard checks if the player has a specific card
func (p *Player) HasCard(card Card) bool {
	for _, handCard := range p.Hand {
		if handCard.IsEqual(card) {
			return true
		}
	}
	return false
}

// HasCards checks if the player has all specified cards
func (p *Player) HasCards(cards []Card) bool {
	for _, card := range cards {
		if !p.HasCard(card) {
			return false
		}
	}
	return true
}

// GetHandSize returns the number of cards in the player's hand
func (p *Player) GetHandSize() int {
	return len(p.Hand)
}

// BidInfo represents a bid made by a player
type BidInfo struct {
	PlayerID string `json:"player_id"`
	Amount   int    `json:"amount"`
	IsPassed bool   `json:"is_passed"`
}

// GameState represents the complete state of a Chinese Bridge game
type GameState struct {
	ID                string            `json:"id"`
	RoomID            string            `json:"room_id"`
	Phase             GamePhase         `json:"phase"`
	Players           [4]*Player        `json:"players"`
	CurrentPlayerTurn PlayerPosition    `json:"current_player_turn"`
	Declarer          *PlayerPosition   `json:"declarer,omitempty"`
	TrumpSuit         *Suit             `json:"trump_suit,omitempty"`
	Contract          int               `json:"contract"`
	CurrentBid        int               `json:"current_bid"`
	BidHistory        []BidInfo         `json:"bid_history"`
	ConsecutivePasses int               `json:"consecutive_passes"`
	CurrentTrick      *Trick            `json:"current_trick,omitempty"`
	Tricks            []Trick           `json:"tricks"`
	Kitty             []Card            `json:"kitty"`
	Scores            map[string]int    `json:"scores"`
	WinnerTeam        *string           `json:"winner_team,omitempty"` // "declarer" or "defenders"
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
}

// NewGameState creates a new game state
func NewGameState(id, roomID string, playerIDs []string, playerNames []string) (*GameState, error) {
	if len(playerIDs) != 4 || len(playerNames) != 4 {
		return nil, fmt.Errorf("exactly 4 players required")
	}

	gameState := &GameState{
		ID:                id,
		RoomID:            roomID,
		Phase:             PhaseWaiting,
		CurrentPlayerTurn: North,
		Contract:          0,
		CurrentBid:        125, // Starting bid
		BidHistory:        make([]BidInfo, 0),
		ConsecutivePasses: 0,
		Tricks:            make([]Trick, 0),
		Kitty:             make([]Card, 0, 8),
		Scores:            make(map[string]int),
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// Initialize players
	for i := 0; i < 4; i++ {
		gameState.Players[i] = NewPlayer(playerIDs[i], playerNames[i], PlayerPosition(i))
		gameState.Scores[playerIDs[i]] = 0
	}

	return gameState, nil
}

// GetPlayer returns a player by ID
func (gs *GameState) GetPlayer(playerID string) *Player {
	for _, player := range gs.Players {
		if player.ID == playerID {
			return player
		}
	}
	return nil
}

// GetPlayerByPosition returns a player by position
func (gs *GameState) GetPlayerByPosition(position PlayerPosition) *Player {
	if int(position) >= 0 && int(position) < 4 {
		return gs.Players[position]
	}
	return nil
}

// GetCurrentPlayer returns the player whose turn it is
func (gs *GameState) GetCurrentPlayer() *Player {
	return gs.GetPlayerByPosition(gs.CurrentPlayerTurn)
}

// NextTurn advances to the next player's turn
func (gs *GameState) NextTurn() {
	gs.CurrentPlayerTurn = gs.CurrentPlayerTurn.GetNextPosition()
	gs.UpdatedAt = time.Now()
}

// DealCards deals cards to all players and sets up the kitty
func (gs *GameState) DealCards(deck *Deck) error {
	if gs.Phase != PhaseWaiting {
		return fmt.Errorf("can only deal cards in waiting phase")
	}

	// Deal 25 cards to each player
	for i := 0; i < 4; i++ {
		cards, err := deck.Deal(25)
		if err != nil {
			return fmt.Errorf("failed to deal cards to player %d: %w", i, err)
		}
		gs.Players[i].AddCards(cards)
	}

	// Remaining 8 cards go to kitty
	kittyCards, err := deck.Deal(8)
	if err != nil {
		return fmt.Errorf("failed to deal kitty cards: %w", err)
	}
	gs.Kitty = kittyCards

	gs.Phase = PhaseBidding
	gs.UpdatedAt = time.Now()
	return nil
}

// PlaceBid places a bid for the current player
func (gs *GameState) PlaceBid(playerID string, bidAmount int) error {
	if gs.Phase != PhaseBidding {
		return fmt.Errorf("not in bidding phase")
	}

	currentPlayer := gs.GetCurrentPlayer()
	if currentPlayer.ID != playerID {
		return fmt.Errorf("not player's turn")
	}

	if currentPlayer.HasPassed {
		return fmt.Errorf("player has already passed and cannot bid")
	}

	// Validate bid amount
	if bidAmount < 95 || bidAmount > 200 {
		return fmt.Errorf("bid must be between 95 and 200")
	}

	if bidAmount >= gs.CurrentBid {
		return fmt.Errorf("bid must be lower than current bid of %d", gs.CurrentBid)
	}

	if (gs.CurrentBid-bidAmount)%5 != 0 {
		return fmt.Errorf("bid must decrease by increments of 5")
	}

	// Record the bid
	gs.BidHistory = append(gs.BidHistory, BidInfo{
		PlayerID: playerID,
		Amount:   bidAmount,
		IsPassed: false,
	})

	gs.CurrentBid = bidAmount
	gs.ConsecutivePasses = 0
	gs.NextTurn()

	return nil
}

// PassBid passes the current player's turn in bidding
func (gs *GameState) PassBid(playerID string) error {
	if gs.Phase != PhaseBidding {
		return fmt.Errorf("not in bidding phase")
	}

	currentPlayer := gs.GetCurrentPlayer()
	if currentPlayer.ID != playerID {
		return fmt.Errorf("not player's turn")
	}

	if currentPlayer.HasPassed {
		return fmt.Errorf("player has already passed")
	}

	// Mark player as passed
	currentPlayer.HasPassed = true
	gs.BidHistory = append(gs.BidHistory, BidInfo{
		PlayerID: playerID,
		Amount:   0,
		IsPassed: true,
	})

	gs.ConsecutivePasses++

	// Check if bidding should end
	if gs.ConsecutivePasses >= 3 {
		// Find the declarer (last player to make a bid)
		for i := len(gs.BidHistory) - 1; i >= 0; i-- {
			if !gs.BidHistory[i].IsPassed {
				declarerPlayer := gs.GetPlayer(gs.BidHistory[i].PlayerID)
				if declarerPlayer != nil {
					gs.Declarer = &declarerPlayer.Position
					gs.Contract = gs.BidHistory[i].Amount
					gs.Phase = PhaseTrumpDeclaration
					gs.CurrentPlayerTurn = declarerPlayer.Position
					break
				}
			}
		}
	} else {
		gs.NextTurn()
	}

	gs.UpdatedAt = time.Now()
	return nil
}

// DeclareTrump declares the trump suit
func (gs *GameState) DeclareTrump(playerID string, trumpSuit Suit) error {
	if gs.Phase != PhaseTrumpDeclaration {
		return fmt.Errorf("not in trump declaration phase")
	}

	if gs.Declarer == nil {
		return fmt.Errorf("no declarer set")
	}

	declarer := gs.GetPlayerByPosition(*gs.Declarer)
	if declarer.ID != playerID {
		return fmt.Errorf("only the declarer can declare trump")
	}

	gs.TrumpSuit = &trumpSuit
	gs.Phase = PhaseKittyExchange
	gs.UpdatedAt = time.Now()

	return nil
}

// ExchangeKitty allows the declarer to exchange cards with the kitty
func (gs *GameState) ExchangeKitty(playerID string, cardsToDiscard []Card) error {
	if gs.Phase != PhaseKittyExchange {
		return fmt.Errorf("not in kitty exchange phase")
	}

	if gs.Declarer == nil {
		return fmt.Errorf("no declarer set")
	}

	declarer := gs.GetPlayerByPosition(*gs.Declarer)
	if declarer.ID != playerID {
		return fmt.Errorf("only the declarer can exchange kitty")
	}

	if len(cardsToDiscard) != 8 {
		return fmt.Errorf("must discard exactly 8 cards")
	}

	// Verify declarer has all cards to discard
	if !declarer.HasCards(cardsToDiscard) {
		return fmt.Errorf("player does not have all specified cards")
	}

	// Add kitty cards to declarer's hand
	declarer.AddCards(gs.Kitty)

	// Remove discarded cards from declarer's hand
	if err := declarer.RemoveCards(cardsToDiscard); err != nil {
		return fmt.Errorf("failed to remove cards from hand: %w", err)
	}

	// Update kitty with discarded cards
	gs.Kitty = cardsToDiscard

	gs.Phase = PhasePlaying
	gs.UpdatedAt = time.Now()

	return nil
}

// StartNewTrick starts a new trick
func (gs *GameState) StartNewTrick() {
	trickID := fmt.Sprintf("%s_trick_%d", gs.ID, len(gs.Tricks)+1)
	gs.CurrentTrick = NewTrick(trickID, gs.CurrentPlayerTurn)
}

// IsGameComplete checks if the game is complete
func (gs *GameState) IsGameComplete() bool {
	// Game is complete when all players have no cards left
	for _, player := range gs.Players {
		if len(player.Hand) > 0 {
			return false
		}
	}
	return true
}

// CalculateFinalScore calculates the final score and determines the winner
func (gs *GameState) CalculateFinalScore() {
	if gs.Declarer == nil {
		return
	}

	// Calculate total points captured by defenders
	defendersPoints := 0
	for _, trick := range gs.Tricks {
		winner := gs.GetPlayer(trick.Winner)
		if winner != nil && winner.Position != *gs.Declarer && winner.Position != gs.Declarer.GetPartnerPosition() {
			defendersPoints += trick.Points
		}
	}

	// Add kitty points to the final trick winner's team
	if len(gs.Tricks) > 0 {
		lastTrick := gs.Tricks[len(gs.Tricks)-1]
		lastWinner := gs.GetPlayer(lastTrick.Winner)
		if lastWinner != nil {
			kittyPoints := 0
			for _, card := range gs.Kitty {
				kittyPoints += card.GetPointValue()
			}

			// If last trick winner is on defenders team, add kitty points to defenders
			if lastWinner.Position != *gs.Declarer && lastWinner.Position != gs.Declarer.GetPartnerPosition() {
				defendersPoints += kittyPoints
			}
		}
	}

	// Determine winner
	if defendersPoints >= gs.Contract {
		gs.WinnerTeam = stringPtr("defenders")
	} else {
		gs.WinnerTeam = stringPtr("declarer")
	}

	gs.Phase = PhaseEnded
	gs.UpdatedAt = time.Now()
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}

// GetTeammates returns the teammate of a given position
func (gs *GameState) GetTeammates(position PlayerPosition) PlayerPosition {
	return position.GetPartnerPosition()
}

// IsOnDeclarerTeam checks if a position is on the declarer's team
func (gs *GameState) IsOnDeclarerTeam(position PlayerPosition) bool {
	if gs.Declarer == nil {
		return false
	}
	return position == *gs.Declarer || position == gs.Declarer.GetPartnerPosition()
}

// GetGameSummary returns a summary of the game state
func (gs *GameState) GetGameSummary() map[string]interface{} {
	summary := map[string]interface{}{
		"id":                  gs.ID,
		"room_id":             gs.RoomID,
		"phase":               gs.Phase.String(),
		"current_player_turn": gs.CurrentPlayerTurn.String(),
		"contract":            gs.Contract,
		"current_bid":         gs.CurrentBid,
		"tricks_played":       len(gs.Tricks),
		"created_at":          gs.CreatedAt,
		"updated_at":          gs.UpdatedAt,
	}

	if gs.Declarer != nil {
		summary["declarer"] = gs.Declarer.String()
	}

	if gs.TrumpSuit != nil {
		summary["trump_suit"] = gs.TrumpSuit.String()
	}

	if gs.WinnerTeam != nil {
		summary["winner_team"] = *gs.WinnerTeam
	}

	return summary
}