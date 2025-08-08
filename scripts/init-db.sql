-- Initialize Chinese Bridge Game Database
-- This script sets up the initial database structure

-- Create extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create initial tables (will be managed by GORM migrations in services)
-- This is just for initial setup

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    google_id VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    avatar VARCHAR(500),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- User statistics table
CREATE TABLE IF NOT EXISTS user_stats (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    games_played INTEGER DEFAULT 0,
    games_won INTEGER DEFAULT 0,
    games_as_declarer INTEGER DEFAULT 0,
    declarer_wins INTEGER DEFAULT 0,
    total_points INTEGER DEFAULT 0,
    average_bid DECIMAL(5,2) DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Rooms table
CREATE TABLE IF NOT EXISTS rooms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    host_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    max_players INTEGER DEFAULT 4,
    current_players INTEGER DEFAULT 0,
    status VARCHAR(50) DEFAULT 'waiting',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Room participants junction table
CREATE TABLE IF NOT EXISTS room_participants (
    room_id UUID REFERENCES rooms(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    position INTEGER,
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (room_id, user_id)
);

-- Games table
CREATE TABLE IF NOT EXISTS games (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    declarer_id UUID REFERENCES users(id),
    trump_suit VARCHAR(20),
    contract INTEGER,
    final_score INTEGER,
    winner_team VARCHAR(20),
    game_data JSONB,
    started_at TIMESTAMP WITH TIME ZONE,
    ended_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Game participants junction table
CREATE TABLE IF NOT EXISTS game_participants (
    game_id UUID REFERENCES games(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    position INTEGER,
    role VARCHAR(20),
    points_captured INTEGER DEFAULT 0,
    PRIMARY KEY (game_id, user_id)
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_users_google_id ON users(google_id);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_rooms_host_id ON rooms(host_id);
CREATE INDEX IF NOT EXISTS idx_rooms_status ON rooms(status);
CREATE INDEX IF NOT EXISTS idx_games_room_id ON games(room_id);
CREATE INDEX IF NOT EXISTS idx_games_declarer_id ON games(declarer_id);
CREATE INDEX IF NOT EXISTS idx_games_started_at ON games(started_at);

-- Insert initial data for development
INSERT INTO users (google_id, email, name, avatar) VALUES 
('dev_user_1', 'dev1@example.com', 'Developer 1', 'https://example.com/avatar1.jpg'),
('dev_user_2', 'dev2@example.com', 'Developer 2', 'https://example.com/avatar2.jpg'),
('dev_user_3', 'dev3@example.com', 'Developer 3', 'https://example.com/avatar3.jpg'),
('dev_user_4', 'dev4@example.com', 'Developer 4', 'https://example.com/avatar4.jpg')
ON CONFLICT (google_id) DO NOTHING;

-- Initialize user stats for development users
INSERT INTO user_stats (user_id, games_played, games_won)
SELECT id, 0, 0 FROM users
ON CONFLICT (user_id) DO NOTHING;