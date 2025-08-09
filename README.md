# Chinese Bridge Game

A modern, real-time multiplayer Chinese Bridge card game built with Go microservices backend and Flutter mobile frontend.

## 🎯 Overview

Chinese Bridge is a traditional Chinese card game similar to Bridge, featuring bidding, trump declaration, and strategic card play. This implementation provides:

- **Real-time multiplayer gameplay** with WebSocket support
- **Google OAuth authentication** for seamless user experience
- **Microservices architecture** for scalability and maintainability
- **Cross-platform mobile app** built with Flutter
- **Cloud-ready deployment** with Docker and Kubernetes support

## 🏗️ Architecture

### Backend (Go Microservices)

- **Auth Service** (Port 8080) - ✅ **COMPLETED** - User authentication and JWT token management
  - Google OAuth 2.0 integration with secure token exchange
  - JWT access/refresh token management with Redis session storage
  - Rate limiting, security headers, and comprehensive error handling
  - Swagger/OpenAPI 3.0 documentation with interactive UI
  - 100% test coverage with unit and integration tests
- **User Service** (Port 8081) - User profiles, statistics, and game history
- **Game Service** (Port 8082) - Game logic, room management, and real-time gameplay
- **WebSocket Service** (Port 8083) - Real-time communication for live gameplay

#### Domain Models

The game engine implements comprehensive domain entities following Chinese Bridge rules:

- **Card System**: 108-card deck (2 standard decks + 4 jokers) with trump hierarchy
- **Formation Recognition**: Singles, pairs, and tractors with validation
- **Game State Management**: Complete game flow from bidding to scoring
- **Trick Management**: Suit-following rules and winner determination

### Frontend (Flutter)

- **Clean Architecture** with BLoC state management
- **Authentication Module** - ✅ **COMPLETED** - Complete Google OAuth integration
  - BLoC state management with comprehensive error handling
  - Repository pattern with local/remote data sources
  - Automatic token refresh with secure local storage
  - Material Design login UI with loading states
  - 100% test coverage with BLoC and widget tests
- **Cross-platform** support (iOS, Android, Web)
- **Real-time UI updates** via WebSocket connections (planned)
- **Material Design** with custom game-specific components

#### Domain Entities

Flutter domain layer includes:

- **Card & Deck Models**: Complete card system with JSON serialization
- **Game Models**: Game state, players, formations, and tricks
- **Room Models**: Room management with participant tracking
- **User Models**: User profiles with statistics and authentication data

### Infrastructure

- **PostgreSQL** - Primary database with GORM ORM for user data, game state, and statistics
- **Redis** - High-performance caching layer for sessions, room states, and real-time game data
- **Kafka** - Event streaming for game actions and notifications
- **Docker** - Containerized services for development and deployment
- **Kubernetes** - Production-ready orchestration

#### Database Layer

- **GORM Models**: Comprehensive database models with relationships and constraints
- **Repository Pattern**: Clean data access layer with interface-based design
- **Migration System**: Automated database schema management with indexing
- **Redis Caching**: Intelligent caching with TTL policies and invalidation strategies
- **Test Coverage**: Full integration tests for database and cache operations

## 📁 Project Structure

```
chinese-bridge-game/
├── cmd/                          # Service entry points
│   ├── auth-service/
│   ├── user-service/
│   └── game-service/
├── internal/                     # Private application code
│   ├── auth/                     # Authentication service
│   ├── user/                     # User management service
│   ├── game/                     # Game logic service
│   │   └── domain/               # ✅ Core game entities (Card, GameState, Formation, Trick)
│   └── common/                   # Shared utilities
│       ├── config/               # Configuration management
│       └── database/             # ✅ Database layer (GORM models, repositories, caching)
├── pkg/                          # Public libraries
│   └── middleware/               # HTTP middleware
├── flutter_app/                  # Flutter mobile application
│   ├── lib/
│   │   ├── core/                 # Core utilities and DI
│   │   └── features/             # Feature-based modules
│   │       ├── authentication/
│   │       └── game/
│   │           └── domain/
│   │               └── entities/ # ✅ Game entities (Card, Game, Room, User)
│   └── test/                     # ✅ Comprehensive unit tests
├── docker/                       # Docker configurations
├── k8s/                          # Kubernetes manifests
├── scripts/                      # Database and setup scripts
├── docker-compose.yml            # Local development setup
├── Makefile                      # Development commands
└── README.md
```

