import 'dart:convert';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:chinese_bridge_game/features/authentication/data/models/auth_models.dart';

abstract class AuthLocalDataSource {
  Future<void> saveAuthData(AuthResponse authResponse);
  Future<void> saveAccessToken(String accessToken);
  Future<String?> getAccessToken();
  Future<String?> getRefreshToken();
  Future<UserInfo?> getUserInfo();
  Future<void> clearAuthData();
  Future<bool> isAuthenticated();
}

class AuthLocalDataSourceImpl implements AuthLocalDataSource {
  final SharedPreferences sharedPreferences;

  static const String _accessTokenKey = 'access_token';
  static const String _refreshTokenKey = 'refresh_token';
  static const String _userInfoKey = 'user_info';
  static const String _tokenExpiryKey = 'token_expiry';

  AuthLocalDataSourceImpl({required this.sharedPreferences});

  @override
  Future<void> saveAuthData(AuthResponse authResponse) async {
    final expiryTime = DateTime.now()
        .add(Duration(seconds: authResponse.expiresIn))
        .millisecondsSinceEpoch;

    await Future.wait([
      sharedPreferences.setString(_accessTokenKey, authResponse.accessToken),
      sharedPreferences.setString(_refreshTokenKey, authResponse.refreshToken),
      sharedPreferences.setString(
          _userInfoKey, jsonEncode(authResponse.user.toJson())),
      sharedPreferences.setInt(_tokenExpiryKey, expiryTime),
    ]);
  }

  @override
  Future<void> saveAccessToken(String accessToken) async {
    await sharedPreferences.setString(_accessTokenKey, accessToken);
  }

  @override
  Future<String?> getAccessToken() async {
    final token = sharedPreferences.getString(_accessTokenKey);
    if (token == null) return null;

    // Check if token is expired
    final expiryTime = sharedPreferences.getInt(_tokenExpiryKey);
    if (expiryTime != null &&
        DateTime.now().millisecondsSinceEpoch >= expiryTime) {
      return null; // Token is expired
    }

    return token;
  }

  @override
  Future<String?> getRefreshToken() async {
    return sharedPreferences.getString(_refreshTokenKey);
  }

  @override
  Future<UserInfo?> getUserInfo() async {
    final userInfoJson = sharedPreferences.getString(_userInfoKey);
    if (userInfoJson == null) return null;

    try {
      final userInfoMap = jsonDecode(userInfoJson) as Map<String, dynamic>;
      return UserInfo.fromJson(userInfoMap);
    } catch (e) {
      // If parsing fails, clear the corrupted data
      await clearAuthData();
      return null;
    }
  }

  @override
  Future<void> clearAuthData() async {
    await Future.wait([
      sharedPreferences.remove(_accessTokenKey),
      sharedPreferences.remove(_refreshTokenKey),
      sharedPreferences.remove(_userInfoKey),
      sharedPreferences.remove(_tokenExpiryKey),
    ]);
  }

  @override
  Future<bool> isAuthenticated() async {
    final accessToken = await getAccessToken();
    final refreshToken = await getRefreshToken();
    return accessToken != null || refreshToken != null;
  }
}
