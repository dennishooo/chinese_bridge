import 'package:chinese_bridge_game/features/authentication/domain/entities/user.dart';

abstract class AuthRepository {
  /// Sign in with Google OAuth
  Future<User> signInWithGoogle();

  /// Sign out the current user
  Future<void> signOut();

  /// Refresh the access token
  Future<String> refreshToken();

  /// Get the current user if authenticated
  Future<User?> getCurrentUser();

  /// Check if user is currently authenticated
  Future<bool> isAuthenticated();

  /// Get the current access token
  Future<String?> getAccessToken();

  /// Get the current refresh token
  Future<String?> getRefreshToken();
}
