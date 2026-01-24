package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/anontalk/backend/internal/database"
	"github.com/gorilla/mux"
)

type AdminHandler struct {
	db *database.DB
}

func NewAdminHandler(db *database.DB) *AdminHandler {
	return &AdminHandler{db: db}
}

type HandleReportRequest struct {
	Action    string `json:"action"`
	AdminNotes string `json:"adminNotes"`
}

type BanUserRequest struct {
	Reason string `json:"reason"`
}

func (h *AdminHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.db.GetStats()
	if err != nil {
		http.Error(w, "Error fetching stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (h *AdminHandler) GetReports(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"reports": []interface{}{},
	})
}

func (h *AdminHandler) HandleReport(w http.ResponseWriter, r *http.Request) {
	reportID := mux.Vars(r)["reportID"]

	var req HandleReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"reportId": reportID,
		"action":   req.Action,
		"message":  "Report processed",
	})
}

func (h *AdminHandler) BanUser(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["userID"]

	var req BanUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.db.BanUser(userID, req.Reason); err != nil {
		http.Error(w, "Error banning user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"userId":  userID,
		"banned":  true,
		"reason":  req.Reason,
		"message": "User banned successfully",
	})
}
