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

- **Auth Service** (Port 8080) - User authentication and JWT token management
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
- **Cross-platform** support (iOS, Android, Web)
- **Real-time UI updates** via WebSocket connections
- **Material Design** with custom game-specific components

#### Domain Entities

Flutter domain layer includes:

- **Card & Deck Models**: Complete card system with JSON serialization
- **Game Models**: Game state, players, formations, and tricks
- **Room Models**: Room management with participant tracking
- **User Models**: User profiles with statistics and authentication data

### Infrastructure

- **PostgreSQL** - Primary database for user data and game state
- **Redis** - Session management and real-time game caching
- **Kafka** - Event streaming for game actions and notifications
- **Docker** - Containerized services for development and deployment
- **Kubernetes** - Production-ready orchestration

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
│       └── database/             # Database connections
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

- **Core Domain Models**: Complete implementation of Chinese Bridge game entities
  - Go backend domain models with comprehensive business logic
  - Flutter frontend domain entities with JSON serialization
  - Full unit test coverage (80%+) for all domain logic
- **Project Structure**: Microservices architecture setup
- **Development Environment**: Docker Compose configuration

### 🚧 In Progress

- Repository patterns and data access layers
- REST API implementations
- Authentication service integration
- WebSocket real-time communication
- Flutter UI components and BLoC state management

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
curl http://localhost:8081/api/v1/health  # User Service
curl http://localhost:8082/api/v1/health  # Game Service

# Test authentication
curl http://localhost:8080/api/v1/auth/google

# Test protected endpoints (requires JWT token)
curl -H "Authorization: Bearer <token>" http://localhost:8081/api/v1/users/profile
```

### Database Access

```bash
# Connect to PostgreSQL
docker exec -it chinese-bridge-postgres psql -U user -d chinese_bridge

# Connect to Redis
docker exec -it chinese-bridge-redis redis-cli

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

## 🔐 Authentication

The game uses Google OAuth 2.0 for authentication:

1. **Setup Google OAuth**:

   - Create a project in Google Cloud Console
   - Enable Google+ API
   - Create OAuth 2.0 credentials
   - Add credentials to `.env` file

2. **JWT Tokens**:
   - Access tokens for API authentication
   - Refresh tokens for session management
   - Automatic token refresh in Flutter app

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

# Run domain entity tests specifically
go test ./internal/game/domain -v

# Run specific service tests
go test ./internal/auth/...
```

### Frontend Tests

```bash
cd flutter_app

# Run all unit tests
flutter test

# Run domain entity tests specifically
flutter test test/features/game/domain/entities/

# Run integration tests
flutter test integration_test/

# Run tests with coverage
flutter test --coverage
```

### Test Coverage

Current test coverage for completed components:

- **Go Domain Entities**: 80%+ coverage
- **Flutter Domain Entities**: 80%+ coverage
- **Core Game Logic**: Comprehensive rule validation tests
- **JSON Serialization**: Full serialization/deseriization tests

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
POST /api/v1/auth/google          # Google OAuth login
POST /api/v1/auth/refresh         # Refresh JWT token
POST /api/v1/auth/logout          # Logout user
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