## 📋 Development Status

### ✅ Completed Features

- **Authentication System**: Complete Google OAuth implementation
  - Go backend auth service with JWT token management
  - Flutter authentication module with BLoC state management
  - Secure session management with Redis caching
  - Comprehensive API documentation with Swagger
  - Rate limiting and security middleware
  - Full unit and integration test coverage (100%)
- **Core Domain Models**: Complete implementation of Chinese Bridge game entities
  - Go backend domain models with comprehensive business logic
  - Flutter frontend domain entities with JSON serialization
  - Full unit test coverage (80%+) for all domain logic
- **Database Layer**: Complete data persistence and caching implementation
  - GORM models with PostgreSQL for all game entities
  - Repository pattern with comprehensive CRUD operations
  - Redis caching layer with intelligent invalidation strategies
  - Migration system with automated schema management
  - Full integration test coverage for database operations
- **Project Structure**: Microservices architecture setup
- **Development Environment**: Docker Compose configuration

### 🚧 In Progress

- User management service implementation
- Game service and room management
- WebSocket real-time communication
- Flutter game UI components

### 📅 Planned

- Complete game flow implementation
- User management and statistics
- Room management and matchmaking
- Production deployment with Kubernetes
- Performance optimization and monitoring

## 🚀 Quick Start

### Prerequisites

- **Go 1.21+**
- **Flutter 3.10+**
- **Docker & Docker Compose**
- **Make** (optional, for convenience commands)

### 1. Clone and Setup

```bash
git clone <repository-url>
cd chinese-bridge-game

# Copy environment configuration
cp .env.example .env

# Edit .env with your configuration
# - Database credentials
# - Google OAuth credentials
# - JWT secrets
```

### 2. Start Infrastructure

```bash
# Start databases and message queue
docker-compose up -d postgres redis kafka

# Or use the Makefile
make setup
```

### 3. Build and Run Backend Services

```bash
# Build all services
make build

# Run services (in separate terminals)
make run-auth    # Terminal 1 - Auth Service (port 8080)
make run-user    # Terminal 2 - User Service (port 8081)
make run-game    # Terminal 3 - Game Service (port 8082)
```

### 4. Run Flutter App

```bash
cd flutter_app
flutter pub get
flutter run
```

## 🔧 Development

### Available Make Commands

```bash
make help           # Show all available commands
make setup          # Setup development environment
make build          # Build all Go services
make test           # Run Go tests
make flutter-get    # Get Flutter dependencies
make flutter-run    # Run Flutter app
make docker-up      # Start all services with Docker
make docker-down    # Stop Docker services
make k8s-up         # Deploy to Kubernetes
make clean          # Clean build artifacts
make lint           # Run code linters
make format         # Format code
```

### Testing Services

```bash
# Test service health
curl http://localhost:8080/api/v1/health  # Auth Service
curl http://localhost:8081/api/v1/health  # User Service (planned)
curl http://localhost:8082/api/v1/health  # Game Service (planned)

# View API documentation
open http://localhost:8080/swagger/index.html  # Auth Service Swagger UI

# Test authentication endpoints
curl http://localhost:8080/api/v1/auth/google/url?state=test
curl -X POST http://localhost:8080/api/v1/auth/google \
  -H "Content-Type: application/json" \
  -d '{"code": "google_auth_code"}'

# Test token refresh
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token": "your_refresh_token"}'

# Test protected endpoints (requires JWT token)
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Authorization: Bearer your_access_token"
```

### Database Access

```bash
# Connect to PostgreSQL
docker exec -it chinese-bridge-postgres psql -U user -d chinese_bridge

# Run database migrations
go run cmd/migrate/main.go

# Connect to Redis
docker exec -it chinese-bridge-redis redis-cli

# View cached data
redis-cli KEYS "session:user:*"
redis-cli KEYS "room:state:*"
redis-cli KEYS "game:state:*"

# View Kafka topics
docker exec -it chinese-bridge-kafka kafka-topics --bootstrap-server localhost:9092 --list
```

