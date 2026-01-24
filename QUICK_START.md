# AnonTalk - Quick Start Guide

## 🚀 Start in 3 Steps

### Step 1: Start Backend (30 seconds)
```bash
cd backend
docker-compose up -d
```

Wait for MySQL and Go server to start. You should see:
```
✓ Creating anontalk-mysql   ... done
✓ Creating anontalk-backend ... done
```

### Step 2: Start Frontend (1 minute)
```bash
# In project root (not backend/)
npm install
cp .env.local.example .env.local
npm run dev
```

Wait for:
```
➜  Local:   http://localhost:3000
```

### Step 3: Open in Browser
- Go to http://localhost:3000
- Click "Get Started"
- Create account → Set profile → Find match!

---

## 📚 Common Commands

### Backend
```bash
# Start with Docker
docker-compose up -d          # Start services
docker-compose logs backend   # View logs
docker-compose logs mysql     # View database logs
docker-compose down           # Stop services

# Development (without Docker)
go run main.go               # Run locally
go mod download              # Install dependencies
```

### Frontend
```bash
npm run dev                  # Start development
npm run build               # Production build
npm run lint                # Check code
npm run type-check          # TypeScript check
npm run format              # Format code
```

### Database
```bash
# Access MySQL
docker exec -it anontalk-mysql mysql -u root -p anontalk

# Common queries
SELECT COUNT(*) FROM users;
SELECT * FROM chat_sessions;
DESCRIBE messages;
```

---

## 🔍 Verify Everything Works

### Backend Running?
```bash
curl http://localhost:8080/api/auth/me
# Should return: "Invalid token" error (that's good!)
```

### Frontend Running?
```bash
# Visit in browser
http://localhost:3000
# Should show landing page
```

### Database Connected?
```bash
docker exec -it anontalk-mysql mysql -u root -p anontalk -e "SELECT 1;"
# Should return: 1
```

---

## 🧪 Quick Test

### Create Account 1
1. Go to http://localhost:3000
2. Click "Get Started"
3. Email: test1@example.com
4. Password: password123
5. Gender: Male
6. Click Sign Up
7. Add profile info → Continue

### Create Account 2 (Incognito)
1. Open incognito window
2. Go to http://localhost:3000
3. Email: test2@example.com
4. Password: password123
5. Gender: Female
6. Repeat steps 1

### Match & Chat
1. Account 1: Click "Find Someone"
2. Account 2: Click "Find Someone"
3. Both should match within 1-2 seconds
4. Start chatting!

---

## 📝 Customization Quick Edits

### Change Brand Color
Edit `/app/globals.css`:
```css
@theme inline {
  --color-primary: #your-color-hex;
}
```

### Change Site Title
Edit `/app/layout.tsx`:
```tsx
export const metadata = {
  title: 'Your App Name',
}
```

### Change Backend URL
Edit `/.env.local`:
```
NEXT_PUBLIC_API_URL=http://your-server:8080
```

### Change Database Credentials
Edit `/backend/.env`:
```
DATABASE_URL=username:password@tcp(host:3306)/dbname
```

---

## 🐛 Troubleshooting

### Port 3000 in use?
```bash
kill -9 $(lsof -t -i:3000)
```

### Port 8080 in use?
```bash
kill -9 $(lsof -t -i:8080)
```

### MySQL won't start?
```bash
docker-compose down
docker volume rm backend_mysql_data  # Clear data
docker-compose up -d
```

### Frontend shows "API Error"?
```bash
# Check backend is running
docker-compose ps
# Check logs
docker-compose logs backend
```

### "Cannot GET /" after npm start?
```bash
# Make sure you're in project root (not /app or /backend)
cd ..
npm run dev
```

---

## 📊 Useful Status Checks

### Is Docker running?
```bash
docker ps
```

### Are containers healthy?
```bash
docker-compose ps
```

### Backend logs (follow in real-time)
```bash
docker-compose logs -f backend
```

### Database access
```bash
docker exec -it anontalk-mysql bash
mysql -u root -p
# Password: root
```

---

## 🔐 Security Notes

### Default Credentials
```
MySQL:
  User: root
  Password: root

JWT Secret:
  your-secret-key (CHANGE IN PRODUCTION!)
```

### Before Going Public
1. Change JWT_SECRET in backend/.env
2. Change MySQL password
3. Set strong CORS origins
4. Enable HTTPS
5. Set NEXT_PUBLIC_API_URL to production URL

---

## 📱 Testing Across Devices

### Same Network
```bash
# Get your machine's IP
ifconfig | grep inet

# On another device, use:
http://YOUR_IP:3000
```

### Mobile
- Use ngrok for tunneling
- Or deploy to cloud

---

## 🎯 Next Steps

- [ ] Edit colors in globals.css
- [ ] Update landing page copy
- [ ] Customize email settings
- [ ] Add more validation
- [ ] Deploy to cloud
- [ ] Set up monitoring
- [ ] Enable analytics

---

## 📞 Quick Help

| Issue | Solution |
|-------|----------|
| Nothing loads | Restart: `docker-compose restart` |
| "Cannot connect" | Check backend: `curl http://localhost:8080/api/auth/me` |
| Database error | Check MySQL: `docker-compose logs mysql` |
| Blank frontend | Clear cache: `Cmd+Shift+R` (Mac) or `Ctrl+Shift+R` (Windows) |
| Matching not working | Check both accounts created successfully |

---

## 📖 For More Details

- **Full Setup**: See `SETUP.md`
- **Project Overview**: See `README.md`
- **Download Guide**: See `DOWNLOAD.md`
- **Project Summary**: See `PROJECT_SUMMARY.md`
- **Backend Docs**: See `backend/README.md`

---

## 🎉 You're Ready!

Your full-stack AnonTalk application is now running locally!

**Frontend**: http://localhost:3000  
**Backend**: http://localhost:8080  
**Database**: localhost:3306  

Happy coding! 🚀
