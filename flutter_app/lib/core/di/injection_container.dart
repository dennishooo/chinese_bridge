import 'package:dio/dio.dart';
import 'package:get_it/get_it.dart';
import 'package:google_sign_in/google_sign_in.dart';
import 'package:injectable/injectable.dart';
import 'package:shared_preferences/shared_preferences.dart';

import 'package:chinese_bridge_game/features/authentication/data/datasources/auth_local_data_source.dart';
import 'package:chinese_bridge_game/features/authentication/data/datasources/auth_remote_data_source.dart';
import 'package:chinese_bridge_game/features/authentication/data/repositories/auth_repository_impl.dart';
import 'package:chinese_bridge_game/features/authentication/domain/repositories/auth_repository.dart';
import 'package:chinese_bridge_game/features/authentication/presentation/bloc/auth_bloc.dart';
import 'package:chinese_bridge_game/core/constants/api_constants.dart';

import 'injection_container.config.dart';

final getIt = GetIt.instance;

@InjectableInit()
Future<void> configureDependencies() async {
  // Register external dependencies
  final sharedPreferences = await SharedPreferences.getInstance();
  getIt.registerSingleton<SharedPreferences>(sharedPreferences);

  final dio = Dio();
  dio.options.baseURL = ApiConstants.baseUrl;
  dio.options.connectTimeout = const Duration(seconds: 30);
  dio.options.receiveTimeout = const Duration(seconds: 30);

  // Add interceptors for logging and authentication
  dio.interceptors.add(LogInterceptor(
    requestBody: true,
    responseBody: true,
  ));

  getIt.registerSingleton<Dio>(dio);

  final googleSignIn = GoogleSignIn(
    scopes: ['email', 'profile'],
    // serverClientId should be configured in production
  );
  getIt.registerSingleton<GoogleSignIn>(googleSignIn);

  // Register authentication dependencies
  getIt.registerLazySingleton<AuthLocalDataSource>(
    () => AuthLocalDataSourceImpl(sharedPreferences: getIt()),
  );

  getIt.registerLazySingleton<AuthRemoteDataSource>(
    () => AuthRemoteDataSourceImpl(
      dio: getIt(),
      googleSignIn: getIt(),
    ),
  );

  getIt.registerLazySingleton<AuthRepository>(
    () => AuthRepositoryImpl(
      remoteDataSource: getIt(),
      localDataSource: getIt(),
    ),
  );

  getIt.registerFactory<AuthBloc>(
    () => AuthBloc(authRepository: getIt()),
  );

  getIt.init();
}

@module
abstract class RegisterModule {
  @singleton
  Dio get dio => getIt<Dio>();

  @singleton
  SharedPreferences get sharedPreferences => getIt<SharedPreferences>();

  @singleton
  GoogleSignIn get googleSignIn => getIt<GoogleSignIn>();
}
