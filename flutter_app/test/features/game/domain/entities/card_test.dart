import 'package:flutter_test/flutter_test.dart';
import 'package:chinese_bridge_game/features/game/domain/entities/card.dart';

void main() {
  group('Card', () {
    group('regular card', () {
      test('should create regular card correctly', () {
        final card = Card.regular(
          suit: Suit.hearts,
          rank: Rank.king,
          deckId: 1,
        );

        expect(card.suit, Suit.hearts);
        expect(card.rank, Rank.king);
        expect(card.deckId, 1);
        expect(card.isJoker, false);
        expect(card.jokerType, null);
      });

      test('should calculate point values correctly', () {
        final king =
            Card.regular(suit: Suit.hearts, rank: Rank.king, deckId: 1);
        final ten = Card.regular(suit: Suit.hearts, rank: Rank.ten, deckId: 1);
        final five =
            Card.regular(suit: Suit.hearts, rank: Rank.five, deckId: 1);
        final ace = Card.regular(suit: Suit.hearts, rank: Rank.ace, deckId: 1);

        expect(king.pointValue, 10);
        expect(ten.pointValue, 10);
        expect(five.pointValue, 5);
        expect(ace.pointValue, 0);
      });

      test('should check same face correctly', () {
        final card1 =
            Card.regular(suit: Suit.hearts, rank: Rank.king, deckId: 1);
        final card2 =
            Card.regular(suit: Suit.hearts, rank: Rank.king, deckId: 2);
        final card3 =
            Card.regular(suit: Suit.spades, rank: Rank.king, deckId: 1);

        expect(card1.isSameFace(card2), true);
        expect(card1.isSameFace(card3), false);
      });

      test('should calculate trump hierarchy correctly', () {
        const trumpSuit = Suit.hearts;
        final trumpKing =
            Card.regular(suit: Suit.hearts, rank: Rank.king, deckId: 1);
        final offSuitKing =
            Card.regular(suit: Suit.spades, rank: Rank.king, deckId: 1);
        final trumpTwo =
            Card.regular(suit: Suit.hearts, rank: Rank.two, deckId: 1);
        final offSuitTwo =
            Card.regular(suit: Suit.spades, rank: Rank.two, deckId: 1);

        expect(trumpKing.getTrumpHierarchy(trumpSuit), 900 + Rank.king.value);
        expect(offSuitKing.getTrumpHierarchy(trumpSuit), 0);
        expect(trumpTwo.getTrumpHierarchy(trumpSuit), 998);
        expect(offSuitTwo.getTrumpHierarchy(trumpSuit), 997);
      });
    });

    group('joker card', () {
      test('should create joker card correctly', () {
        final joker = Card.joker(
          jokerType: JokerType.bigJoker,
          deckId: 1,
        );

        expect(joker.suit, null);
        expect(joker.rank, null);
        expect(joker.deckId, 1);
        expect(joker.isJoker, true);
        expect(joker.jokerType, JokerType.bigJoker);
      });

      test('should have zero point value', () {
        final bigJoker = Card.joker(jokerType: JokerType.bigJoker, deckId: 1);
        final smallJoker =
            Card.joker(jokerType: JokerType.smallJoker, deckId: 1);

        expect(bigJoker.pointValue, 0);
        expect(smallJoker.pointValue, 0);
      });

      test('should check same face correctly', () {
        final joker1 = Card.joker(jokerType: JokerType.bigJoker, deckId: 1);
        final joker2 = Card.joker(jokerType: JokerType.bigJoker, deckId: 2);
        final joker3 = Card.joker(jokerType: JokerType.smallJoker, deckId: 1);

        expect(joker1.isSameFace(joker2), true);
        expect(joker1.isSameFace(joker3), false);
      });

      test('should have highest trump hierarchy', () {
        const trumpSuit = Suit.hearts;
        final bigJoker = Card.joker(jokerType: JokerType.bigJoker, deckId: 1);
        final smallJoker =
            Card.joker(jokerType: JokerType.smallJoker, deckId: 1);

        expect(bigJoker.getTrumpHierarchy(trumpSuit), 1000);
        expect(smallJoker.getTrumpHierarchy(trumpSuit), 999);
      });
    });

    group('JSON serialization', () {
      test('should serialize and deserialize regular card', () {
        final originalCard = Card.regular(
          suit: Suit.hearts,
          rank: Rank.king,
          deckId: 1,
        );

        final json = originalCard.toJson();
        final deserializedCard = Card.fromJson(json);

        expect(deserializedCard, originalCard);
      });

      test('should serialize and deserialize joker card', () {
        final originalJoker = Card.joker(
          jokerType: JokerType.bigJoker,
          deckId: 1,
        );

        final json = originalJoker.toJson();
        final deserializedJoker = Card.fromJson(json);

        expect(deserializedJoker, originalJoker);
      });
    });

    group('equality', () {
      test('should be equal when all properties match', () {
        final card1 =
            Card.regular(suit: Suit.hearts, rank: Rank.king, deckId: 1);
        final card2 =
            Card.regular(suit: Suit.hearts, rank: Rank.king, deckId: 1);

        expect(card1, card2);
      });

      test('should not be equal when properties differ', () {
        final card1 =
            Card.regular(suit: Suit.hearts, rank: Rank.king, deckId: 1);
        final card2 =
            Card.regular(suit: Suit.hearts, rank: Rank.king, deckId: 2);

        expect(card1, isNot(card2));
      });
    });
  });

  group('Deck', () {
    test('should create new deck with 108 cards', () {
      final deck = Deck.newDeck();

      expect(deck.cards.length, 108);
      expect(deck.remaining, 108);
    });

    test('should have valid composition', () {
      final deck = Deck.newDeck();

      expect(deck.isValidComposition, true);
    });

    test('should have correct total points (200)', () {
      final deck = Deck.newDeck();

      expect(deck.totalPoints, 200);
    });

    test('should serialize and deserialize correctly', () {
      final originalDeck = Deck.newDeck();

      final json = originalDeck.toJson();
      final deserializedDeck = Deck.fromJson(json);

      expect(deserializedDeck.cards.length, originalDeck.cards.length);
      expect(deserializedDeck.totalPoints, originalDeck.totalPoints);
    });
  });

  group('Enums', () {
    group('Suit', () {
      test('should convert to and from JSON correctly', () {
        for (final suit in Suit.values) {
          final json = suit.toJson();
          final fromJson = SuitExtension.fromJson(json);
          expect(fromJson, suit);
        }
      });

      test('should have correct names', () {
        expect(Suit.spades.name, 'Spades');
        expect(Suit.hearts.name, 'Hearts');
        expect(Suit.clubs.name, 'Clubs');
        expect(Suit.diamonds.name, 'Diamonds');
      });
    });

    group('Rank', () {
      test('should convert to and from JSON correctly', () {
        for (final rank in Rank.values) {
          final json = rank.toJson();
          final fromJson = RankExtension.fromJson(json);
          expect(fromJson, rank);
        }
      });

      test('should have correct names', () {
        expect(Rank.two.name, '2');
        expect(Rank.jack.name, 'J');
        expect(Rank.queen.name, 'Q');
        expect(Rank.king.name, 'K');
        expect(Rank.ace.name, 'A');
      });

      test('should have correct values', () {
        expect(Rank.two.value, 2);
        expect(Rank.three.value, 3);
        expect(Rank.ace.value, 14);
      });
    });

    group('JokerType', () {
      test('should convert to and from JSON correctly', () {
        for (final jokerType in JokerType.values) {
          final json = jokerType.toJson();
          final fromJson = JokerTypeExtension.fromJson(json);
          expect(fromJson, jokerType);
        }
      });

      test('should have correct names', () {
        expect(JokerType.bigJoker.name, 'Big Joker');
        expect(JokerType.smallJoker.name, 'Small Joker');
      });
    });
  });
}
