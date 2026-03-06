package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"jobQueue/internal/handlers"
	"jobQueue/internal/redis"
	"jobQueue/internal/service"
)

func main() {

	rdb := redis.NewRedisClient()

	jobService := service.NewJobService(rdb)

	jobHandler := handlers.NewJobHandler(jobService)

	r := mux.NewRouter()
	r.HandleFunc("/v1/jobs", jobHandler.CreateJob).Methods("POST")
	r.HandleFunc("/v1/jobs/{jobId}", jobHandler.GetJob).Methods("GET")

	log.Println("API server running on :8080")

	http.ListenAndServe(":8080", r)
}
