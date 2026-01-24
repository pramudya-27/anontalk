# AnonTalk Backend (Go + MySQL)

A high-performance anonymous chat backend built with Go, featuring real-time matchmaking and WebSocket support.

## Quick Start

### Option 1: Using Docker (Recommended)

```bash
cd backend
docker-compose up -d
```

The backend will be available at `http://localhost:8080`

### Option 2: Local Setup

1. **Install Dependencies**
   ```bash
   go mod download
   ```

2. **Set up MySQL Database**
   ```bash
   mysql -u root -p < migrations/001_init.sql
   ```

3. **Configure Environment**
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials
   ```

4. **Run the Server**
   ```bash
   go run main.go
   ```

## Database Schema

The application uses MySQL with the following tables:
- `users` - User accounts with authentication
- `user_profiles` - Extended user profile information
- `chat_sessions` - Active chat sessions between users
- `messages` - Chat message history
- `matchmaking_queue` - Waiting queue for matchmaking
- `reports` - User reports for moderation
- `admin_users` - Admin/moderator accounts
- `blocked_users` - User blocking records

## API Endpoints

### Authentication
- `POST /api/auth/signup` - Create new account
- `POST /api/auth/login` - Login user
- `POST /api/auth/logout` - Logout user
- `GET /api/auth/me` - Get current user

### User Management
- `GET /api/users/profile` - Get user profile
- `PUT /api/users/profile` - Update profile
- `PUT /api/users/preferences` - Update chat preferences

### Matchmaking
- `POST /api/matchmaking/join-queue` - Join matchmaking queue
- `POST /api/matchmaking/leave-queue` - Leave queue
- `GET /api/matchmaking/status` - Get queue status

### Chat
- `GET /api/chat/sessions` - Get user's chat sessions
- `GET /api/chat/sessions/{sessionID}/messages` - Get messages
- `POST /api/chat/send-message` - Send a message
- `POST /api/chat/sessions/{sessionID}/end` - End chat session
- `WS /ws/chat/{sessionID}` - WebSocket for real-time chat

### Reporting & Admin
- `POST /api/reports/create` - Report a user
- `GET /api/admin/stats` - Get platform statistics
- `GET /api/admin/reports` - List all reports
- `POST /api/admin/reports/{reportID}/action` - Handle report
- `POST /api/admin/users/{userID}/ban` - Ban a user

## Architecture

```
backend/
├── main.go                 # Server entry point
├── internal/
│   ├── database/          # MySQL database layer
│   ├── handlers/          # API route handlers
│   └── middleware/        # Auth & CORS middleware
├── migrations/            # Database migrations
├── docker-compose.yml     # Docker compose configuration
├── Dockerfile            # Docker image definition
└── go.mod               # Go module dependencies
```

## Development

### Running Tests
```bash
go test ./...
```

### Database Migrations
Add new migrations to the `migrations/` folder with incremental numbering.

## Security Considerations

- All passwords are hashed with bcrypt
- JWT tokens expire after 24 hours
- Row-level security is enforced at the database level
- CORS is configured for frontend communication
- SQL injection protection via parameterized queries

## Performance

- Connection pooling (max 25 connections, 5 idle)
- Efficient database indexing
- WebSocket support for real-time updates
- Optimized matchmaking algorithm

## Environment Variables

```
DATABASE_URL=root:password@tcp(localhost:3306)/anontalk
JWT_SECRET=your-secret-key
SERVER_PORT=8080
NODE_ENV=development
```

## Troubleshooting

**MySQL Connection Error**
- Ensure MySQL is running
- Check DATABASE_URL in .env
- Verify database exists: `CREATE DATABASE anontalk;`

**Port Already in Use**
- Change SERVER_PORT in .env
- Or kill existing process: `lsof -ti:8080 | xargs kill`

**WebSocket Connection Failed**
- Ensure frontend is configured to connect to correct backend URL
- Check CORS settings in middleware

## License

MIT
