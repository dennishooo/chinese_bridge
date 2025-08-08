import 'package:equatable/equatable.dart';

enum Suit { spades, hearts, clubs, diamonds }

extension SuitExtension on Suit {
  String get name {
    switch (this) {
      case Suit.spades:
        return 'Spades';
      case Suit.hearts:
        return 'Hearts';
      case Suit.clubs:
        return 'Clubs';
      case Suit.diamonds:
        return 'Diamonds';
    }
  }

  String toJson() => name;

  static Suit fromJson(String json) {
    switch (json) {
      case 'Spades':
        return Suit.spades;
      case 'Hearts':
        return Suit.hearts;
      case 'Clubs':
        return Suit.clubs;
      case 'Diamonds':
        return Suit.diamonds;
      default:
        throw ArgumentError('Invalid suit: $json');
    }
  }
}

enum Rank {
  two,
  three,
  four,
  five,
  six,
  seven,
  eight,
  nine,
  ten,
  jack,
  queen,
  king,
  ace
}

extension RankExtension on Rank {
  String get name {
    switch (this) {
      case Rank.two:
        return '2';
      case Rank.three:
        return '3';
      case Rank.four:
        return '4';
      case Rank.five:
        return '5';
      case Rank.six:
        return '6';
      case Rank.seven:
        return '7';
      case Rank.eight:
        return '8';
      case Rank.nine:
        return '9';
      case Rank.ten:
        return '10';
      case Rank.jack:
        return 'J';
      case Rank.queen:
        return 'Q';
      case Rank.king:
        return 'K';
      case Rank.ace:
        return 'A';
    }
  }

  int get pointValue {
    switch (this) {
      case Rank.five:
        return 5;
      case Rank.ten:
      case Rank.king:
        return 10;
      default:
        return 0;
    }
  }

  int get value {
    return index + 2; // Two = 2, Three = 3, ..., Ace = 14
  }

  String toJson() => name;

  static Rank fromJson(String json) {
    switch (json) {
      case '2':
        return Rank.two;
      case '3':
        return Rank.three;
      case '4':
        return Rank.four;
      case '5':
        return Rank.five;
      case '6':
        return Rank.six;
      case '7':
        return Rank.seven;
      case '8':
        return Rank.eight;
      case '9':
        return Rank.nine;
      case '10':
        return Rank.ten;
      case 'J':
        return Rank.jack;
      case 'Q':
        return Rank.queen;
      case 'K':
        return Rank.king;
      case 'A':
        return Rank.ace;
      default:
        throw ArgumentError('Invalid rank: $json');
    }
  }
}

enum JokerType { bigJoker, smallJoker }

extension JokerTypeExtension on JokerType {
  String get name {
    switch (this) {
      case JokerType.bigJoker:
        return 'Big Joker';
      case JokerType.smallJoker:
        return 'Small Joker';
    }
  }

  String toJson() => name;

  static JokerType fromJson(String json) {
    switch (json) {
      case 'Big Joker':
        return JokerType.bigJoker;
      case 'Small Joker':
        return JokerType.smallJoker;
      default:
        throw ArgumentError('Invalid joker type: $json');
    }
  }
}

class Card extends Equatable {
  final Suit? suit;
  final Rank? rank;
  final int deckId;
  final bool isJoker;
  final JokerType? jokerType;

  const Card({
    this.suit,
    this.rank,
    required this.deckId,
    this.isJoker = false,
    this.jokerType,
  });

  // Factory constructor for regular cards
  factory Card.regular({
    required Suit suit,
    required Rank rank,
    required int deckId,
  }) {
    return Card(
      suit: suit,
      rank: rank,
      deckId: deckId,
      isJoker: false,
    );
  }

  // Factory constructor for jokers
  factory Card.joker({
    required JokerType jokerType,
    required int deckId,
  }) {
    return Card(
      deckId: deckId,
      isJoker: true,
      jokerType: jokerType,
    );
  }

  // Factory constructor from JSON
  factory Card.fromJson(Map<String, dynamic> json) {
    final isJoker = json['is_joker'] as bool? ?? false;

    if (isJoker) {
      return Card.joker(
        jokerType: JokerTypeExtension.fromJson(json['joker_type'] as String),
        deckId: json['deck_id'] as int,
      );
    } else {
      return Card.regular(
        suit: SuitExtension.fromJson(json['suit'] as String),
        rank: RankExtension.fromJson(json['rank'] as String),
        deckId: json['deck_id'] as int,
      );
    }
  }

  // Convert to JSON
  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{
      'deck_id': deckId,
      'is_joker': isJoker,
    };

