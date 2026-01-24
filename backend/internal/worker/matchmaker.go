package worker

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/anontalk/backend/internal/database"
	"github.com/anontalk/backend/internal/redis_client"
	"github.com/google/uuid"
)

type Matchmaker struct {
	db    *database.DB
	redis *redis_client.RedisClient
}

func NewMatchmaker(db *database.DB, redis *redis_client.RedisClient) *Matchmaker {
	return &Matchmaker{db: db, redis: redis}
}

func (m *Matchmaker) Run() {
	fmt.Println("Matchmaking worker started...")
	for {
		// 1. Pop user from queue
		val, err := m.redis.Client.LPop(redis_client.Ctx, "matchmaking_queue").Result()
		if err != nil {
			// Queue empty or error
			time.Sleep(1 * time.Second)
			continue
		}
		user1ID := val

		// 2. Try to find a match
		// Simple approach: Pop another user. Ideally, we look for compatible match.
		// For improved matching (gender prefs), we would need to inspect queue items without popping,
		// or maintain multiple queues (e.g., 'male_seeking_female', etc).
		
		// Let's implement the logic requested: "Worker ambil dua session dari queue. Worker cek preferensi gender dua arah."
		// Note: LPOP is destructive. If we pop user2 and they don't match user1, we must returning them to queue.

		val2, err := m.redis.Client.LPop(redis_client.Ctx, "matchmaking_queue").Result()
		if err != nil {
			// Only one user in queue, push back user1 and wait
			m.redis.Client.RPush(redis_client.Ctx, "matchmaking_queue", user1ID)
			time.Sleep(1 * time.Second)
			continue
		}
		user2ID := val2

		if user1ID == user2ID {
			// Should not happen unless bad data
			continue 
		}

		// 3. Check preferences
		if m.checkMatch(user1ID, user2ID) {
			// 4. Create Session
			sessionID := uuid.New().String()
			err := m.db.CreateChatSession(sessionID, user1ID, user2ID)
			if err != nil {
				fmt.Printf("Error creating session: %v\n", err)
				// Refund users to queue?
				m.redis.Client.RPush(redis_client.Ctx, "matchmaking_queue", user1ID)
				m.redis.Client.RPush(redis_client.Ctx, "matchmaking_queue", user2ID)
			} else {
				fmt.Printf("Matched %s and %s in session %s\n", user1ID, user2ID, sessionID)
				// Cleanup redis data
				m.redis.Client.Del(redis_client.Ctx, "user_pref:"+user1ID)
				m.redis.Client.Del(redis_client.Ctx, "user_pref:"+user2ID)
			}
		} else {
			// No match, push back.
			// To avoid infinite loop of popping same pair, maybe push to back?
			m.redis.Client.RPush(redis_client.Ctx, "matchmaking_queue", user1ID)
			m.redis.Client.RPush(redis_client.Ctx, "matchmaking_queue", user2ID)
			
			// Sleep briefly to allow other permutations if queue is small
			time.Sleep(100 * time.Millisecond)
		}
	}
}

type UserPref struct {
	Gender          string   `json:"gender"`
	PreferredGender []string `json:"preferredGender"`
	Interests       []string `json:"interests"`
}

func (m *Matchmaker) checkMatch(u1, u2 string) bool {
	// Fetch prefs from Redis
	p1Str, err := m.redis.Client.Get(redis_client.Ctx, "user_pref:"+u1).Result()
	if err != nil {
		return false // Prefs missing, maybe stale queue item
	}
	p2Str, err := m.redis.Client.Get(redis_client.Ctx, "user_pref:"+u2).Result()
	if err != nil {
		return false
	}

	var p1, p2 UserPref
	json.Unmarshal([]byte(p1Str), &p1)
	json.Unmarshal([]byte(p2Str), &p2)

	// Check U1 pref for U2 gender
	match1 := false
	for _, g := range p1.PreferredGender {
		if g == "any" || g == p2.Gender {
			match1 = true
			break
		}
	}

	// Check U2 pref for U1 gender
	match2 := false
	for _, g := range p2.PreferredGender {
		if g == "any" || g == p1.Gender {
			match2 = true
			break
		}
	}

	return match1 && match2
}
