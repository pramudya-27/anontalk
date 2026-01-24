package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/anontalk/backend/internal/database"
	"github.com/google/uuid"
)

type ReportHandler struct {
	db *database.DB
}

func NewReportHandler(db *database.DB) *ReportHandler {
	return &ReportHandler{db: db}
}

type CreateReportRequest struct {
	ReportedUserID string `json:"reportedUserId"`
	Reason         string `json:"reason"`
	Description    string `json:"description"`
}

func (h *ReportHandler) CreateReport(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	var req CreateReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	reportID := uuid.New().String()
	if err := h.db.CreateReport(reportID, userID, req.ReportedUserID, req.Reason, req.Description); err != nil {
		http.Error(w, "Error creating report", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      reportID,
		"status":  "pending",
		"message": "Report submitted",
	})
}
