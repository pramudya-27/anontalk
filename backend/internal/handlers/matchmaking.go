package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/anontalk/backend/internal/database"
	"github.com/anontalk/backend/internal/redis_client"
)

type MatchHandler struct {
	db    *database.DB
	redis *redis_client.RedisClient
}

func NewMatchHandler(db *database.DB, redis *redis_client.RedisClient) *MatchHandler {
	return &MatchHandler{db: db, redis: redis}
}

type JoinQueueRequest struct {
	Gender          string   `json:"gender"`
	PreferredGender []string `json:"preferredGender"`
	Interests       []string `json:"interests"`
}

func (h *MatchHandler) JoinQueue(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	var req JoinQueueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// 1. Store preferences in Redis
	pref := map[string]interface{}{
		"gender":          req.Gender,
		"preferredGender": req.PreferredGender,
		"interests":       req.Interests,
	}
	prefJSON, _ := json.Marshal(pref)

	err := h.redis.Client.Set(redis_client.Ctx, "user_pref:"+userID, prefJSON, 0).Err()
	if err != nil {
		http.Error(w, "redis error", http.StatusInternalServerError)
		return
	}

	// 2. Add to Queue (prevent duplicates)
	h.redis.Client.LRem(redis_client.Ctx, "matchmaking_queue", 0, userID)
	err = h.redis.Client.RPush(redis_client.Ctx, "matchmaking_queue", userID).Err()
	if err != nil {
		http.Error(w, "redis error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Joined queue"})
}

func (h *MatchHandler) LeaveQueue(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	h.redis.Client.LRem(redis_client.Ctx, "matchmaking_queue", 0, userID)
	h.redis.Client.Del(redis_client.Ctx, "user_pref:"+userID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Left queue"})
}

func (h *MatchHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	// Check if user has active session
	session, err := h.db.GetActiveSessionByUserID(userID)
	if err != nil {
		log.Printf("Error getting match status: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if session != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"match":   true,
			"session": session,
		})
	} else {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"match": false,
		})
	}
}
