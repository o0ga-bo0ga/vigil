package api

import (
	"encoding/json"
	"net/http"

	"github.com/o0ga-bo0ga/vigil/internal/service"
)

type Handler struct {
	service *service.JobService
}

func NewHandler(service *service.JobService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var req service.CreateJobRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Tenant == "" {
		http.Error(w, "Missing required field: name/tenant", http.StatusUnprocessableEntity)
		return
	}

	job, err := h.service.CreateJob(req)
	if err != nil {
		http.Error(w, "Failed to create job: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(job); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *Handler) ListJobs(w http.ResponseWriter, r *http.Request) {
	tenant := r.URL.Query().Get("tenant")

	jobs, err := h.service.ListJobs(tenant)

	if err != nil {
		http.Error(w, "Failed to list jobs: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}