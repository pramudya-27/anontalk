# AnonTalk - Anonymous Chat Platform

A modern anonymous chat application where users match based on gender preferences and chat in real-time. Built with Next.js/React frontend and Go backend with MySQL database.

## Features

✅ **Anonymous Matchmaking**
- Smart matching based on gender preferences
- Real-time queue system
- Instant pairing algorithm

✅ **Real-Time Chat**
- WebSocket-based messaging
- Message history
- Session management

✅ **User Management**
- Email-based authentication
- User profiles with interests
- Privacy-focused design
- User preferences and blocking

✅ **Safety & Moderation**
- Report inappropriate behavior
- Admin dashboard
- User banning system
- Report management

✅ **Full-Stack Architecture**
- Frontend: Next.js 16 with React 19
- Backend: Go with MySQL
- Scalable and performant
- Docker support

## Tech Stack

### Frontend
- **Next.js 16** - React framework with App Router
- **React 19** - UI library
- **TypeScript** - Type safety
- **Tailwind CSS** - Styling
- **shadcn/ui** - UI components

### Backend
- **Go 1.21** - Server runtime
- **MySQL 8.0** - Database
- **Gorilla Mux** - HTTP routing
- **Gorilla WebSocket** - Real-time communication
- **JWT** - Authentication
- **bcrypt** - Password hashing

### DevOps
- **Docker** - Containerization
- **Docker Compose** - Multi-container orchestration
- **MySQL** - Relational database

## Quick Start

### Prerequisites
- Docker and Docker Compose (recommended)
- Or: Node.js 18+, Go 1.21+, MySQL 8.0+

### Jalankan lokal dengan Make

Prasyarat: Docker Compose, Node.js/npm, Go, dan Make.

```bash
make setup  # install dependency dan siapkan .env.local
make dev    # jalankan MySQL, Redis, backend, dan frontend
```

Buka http://localhost:3000. Backend tersedia di http://localhost:8080.
Gunakan `make clean` untuk menghentikan MySQL dan Redis.

### Start with Docker

```bash
# Clone the repository
git clone <repo-url>
cd anontalk

# Start backend and database
cd backend
docker-compose up -d

# In another terminal, start frontend
npm install
cp .env.local.example .env.local
npm run dev
```

Visit http://localhost:3000 to access the application.

### Manual Setup

See [SETUP.md](./SETUP.md) for detailed local installation instructions.

## Project Structure

```
anontalk/
├── app/                    # Next.js application
│   ├── auth/              # Authentication pages
│   ├── chat/              # Chat interface
│   ├── dashboard/         # User dashboard
│   ├── admin/             # Admin panel
│   ├── api/               # API routes (legacy, using Go backend)
│   ├── globals.css        # Global styles
│   └── layout.tsx         # Root layout
├── backend/               # Go backend
│   ├── main.go            # Server entry point
│   ├── internal/
│   │   ├── database/      # MySQL operations
│   │   ├── handlers/      # API handlers
│   │   └── middleware/    # Auth & CORS
│   ├── migrations/        # Database schema
│   ├── docker-compose.yml # Container setup
│   └── Dockerfile         # Backend image
├── components/            # React components
├── lib/                   # Utilities
│   └── api-client.ts      # Go backend client
├── scripts/               # Database scripts
├── SETUP.md              # Setup guide
└── README.md             # This file
```

## API Endpoints

All endpoints require authentication (JWT token) except signup and login.

### Authentication
- `POST /api/auth/signup` - Create account
- `POST /api/auth/login` - Login
- `POST /api/auth/logout` - Logout
- `GET /api/auth/me` - Get current user

### Matchmaking
- `POST /api/matchmaking/join-queue` - Join matching queue
- `POST /api/matchmaking/leave-queue` - Leave queue
- `GET /api/matchmaking/status` - Get status

### Chat
- `POST /api/chat/send-message` - Send message
- `GET /api/chat/sessions/{id}/messages` - Get messages
- `POST /api/chat/sessions/{id}/end` - End session
- `WS /ws/chat/{sessionId}` - WebSocket connection

### Moderation
- `POST /api/reports/create` - Report user
- `GET /api/admin/stats` - Get statistics
- `GET /api/admin/reports` - List reports
- `POST /api/admin/reports/{id}/action` - Process report
- `POST /api/admin/users/{id}/ban` - Ban user

