import 'package:equatable/equatable.dart';
import 'package:json_annotation/json_annotation.dart';

part 'auth_models.g.dart';

@JsonSerializable()
class AuthResponse extends Equatable {
  @JsonKey(name: 'access_token')
  final String accessToken;

  @JsonKey(name: 'refresh_token')
  final String refreshToken;

  @JsonKey(name: 'token_type')
  final String tokenType;

  @JsonKey(name: 'expires_in')
  final int expiresIn;

  final UserInfo user;

  const AuthResponse({
    required this.accessToken,
    required this.refreshToken,
    required this.tokenType,
    required this.expiresIn,
    required this.user,
  });

  factory AuthResponse.fromJson(Map<String, dynamic> json) =>
      _$AuthResponseFromJson(json);

  Map<String, dynamic> toJson() => _$AuthResponseToJson(this);

  @override
  List<Object> get props =>
      [accessToken, refreshToken, tokenType, expiresIn, user];
}

@JsonSerializable()
class UserInfo extends Equatable {
  final String id;

  @JsonKey(name: 'google_id')
  final String googleId;

  final String email;
  final String name;
  final String avatar;

  const UserInfo({
    required this.id,
    required this.googleId,
    required this.email,
    required this.name,
    required this.avatar,
  });

  factory UserInfo.fromJson(Map<String, dynamic> json) =>
      _$UserInfoFromJson(json);

  Map<String, dynamic> toJson() => _$UserInfoToJson(this);

  @override
  List<Object> get props => [id, googleId, email, name, avatar];
}

@JsonSerializable()
class TokenResponse extends Equatable {
  @JsonKey(name: 'access_token')
  final String accessToken;

  @JsonKey(name: 'token_type')
  final String tokenType;

  @JsonKey(name: 'expires_in')
  final int expiresIn;

  const TokenResponse({
    required this.accessToken,
    required this.tokenType,
    required this.expiresIn,
  });

  factory TokenResponse.fromJson(Map<String, dynamic> json) =>
      _$TokenResponseFromJson(json);

  Map<String, dynamic> toJson() => _$TokenResponseToJson(this);

  @override
  List<Object> get props => [accessToken, tokenType, expiresIn];
}

@JsonSerializable()
class GoogleOAuthRequest extends Equatable {
  final String code;
  final String? state;

  const GoogleOAuthRequest({
    required this.code,
    this.state,
  });

  factory GoogleOAuthRequest.fromJson(Map<String, dynamic> json) =>
      _$GoogleOAuthRequestFromJson(json);

  Map<String, dynamic> toJson() => _$GoogleOAuthRequestToJson(this);

  @override
  List<Object?> get props => [code, state];
}

@JsonSerializable()
class RefreshTokenRequest extends Equatable {
  @JsonKey(name: 'refresh_token')
  final String refreshToken;

  const RefreshTokenRequest({
    required this.refreshToken,
  });

  factory RefreshTokenRequest.fromJson(Map<String, dynamic> json) =>
      _$RefreshTokenRequestFromJson(json);

  Map<String, dynamic> toJson() => _$RefreshTokenRequestToJson(this);

  @override
  List<Object> get props => [refreshToken];
}

@JsonSerializable()
class ErrorResponse extends Equatable {
  final String code;
  final String message;
  final String? details;

  @JsonKey(name: 'trace_id')
  final String traceId;

  const ErrorResponse({
    required this.code,
    required this.message,
    this.details,
    required this.traceId,
  });

  factory ErrorResponse.fromJson(Map<String, dynamic> json) =>
      _$ErrorResponseFromJson(json);

  Map<String, dynamic> toJson() => _$ErrorResponseToJson(this);

  @override
  List<Object?> get props => [code, message, details, traceId];
}
