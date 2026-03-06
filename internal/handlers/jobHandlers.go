package handlers

import (
	"encoding/json"
	"net/http"

	"jobQueue/internal/service"

	"github.com/gorilla/mux"
)

type JobHandler struct {
	Service *service.JobService
}

func NewJobHandler(service *service.JobService) *JobHandler {
	return &JobHandler{Service: service}
}

func (h *JobHandler) CreateJob(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Type    string `json:"type"`
		Payload string `json:"payload"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	jobId, err := h.Service.CreateJob(req.Type, req.Payload)
	if err != nil {
		http.Error(w, "Failed to create job", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"jobId": jobId})
}

func (h *JobHandler) GetJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobId := vars["jobId"]

	job, err := h.Service.GetJob(jobId)
	if err != nil {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(job)
}
