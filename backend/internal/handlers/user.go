package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/anontalk/backend/internal/database"
	"github.com/google/uuid"
)

type UserHandler struct {
	db *database.DB
}

func NewUserHandler(db *database.DB) *UserHandler {
	return &UserHandler{db: db}
}

type ProfileRequest struct {
	DisplayName string `json:"displayName"`
	AgeRange    string `json:"ageRange"`
	LookingFor  string `json:"lookingFor"`
}

type PreferencesRequest struct {
	PreferredGender []string `json:"preferredGender"`
	Interests       []string `json:"interests"`
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	user, err := h.db.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	var req ProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	profileID := uuid.New().String()
	if err := h.db.UpdateUserProfile(userID, req.DisplayName, req.AgeRange, req.LookingFor); err != nil {
		http.Error(w, "Error updating profile", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Profile updated",
		"id":      profileID,
	})
}

func (h *UserHandler) UpdatePreferences(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	var req PreferencesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.db.UpdateUserPreferences(userID, req.PreferredGender, req.Interests); err != nil {
		http.Error(w, "Error updating preferences", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Preferences updated"})
}
