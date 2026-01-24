package main

import (
	"flag"
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"

	"github.com/anontalk/backend/internal/database"
	"github.com/anontalk/backend/internal/handlers"
	"github.com/anontalk/backend/internal/middleware"
	"github.com/anontalk/backend/internal/redis_client"
	"github.com/anontalk/backend/internal/worker"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	// Register common MIME types to ensure they are served correctly,
	// especially on Windows where registry might be incomplete.
	mime.AddExtensionType(".css", "text/css")
	mime.AddExtensionType(".js", "application/javascript")
	mime.AddExtensionType(".png", "image/png")
	mime.AddExtensionType(".jpg", "image/jpeg")
	mime.AddExtensionType(".jpeg", "image/jpeg")
	mime.AddExtensionType(".gif", "image/gif")
	mime.AddExtensionType(".svg", "image/svg+xml")
	mime.AddExtensionType(".wav", "audio/wav")
	mime.AddExtensionType(".mp3", "audio/mpeg")
	mime.AddExtensionType(".pdf", "application/pdf")
}

func main() {
	mode := flag.String("mode", "all", "Application mode: 'api', 'worker', or 'all'")
	flag.Parse()

	godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	db, err := database.NewDB(dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize Redis Client
	redisClient := redis_client.NewRedisClient()
	log.Println("Connected to Redis")

	// --- WORKER MODE ---
	if *mode == "worker" || *mode == "all" {
		log.Println("Matchmaking worker started...")
		matchmaker := worker.NewMatchmaker(db, redisClient)
		go matchmaker.Run()
	}

	// --- API MODE ---
	if *mode == "api" || *mode == "all" {
		// Initialize Hub
		hub := handlers.NewHub()
		go hub.Run()

		// Handlers
		authHandler := handlers.NewAuthHandler(db)
		matchHandler := handlers.NewMatchHandler(db, redisClient)
		chatHandler := handlers.NewChatHandler(db, hub)
		reportHandler := handlers.NewReportHandler(db)
		adminHandler := handlers.NewAdminHandler(db)
		userHandler := handlers.NewUserHandler(db)

		router := mux.NewRouter()

		// Auth routes
		router.HandleFunc("/api/auth/signup", authHandler.Signup).Methods(http.MethodPost)
		router.HandleFunc("/api/auth/login", authHandler.Login).Methods(http.MethodPost)
		router.HandleFunc("/api/auth/logout", middleware.AuthRequired(authHandler.Logout, db)).Methods(http.MethodPost)
		router.HandleFunc("/api/auth/me", middleware.AuthRequired(authHandler.GetUser, db)).Methods(http.MethodGet)

		// User routes
		router.HandleFunc("/api/users/profile", middleware.AuthRequired(userHandler.GetProfile, db)).Methods(http.MethodGet)
		router.HandleFunc("/api/users/profile", middleware.AuthRequired(userHandler.UpdateProfile, db)).Methods(http.MethodPut)
		router.HandleFunc("/api/users/preferences", middleware.AuthRequired(userHandler.UpdatePreferences, db)).Methods(http.MethodPut)

		// Matchmaking routes
		router.HandleFunc("/api/match/join", middleware.AuthRequired(matchHandler.JoinQueue, db)).Methods(http.MethodPost)
		router.HandleFunc("/api/match/status", middleware.AuthRequired(matchHandler.GetStatus, db)).Methods(http.MethodGet)
		router.HandleFunc("/api/match/leave", middleware.AuthRequired(matchHandler.LeaveQueue, db)).Methods(http.MethodPost)

		// Chat routes
		router.HandleFunc("/ws/chat", func(w http.ResponseWriter, r *http.Request) {
			handlers.HandleWebSocket(hub, w, r, db)
		})
		router.HandleFunc("/api/chat/messages", middleware.AuthRequired(chatHandler.GetMessages, db)).Methods(http.MethodGet)
		router.HandleFunc("/api/chat/send", middleware.AuthRequired(chatHandler.SendMessage, db)).Methods(http.MethodPost)
		router.HandleFunc("/api/chat/upload", middleware.AuthRequired(chatHandler.UploadFile, db)).Methods(http.MethodPost)
		router.HandleFunc("/api/chat/sessions", middleware.AuthRequired(chatHandler.GetSessions, db)).Methods(http.MethodGet)
		router.HandleFunc("/api/chat/sessions/{sessionID}/end", middleware.AuthRequired(chatHandler.EndSession, db)).Methods(http.MethodPost)

		// Serve uploaded files
		router.PathPrefix("/api/uploads/").Handler(http.StripPrefix("/api/uploads/", http.FileServer(http.Dir("./uploads"))))

		// Reporting routes
		router.HandleFunc("/api/report", middleware.AuthRequired(reportHandler.CreateReport, db)).Methods(http.MethodPost)

		// Admin routes
		router.HandleFunc("/api/admin/dashboard", middleware.AuthRequired(adminHandler.GetStats, db)).Methods(http.MethodGet)
		router.HandleFunc("/api/admin/reports", middleware.AuthRequired(adminHandler.GetReports, db)).Methods(http.MethodGet)
		router.HandleFunc("/api/admin/reports/{reportID}/action", middleware.AuthRequired(adminHandler.HandleReport, db)).Methods(http.MethodPost)
		router.HandleFunc("/api/admin/users/{userID}/ban", middleware.AuthRequired(adminHandler.BanUser, db)).Methods(http.MethodPost)

		// CORS middleware
		corsMiddleware := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

				if r.Method == "OPTIONS" {
					w.WriteHeader(http.StatusOK)
					return
				}

				next.ServeHTTP(w, r)
			})
		}
		router.Use(corsMiddleware)
		router.Use(middleware.JSONMiddleware)

		port := os.Getenv("PORT")
		if port == "" {
			port = os.Getenv("SERVER_PORT")
		}
		if port == "" {
			port = "8080"
		}

		fmt.Printf("Server running on http://0.0.0.0:%s\n", port)
		log.Fatal(http.ListenAndServe(":"+port, router))
	} else {
		// If only worker is running, block here so main() doesn't exit
		select {}
	}
}
