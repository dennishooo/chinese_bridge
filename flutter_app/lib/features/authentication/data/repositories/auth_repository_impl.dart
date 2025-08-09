import 'package:chinese_bridge_game/features/authentication/domain/entities/user.dart';
import 'package:chinese_bridge_game/features/authentication/domain/repositories/auth_repository.dart';
import 'package:chinese_bridge_game/features/authentication/data/datasources/auth_local_data_source.dart';
import 'package:chinese_bridge_game/features/authentication/data/datasources/auth_remote_data_source.dart';
import 'package:chinese_bridge_game/features/authentication/data/models/auth_models.dart';

class AuthRepositoryImpl implements AuthRepository {
  final AuthRemoteDataSource remoteDataSource;
  final AuthLocalDataSource localDataSource;

  AuthRepositoryImpl({
    required this.remoteDataSource,
    required this.localDataSource,
  });

  @override
  Future<User> signInWithGoogle() async {
    try {
      // Get Google auth code
      final authCode = await (remoteDataSource as AuthRemoteDataSourceImpl)
          .getGoogleAuthCode();

      // Exchange auth code for tokens
      final authResponse = await remoteDataSource.signInWithGoogle(authCode);

      // Save auth data locally
      await localDataSource.saveAuthData(authResponse);

      // Convert UserInfo to User entity
      return _userInfoToUser(authResponse.user);
    } catch (e) {
      throw Exception('Sign in failed: $e');
    }
  }

  @override
  Future<void> signOut() async {
    try {
      // Sign out from server
      await remoteDataSource.signOut();
    } finally {
      // Always clear local data, even if server sign out fails
      await localDataSource.clearAuthData();
    }
  }

  @override
  Future<String> refreshToken() async {
    try {
      final refreshToken = await localDataSource.getRefreshToken();
      if (refreshToken == null) {
        throw Exception('No refresh token available');
      }

      final tokenResponse = await remoteDataSource.refreshToken(refreshToken);

      // Save new access token
      await localDataSource.saveAccessToken(tokenResponse.accessToken);

      return tokenResponse.accessToken;
    } catch (e) {
      // If refresh fails, clear auth data and require re-login
      await localDataSource.clearAuthData();
      throw Exception('Token refresh failed: $e');
    }
  }

  @override
  Future<User?> getCurrentUser() async {
    try {
      final userInfo = await localDataSource.getUserInfo();
      if (userInfo == null) return null;

      return _userInfoToUser(userInfo);
    } catch (e) {
      return null;
    }
  }

  @override
  Future<bool> isAuthenticated() async {
    try {
      final isAuth = await localDataSource.isAuthenticated();
      if (!isAuth) return false;

      // Try to get a valid access token
      final accessToken = await getAccessToken();
      return accessToken != null;
    } catch (e) {
      return false;
    }
  }

  @override
  Future<String?> getAccessToken() async {
    try {
      String? accessToken = await localDataSource.getAccessToken();

      // If access token is null or expired, try to refresh
      if (accessToken == null) {
        final refreshToken = await localDataSource.getRefreshToken();
        if (refreshToken != null) {
          try {
            accessToken = await this.refreshToken();
          } catch (e) {
            // Refresh failed, user needs to re-authenticate
            return null;
          }
        }
      }

      return accessToken;
    } catch (e) {
      return null;
    }
  }

  @override
  Future<String?> getRefreshToken() async {
    try {
      return await localDataSource.getRefreshToken();
    } catch (e) {
      return null;
    }
  }

  /// Convert UserInfo data model to User domain entity
  User _userInfoToUser(UserInfo userInfo) {
    return User(
      id: userInfo.id,
      googleId: userInfo.googleId,
      email: userInfo.email,
      name: userInfo.name,
      avatar: userInfo.avatar,
    );
  }
}
