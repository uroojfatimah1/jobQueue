package service

import (
	"fmt"
	"strconv"
	"time"

	"jobQueue/internal/model"
	redisClient "jobQueue/internal/redis"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type JobService struct {
	Rdb *redis.Client
}

func NewJobService(rdb *redis.Client) *JobService {
	return &JobService{Rdb: rdb}
}

func (s *JobService) CreateJob(jobType string, payload string) (string, error) {
	jobId := uuid.New().String()
	job := models.Job{
		ID:        jobId,
		Type:      jobType,
		Payload:   payload,
		Status:    "queued",
		Attempts:  0,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	key := fmt.Sprintf("job:%s", jobId)

	// Store metadata
	err := s.Rdb.HSet(redisClient.Ctx, key, map[string]interface{}{
		"type":      job.Type,
		"payload":   job.Payload,
		"status":    job.Status,
		"attempts":  job.Attempts,
		"createdAt": job.CreatedAt,
	}).Err()
	if err != nil {
		return "", fmt.Errorf("failed to store job: %v", err)
	}

	// Push jobId to queue
	err = s.Rdb.RPush(redisClient.Ctx, "jobs:queue", jobId).Err()
	if err != nil {
		return "", fmt.Errorf("failed to push job to queue: %v", err)
	}

	return jobId, nil
}

func (s *JobService) GetJob(jobId string) (*models.Job, error) {
	key := fmt.Sprintf("job:%s", jobId)

	data, err := s.Rdb.HGetAll(redisClient.Ctx, key).Result()
	if err != nil || len(data) == 0 {
		return nil, fmt.Errorf("job not found")
	}

	attempts, _ := strconv.Atoi(data["attempts"])

	job := &models.Job{
		ID:        jobId,
		Type:      data["type"],
		Payload:   data["payload"],
		Status:    data["status"],
		Attempts:  attempts,
		CreatedAt: data["createdAt"],
	}

	return job, nil
}
