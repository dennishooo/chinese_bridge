import 'package:equatable/equatable.dart';

class User extends Equatable {
  final String id;
  final String googleId;
  final String email;
  final String name;
  final String avatar;

  const User({
    required this.id,
    required this.googleId,
    required this.email,
    required this.name,
    required this.avatar,
  });

  @override
  List<Object> get props => [id, googleId, email, name, avatar];

  User copyWith({
    String? id,
    String? googleId,
    String? email,
    String? name,
    String? avatar,
  }) {
    return User(
      id: id ?? this.id,
      googleId: googleId ?? this.googleId,
      email: email ?? this.email,
      name: name ?? this.name,
      avatar: avatar ?? this.avatar,
    );
  }
}
