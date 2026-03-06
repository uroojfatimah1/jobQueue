**Redis-Based Background Job Queue Service**

**Overview**

This project implements a simple asynchronous job queue using Go and Redis.
Clients can create jobs via an API, and a separate worker process asynchronously processes them. The system tracks job metadata, retries failed jobs, and provides job status through a REST API.
---
**Features**

* **Create Job API (POST /v1/jobs)**

  * Accepts type and payload in JSON body.

  * Generates a unique jobId.

  * Stores job metadata in Redis.

  * Pushes job into the Redis queue.

* **Job Status API (GET /v1/jobs/{jobId})**

    * Returns job metadata: type, payload, status, attempts, timestamps.

    * Statuses: queued, processing, completed, failed.

* **Worker Service**

    * Continuously listens to jobs:queue.

    * Marks jobs as processing while executing.

    * Simulates job execution with configurable duration.

    * Updates job status to completed or failed.

    * Retries failed jobs up to 3 times; moves permanently failed jobs to jobs:failed.

* **Redis Structure**

    * job:{jobId} → Redis hash storing job metadata.

    * jobs:queue → Redis list for queued jobs.

    * jobs:failed → Redis list for failed jobs.

* **Logging**

    * Logs job creation, processing, completion, retries, and failures.

* **Simulated Job Failures**

  * 1 in 4 jobs randomly fails to test retry logic.
---
***Tech Stack***

* Go – backend and worker implementation.

* Redis – queue and job metadata storage.

* Mux Router – HTTP routing for job APIs.
---
***Getting Started***
**Prerequisites**

* Go 1.20+ installed.

* Redis server running locally (default localhost:6379).
---
**Installation**

**Clone the repository:**

git clone git@github.com:uroojfatimah1/jobQueue.git
cd jobQueue

**Install dependencies:**

```go mod tidy```

---
**Run API Server**

`go run cmd/server/main.go`
* Exposes endpoints:

    * ```POST /v1/jobs``` → Create a job

    * ```GET /v1/jobs/{jobId}``` → Get job status
---
**Run Worker**

```go run cmd/worker/main.go```

* Continuously processes jobs from jobs:queue.

* Handles retries and failed queue automatically.
---
**API Examples:**

**Create Job**

```
curl -X POST http://localhost:8080/v1/jobs \
-H "Content-Type: application/json" \
-d '{"type":"email","payload":"send farewell mail"}'
```

**Response:**

```
{
"jobId": "123e4567-e89b-12d3-a456-426614174000"
}
```

**Get Job Status**

`curl http://localhost:8080/v1/jobs/123e4567-e89b-12d3-a456-426614174000`

**Response:**

```
{
"id": "123e4567-e89b-12d3-a456-426614174000",
"type": "email",
"payload": "send farewell mail",
"status": "processing",
"attempts": 1,
"createdAt": "2026-03-06T13:00:00Z"
}
```
---
**Job Lifecycle**

* **Queued** – Job added to Redis jobs:queue.
* **Processing** – Worker picks up job and executes.
* **Completed** – Execution finished successfully.
* **Failed** – Execution failed; retried up to 3 times.
* **Failed Queue** – Jobs exceeding retry limit moved to jobs:failed.