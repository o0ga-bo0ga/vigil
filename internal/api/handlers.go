package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/o0ga-bo0ga/vigil/internal/service"
	"github.com/o0ga-bo0ga/vigil/internal/store"
)

type Handler struct {
	service *service.JobService
	hub     *Hub
}

func NewHandler(service *service.JobService, hub *Hub) *Handler {
	return &Handler{service: service, hub: hub}
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

	jsonJob, err := json.Marshal(job)
	if err != nil {
		log.Println("Failed to marshal job for broadcast: ", err)
	} else {
		h.hub.Broadcast(jsonJob)
	}
}

func (h *Handler) ListJobs(w http.ResponseWriter, r *http.Request) {
	tenant := r.URL.Query().Get("tenant")
	status := r.URL.Query().Get("status")
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 0
	}
	filter := store.ListJobsFilter{
		Tenant: tenant,
		Status: status,
		Limit: limit,
	}

	jobs, err := h.service.ListJobs(filter)

	if err != nil {
		http.Error(w, "Failed to list jobs: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}

func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	tenant := r.URL.Query().Get("tenant")
	stats, err := h.service.GetStats(tenant)

	if err != nil {
		http.Error(w, "Failed to get stats: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}