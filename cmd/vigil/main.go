package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/o0ga-bo0ga/vigil/internal/api"
	"github.com/o0ga-bo0ga/vigil/internal/service"
	"github.com/o0ga-bo0ga/vigil/internal/static"
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

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "vigil.db"
	}
	s, err := store.NewSQLiteStore(dbPath)
	if err != nil {
		log.Fatal(err)
	}

	hub := api.NewHub()
	go hub.Run()
	svc := service.NewJobService(s)
	h := api.NewHandler(svc, hub)

	r.Handle("/static/*",
				http.StripPrefix("/static/",
								http.FileServer(http.FS(static.Files))))

	r.Get("/ws", hub.HandleWebSocket)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		data, _ := static.Files.ReadFile("index.html")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(data)
	})

	r.Route("/api", func(r chi.Router) {
		r.Use(api.AuthMiddleware)
		r.Post("/events", h.CreateEvent)
		r.Get("/jobs", h.ListJobs)
		r.Get("/stats", h.GetStats)
	})

	log.Printf("Starting server on %s...", port)
	if err := http.ListenAndServe(port, r); err != nil{
		log.Fatal(err)
	}
}