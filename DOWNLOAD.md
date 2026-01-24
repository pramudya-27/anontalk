# AnonTalk - Download & Local Setup

This complete full-stack application is ready to download and run on your local machine!

## What You're Getting

✅ **Complete Frontend (Next.js)**
- Landing page with feature showcase
- Authentication pages (signup, login)
- User profile setup
- Real-time chat interface
- Dashboard with matchmaking
- Admin panel for moderation

✅ **Complete Backend (Go)**
- RESTful API with 20+ endpoints
- MySQL database with optimized schema
- Real-time WebSocket support
- JWT authentication
- Admin moderation system
- Complete error handling

✅ **Database**
- MySQL schema with migrations
- Optimized indexes
- Proper foreign keys and constraints
- Ready-to-use test data scripts

✅ **DevOps**
- Docker & Docker Compose configuration
- Multi-container setup
- Easy one-command deployment
- Development environment ready

## Quick Download & Start (5 minutes)

### Step 1: Download the Project
```bash
# Click the Download ZIP button in v0
# Extract the ZIP file
cd anontalk
```

### Step 2: Start Backend with Docker
```bash
cd backend
docker-compose up -d
```

You should see:
- ✅ MySQL container running (port 3306)
- ✅ Go backend running (port 8080)

### Step 3: Start Frontend
```bash
# In project root directory
npm install
cp .env.local.example .env.local
npm run dev
```

The frontend will run on http://localhost:3000

### Step 4: Test It Out!
- Open http://localhost:3000
- Click "Get Started"
- Create an account (use any email format: test@example.com)
- Set up your profile
- Open another browser tab in incognito and create another account
- Both should be able to match and chat!

## File Structure After Download

```
anontalk/
├── app/                    # Next.js Pages & Components
├── backend/               # Go Backend Server
│   ├── main.go           # Entry point
│   ├── go.mod            # Dependencies
│   ├── docker-compose.yml # Containers config
│   ├── Dockerfile        # Backend image
│   ├── migrations/       # Database schemas
│   ├── internal/         # Source code
│   └── README.md         # Backend docs
├── components/           # React UI Components
├── lib/                  # Utilities & API Client
├── public/              # Static assets
├── SETUP.md             # Detailed setup guide
├── README.md            # Project overview
├── DOWNLOAD.md          # This file
└── .env.local.example   # Environment template
```

## System Requirements

### Minimum
- 4GB RAM
- 2GB disk space
- Internet connection

### Recommended
- 8GB RAM
- 5GB disk space
- Docker installed
- Modern browser (Chrome, Firefox, Safari, Edge)

## Prerequisites

### Option A: Docker Setup (Easiest)
- Docker Desktop installed
- 10-15 MB download for MySQL image
- That's it!

### Option B: Local Setup
- **Node.js 18+**: https://nodejs.org
- **Go 1.21+**: https://golang.org
- **MySQL 8.0+**: https://www.mysql.com
- **Git** (optional): https://git-scm.com

## Installation Checklist

- [ ] Downloaded and extracted the ZIP file
- [ ] Opened terminal/cmd in project directory
- [ ] Ran `cd backend && docker-compose up -d` (or local setup)
- [ ] Ran `npm install` in project root
- [ ] Ran `npm run dev`
- [ ] Opened http://localhost:3000
- [ ] Created test account
- [ ] Successfully logged in

## Troubleshooting Quick Fixes

### "Cannot connect to backend"
```bash
# Make sure backend is running
docker-compose ps  # Should show 2 running containers

# Or check local ports
lsof -i :8080  # Backend
lsof -i :3000  # Frontend
```

### "Port already in use"
```bash
# Kill process on port 8080
kill -9 $(lsof -t -i:8080)

# Kill process on port 3000
kill -9 $(lsof -t -i:3000)
```

### "MySQL connection error"
```bash
# Restart containers
docker-compose restart

# Or check logs
docker-compose logs mysql
```

### "npm install fails"
```bash
# Clear cache and retry
npm cache clean --force
rm -rf node_modules package-lock.json
npm install
```

## Features to Try

1. **Sign Up & Profile**
   - Create account
   - Set gender preference
   - Add interests

2. **Matchmaking**
   - Click "Find Someone"
   - Open incognito with another account
   - Watch real-time matching

3. **Chat**
   - Send messages
   - See real-time delivery
   - View message history

4. **Reporting**
   - Report a user
   - Check admin panel

5. **Admin Dashboard**
   - View stats
   - See reports
   - Ban users

## Configuration

### Backend (.env)
```
DATABASE_URL=root:password@tcp(localhost:3306)/anontalk
JWT_SECRET=your-secret-key-change-in-production
SERVER_PORT=8080
NODE_ENV=development
```

### Frontend (.env.local)
```
NEXT_PUBLIC_API_URL=http://localhost:8080
```

## Development Tips

### Add a New Feature
1. Create component in `/components`
2. Add page in `/app` if needed
3. Create API endpoint in backend `/internal/handlers`
4. Use `apiClient` to call backend
5. Test and commit

### Debug Issues
```bash
# Frontend logs
npm run dev  # Check console output

# Backend logs
docker-compose logs backend  # or `go run main.go`

# Database
docker exec -it anontalk-mysql mysql -u root -p anontalk
```

## Next Steps

### Local Development
1. Customize UI in `/app` and `/components`
2. Add new features in backend `/internal/handlers`
3. Update database schema if needed
4. Test thoroughly

### Deployment
1. Build Docker images: `docker build -t anontalk-backend ./backend`
2. Push to Docker Hub or registry
3. Deploy to cloud (AWS, Google Cloud, DigitalOcean, etc)
4. Set up production environment variables
5. Configure domain and SSL

## Documentation

- **Full Setup Guide**: See `SETUP.md`
- **Backend API Docs**: See `backend/README.md`
- **Project Overview**: See `README.md`

## Support

If you encounter issues:
1. Check `SETUP.md` troubleshooting section
2. Review logs: `docker-compose logs`
3. Check network connectivity
4. Verify all ports are available
5. Try restarting containers

## Security Notes

⚠️ **Development Only Configuration**
- Default JWT secret (change in production)
- CORS allows all origins (restrict in production)
- No HTTPS (use HTTPS in production)
- Default database credentials (change in production)

See README.md for production security checklist.

## What to Customize

Before going public, customize:
- [ ] JWT_SECRET in backend/.env
- [ ] Database password
- [ ] CORS origins
- [ ] Email configuration
- [ ] Brand colors and text
- [ ] Rate limiting
- [ ] Content moderation rules
- [ ] Terms of service
- [ ] Privacy policy

## License

MIT - Feel free to modify and deploy!

## Questions?

1. Check the documentation files
2. Review the code comments
3. Check API tests in setup guide
4. Review backend/README.md for API details

Happy coding! 🚀
