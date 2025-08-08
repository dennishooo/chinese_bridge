import 'package:get_it/get_it.dart';
import 'package:injectable/injectable.dart';
import 'package:dio/dio.dart';
import 'package:shared_preferences/shared_preferences.dart';

import 'injection_container.config.dart';

final getIt = GetIt.instance;

@InjectableInit()
Future<void> configureDependencies() async {
  // Register external dependencies
  final sharedPreferences = await SharedPreferences.getInstance();
  getIt.registerSingleton<SharedPreferences>(sharedPreferences);

  final dio = Dio();
  dio.options.baseUrl = 'http://localhost:8080/api/v1';
  dio.options.connectTimeout = const Duration(seconds: 5);
  dio.options.receiveTimeout = const Duration(seconds: 3);
  getIt.registerSingleton<Dio>(dio);

  getIt.init();
}

@module
abstract class RegisterModule {
  @singleton
  Dio get dio => getIt<Dio>();

  @singleton
  SharedPreferences get sharedPreferences => getIt<SharedPreferences>();
}
