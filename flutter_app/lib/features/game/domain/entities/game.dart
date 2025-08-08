import 'package:equatable/equatable.dart';
import 'card.dart';
import 'user.dart';

enum GamePhase {
  waiting,
  dealing,
  bidding,
  trumpDeclaration,
  kittyExchange,
  playing,
  ended
}

extension GamePhaseExtension on GamePhase {
  String get name {
    switch (this) {
      case GamePhase.waiting:
        return 'Waiting';
      case GamePhase.dealing:
        return 'Dealing';
      case GamePhase.bidding:
        return 'Bidding';
      case GamePhase.trumpDeclaration:
        return 'Trump Declaration';
      case GamePhase.kittyExchange:
        return 'Kitty Exchange';
      case GamePhase.playing:
        return 'Playing';
      case GamePhase.ended:
        return 'Ended';
    }
  }

  String toJson() => name;

  static GamePhase fromJson(String json) {
    switch (json) {
      case 'Waiting':
        return GamePhase.waiting;
      case 'Dealing':
        return GamePhase.dealing;
      case 'Bidding':
        return GamePhase.bidding;
      case 'Trump Declaration':
        return GamePhase.trumpDeclaration;
      case 'Kitty Exchange':
        return GamePhase.kittyExchange;
      case 'Playing':
        return GamePhase.playing;
      case 'Ended':
        return GamePhase.ended;
      default:
        throw ArgumentError('Invalid game phase: $json');
    }
  }
}

enum PlayerPosition { north, east, south, west }

extension PlayerPositionExtension on PlayerPosition {
  String get name {
    switch (this) {
      case PlayerPosition.north:
        return 'North';
      case PlayerPosition.east:
        return 'East';
      case PlayerPosition.south:
        return 'South';
      case PlayerPosition.west:
        return 'West';
    }
  }

  PlayerPosition get nextPosition {
    switch (this) {
      case PlayerPosition.north:
        return PlayerPosition.east;
      case PlayerPosition.east:
        return PlayerPosition.south;
      case PlayerPosition.south:
        return PlayerPosition.west;
      case PlayerPosition.west:
        return PlayerPosition.north;
    }
  }

  PlayerPosition get partnerPosition {
    switch (this) {
      case PlayerPosition.north:
        return PlayerPosition.south;
      case PlayerPosition.east:
        return PlayerPosition.west;
      case PlayerPosition.south:
        return PlayerPosition.north;
      case PlayerPosition.west:
        return PlayerPosition.east;
    }
  }

  String toJson() => name;

  static PlayerPosition fromJson(String json) {
    switch (json) {
      case 'North':
        return PlayerPosition.north;
      case 'East':
        return PlayerPosition.east;
      case 'South':
        return PlayerPosition.south;
      case 'West':
        return PlayerPosition.west;
      default:
        throw ArgumentError('Invalid player position: $json');
    }
  }
}

enum FormationType { single, pair, tractor }

extension FormationTypeExtension on FormationType {
  String get name {
    switch (this) {
      case FormationType.single:
        return 'Single';
      case FormationType.pair:
        return 'Pair';
      case FormationType.tractor:
        return 'Tractor';
    }
  }

  String toJson() => name;

  static FormationType fromJson(String json) {
    switch (json) {
      case 'Single':
        return FormationType.single;
      case 'Pair':
        return FormationType.pair;
      case 'Tractor':
        return FormationType.tractor;
      default:
        throw ArgumentError('Invalid formation type: $json');
    }
  }
}

class Formation extends Equatable {
  final FormationType type;
  final List<Card> cards;
  final Suit suit;

  const Formation({
    required this.type,
    required this.cards,
    required this.suit,
  });

  // Factory constructor from JSON
  factory Formation.fromJson(Map<String, dynamic> json) {
    final cardsList = json['cards'] as List<dynamic>;
    final cards = cardsList
        .map((cardJson) => Card.fromJson(cardJson as Map<String, dynamic>))
        .toList();

    return Formation(
      type: FormationTypeExtension.fromJson(json['type'] as String),
      cards: cards,
      suit: SuitExtension.fromJson(json['suit'] as String),
    );
  }

  // Convert to JSON
  Map<String, dynamic> toJson() {
    return {
      'type': type.toJson(),
      'cards': cards.map((card) => card.toJson()).toList(),
      'suit': suit.toJson(),
    };
  }

