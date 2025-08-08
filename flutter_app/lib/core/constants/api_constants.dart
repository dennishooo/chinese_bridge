class ApiConstants {
  static const String baseUrl = 'http://localhost:8080/api/v1';
  static const String authBaseUrl = 'http://localhost:8080/api/v1/auth';
  static const String userBaseUrl = 'http://localhost:8081/api/v1/users';
  static const String gameBaseUrl = 'http://localhost:8082/api/v1/games';
  static const String wsUrl = 'ws://localhost:8083/ws';

  // Auth endpoints
  static const String googleLogin = '/google';
  static const String refreshToken = '/refresh';
  static const String logout = '/logout';

  // User endpoints
  static const String profile = '/profile';
  static const String stats = '/stats';
  static const String history = '/history';

  // Game endpoints
  static const String startGame = '/start';
  static const String placeBid = '/bid';
  static const String declareTrump = '/trump';
  static const String exchangeKitty = '/kitty';
  static const String playCards = '/play';
}