## 🎮 Game Rules

### Chinese Bridge Basics

1. **Players**: 4 players in partnerships (North-South vs East-West)
2. **Deck**: Standard 52-card deck
3. **Bidding**: Players bid on the number of tricks they can take
4. **Trump**: Winning bidder declares trump suit
5. **Kitty**: Special cards that can be exchanged
6. **Play**: Trick-taking gameplay with trump suit advantages

### Game Flow

1. **Room Creation**: Host creates a game room
2. **Player Joining**: 3 other players join the room
3. **Bidding Phase**: Players bid on contract
4. **Trump Declaration**: Winning bidder declares trump
5. **Kitty Exchange**: Declarer exchanges cards with kitty
6. **Card Play**: 13 tricks of card play
7. **Scoring**: Points calculated based on contract success

## 🗄️ Database Schema

### Core Entities

The database implements a comprehensive schema for the Chinese Bridge game:

#### User Management

- **Users**: Player profiles with Google OAuth integration
- **UserStats**: Game statistics and performance metrics
- **Sessions**: JWT token management with expiration

#### Game Management

- **Rooms**: Game room creation and participant management
- **RoomParticipants**: Junction table for room membership
- **Games**: Individual game instances with complete state
- **GameParticipants**: Player participation in specific games

### Caching Strategy

Redis caching layer provides high-performance data access:

#### Cache Types

- **User Sessions** (24h TTL): Authentication and user state
- **Room States** (30min TTL): Active room information and participants
- **Game States** (2h TTL): Real-time game data and player actions
- **Leaderboards** (5min TTL): Player rankings and statistics
- **WebSocket Connections** (1h TTL): Active player connections
- **Matchmaking Queue**: Player queue for game matching

#### Cache Invalidation

- **Automatic TTL**: Time-based expiration for all cache entries
- **Manual Invalidation**: Event-driven cache clearing on data changes
- **Cascade Invalidation**: Related data cleanup (e.g., user changes invalidate leaderboard)
- **Periodic Cleanup**: Background processes for expired entry removal

## 🔐 Authentication

The game uses a comprehensive Google OAuth 2.0 authentication system with JWT token management:

### Backend Authentication Service

- **Google OAuth Integration**: Server-side OAuth flow with secure token exchange
- **JWT Token Management**: Access tokens (1h) and refresh tokens (7 days) with automatic rotation
- **Session Management**: Redis-based session storage with TTL policies
- **Security Features**: Rate limiting (5 req/sec per IP), security headers, CORS protection
- **API Documentation**: Complete Swagger/OpenAPI 3.0 documentation at `/swagger/index.html`

### Flutter Authentication Module

- **Clean Architecture**: Repository pattern with separate local/remote data sources
- **BLoC State Management**: Comprehensive state management with proper error handling
- **Automatic Token Refresh**: Seamless token renewal with fallback to re-authentication
- **Secure Storage**: Encrypted local storage with token expiration checking
- **Responsive UI**: Material Design login screen with loading states and error handling

### Setup Instructions

1. **Google Cloud Console Setup**:

   ```bash
   # 1. Create a new project in Google Cloud Console
   # 2. Enable Google+ API and OAuth2 API
   # 3. Create OAuth 2.0 credentials (Web application)
   # 4. Add authorized redirect URIs:
   #    - http://localhost:8080/api/v1/auth/google (for development)
   #    - https://yourdomain.com/api/v1/auth/google (for production)
   ```

2. **Environment Configuration**:

   ```bash
   # Copy the example environment file
   cp .env.example .env

   # Add your Google OAuth credentials
   GOOGLE_CLIENT_ID=your-google-client-id
   GOOGLE_CLIENT_SECRET=your-google-client-secret
   GOOGLE_REDIRECT_URL=http://localhost:8080/api/v1/auth/google

   # Set a secure JWT secret (use a strong random string)
   JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
   ```

