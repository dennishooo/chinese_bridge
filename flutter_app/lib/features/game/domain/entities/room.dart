import 'package:equatable/equatable.dart';
import 'user.dart';

enum RoomStatus { waiting, full, inGame, ended }

extension RoomStatusExtension on RoomStatus {
  String get name {
    switch (this) {
      case RoomStatus.waiting:
        return 'waiting';
      case RoomStatus.full:
        return 'full';
      case RoomStatus.inGame:
        return 'in_game';
      case RoomStatus.ended:
        return 'ended';
    }
  }

  String toJson() => name;

  static RoomStatus fromJson(String json) {
    switch (json) {
      case 'waiting':
        return RoomStatus.waiting;
      case 'full':
        return RoomStatus.full;
      case 'in_game':
        return RoomStatus.inGame;
      case 'ended':
        return RoomStatus.ended;
      default:
        throw ArgumentError('Invalid room status: $json');
    }
  }
}

class RoomParticipant extends Equatable {
  final String roomId;
  final String userId;
  final int position; // 0-3 for seating position
  final DateTime joinedAt;
  final User user;

  const RoomParticipant({
    required this.roomId,
    required this.userId,
    required this.position,
    required this.joinedAt,
    required this.user,
  });

  // Factory constructor from JSON
  factory RoomParticipant.fromJson(Map<String, dynamic> json) {
    return RoomParticipant(
      roomId: json['room_id'] as String,
      userId: json['user_id'] as String,
      position: json['position'] as int,
      joinedAt: DateTime.parse(json['joined_at'] as String),
      user: User.fromJson(json['user'] as Map<String, dynamic>),
    );
  }

  // Convert to JSON
  Map<String, dynamic> toJson() {
    return {
      'room_id': roomId,
      'user_id': userId,
      'position': position,
      'joined_at': joinedAt.toIso8601String(),
      'user': user.toJson(),
    };
  }

  // Copy with new values
  RoomParticipant copyWith({
    String? roomId,
    String? userId,
    int? position,
    DateTime? joinedAt,
    User? user,
  }) {
    return RoomParticipant(
      roomId: roomId ?? this.roomId,
      userId: userId ?? this.userId,
      position: position ?? this.position,
      joinedAt: joinedAt ?? this.joinedAt,
      user: user ?? this.user,
    );
  }

  @override
  List<Object?> get props => [roomId, userId, position, joinedAt, user];
}

class Room extends Equatable {
  final String id;
  final String name;
  final String hostId;
  final int maxPlayers;
  final int currentPlayers;
  final RoomStatus status;
  final DateTime createdAt;
  final DateTime updatedAt;
  final User host;
  final List<RoomParticipant> participants;

  const Room({
    required this.id,
    required this.name,
    required this.hostId,
    this.maxPlayers = 4,
    this.currentPlayers = 0,
    this.status = RoomStatus.waiting,
    required this.createdAt,
    required this.updatedAt,
    required this.host,
    required this.participants,
  });

  // Factory constructor from JSON
  factory Room.fromJson(Map<String, dynamic> json) {
    final participantsList = json['participants'] as List<dynamic>? ?? [];
    final participants = participantsList
        .map((participantJson) =>
            RoomParticipant.fromJson(participantJson as Map<String, dynamic>))
        .toList();

    return Room(
      id: json['id'] as String,
      name: json['name'] as String,
      hostId: json['host_id'] as String,
      maxPlayers: json['max_players'] as int? ?? 4,
      currentPlayers: json['current_players'] as int? ?? 0,
      status:
          RoomStatusExtension.fromJson(json['status'] as String? ?? 'waiting'),
      createdAt: DateTime.parse(json['created_at'] as String),
      updatedAt: DateTime.parse(json['updated_at'] as String),
      host: User.fromJson(json['host'] as Map<String, dynamic>),
      participants: participants,
    );
  }

  // Convert to JSON
  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'host_id': hostId,
      'max_players': maxPlayers,
      'current_players': currentPlayers,
      'status': status.toJson(),
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt.toIso8601String(),
      'host': host.toJson(),
      'participants':
          participants.map((participant) => participant.toJson()).toList(),
    };
  }

  // Check if room is full
  bool get isFull => currentPlayers >= maxPlayers;

  // Check if room can start game
  bool get canStartGame =>
      currentPlayers == maxPlayers && status == RoomStatus.full;

  // Check if user is host
  bool isHost(String userId) => hostId == userId;

  // Check if user is participant
  bool isParticipant(String userId) {
    return participants.any((participant) => participant.userId == userId);
  }

  // Get participant by user ID
  RoomParticipant? getParticipant(String userId) {
    try {
      return participants
          .firstWhere((participant) => participant.userId == userId);
    } catch (e) {
      return null;
    }
  }

  // Get participant by position
  RoomParticipant? getParticipantByPosition(int position) {
    try {
      return participants
          .firstWhere((participant) => participant.position == position);
    } catch (e) {
      return null;
    }
  }

  // Get available positions
  List<int> get availablePositions {
    final occupiedPositions = participants.map((p) => p.position).toSet();
    return List.generate(maxPlayers, (index) => index)
        .where((position) => !occupiedPositions.contains(position))
        .toList();
  }

  // Get next available position
  int? get nextAvailablePosition {
    final available = availablePositions;
    return available.isEmpty ? null : available.first;
  }

  // Copy with new values
  Room copyWith({
    String? id,
    String? name,
    String? hostId,
    int? maxPlayers,
    int? currentPlayers,
    RoomStatus? status,
    DateTime? createdAt,
    DateTime? updatedAt,
    User? host,
    List<RoomParticipant>? participants,
  }) {
    return Room(
      id: id ?? this.id,
      name: name ?? this.name,
      hostId: hostId ?? this.hostId,
      maxPlayers: maxPlayers ?? this.maxPlayers,
      currentPlayers: currentPlayers ?? this.currentPlayers,
      status: status ?? this.status,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
      host: host ?? this.host,
      participants: participants ?? this.participants,
    );
  }

  @override
  List<Object?> get props => [
        id,
        name,
        hostId,
        maxPlayers,
        currentPlayers,
        status,
        createdAt,
        updatedAt,
        host,
        participants,
      ];
}
