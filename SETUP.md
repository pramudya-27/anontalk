# AnonTalk - Complete Setup Guide

## Project Structure

```
anontalk/
├── backend/           # Go backend with MySQL
├── app/              # Next.js frontend
├── lib/              # Shared utilities
├── components/       # UI components
└── scripts/          # Database migrations
```

## Quick Start

### Option 1: Using Docker (Recommended)

1. **Start the backend with Docker:**
   ```bash
   cd backend
   docker-compose up -d
   ```

2. **Install frontend dependencies:**
   ```bash
   npm install
   ```

3. **Set up frontend environment:**
   ```bash
   cp .env.local.example .env.local
   ```

4. **Run the frontend:**
   ```bash
   npm run dev
   ```

5. **Access the application:**
   - Frontend: http://localhost:3000
   - Backend: http://localhost:8080

### Option 2: Local Setup (Without Docker)

#### Backend Setup:

1. **Install MySQL:**
   - On macOS: `brew install mysql`
   - On Windows: Download from https://dev.mysql.com/downloads/mysql/
   - On Linux: `sudo apt-get install mysql-server`

2. **Start MySQL service:**
   ```bash
   # macOS
   brew services start mysql
   
   # Linux
   sudo systemctl start mysql
   
   # Windows (if installed as service)
   net start MySQL80
   ```

3. **Create database and user:**
   ```bash
   mysql -u root -p
   CREATE DATABASE anontalk;
   CREATE USER 'anontalk'@'localhost' IDENTIFIED BY 'password';
   GRANT ALL PRIVILEGES ON anontalk.* TO 'anontalk'@'localhost';
   FLUSH PRIVILEGES;
   ```

4. **Run migrations:**
   ```bash
   mysql -u anontalk -p anontalk < backend/migrations/001_init.sql
   ```

5. **Install Go dependencies:**
   ```bash
   cd backend
   go mod download
   ```

6. **Set up backend environment:**
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials
   ```

7. **Run the backend:**
   ```bash
   go run main.go
   ```

#### Frontend Setup:

1. **Install Node dependencies:**
   ```bash
   npm install
   ```

2. **Set up environment:**
   ```bash
   cp .env.local.example .env.local
   ```

3. **Run the development server:**
   ```bash
   npm run dev
   ```

## Default Credentials

### Database
- **Username:** anontalk
- **Password:** password
- **Database:** anontalk
- **Host:** localhost:3306

### Backend API
- **Base URL:** http://localhost:8080
- **Default Port:** 8080

### Frontend
- **URL:** http://localhost:3000
- **Default Port:** 3000

## Testing the Application

### 1. Create an Account
- Visit http://localhost:3000
- Click "Get Started"
- Fill in email, password, and gender
- Click "Sign Up"

### 2. Set Up Profile
- Add display name, age range, and interests
- Select preferred gender for matching

### 3. Start Matching
- Click "Find Someone to Chat"
- Wait for a match (in test, you can open another incognito window and repeat steps 1-2)

### 4. Chat
- Once matched, you'll see the chat interface
- Messages are sent via WebSocket in real-time
- Click "End Chat" to disconnect

## API Testing with cURL

### Sign Up
```bash
curl -X POST http://localhost:8080/api/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "gender": "male",
    "preferredGender": ["female"],
    "interests": ["travel", "music"]
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

### Join Matchmaking Queue
```bash
curl -X POST http://localhost:8080/api/matchmaking/join-queue \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "gender": "male",
    "preferredGender": ["female"],
    "interests": ["travel", "music"]
  }'
```

## Troubleshooting

### MySQL Connection Refused
- Ensure MySQL service is running: `mysql.server status`
- Check credentials in backend/.env
- Verify database exists: `mysql -u anontalk -p anontalk -e "SELECT 1;"`

### Backend Port Already in Use
```bash
# Find process using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>
```

### Frontend Cannot Connect to Backend
- Ensure backend is running on port 8080
- Check NEXT_PUBLIC_API_URL in .env.local
- Ensure CORS is enabled (it is by default in backend)

### WebSocket Connection Failed
- Check backend is running
- Verify firewall allows connections to port 8080
- Check browser console for specific error messages

## Development Tips

### Frontend
- Modify pages in `/app` directory
- Create new API routes in `/app/api`
- Components are in `/components` directory
- Run `npm run dev` for hot reload

### Backend
- API handlers are in `/backend/internal/handlers`
- Database operations in `/backend/internal/database`
- Middleware in `/backend/internal/middleware`
- Run `go run main.go` for development

### Database
- View tables: `SHOW TABLES;`
- Inspect schema: `DESCRIBE users;`
- Query data: `SELECT * FROM users;`

## Production Deployment

### Backend
```bash
# Build Docker image
cd backend
docker build -t anontalk-backend .

# Push to Docker Hub or private registry
docker push your-registry/anontalk-backend:latest
```

### Frontend
```bash
# Build for production
npm run build

# Deploy to Vercel, Netlify, or your own server
npm run start
```

## Security Checklist

- [ ] Change JWT_SECRET in backend/.env
- [ ] Use strong MySQL password
- [ ] Set up HTTPS for production
- [ ] Enable rate limiting
- [ ] Configure proper CORS origins
- [ ] Set up user input validation
- [ ] Enable database backups
- [ ] Monitor for suspicious activity
- [ ] Keep dependencies updated
- [ ] Use environment variables for secrets

## Support

For issues or questions:
1. Check the logs: `docker-compose logs backend`
2. Check browser console for frontend errors
3. Verify all services are running
4. Review the API documentation in backend/README.md

## Next Steps

- [ ] Customize the UI theme
- [ ] Add email verification
- [ ] Implement rate limiting
- [ ] Set up monitoring and logging
- [ ] Add analytics
- [ ] Deploy to production
