package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/o0ga-bo0ga/vigil/internal/api"
	"github.com/o0ga-bo0ga/vigil/internal/service"
	"github.com/o0ga-bo0ga/vigil/internal/store"
)

func healthHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "OK"}`))
}

func main(){

	godotenv.Load()

	var port string

	port = os.Getenv("PORT")
	if port == ""{
		port = "8080"
	}
	port = ":" + port

	r := chi.NewRouter()

	r.Get("/health", healthHandler)

	s := store.NewMemoryStore()
	svc := service.NewJobService(s)
	h := api.NewHandler(svc)

	r.Route("/api", func(r chi.Router) {
		r.Use(api.AuthMiddleware)
		r.Post("/events", h.CreateEvent)
		r.Get("/jobs", h.ListJobs)
	})

	log.Printf("Starting server on %s...", port)
	if err := http.ListenAndServe(port, r); err != nil{
		log.Fatal(err)
	}
}