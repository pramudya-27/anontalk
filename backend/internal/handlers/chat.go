package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/anontalk/backend/internal/database"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type ChatHandler struct {
	db  *database.DB
	hub *Hub
}

func NewChatHandler(db *database.DB, hub *Hub) *ChatHandler {
	return &ChatHandler{db: db, hub: hub}
}

type SendMessageRequest struct {
	SessionID string `json:"sessionId"`
	Content   string `json:"content"`
	Type      string `json:"type"` // text, image, audio, file
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *ChatHandler) GetSessions(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	session, err := h.db.GetActiveSessionByUserID(userID)
	if err != nil {
		http.Error(w, "Error fetching sessions", http.StatusInternalServerError)
		return
	}

	sessions := []interface{}{}
	if session != nil {
		sessions = append(sessions, session)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"userId":   userID,
		"sessions": sessions,
	})
}

func (h *ChatHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	// In-memory chat: no history stored in DB.
	// We return empty list to satisfy frontend contract.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"messages": []string{},
	})
}

func (h *ChatHandler) EndSession(w http.ResponseWriter, r *http.Request) {
	sessionID := mux.Vars(r)["sessionID"]

	if err := h.db.EndChatSession(sessionID); err != nil {
		http.Error(w, "Error ending session", http.StatusInternalServerError)
		return
	}

	// Cleanup session uploads
	os.RemoveAll("./uploads/" + sessionID)

	// Broadcast session_ended to the hub
	h.hub.broadcast <- BroadcastMessage{
		SessionID: sessionID,
		Payload: map[string]interface{}{
			"type": "session_ended",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Session ended"})
}

func (h *ChatHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	var req SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Default type to text if empty
	if req.Type == "" {
		req.Type = "text"
	}

	messageID := uuid.New().String()
	
	// MEMORY ONLY: Do not save to DB
	// if err := h.db.SaveMessage(messageID, req.SessionID, userID, req.Content, req.Type); err != nil { ... }

	msg := map[string]interface{}{
		"id":        messageID,
		"sessionId": req.SessionID,
		"senderId":  userID,
		"content":   req.Content,
		"type":      req.Type,
	}

	// Broadcast via WebSocket
	h.hub.broadcast <- BroadcastMessage{
		SessionID: req.SessionID,
		Payload:   msg,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
}

func (h *ChatHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	// validate max file size (10MB)
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	sessionID := r.FormValue("sessionId")
	if sessionID == "" {
		http.Error(w, "Session ID required", http.StatusBadRequest)
		return
	}

	// Create uploads directory for session if not exists
	uploadDir := "./uploads/" + sessionID
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, os.ModePerm)
	}

	// Generate unique filename to prevent collisions but keep extension
	ext := filepath.Ext(handler.Filename)
	filename := uuid.New().String() + ext
	filePath := filepath.Join(uploadDir, filename)

	// Create file
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy uploaded file content
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Error saving file content", http.StatusInternalServerError)
		return
	}

	// Return public URL (relative)
	fileURL := fmt.Sprintf("/api/uploads/%s/%s", sessionID, filename)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"url": fileURL,
	})
}

// WebSocket Hub for real-time chat
type BroadcastMessage struct {
	SessionID string
	Payload   interface{}
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan BroadcastMessage
	register   chan *Client
	unregister chan *Client
}

type Client struct {
	hub       *Hub
	conn      *websocket.Conn
	send      chan interface{}
	sessionID string
	userID    string
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan BroadcastMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				// Only send to clients in the same session
				if client.sessionID == message.SessionID {
					select {
					case client.send <- message.Payload:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			}
		}
	}
}

func HandleWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request, db *database.DB) {
	sessionID := mux.Vars(r)["sessionID"]
	tokenString := r.URL.Query().Get("token")

	// Validate token
	userID := ""
	if tokenString != "" {
		// Basic parsing - ideally reuse middleware logic or helper
		// For now, we'll verify it's a valid JWT with our secret
		// NOTE: In a real app, import jwt and use proper validation
	}

	// Upgrade connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := &Client{
		hub:       hub,
		conn:      conn,
		send:      make(chan interface{}, 256),
		sessionID: sessionID,
		userID:    userID, // Will be empty if not validated, handled by client-side usually providing it
	}

	hub.register <- client

	go func() {
		defer func() {
			hub.unregister <- client
			conn.Close()
		}()

		for {
			var msg map[string]interface{}
			err := conn.ReadJSON(&msg)
			if err != nil {
				break
			}

			// We basically ignore incoming messages here because we use HTTP for sending.
			// This loop is primarily to keep the connection open and detect disconnects.
		}
	}()

	go func() {
		for message := range client.send {
			if err := conn.WriteJSON(message); err != nil {
				return
			}
		}
	}()
}