  // Get total point value of all cards in the formation
  int get pointValue {
    return cards.fold(0, (sum, card) => sum + card.pointValue);
  }

  // Get highest card in the formation
  Card getHighestCard(Suit trumpSuit) {
    if (cards.isEmpty) throw StateError('Formation has no cards');

    Card highest = cards.first;
    for (final card in cards.skip(1)) {
      if (card.getTrumpHierarchy(trumpSuit) >
          highest.getTrumpHierarchy(trumpSuit)) {
        highest = card;
      } else if (card.getTrumpHierarchy(trumpSuit) ==
          highest.getTrumpHierarchy(trumpSuit)) {
        // Same trump hierarchy, compare suit hierarchy
        if (card.suitHierarchy > highest.suitHierarchy) {
          highest = card;
        }
      }
    }
    return highest;
  }

  // Check if formation contains trump cards
  bool isTrump(Suit trumpSuit) {
    return cards.any((card) => card.isTrump(trumpSuit));
  }

  // Validate formation is legal
  bool get isValid {
    switch (type) {
      case FormationType.single:
        return cards.length == 1;
      case FormationType.pair:
        return cards.length == 2 && cards[0].isSameFace(cards[1]);
      case FormationType.tractor:
        return cards.length >= 4 && cards.length % 2 == 0;
    }
  }

  @override
  String toString() {
    return '${type.name}: ${cards.map((c) => c.toString()).join(', ')}';
  }

  @override
  List<Object?> get props => [type, cards, suit];
}

class Player extends Equatable {
  final String id;
  final String name;
  final PlayerPosition position;
  final List<Card> hand;
  final bool hasPassed;

  const Player({
    required this.id,
    required this.name,
    required this.position,
    required this.hand,
    this.hasPassed = false,
  });

  // Factory constructor from JSON
  factory Player.fromJson(Map<String, dynamic> json) {
    final handList = json['hand'] as List<dynamic>;
    final hand = handList
        .map((cardJson) => Card.fromJson(cardJson as Map<String, dynamic>))
        .toList();

    return Player(
      id: json['id'] as String,
      name: json['name'] as String,
      position: PlayerPositionExtension.fromJson(json['position'] as String),
      hand: hand,
      hasPassed: json['has_passed'] as bool? ?? false,
    );
  }

  // Convert to JSON
  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'position': position.toJson(),
      'hand': hand.map((card) => card.toJson()).toList(),
      'has_passed': hasPassed,
    };
  }

  // Get hand size
  int get handSize => hand.length;

  // Check if player has a specific card
  bool hasCard(Card card) {
    return hand.any((handCard) => handCard == card);
  }

  // Check if player has all specified cards
  bool hasCards(List<Card> cards) {
    return cards.every((card) => hasCard(card));
  }

  // Copy with new values
  Player copyWith({
    String? id,
    String? name,
    PlayerPosition? position,
    List<Card>? hand,
    bool? hasPassed,
  }) {
    return Player(
      id: id ?? this.id,
      name: name ?? this.name,
      position: position ?? this.position,
      hand: hand ?? this.hand,
      hasPassed: hasPassed ?? this.hasPassed,
    );
  }

  @override
  List<Object?> get props => [id, name, position, hand, hasPassed];
}

class BidInfo extends Equatable {
  final String playerId;
  final int amount;
  final bool isPassed;

  const BidInfo({
    required this.playerId,
    required this.amount,
    required this.isPassed,
  });

  // Factory constructor from JSON
  factory BidInfo.fromJson(Map<String, dynamic> json) {
    return BidInfo(
      playerId: json['player_id'] as String,
      amount: json['amount'] as int,
      isPassed: json['is_passed'] as bool,
    );
  }

  // Convert to JSON
  Map<String, dynamic> toJson() {
    return {
      'player_id': playerId,
      'amount': amount,
      'is_passed': isPassed,
    };
  }

  @override
  List<Object?> get props => [playerId, amount, isPassed];
}

class Trick extends Equatable {
  final String id;
  final PlayerPosition leader;
  final Map<PlayerPosition, Formation> plays;
  final String? winner;
  final int points;
  final Suit? trumpSuit;
  final Suit? ledSuit;
  final bool isComplete;
  final DateTime createdAt;
  final DateTime? completedAt;