3. **Flutter Configuration**:
   ```dart
   // Update flutter_app/lib/core/di/injection_container.dart
   // Add your Google Client ID for mobile apps
   final googleSignIn = GoogleSignIn(
     scopes: ['email', 'profile'],
     serverClientId: 'your-google-client-id', // Same as backend
   );
   ```

### Authentication Flow

1. **User Login**: User taps "Sign in with Google" in Flutter app
2. **Google OAuth**: App opens Google sign-in flow and gets authorization code
3. **Token Exchange**: Flutter sends auth code to backend `/auth/google` endpoint
4. **JWT Generation**: Backend validates with Google and generates JWT tokens
5. **Session Storage**: Backend stores session in Redis, Flutter stores tokens locally
6. **API Access**: Flutter includes JWT in Authorization header for protected endpoints
7. **Token Refresh**: Automatic token refresh when access token expires

### API Endpoints

```bash
# Get Google OAuth URL (for web flows)
GET /api/v1/auth/google/url?state=random_state

# Exchange Google auth code for JWT tokens
POST /api/v1/auth/google
Content-Type: application/json
{
  "code": "google_auth_code",
  "state": "optional_state"
}

# Refresh expired access token
POST /api/v1/auth/refresh
Content-Type: application/json
{
  "refresh_token": "your_refresh_token"
}

# Logout and invalidate all sessions
POST /api/v1/auth/logout
Authorization: Bearer your_access_token
```

### Testing Authentication

```bash
# Start the auth service
make run-auth

# Test health endpoint
curl http://localhost:8080/api/v1/health

# View API documentation
open http://localhost:8080/swagger/index.html

# Test authentication flow (requires valid Google auth code)
curl -X POST http://localhost:8080/api/v1/auth/google \
  -H "Content-Type: application/json" \
  -d '{"code": "your_google_auth_code"}'
```

## 🚢 Deployment

### Docker Deployment

```bash
# Build Docker images
make docker-build

# Deploy with Docker Compose
docker-compose up -d
```

### Kubernetes Deployment

```bash
# Deploy to Kubernetes cluster
make k8s-up

# Check deployment status
kubectl get all -n chinese-bridge

# View logs
kubectl logs -f deployment/auth-service -n chinese-bridge
```

### Environment Variables

Key environment variables to configure:

```bash
# Database
DATABASE_URL=postgres://user:password@localhost:5432/chinese_bridge
REDIS_URL=redis://localhost:6379

# Authentication
JWT_SECRET=your-secret-key
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret

# Services
KAFKA_URL=localhost:9092
ENVIRONMENT=development
```

## 🧪 Testing

### Backend Tests

```bash
# Run all Go tests
make test

# Run tests with coverage
go test -cover ./...

# Run authentication service tests
go test ./internal/auth/service -v
go test ./internal/auth/handler -v
go test ./internal/auth/repository -v

# Run domain entity tests specifically
go test ./internal/game/domain -v

# Run database layer tests
go test ./internal/common/database -v

# Run Redis cache tests (requires Redis running)
go test ./internal/common/database -run TestRedisCache -v

# Run tests with coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Frontend Tests

```bash
cd flutter_app

# Run all unit tests
flutter test

# Run authentication BLoC tests specifically
flutter test test/features/authentication/presentation/bloc/

# Run domain entity tests specifically
flutter test test/features/game/domain/entities/

# Run integration tests
flutter test integration_test/

# Run tests with coverage
flutter test --coverage

