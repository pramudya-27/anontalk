package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	
	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	conn *sql.DB
}

func NewDB(dsn string) (*DB, error) {
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = conn.Ping()
	if err != nil {
		return nil, err
	}

	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(5)

	return &DB{conn: conn}, nil
}

func (d *DB) Close() error {
	return d.conn.Close()
}

func (d *DB) GetConnection() *sql.DB {
	return d.conn
}

// User methods
func (d *DB) CreateUser(id, email, passwordHash, gender string, interests []string, preferredGender []string) error {
	interestsJSON, _ := json.Marshal(interests)
	preferredGenderJSON, _ := json.Marshal(preferredGender)

	query := `
		INSERT INTO users (id, email, password_hash, gender, interests, preferred_gender, is_verified)
		VALUES (?, ?, ?, ?, ?, ?, true)
	`
	_, err := d.conn.Exec(query, id, email, passwordHash, gender, interestsJSON, preferredGenderJSON)
	return err
}

func (d *DB) GetUserByEmail(email string) (map[string]interface{}, error) {
	query := `
		SELECT id, email, password_hash, gender, interests, preferred_gender, is_banned
		FROM users WHERE email = ?
	`
	row := d.conn.QueryRow(query, email)

	var user map[string]interface{}
	var id, email2, passwordHash, gender string
	var interests, preferredGender sql.NullString
	var isBanned bool

	err := row.Scan(&id, &email2, &passwordHash, &gender, &interests, &preferredGender, &isBanned)
	if err != nil {
		return nil, err
	}

	user = map[string]interface{}{
		"id":               id,
		"email":            email2,
		"passwordHash":     passwordHash,
		"gender":           gender,
		"isBanned":         isBanned,
	}

	return user, nil
}

func (d *DB) GetUserByID(userID string) (map[string]interface{}, error) {
	query := `
		SELECT u.id, u.email, u.gender, u.interests, u.preferred_gender, u.is_verified, u.is_banned,
		       p.display_name, p.age_range, p.looking_for
		FROM users u
		LEFT JOIN user_profiles p ON u.id = p.user_id
		WHERE u.id = ?
	`
	row := d.conn.QueryRow(query, userID)

	var id, email, gender string
	var interests, preferredGender sql.NullString // captured as JSON string
	var isVerified, isBanned bool
	var displayName, ageRange, lookingFor sql.NullString

	err := row.Scan(&id, &email, &gender, &interests, &preferredGender, &isVerified, &isBanned, &displayName, &ageRange, &lookingFor)
	if err != nil {
		return nil, err
	}

	// Parse JSON fields
	var interestsList []string
	if interests.Valid {
		json.Unmarshal([]byte(interests.String), &interestsList)
	}

	var preferredGenderList []string
	if preferredGender.Valid {
		json.Unmarshal([]byte(preferredGender.String), &preferredGenderList)
	}

	user := map[string]interface{}{
		"id":              id,
		"email":           email,
		"gender":          gender,
		"isVerified":      isVerified,
		"isBanned":        isBanned,
		"displayName":     displayName.String,
		"ageRange":        ageRange.String,
		"lookingFor":      lookingFor.String,
		"interests":       interestsList,
		"preferredGender": preferredGenderList,
	}

	return user, nil
}

func (d *DB) UpdateUserProfile(userID, displayName, ageRange, lookingFor string) error {
	query := `
		INSERT INTO user_profiles (id, user_id, display_name, age_range, looking_for)
		VALUES (UUID(), ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE display_name = VALUES(display_name), age_range = VALUES(age_range), looking_for = VALUES(looking_for)
	`
	_, err := d.conn.Exec(query, userID, displayName, ageRange, lookingFor)
	return err
}

func (d *DB) UpdateUserPreferences(userID string, preferredGender, interests []string) error {
	preferredGenderJSON, _ := json.Marshal(preferredGender)
	interestsJSON, _ := json.Marshal(interests)

	query := `
		UPDATE users SET preferred_gender = ?, interests = ?
		WHERE id = ?
	`
	_, err := d.conn.Exec(query, preferredGenderJSON, interestsJSON, userID)
	return err
}

// Chat session methods
func (d *DB) CreateChatSession(sessionID, user1ID, user2ID string) error {
	query := `
		INSERT INTO chat_sessions (id, user1_id, user2_id, is_active)
		VALUES (?, ?, ?, true)
	`
	_, err := d.conn.Exec(query, sessionID, user1ID, user2ID)
	return err
}

func (d *DB) GetChatSession(sessionID string) (map[string]interface{}, error) {
	query := `
		SELECT id, user1_id, user2_id, started_at, ended_at, is_active
		FROM chat_sessions WHERE id = ?
	`
	row := d.conn.QueryRow(query, sessionID)

	var id, user1ID, user2ID string
	var startedAt, endedAt sql.NullTime
	var isActive bool

	err := row.Scan(&id, &user1ID, &user2ID, &startedAt, &endedAt, &isActive)
	if err != nil {
		return nil, err
	}

	session := map[string]interface{}{
		"id":       id,
		"user1Id":  user1ID,
		"user2Id":  user2ID,
		"isActive": isActive,
	}

	return session, nil
}

