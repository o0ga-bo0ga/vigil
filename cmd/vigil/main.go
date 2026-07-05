package main

import (
	"net/http"
	"os"
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
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

	log.Printf("Starting server on %s...", port)
	if err := http.ListenAndServe(port, r); err != nil{
		log.Fatal(err)
	}
}