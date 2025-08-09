import 'package:dio/dio.dart';
import 'package:google_sign_in/google_sign_in.dart';
import 'package:chinese_bridge_game/features/authentication/data/models/auth_models.dart';
import 'package:chinese_bridge_game/core/constants/api_constants.dart';

abstract class AuthRemoteDataSource {
  Future<AuthResponse> signInWithGoogle(String authCode);
  Future<TokenResponse> refreshToken(String refreshToken);
  Future<void> signOut();
}

class AuthRemoteDataSourceImpl implements AuthRemoteDataSource {
  final Dio dio;
  final GoogleSignIn googleSignIn;

  AuthRemoteDataSourceImpl({
    required this.dio,
    required this.googleSignIn,
  });

  @override
  Future<AuthResponse> signInWithGoogle(String authCode) async {
    try {
      final response = await dio.post(
        '${ApiConstants.baseUrl}/auth/google',
        data: GoogleOAuthRequest(code: authCode).toJson(),
      );

      if (response.statusCode == 200) {
        return AuthResponse.fromJson(response.data);
      } else {
        throw Exception('Failed to authenticate with Google');
      }
    } on DioException catch (e) {
      if (e.response?.data != null) {
        final errorResponse = ErrorResponse.fromJson(e.response!.data);
        throw Exception('Authentication failed: ${errorResponse.message}');
      }
      throw Exception('Network error: ${e.message}');
    } catch (e) {
      throw Exception('Unexpected error: $e');
    }
  }

  @override
  Future<TokenResponse> refreshToken(String refreshToken) async {
    try {
      final response = await dio.post(
        '${ApiConstants.baseUrl}/auth/refresh',
        data: RefreshTokenRequest(refreshToken: refreshToken).toJson(),
      );

      if (response.statusCode == 200) {
        return TokenResponse.fromJson(response.data);
      } else {
        throw Exception('Failed to refresh token');
      }
    } on DioException catch (e) {
      if (e.response?.data != null) {
        final errorResponse = ErrorResponse.fromJson(e.response!.data);
        throw Exception('Token refresh failed: ${errorResponse.message}');
      }
      throw Exception('Network error: ${e.message}');
    } catch (e) {
      throw Exception('Unexpected error: $e');
    }
  }

  @override
  Future<void> signOut() async {
    try {
      await dio.post('${ApiConstants.baseUrl}/auth/logout');
      await googleSignIn.signOut();
    } on DioException catch (e) {
      // Log the error but don't throw - we still want to sign out locally
      print('Server sign out failed: ${e.message}');
      await googleSignIn.signOut();
    } catch (e) {
      print('Sign out error: $e');
      await googleSignIn.signOut();
    }
  }

  /// Get Google OAuth authorization code
  Future<String> getGoogleAuthCode() async {
    try {
      final GoogleSignInAccount? account = await googleSignIn.signIn();
      if (account == null) {
        throw Exception('Google sign in was cancelled');
      }

      final GoogleSignInAuthentication auth = await account.authentication;
      if (auth.serverAuthCode == null) {
        throw Exception('Failed to get server auth code');
      }

      return auth.serverAuthCode!;
    } catch (e) {
      throw Exception('Google sign in failed: $e');
    }
  }
}