# Generate code (for JSON serialization)
flutter packages pub run build_runner build
```

### Test Coverage

Current test coverage for completed components:

- **Authentication Service**: 100% coverage (service, handler, repository layers)
- **Flutter Authentication Module**: 100% coverage (BLoC, repository, data sources)
- **Go Domain Entities**: 80%+ coverage
- **Flutter Domain Entities**: 80%+ coverage
- **Database Layer**: 90%+ coverage with integration tests
- **Redis Caching**: Full coverage including TTL expiration tests
- **Core Game Logic**: Comprehensive rule validation tests
- **JSON Serialization**: Full serialization/deserialization tests

## 📊 Monitoring

### Health Checks

All services provide health check endpoints:

- `GET /api/v1/health` - Service health status
- `GET /api/v1/ready` - Service readiness status

### Logging

- **Structured logging** with JSON format
- **Log levels**: DEBUG, INFO, WARN, ERROR
- **Centralized logging** via Docker/Kubernetes

### Metrics

- **Service metrics** via Prometheus (planned)
- **Application metrics** for game statistics
- **Infrastructure metrics** for system monitoring

## 🔧 Troubleshooting

### Authentication Issues

**Problem**: Google OAuth login fails with "invalid_client" error

```bash
# Solution: Check your Google OAuth configuration
# 1. Verify GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET in .env
# 2. Ensure redirect URI matches exactly in Google Cloud Console
# 3. Check that Google+ API is enabled in your project
```

**Problem**: JWT token validation fails

```bash
# Solution: Check JWT secret configuration
# 1. Ensure JWT_SECRET is set in .env and matches across services
# 2. Verify token hasn't expired (access tokens expire in 1 hour)
# 3. Try refreshing the token using the refresh endpoint
```

**Problem**: Redis connection errors

```bash
# Solution: Ensure Redis is running
docker-compose up -d redis

# Check Redis connectivity
docker exec -it chinese-bridge-redis redis-cli ping
```

**Problem**: Database migration errors

```bash
# Solution: Reset database and run migrations
docker-compose down -v  # Remove volumes
docker-compose up -d postgres
make run-auth  # Migrations run automatically on service start
```

**Problem**: Flutter build errors after adding authentication

```bash
# Solution: Regenerate code and clean build
cd flutter_app
flutter packages pub run build_runner build --delete-conflicting-outputs
flutter clean
flutter pub get
flutter run
```

### Common Development Issues

**Problem**: Port conflicts when starting services

```bash
# Solution: Check what's running on ports 8080-8083
lsof -i :8080
# Kill conflicting processes or change ports in docker-compose.yml
```

**Problem**: CORS errors in Flutter web development

```bash
# Solution: The auth service includes CORS headers, but for development:
flutter run -d chrome --web-renderer html --web-port 3000
```

## 🤝 Contributing

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/amazing-feature`
3. **Commit changes**: `git commit -m 'Add amazing feature'`
4. **Push to branch**: `git push origin feature/amazing-feature`
5. **Open a Pull Request**

### Code Standards

- **Go**: Follow Go conventions, use `gofmt` and `golint`
- **Flutter**: Follow Dart conventions, use `dart format`
- **Git**: Use conventional commit messages
- **Testing**: Maintain test coverage above 80%

## 📝 API Documentation

### Authentication Endpoints

```
GET  /api/v1/auth/google/url      # Get Google OAuth authorization URL
POST /api/v1/auth/google          # Exchange Google auth code for JWT tokens
POST /api/v1/auth/refresh         # Refresh expired access token
POST /api/v1/auth/logout          # Logout user and invalidate sessions
GET  /api/v1/health               # Service health check
GET  /api/v1/ready                # Service readiness check
GET  /swagger/index.html          # Interactive API documentation
```

### User Endpoints

```
GET  /api/v1/users/profile        # Get user profile
PUT  /api/v1/users/profile        # Update user profile
GET  /api/v1/users/stats          # Get user statistics
GET  /api/v1/users/history        # Get game history
```

### Game Endpoints

```
POST /api/v1/rooms/:roomId/start  # Start game in room
GET  /api/v1/games/:gameId        # Get game state
POST /api/v1/games/:gameId/bid    # Place bid
POST /api/v1/games/:gameId/trump  # Declare trump
POST /api/v1/games/:gameId/kitty  # Exchange kitty
POST /api/v1/games/:gameId/play   # Play cards
```

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- **Chinese Bridge Community** for game rules and traditions
- **Go Community** for excellent microservices patterns
- **Flutter Team** for the amazing cross-platform framework
- **Open Source Contributors** for the tools and libraries used

---

**Happy Gaming! 🎮**

For questions or support, please open an issue or contact the development team.
