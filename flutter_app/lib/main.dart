import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:chinese_bridge_game/core/di/injection_container.dart';
import 'package:chinese_bridge_game/features/authentication/presentation/bloc/auth_bloc.dart';
import 'package:chinese_bridge_game/features/authentication/presentation/pages/login_page.dart';
import 'package:chinese_bridge_game/features/game/presentation/pages/game_lobby_page.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();

  // Initialize dependency injection
  await configureDependencies();

  runApp(const ChineseBridgeApp());
}

class ChineseBridgeApp extends StatelessWidget {
  const ChineseBridgeApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MultiBlocProvider(
      providers: [
        BlocProvider<AuthBloc>(
          create: (context) => getIt<AuthBloc>()..add(AuthStarted()),
        ),
      ],
      child: MaterialApp(
        title: 'Chinese Bridge Game',
        theme: ThemeData(
          primarySwatch: Colors.blue,
          useMaterial3: true,
        ),
        home: const AuthWrapper(),
      ),
    );
  }
}

class AuthWrapper extends StatelessWidget {
  const AuthWrapper({super.key});

  @override
  Widget build(BuildContext context) {
    return BlocBuilder<AuthBloc, AuthState>(
      builder: (context, state) {
        if (state is AuthLoading) {
          return const Scaffold(
            body: Center(
              child: CircularProgressIndicator(),
            ),
          );
        } else if (state is AuthAuthenticated) {
          return const GameLobbyPage();
        } else {
          return const LoginPage();
        }
      },
    );
  }
}