func (d *DB) GetActiveSessionByUserID(userID string) (map[string]interface{}, error) {
	query := `
		SELECT id, user1_id, user2_id, started_at, is_active
		FROM chat_sessions 
		WHERE (user1_id = ? OR user2_id = ?) AND is_active = true
		LIMIT 1
	`
	row := d.conn.QueryRow(query, userID, userID)

	var id, user1ID, user2ID string
	var startedAt sql.NullTime
	var isActive bool

	err := row.Scan(&id, &user1ID, &user2ID, &startedAt, &isActive)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	session := map[string]interface{}{
		"id":       id,
		"user1Id":  user1ID,
		"user2Id":  user2ID,
		"isActive": isActive,
	}

	return session, nil
}

func (d *DB) EndChatSession(sessionID string) error {
	query := `
		UPDATE chat_sessions SET is_active = false, ended_at = NOW()
		WHERE id = ?
	`
	_, err := d.conn.Exec(query, sessionID)
	return err
}

// Message methods
func (d *DB) SaveMessage(messageID, sessionID, senderID, content, msgType string) error {
	query := `
		INSERT INTO messages (id, session_id, sender_id, content, type)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := d.conn.Exec(query, messageID, sessionID, senderID, content, msgType)
	return err
}

func (d *DB) GetMessages(sessionID string, limit, offset int) ([]map[string]interface{}, error) {
	query := `
		SELECT id, sender_id, content, type, created_at
		FROM messages WHERE session_id = ?
		ORDER BY created_at ASC
		LIMIT ? OFFSET ?
	`
	rows, err := d.conn.Query(query, sessionID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []map[string]interface{}
	for rows.Next() {
		var id, senderID, content, msgType string
		var createdAt sql.NullTime

		err := rows.Scan(&id, &senderID, &content, &msgType, &createdAt)
		if err != nil {
			return nil, err
		}

		message := map[string]interface{}{
			"id":       id,
			"senderId": senderID,
			"content":  content,
			"type":     msgType,
		}

		messages = append(messages, message)
	}

	return messages, nil
}

// Matchmaking queue methods
func (d *DB) AddToQueue(queueID, userID, gender string, preferredGender, interests []string) error {
	preferredGenderJSON, _ := json.Marshal(preferredGender)
	interestsJSON, _ := json.Marshal(interests)

	// Use ON DUPLICATE KEY UPDATE to handle cases where user is already in queue
	// This refreshes their spot in the queue and updates preferences
	query := `
		INSERT INTO matchmaking_queue (id, user_id, gender, preferred_gender, interests)
		VALUES (?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE 
			id = VALUES(id),
			gender = VALUES(gender),
			preferred_gender = VALUES(preferred_gender),
			interests = VALUES(interests),
			joined_at = CURRENT_TIMESTAMP
	`
	_, err := d.conn.Exec(query, queueID, userID, gender, preferredGenderJSON, interestsJSON)
	return err
}

func (d *DB) RemoveFromQueue(userID string) error {
	query := `DELETE FROM matchmaking_queue WHERE user_id = ?`
	_, err := d.conn.Exec(query, userID)
	return err
}

func (d *DB) FindMatch(userID, gender string, preferredGender []string) (map[string]interface{}, error) {
	// Simple matching logic - can be enhanced
	query := `
		SELECT user_id, gender, interests FROM matchmaking_queue
		WHERE user_id != ? AND gender IN (?)
		ORDER BY joined_at ASC
		LIMIT 1
	`

	var preferredGenderStr string
	for i, g := range preferredGender {
		if i > 0 {
			preferredGenderStr += ","
		}
		preferredGenderStr += "'" + g + "'"
	}

	query = fmt.Sprintf(`
		SELECT user_id, gender, interests FROM matchmaking_queue
		WHERE user_id != ? AND gender IN (%s)
		ORDER BY joined_at ASC
		LIMIT 1
	`, preferredGenderStr)

	row := d.conn.QueryRow(query, userID)

	var matchUserID, matchGender string
	var interests sql.NullString

	err := row.Scan(&matchUserID, &matchGender, &interests)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	match := map[string]interface{}{
		"userId": matchUserID,
		"gender": matchGender,
	}

	return match, nil
}

// Report methods
func (d *DB) CreateReport(reportID, reporterID, reportedUserID, reason, description string) error {
	query := `
		INSERT INTO reports (id, reporter_id, reported_user_id, reason, description, status)
		VALUES (?, ?, ?, ?, ?, 'pending')
	`
	_, err := d.conn.Exec(query, reportID, reporterID, reportedUserID, reason, description)
	return err
}

func (d *DB) BanUser(userID, reason string) error {
	query := `
		UPDATE users SET is_banned = true, ban_reason = ?, banned_at = NOW()
		WHERE id = ?
	`
	_, err := d.conn.Exec(query, reason, userID)
	return err
}

// Admin methods
func (d *DB) GetStats() (map[string]interface{}, error) {
	var totalUsers, activeChats, totalMessages int

	d.conn.QueryRow("SELECT COUNT(*) FROM users").Scan(&totalUsers)
	d.conn.QueryRow("SELECT COUNT(*) FROM chat_sessions WHERE is_active = true").Scan(&activeChats)
	d.conn.QueryRow("SELECT COUNT(*) FROM messages").Scan(&totalMessages)

	stats := map[string]interface{}{
		"totalUsers":    totalUsers,
		"activeChats":   activeChats,
		"totalMessages": totalMessages,
	}

	return stats, nil
}