  const Trick({
    required this.id,
    required this.leader,
    required this.plays,
    this.winner,
    this.points = 0,
    this.trumpSuit,
    this.ledSuit,
    this.isComplete = false,
    required this.createdAt,
    this.completedAt,
  });

  // Factory constructor from JSON
  factory Trick.fromJson(Map<String, dynamic> json) {
    final playsMap = json['plays'] as Map<String, dynamic>? ?? {};
    final plays = <PlayerPosition, Formation>{};

    for (final entry in playsMap.entries) {
      final position = PlayerPositionExtension.fromJson(entry.key);
      final formation = Formation.fromJson(entry.value as Map<String, dynamic>);
      plays[position] = formation;
    }

    return Trick(
      id: json['id'] as String,
      leader: PlayerPositionExtension.fromJson(json['leader'] as String),
      plays: plays,
      winner: json['winner'] as String?,
      points: json['points'] as int? ?? 0,
      trumpSuit: json['trump_suit'] != null
          ? SuitExtension.fromJson(json['trump_suit'] as String)
          : null,
      ledSuit: json['led_suit'] != null
          ? SuitExtension.fromJson(json['led_suit'] as String)
          : null,
      isComplete: json['is_complete'] as bool? ?? false,
      createdAt: DateTime.parse(json['created_at'] as String),
      completedAt: json['completed_at'] != null
          ? DateTime.parse(json['completed_at'] as String)
          : null,
    );
  }

  // Convert to JSON
  Map<String, dynamic> toJson() {
    final playsJson = <String, dynamic>{};
    for (final entry in plays.entries) {
      playsJson[entry.key.toJson()] = entry.value.toJson();
    }

    return {
      'id': id,
      'leader': leader.toJson(),
      'plays': playsJson,
      'winner': winner,
      'points': points,
      'trump_suit': trumpSuit?.toJson(),
      'led_suit': ledSuit?.toJson(),
      'is_complete': isComplete,
      'created_at': createdAt.toIso8601String(),
      'completed_at': completedAt?.toIso8601String(),
    };
  }

  @override
  List<Object?> get props => [
        id,
        leader,
        plays,
        winner,
        points,
        trumpSuit,
        ledSuit,
        isComplete,
        createdAt,
        completedAt,
      ];
}

class Game extends Equatable {
  final String id;
  final String roomId;
  final GamePhase phase;
  final List<Player> players;
  final PlayerPosition currentPlayerTurn;
  final PlayerPosition? declarer;
  final Suit? trumpSuit;
  final int contract;
  final int currentBid;
  final List<BidInfo> bidHistory;
  final int consecutivePasses;
  final Trick? currentTrick;
  final List<Trick> tricks;
  final List<Card> kitty;
  final Map<String, int> scores;
  final String? winnerTeam;
  final DateTime createdAt;
  final DateTime updatedAt;

  const Game({
    required this.id,
    required this.roomId,
    required this.phase,
    required this.players,
    required this.currentPlayerTurn,
    this.declarer,
    this.trumpSuit,
    this.contract = 0,
    this.currentBid = 125,
    required this.bidHistory,
    this.consecutivePasses = 0,
    this.currentTrick,
    required this.tricks,
    required this.kitty,
    required this.scores,
    this.winnerTeam,
    required this.createdAt,
    required this.updatedAt,
  });

  // Factory constructor from JSON
  factory Game.fromJson(Map<String, dynamic> json) {
    final playersList = json['players'] as List<dynamic>;
    final players = playersList
        .map(
            (playerJson) => Player.fromJson(playerJson as Map<String, dynamic>))
        .toList();

    final bidHistoryList = json['bid_history'] as List<dynamic>;
    final bidHistory = bidHistoryList
        .map((bidJson) => BidInfo.fromJson(bidJson as Map<String, dynamic>))
        .toList();

    final tricksList = json['tricks'] as List<dynamic>;
    final tricks = tricksList
        .map((trickJson) => Trick.fromJson(trickJson as Map<String, dynamic>))
        .toList();

    final kittyList = json['kitty'] as List<dynamic>;
    final kitty = kittyList
        .map((cardJson) => Card.fromJson(cardJson as Map<String, dynamic>))
        .toList();

    final scoresMap = json['scores'] as Map<String, dynamic>;
    final scores = scoresMap.map((key, value) => MapEntry(key, value as int));

    return Game(
      id: json['id'] as String,
      roomId: json['room_id'] as String,
      phase: GamePhaseExtension.fromJson(json['phase'] as String),
      players: players,
      currentPlayerTurn: PlayerPositionExtension.fromJson(
          json['current_player_turn'] as String),
      declarer: json['declarer'] != null
          ? PlayerPositionExtension.fromJson(json['declarer'] as String)
          : null,
      trumpSuit: json['trump_suit'] != null
          ? SuitExtension.fromJson(json['trump_suit'] as String)
          : null,
      contract: json['contract'] as int? ?? 0,
      currentBid: json['current_bid'] as int? ?? 125,
      bidHistory: bidHistory,
      consecutivePasses: json['consecutive_passes'] as int? ?? 0,
      currentTrick: json['current_trick'] != null
          ? Trick.fromJson(json['current_trick'] as Map<String, dynamic>)
          : null,
      tricks: tricks,
      kitty: kitty,
      scores: scores,
      winnerTeam: json['winner_team'] as String?,
      createdAt: DateTime.parse(json['created_at'] as String),
      updatedAt: DateTime.parse(json['updated_at'] as String),
    );
  }

