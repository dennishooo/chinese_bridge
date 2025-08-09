import 'package:bloc_test/bloc_test.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/annotations.dart';
import 'package:mockito/mockito.dart';

import 'package:chinese_bridge_game/features/authentication/domain/entities/user.dart';
import 'package:chinese_bridge_game/features/authentication/domain/repositories/auth_repository.dart';
import 'package:chinese_bridge_game/features/authentication/presentation/bloc/auth_bloc.dart';

import 'auth_bloc_test.mocks.dart';

@GenerateMocks([AuthRepository])
void main() {
  group('AuthBloc', () {
    late MockAuthRepository mockAuthRepository;
    late AuthBloc authBloc;

    setUp(() {
      mockAuthRepository = MockAuthRepository();
      authBloc = AuthBloc(authRepository: mockAuthRepository);
    });

    tearDown(() {
      authBloc.close();
    });

    test('initial state is AuthInitial', () {
      expect(authBloc.state, AuthInitial());
    });

    group('AuthStarted', () {
      const tUser = User(
        id: '1',
        googleId: 'google123',
        email: 'test@example.com',
        name: 'Test User',
        avatar: 'https://example.com/avatar.jpg',
      );

      blocTest<AuthBloc, AuthState>(
        'emits [AuthLoading, AuthAuthenticated] when user is authenticated',
        build: () {
          when(mockAuthRepository.isAuthenticated())
              .thenAnswer((_) async => true);
          when(mockAuthRepository.getCurrentUser())
              .thenAnswer((_) async => tUser);
          return authBloc;
        },
        act: (bloc) => bloc.add(AuthStarted()),
        expect: () => [
          AuthLoading(),
          const AuthAuthenticated(user: tUser),
        ],
      );

      blocTest<AuthBloc, AuthState>(
        'emits [AuthLoading, AuthUnauthenticated] when user is not authenticated',
        build: () {
          when(mockAuthRepository.isAuthenticated())
              .thenAnswer((_) async => false);
          return authBloc;
        },
        act: (bloc) => bloc.add(AuthStarted()),
        expect: () => [
          AuthLoading(),
          AuthUnauthenticated(),
        ],
      );

      blocTest<AuthBloc, AuthState>(
        'emits [AuthLoading, AuthError] when checking authentication fails',
        build: () {
          when(mockAuthRepository.isAuthenticated())
              .thenThrow(Exception('Network error'));
          return authBloc;
        },
        act: (bloc) => bloc.add(AuthStarted()),
        expect: () => [
          AuthLoading(),
          const AuthError(
              message:
                  'Failed to check authentication status: Exception: Network error'),
        ],
      );
    });

    group('AuthGoogleSignInRequested', () {
      const tUser = User(
        id: '1',
        googleId: 'google123',
        email: 'test@example.com',
        name: 'Test User',
        avatar: 'https://example.com/avatar.jpg',
      );

      blocTest<AuthBloc, AuthState>(
        'emits [AuthLoading, AuthAuthenticated] when Google sign in succeeds',
        build: () {
          when(mockAuthRepository.signInWithGoogle())
              .thenAnswer((_) async => tUser);
          return authBloc;
        },
        act: (bloc) => bloc.add(AuthGoogleSignInRequested()),
        expect: () => [
          AuthLoading(),
          const AuthAuthenticated(user: tUser),
        ],
      );

      blocTest<AuthBloc, AuthState>(
        'emits [AuthLoading, AuthError] when Google sign in fails',
        build: () {
          when(mockAuthRepository.signInWithGoogle())
              .thenThrow(Exception('Sign in failed'));
          return authBloc;
        },
        act: (bloc) => bloc.add(AuthGoogleSignInRequested()),
        expect: () => [
          AuthLoading(),
          const AuthError(
              message: 'Google sign in failed: Exception: Sign in failed'),
        ],
      );
    });

    group('AuthSignOutRequested', () {
      blocTest<AuthBloc, AuthState>(
        'emits [AuthLoading, AuthUnauthenticated] when sign out succeeds',
        build: () {
          when(mockAuthRepository.signOut()).thenAnswer((_) async {});
          return authBloc;
        },
        act: (bloc) => bloc.add(AuthSignOutRequested()),
        expect: () => [
          AuthLoading(),
          AuthUnauthenticated(),
        ],
      );

      blocTest<AuthBloc, AuthState>(
        'emits [AuthLoading, AuthUnauthenticated] even when sign out fails',
        build: () {
          when(mockAuthRepository.signOut())
              .thenThrow(Exception('Sign out failed'));
          return authBloc;
        },
        act: (bloc) => bloc.add(AuthSignOutRequested()),
        expect: () => [
          AuthLoading(),
          AuthUnauthenticated(),
        ],
      );
    });

    group('AuthTokenRefreshRequested', () {
      const tUser = User(
        id: '1',
        googleId: 'google123',
        email: 'test@example.com',
        name: 'Test User',
        avatar: 'https://example.com/avatar.jpg',
      );

      blocTest<AuthBloc, AuthState>(
        'emits [AuthAuthenticated] when token refresh succeeds',
        build: () {
          when(mockAuthRepository.refreshToken())
              .thenAnswer((_) async => 'new_access_token');
          when(mockAuthRepository.getCurrentUser())
              .thenAnswer((_) async => tUser);
          return authBloc;
        },
        act: (bloc) => bloc.add(AuthTokenRefreshRequested()),
        expect: () => [
          const AuthAuthenticated(user: tUser),
        ],
      );

      blocTest<AuthBloc, AuthState>(
        'emits [AuthError] when token refresh fails',
        build: () {
          when(mockAuthRepository.refreshToken())
              .thenThrow(Exception('Token refresh failed'));
          return authBloc;
        },
        act: (bloc) => bloc.add(AuthTokenRefreshRequested()),
        expect: () => [
          const AuthError(
              message: 'Token refresh failed: Exception: Token refresh failed'),
        ],
      );
    });
  });
}