    if (isJoker) {
      json['joker_type'] = jokerType!.toJson();
    } else {
      json['suit'] = suit!.toJson();
      json['rank'] = rank!.toJson();
    }

    return json;
  }

  // Get point value for scoring
  int get pointValue {
    if (isJoker) return 0;
    return rank!.pointValue;
  }

  // Check if two cards have the same face value (ignoring deck ID)
  bool isSameFace(Card other) {
    if (isJoker && other.isJoker) {
      return jokerType == other.jokerType;
    }
    if (isJoker || other.isJoker) {
      return false;
    }
    return suit == other.suit && rank == other.rank;
  }

  // Get trump hierarchy value for card comparison
  int getTrumpHierarchy(Suit trumpSuit) {
    if (isJoker) {
      switch (jokerType!) {
        case JokerType.bigJoker:
          return 1000; // Highest trump
        case JokerType.smallJoker:
          return 999; // Second highest trump
      }
    }

    // Trump 2s are permanent trumps
    if (rank == Rank.two) {
      if (suit == trumpSuit) {
        return 998; // Trump suit 2s
      }
      return 997; // Off-suit 2s
    }

    // Trump suit cards (excluding 2s)
    if (suit == trumpSuit) {
      return 900 + rank!.value; // Trump suit cards ranked by rank
    }

    // Off-suit cards have no trump hierarchy
    return 0;
  }

  // Get suit hierarchy value
  int get suitHierarchy {
    if (isJoker) return 0;
    return rank!.value;
  }

  // Check if card is trump
  bool isTrump(Suit trumpSuit) {
    return getTrumpHierarchy(trumpSuit) > 0;
  }

  @override
  String toString() {
    if (isJoker) {
      return '${jokerType!.name} (Deck $deckId)';
    }
    return '${rank!.name} of ${suit!.name} (Deck $deckId)';
  }

  @override
  List<Object?> get props => [suit, rank, deckId, isJoker, jokerType];
}

class Deck extends Equatable {
  final List<Card> cards;

  const Deck({required this.cards});

  // Factory constructor to create a new deck with 108 cards
  factory Deck.newDeck() {
    final cards = <Card>[];

    // Add two standard 52-card decks
    for (int deckId = 1; deckId <= 2; deckId++) {
      // Add regular cards for each suit
      for (final suit in Suit.values) {
        for (final rank in Rank.values) {
          cards.add(Card.regular(suit: suit, rank: rank, deckId: deckId));
        }
      }
    }

    // Add 4 jokers (2 big, 2 small)
    for (int deckId = 1; deckId <= 2; deckId++) {
      cards.add(Card.joker(jokerType: JokerType.bigJoker, deckId: deckId));
      cards.add(Card.joker(jokerType: JokerType.smallJoker, deckId: deckId));
    }

    return Deck(cards: cards);
  }

  // Factory constructor from JSON
  factory Deck.fromJson(Map<String, dynamic> json) {
    final cardsList = json['cards'] as List<dynamic>;
    final cards = cardsList
        .map((cardJson) => Card.fromJson(cardJson as Map<String, dynamic>))
        .toList();
    return Deck(cards: cards);
  }

  // Convert to JSON
  Map<String, dynamic> toJson() {
    return {
      'cards': cards.map((card) => card.toJson()).toList(),
    };
  }

  // Get total point value of all cards
  int get totalPoints {
    return cards.fold(0, (sum, card) => sum + card.pointValue);
  }

  // Get number of remaining cards
  int get remaining => cards.length;

  // Validate deck composition
  bool get isValidComposition {
    if (cards.length != 108) return false;

    // Count cards by type
    final suitCounts = <Suit, Map<Rank, int>>{};
    final jokerCounts = <JokerType, int>{};

    for (final suit in Suit.values) {
      suitCounts[suit] = <Rank, int>{};
      for (final rank in Rank.values) {
        suitCounts[suit]![rank] = 0;
      }
    }

    for (final jokerType in JokerType.values) {
      jokerCounts[jokerType] = 0;
    }

    for (final card in cards) {
      if (card.isJoker) {
        jokerCounts[card.jokerType!] = jokerCounts[card.jokerType!]! + 1;
      } else {
        suitCounts[card.suit!]![card.rank!] =
            suitCounts[card.suit!]![card.rank!]! + 1;
      }
    }

    // Validate each rank appears exactly twice in each suit
    for (final suit in Suit.values) {
      for (final rank in Rank.values) {
        if (suitCounts[suit]![rank] != 2) {
          return false;
        }
      }
    }

    // Validate jokers
    if (jokerCounts[JokerType.bigJoker] != 2) return false;
    if (jokerCounts[JokerType.smallJoker] != 2) return false;

    return true;
  }

  @override
  List<Object?> get props => [cards];
}