  // Convert to JSON
  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'room_id': roomId,
      'phase': phase.toJson(),
      'players': players.map((player) => player.toJson()).toList(),
      'current_player_turn': currentPlayerTurn.toJson(),
      'declarer': declarer?.toJson(),
      'trump_suit': trumpSuit?.toJson(),
      'contract': contract,
      'current_bid': currentBid,
      'bid_history': bidHistory.map((bid) => bid.toJson()).toList(),
      'consecutive_passes': consecutivePasses,
      'current_trick': currentTrick?.toJson(),
      'tricks': tricks.map((trick) => trick.toJson()).toList(),
      'kitty': kitty.map((card) => card.toJson()).toList(),
      'scores': scores,
      'winner_team': winnerTeam,
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt.toIso8601String(),
    };
  }

  // Get player by ID
  Player? getPlayer(String playerId) {
    try {
      return players.firstWhere((player) => player.id == playerId);
    } catch (e) {
      return null;
    }
  }

  // Get player by position
  Player? getPlayerByPosition(PlayerPosition position) {
    try {
      return players.firstWhere((player) => player.position == position);
    } catch (e) {
      return null;
    }
  }

  // Get current player
  Player? get currentPlayer => getPlayerByPosition(currentPlayerTurn);

  // Check if game is complete
  bool get isComplete => players.every((player) => player.hand.isEmpty);

  // Check if player is on declarer's team
  bool isOnDeclarerTeam(PlayerPosition position) {
    if (declarer == null) return false;
    return position == declarer || position == declarer!.partnerPosition;
  }

  // Copy with new values
  Game copyWith({
    String? id,
    String? roomId,
    GamePhase? phase,
    List<Player>? players,
    PlayerPosition? currentPlayerTurn,
    PlayerPosition? declarer,
    Suit? trumpSuit,
    int? contract,
    int? currentBid,
    List<BidInfo>? bidHistory,
    int? consecutivePasses,
    Trick? currentTrick,
    List<Trick>? tricks,
    List<Card>? kitty,
    Map<String, int>? scores,
    String? winnerTeam,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) {
    return Game(
      id: id ?? this.id,
      roomId: roomId ?? this.roomId,
      phase: phase ?? this.phase,
      players: players ?? this.players,
      currentPlayerTurn: currentPlayerTurn ?? this.currentPlayerTurn,
      declarer: declarer ?? this.declarer,
      trumpSuit: trumpSuit ?? this.trumpSuit,
      contract: contract ?? this.contract,
      currentBid: currentBid ?? this.currentBid,
      bidHistory: bidHistory ?? this.bidHistory,
      consecutivePasses: consecutivePasses ?? this.consecutivePasses,
      currentTrick: currentTrick ?? this.currentTrick,
      tricks: tricks ?? this.tricks,
      kitty: kitty ?? this.kitty,
      scores: scores ?? this.scores,
      winnerTeam: winnerTeam ?? this.winnerTeam,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }

  @override
  List<Object?> get props => [
        id,
        roomId,
        phase,
        players,
        currentPlayerTurn,
        declarer,
        trumpSuit,
        contract,
        currentBid,
        bidHistory,
        consecutivePasses,
        currentTrick,
        tricks,
        kitty,
        scores,
        winnerTeam,
        createdAt,
        updatedAt,
      ];
}
