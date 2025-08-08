import 'package:flutter_test/flutter_test.dart';
import 'package:chinese_bridge_game/features/game/domain/entities/user.dart';

void main() {
  group('UserStats', () {
    final now = DateTime.now();

    test('should create UserStats with default values', () {
      final stats = UserStats(
        userId: 'user1',
        createdAt: now,
        updatedAt: now,
      );

      expect(stats.userId, 'user1');
      expect(stats.gamesPlayed, 0);
      expect(stats.gamesWon, 0);
      expect(stats.gamesAsDeclarer, 0);
      expect(stats.declarerWins, 0);
      expect(stats.totalPoints, 0);
      expect(stats.averageBid, 0.0);
    });

    test('should calculate win rate correctly', () {
      final stats = UserStats(
        userId: 'user1',
        gamesPlayed: 10,
        gamesWon: 6,
        createdAt: now,
        updatedAt: now,
      );

      expect(stats.winRate, 0.6);
    });

    test('should return 0 win rate when no games played', () {
      final stats = UserStats(
        userId: 'user1',
        createdAt: now,
        updatedAt: now,
      );

      expect(stats.winRate, 0.0);
    });

    test('should calculate declarer win rate correctly', () {
      final stats = UserStats(
        userId: 'user1',
        gamesAsDeclarer: 5,
        declarerWins: 3,
        createdAt: now,
        updatedAt: now,
      );

      expect(stats.declarerWinRate, 0.6);
    });

    test('should return 0 declarer win rate when no games as declarer', () {
      final stats = UserStats(
        userId: 'user1',
        createdAt: now,
        updatedAt: now,
      );

      expect(stats.declarerWinRate, 0.0);
    });

    test('should calculate average points per game correctly', () {
      final stats = UserStats(
        userId: 'user1',
        gamesPlayed: 4,
        totalPoints: 400,
        createdAt: now,
        updatedAt: now,
      );

      expect(stats.averagePointsPerGame, 100.0);
    });

    test('should return 0 average points when no games played', () {
      final stats = UserStats(
        userId: 'user1',
        createdAt: now,
        updatedAt: now,
      );

      expect(stats.averagePointsPerGame, 0.0);
    });

    test('should serialize and deserialize correctly', () {
      final originalStats = UserStats(
        userId: 'user1',
        gamesPlayed: 10,
        gamesWon: 6,
        gamesAsDeclarer: 3,
        declarerWins: 2,
        totalPoints: 500,
        averageBid: 115.5,
        createdAt: now,
        updatedAt: now,
      );

      final json = originalStats.toJson();
      final deserializedStats = UserStats.fromJson(json);

      expect(deserializedStats, originalStats);
    });

    test('should copy with new values', () {
      final originalStats = UserStats(
        userId: 'user1',
        gamesPlayed: 10,
        createdAt: now,
        updatedAt: now,
      );

      final copiedStats = originalStats.copyWith(
        gamesPlayed: 15,
        gamesWon: 8,
      );

      expect(copiedStats.userId, 'user1');
      expect(copiedStats.gamesPlayed, 15);
      expect(copiedStats.gamesWon, 8);
      expect(copiedStats.createdAt, now);
    });
  });

  group('User', () {
    final now = DateTime.now();
    final stats = UserStats(
      userId: 'user1',
      gamesPlayed: 10,
      gamesWon: 6,
      createdAt: now,
      updatedAt: now,
    );

    test('should create User correctly', () {
      final user = User(
        id: 'user1',
        googleId: 'google123',
        email: 'test@example.com',
        name: 'Test User',
        avatar: 'https://example.com/avatar.jpg',
        createdAt: now,
        updatedAt: now,
        stats: stats,
      );

      expect(user.id, 'user1');
      expect(user.googleId, 'google123');
      expect(user.email, 'test@example.com');
      expect(user.name, 'Test User');
      expect(user.avatar, 'https://example.com/avatar.jpg');
      expect(user.stats, stats);
    });

    test('should return display name correctly', () {
      final userWithName = User(
        id: 'user1',
        googleId: 'google123',
        email: 'test@example.com',
        name: 'Test User',
        createdAt: now,
        updatedAt: now,
      );

      final userWithoutName = User(
        id: 'user2',
        googleId: 'google456',
        email: 'test@example.com',
        name: '',
        createdAt: now,
        updatedAt: now,
      );

      expect(userWithName.displayName, 'Test User');
      expect(userWithoutName.displayName, 'test');
    });

    test('should return correct initials', () {
      final userTwoNames = User(
        id: 'user1',
        googleId: 'google123',
        email: 'test@example.com',
        name: 'John Doe',
        createdAt: now,
        updatedAt: now,
      );

      final userOneName = User(
        id: 'user2',
        googleId: 'google456',
        email: 'test@example.com',
        name: 'John',
        createdAt: now,
        updatedAt: now,
      );

      final userNoName = User(
        id: 'user3',
        googleId: 'google789',
        email: 'test@example.com',
        name: '',
        createdAt: now,
        updatedAt: now,
      );

      expect(userTwoNames.initials, 'JD');
      expect(userOneName.initials, 'J');
      expect(userNoName.initials, 'T'); // From email
    });

    test('should check avatar presence correctly', () {
      final userWithAvatar = User(
        id: 'user1',
        googleId: 'google123',
        email: 'test@example.com',
        name: 'Test User',
        avatar: 'https://example.com/avatar.jpg',
        createdAt: now,
        updatedAt: now,
      );

      final userWithoutAvatar = User(
        id: 'user2',
        googleId: 'google456',
        email: 'test@example.com',
        name: 'Test User',
        createdAt: now,
        updatedAt: now,
      );

      final userWithEmptyAvatar = User(
        id: 'user3',
        googleId: 'google789',
        email: 'test@example.com',
        name: 'Test User',
        avatar: '',
        createdAt: now,
        updatedAt: now,
      );

      expect(userWithAvatar.hasAvatar, true);
      expect(userWithoutAvatar.hasAvatar, false);
      expect(userWithEmptyAvatar.hasAvatar, false);
    });

    test('should serialize and deserialize correctly', () {
      final originalUser = User(
        id: 'user1',
        googleId: 'google123',
        email: 'test@example.com',
        name: 'Test User',
        avatar: 'https://example.com/avatar.jpg',
        createdAt: now,
        updatedAt: now,
        stats: stats,
      );

      final json = originalUser.toJson();
      final deserializedUser = User.fromJson(json);

      expect(deserializedUser, originalUser);
    });

    test('should serialize and deserialize without stats', () {
      final originalUser = User(
        id: 'user1',
        googleId: 'google123',
        email: 'test@example.com',
        name: 'Test User',
        createdAt: now,
        updatedAt: now,
      );

      final json = originalUser.toJson();
      final deserializedUser = User.fromJson(json);

      expect(deserializedUser, originalUser);
      expect(deserializedUser.stats, null);
    });

    test('should copy with new values', () {
      final originalUser = User(
        id: 'user1',
        googleId: 'google123',
        email: 'test@example.com',
        name: 'Test User',
        createdAt: now,
        updatedAt: now,
      );

      final copiedUser = originalUser.copyWith(
        name: 'Updated User',
        avatar: 'https://example.com/new-avatar.jpg',
      );

      expect(copiedUser.id, 'user1');
      expect(copiedUser.name, 'Updated User');
      expect(copiedUser.avatar, 'https://example.com/new-avatar.jpg');
      expect(copiedUser.email, 'test@example.com');
    });

    test('should be equal when all properties match', () {
      final user1 = User(
        id: 'user1',
        googleId: 'google123',
        email: 'test@example.com',
        name: 'Test User',
        createdAt: now,
        updatedAt: now,
      );

      final user2 = User(
        id: 'user1',
        googleId: 'google123',
        email: 'test@example.com',
        name: 'Test User',
        createdAt: now,
        updatedAt: now,
      );

      expect(user1, user2);
    });

    test('should not be equal when properties differ', () {
      final user1 = User(
        id: 'user1',
        googleId: 'google123',
        email: 'test@example.com',
        name: 'Test User',
        createdAt: now,
        updatedAt: now,
      );

      final user2 = User(
        id: 'user2',
        googleId: 'google123',
        email: 'test@example.com',
        name: 'Test User',
        createdAt: now,
        updatedAt: now,
      );

      expect(user1, isNot(user2));
    });
  });
}
