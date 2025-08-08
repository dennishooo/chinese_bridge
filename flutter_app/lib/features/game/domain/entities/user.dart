import 'package:equatable/equatable.dart';

class UserStats extends Equatable {
  final String userId;
  final int gamesPlayed;
  final int gamesWon;
  final int gamesAsDeclarer;
  final int declarerWins;
  final int totalPoints;
  final double averageBid;
  final DateTime createdAt;
  final DateTime updatedAt;

  const UserStats({
    required this.userId,
    this.gamesPlayed = 0,
    this.gamesWon = 0,
    this.gamesAsDeclarer = 0,
    this.declarerWins = 0,
    this.totalPoints = 0,
    this.averageBid = 0.0,
    required this.createdAt,
    required this.updatedAt,
  });

  // Factory constructor from JSON
  factory UserStats.fromJson(Map<String, dynamic> json) {
    return UserStats(
      userId: json['user_id'] as String,
      gamesPlayed: json['games_played'] as int? ?? 0,
      gamesWon: json['games_won'] as int? ?? 0,
      gamesAsDeclarer: json['games_as_declarer'] as int? ?? 0,
      declarerWins: json['declarer_wins'] as int? ?? 0,
      totalPoints: json['total_points'] as int? ?? 0,
      averageBid: (json['average_bid'] as num?)?.toDouble() ?? 0.0,
      createdAt: DateTime.parse(json['created_at'] as String),
      updatedAt: DateTime.parse(json['updated_at'] as String),
    );
  }

  // Convert to JSON
  Map<String, dynamic> toJson() {
    return {
      'user_id': userId,
      'games_played': gamesPlayed,
      'games_won': gamesWon,
      'games_as_declarer': gamesAsDeclarer,
      'declarer_wins': declarerWins,
      'total_points': totalPoints,
      'average_bid': averageBid,
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt.toIso8601String(),
    };
  }

  // Calculate win rate
  double get winRate {
    if (gamesPlayed == 0) return 0.0;
    return gamesWon / gamesPlayed;
  }

  // Calculate declarer win rate
  double get declarerWinRate {
    if (gamesAsDeclarer == 0) return 0.0;
    return declarerWins / gamesAsDeclarer;
  }

  // Calculate average points per game
  double get averagePointsPerGame {
    if (gamesPlayed == 0) return 0.0;
    return totalPoints / gamesPlayed;
  }

  // Copy with new values
  UserStats copyWith({
    String? userId,
    int? gamesPlayed,
    int? gamesWon,
    int? gamesAsDeclarer,
    int? declarerWins,
    int? totalPoints,
    double? averageBid,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) {
    return UserStats(
      userId: userId ?? this.userId,
      gamesPlayed: gamesPlayed ?? this.gamesPlayed,
      gamesWon: gamesWon ?? this.gamesWon,
      gamesAsDeclarer: gamesAsDeclarer ?? this.gamesAsDeclarer,
      declarerWins: declarerWins ?? this.declarerWins,
      totalPoints: totalPoints ?? this.totalPoints,
      averageBid: averageBid ?? this.averageBid,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }

  @override
  List<Object?> get props => [
        userId,
        gamesPlayed,
        gamesWon,
        gamesAsDeclarer,
        declarerWins,
        totalPoints,
        averageBid,
        createdAt,
        updatedAt,
      ];
}

class User extends Equatable {
  final String id;
  final String googleId;
  final String email;
  final String name;
  final String? avatar;
  final DateTime createdAt;
  final DateTime updatedAt;
  final UserStats? stats;

  const User({
    required this.id,
    required this.googleId,
    required this.email,
    required this.name,
    this.avatar,
    required this.createdAt,
    required this.updatedAt,
    this.stats,
  });

  // Factory constructor from JSON
  factory User.fromJson(Map<String, dynamic> json) {
    return User(
      id: json['id'] as String,
      googleId: json['google_id'] as String,
      email: json['email'] as String,
      name: json['name'] as String,
      avatar: json['avatar'] as String?,
      createdAt: DateTime.parse(json['created_at'] as String),
      updatedAt: DateTime.parse(json['updated_at'] as String),
      stats: json['stats'] != null
          ? UserStats.fromJson(json['stats'] as Map<String, dynamic>)
          : null,
    );
  }

  // Convert to JSON
  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'google_id': googleId,
      'email': email,
      'name': name,
      'avatar': avatar,
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt.toIso8601String(),
      'stats': stats?.toJson(),
    };
  }

  // Get display name (fallback to email if name is empty)
  String get displayName {
    return name.isNotEmpty ? name : email.split('@').first;
  }

  // Get initials for avatar
  String get initials {
    final nameParts = displayName.split(' ');
    if (nameParts.length >= 2) {
      return '${nameParts[0][0]}${nameParts[1][0]}'.toUpperCase();
    } else if (nameParts.isNotEmpty) {
      return nameParts[0][0].toUpperCase();
    }
    return 'U';
  }

  // Check if user has avatar
  bool get hasAvatar => avatar != null && avatar!.isNotEmpty;

  // Copy with new values
  User copyWith({
    String? id,
    String? googleId,
    String? email,
    String? name,
    String? avatar,
    DateTime? createdAt,
    DateTime? updatedAt,
    UserStats? stats,
  }) {
    return User(
      id: id ?? this.id,
      googleId: googleId ?? this.googleId,
      email: email ?? this.email,
      name: name ?? this.name,
      avatar: avatar ?? this.avatar,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
      stats: stats ?? this.stats,
    );
  }

  @override
  List<Object?> get props => [
        id,
        googleId,
        email,
        name,
        avatar,
        createdAt,
        updatedAt,
        stats,
      ];
}