See [backend/README.md](./backend/README.md) for detailed API documentation.

## Database Schema

### Core Tables
- **users** - User accounts with authentication
- **user_profiles** - Extended profile information
- **chat_sessions** - Active/historical chats
- **messages** - Chat message history
- **matchmaking_queue** - Active queue for matching

### Moderation Tables
- **reports** - User reports
- **admin_users** - Admin account information
- **blocked_users** - User blocking records

See [backend/migrations/001_init.sql](./backend/migrations/001_init.sql) for full schema.

## Authentication Flow

1. User signs up with email, password, and gender
2. Backend creates user account with bcrypt-hashed password
3. JWT token issued (24-hour expiration)
4. Token stored in localStorage
5. All subsequent requests include `Authorization: Bearer <token>` header
6. Middleware validates token on each request
7. Banned users are rejected at middleware level

## Matching Algorithm

1. User joins matchmaking queue with preferences
2. Backend searches for users with:
   - Matching gender preference
   - Opposite gender (if applicable)
   - Similar interests (future enhancement)
3. First available match creates chat session
4. Both users notified of match
5. Session persists until either user leaves

## Real-Time Communication

- WebSocket connection established after matching
- Messages broadcast to both users in session
- Connection automatic reconnect on disconnect
- Message history available via REST API
- Session history stored in database

## Security Features

- **Passwords**: bcrypt hashing with salt
- **Authentication**: JWT tokens with expiration
- **Database**: Parameterized queries prevent SQL injection
- **API**: CORS protection
- **Privacy**: Anonymous user identities
- **Moderation**: Report and ban system
- **HTTPS Ready**: Configured for production SSL

## Performance Optimizations

- Connection pooling (25 max, 5 idle)
- Database query indexing
- Efficient WebSocket message handling
- Lazy loading of components
- Code splitting in Next.js
- Static generation where applicable

## Development Workflow

### Frontend Development
```bash
npm run dev        # Start dev server
npm run build      # Build for production
npm run lint       # Run linter
npm run type-check # TypeScript check
```

### Backend Development
```bash
cd backend
go run main.go     # Run locally
go test ./...      # Run tests
docker-compose up  # Docker setup
```

## Production Deployment

### Docker Deployment
```bash
# Build images
docker build -t anontalk-frontend .
docker build -t anontalk-backend ./backend

# Push to registry
docker push your-registry/anontalk-frontend
docker push your-registry/anontalk-backend

# Deploy using docker-compose or Kubernetes
```

### Environment Variables (Production)
```
NEXT_PUBLIC_API_URL=https://api.yourdomain.com
DATABASE_URL=user:pass@tcp(db-host:3306)/anontalk
JWT_SECRET=use-a-strong-random-string
SERVER_PORT=8080
NODE_ENV=production
```

## Monitoring & Logging

- Error tracking via server logs
- API request/response logging
- Database query logging (optional)
- WebSocket connection logging
- Admin stats endpoint for metrics

## Future Enhancements

- [ ] Email verification
- [ ] Rate limiting
- [ ] Advanced matching algorithm
- [ ] User interests matching
- [ ] Message encryption
- [ ] Multimedia support
- [ ] Push notifications
- [ ] Mobile app (React Native)
- [ ] Analytics dashboard
- [ ] A/B testing framework

## Contributing

1. Create a feature branch
2. Make your changes
3. Submit a pull request
4. Code will be reviewed and merged

## Testing

### Frontend
```bash
npm test           # Run Jest tests
npm run e2e        # Run E2E tests
```

### Backend
```bash
cd backend
go test ./...      # Run all tests
go test ./... -v   # Verbose output
```

## Troubleshooting

See [SETUP.md](./SETUP.md) troubleshooting section for common issues.

## License

MIT

## Support

- Backend Documentation: [backend/README.md](./backend/README.md)
- Setup Guide: [SETUP.md](./SETUP.md)
- API Testing: See SETUP.md for cURL examples

## Authors

- Built with ❤️ for anonymous, safe conversations

---

**Note:** This platform prioritizes user safety and privacy. All conversations are anonymous and temporary unless explicitly saved. Users are expected to follow community guidelines and treat others respectfully.
