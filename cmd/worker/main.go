package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v8"

	"jobQueue/internal/service"
)

var ctx = context.Background()

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	jobService := service.NewJobService(rdb)

	log.Println("Job Worker started...")

	rand.Seed(time.Now().UnixNano())

	for {
		processJob(jobService, rdb)
	}
}

func processJob(jobService *service.JobService, rdb *redis.Client) {

	// Wait for a job from the queue
	jobIds, err := rdb.BLPop(ctx, 0, "jobs:queue").Result()
	if err != nil {
		log.Println("Failed to pop from queue:", err)
		return
	}

	if len(jobIds) < 2 {
		return
	}

	jobId := jobIds[1]

	// Get job metadata
	job, err := jobService.GetJob(jobId)
	if err != nil {
		log.Println("Job not found:", jobId)
		return
	}

	key := fmt.Sprintf("job:%s", jobId)

	// Update status to processing
	rdb.HSet(ctx, key, "status", "processing", "updatedAt", time.Now().Format(time.RFC3339))
	log.Println("Processing job:", jobId, "Type:", job.Type, "Payload:", job.Payload)

	// Simulate execution duration based on a job type
	duration := jobDuration(job.Type)
	time.Sleep(duration)

	// Randomly fail 1 in 4 jobs
	if rand.Intn(4) == 0 {
		job.Attempts++
		rdb.HSet(ctx, key,
			"status", "failed",
			"attempts", job.Attempts,
			"updatedAt", time.Now().Format(time.RFC3339),
		)
		log.Println("Job failed:", jobId, "Attempt:", job.Attempts)

		if job.Attempts < 3 {
			// Retry: push back to queue
			rdb.RPush(ctx, "jobs:queue", jobId)
			log.Println("Job re-queued for retry:", jobId)
		} else {
			// Move to the failed queue after 3 attempts
			rdb.RPush(ctx, "jobs:failed", jobId)
			log.Println("Job moved to failed queue:", jobId)
		}
		return
	}

	// Success
	job.Attempts++
	rdb.HSet(ctx, key,
		"status", "completed",
		"attempts", job.Attempts,
		"updatedAt", time.Now().Format(time.RFC3339),
	)
	log.Println("Job completed:", jobId)
}

// Simulate job duration per type
func jobDuration(jobType string) time.Duration {
	switch jobType {
	case "email":
		return 45 * time.Second
	case "report":
		return 35 * time.Second
	case "image":
		return 10 * time.Second
	default:
		return 15 * time.Second
	}
}
