class ApiConstants {
  static const String baseUrl = 'http://localhost:8080/api/v1';

  // Auth endpoints
  static const String authGoogle = '/auth/google';
  static const String authGoogleUrl = '/auth/google/url';
  static const String authRefresh = '/auth/refresh';
  static const String authLogout = '/auth/logout';

  // User endpoints
  static const String userProfile = '/users/profile';
  static const String userStats = '/users/stats';
  static const String userHistory = '/users/history';

  // Room endpoints
  static const String rooms = '/rooms';
  static const String roomJoin = '/rooms/{id}/join';
  static const String roomLeave = '/rooms/{id}/leave';

  // Game endpoints
  static const String gameStart = '/game/{roomId}/start';
  static const String gameBid = '/game/{gameId}/bid';
  static const String gameTrump = '/game/{gameId}/trump';
  static const String gameKitty = '/game/{gameId}/kitty';
  static const String gamePlay = '/game/{gameId}/play';

  // WebSocket endpoints
  static const String wsUrl = 'ws://localhost:8080/ws';
}
