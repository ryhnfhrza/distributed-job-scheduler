# 📦 Distributed Job Scheduler

A scalable distributed background job processing system built with Go, MySQL, Redis, and Docker.  
This project allows users to create scheduled jobs through an API, queue them using Redis, and execute them asynchronously using distributed workers.

Designed with a microservice architecture consisting of:

- **API Service** → Manage jobs through REST API
- **Scheduler Service** → Move scheduled jobs from MySQL to Redis queue
- **Worker Service** → Consume and execute jobs asynchronously

Supports:

- Scheduled jobs
- Email jobs
- Webhook jobs
- Retry mechanism
- Job logs
- Dockerized deployment
- Database migration
- Distributed worker processing

---

# 🚀 Features

- ✅ Create scheduled background jobs
- ✅ Update & delete pending jobs
- ✅ Distributed worker architecture
- ✅ Redis queue processing
- ✅ Email job execution
- ✅ Webhook job execution
- ✅ Retry mechanism with retry delay
- ✅ Job execution logs
- ✅ MySQL persistence
- ✅ Docker & Docker Compose support
- ✅ Database migration support
- ✅ UTC-based scheduling system
- ✅ Clean architecture structure
- ✅ RESTful API

---

# 🛠️ Tech Stack

- Go (Golang)
- MySQL
- Redis
- Docker
- Docker Compose
- Golang Migrate
- REST API
- Microservices Architecture

---

# 📋 Architecture

```text
                +------------------+
                |   API SERVICE    |
                |------------------|
                | Create Job       |
                | Update Job       |
                | Delete Job       |
                | Get Jobs         |
                +---------+--------+
                          |
                          v
                    +-----------+
                    |  MySQL DB |
                    +-----------+
                          |
                          v
                +------------------+
                | SCHEDULER SERVICE|
                |------------------|
                | Fetch due jobs   |
                | Mark as queued   |
                | Push to Redis    |
                +---------+--------+
                          |
                          v
                    +-----------+
                    |   Redis   |
                    | Job Queue |
                    +-----------+
                          |
                          v
                +------------------+
                |  WORKER SERVICE  |
                |------------------|
                | Consume Queue    |
                | Execute Job      |
                | Retry Failed Job |
                | Insert Logs      |
                +------------------+

# 📦 Prerequisites

Before running this project, make sure you have installed:

- Go >= 1.22
- Docker
- Docker Compose
- Git

Optional for local development:

- MySQL
- Redis

---

# ⚙️ Installation & Setup

## 1️⃣ Clone Repository

```bash
git clone https://github.com/yourusername/distributed-job-scheduler.git

cd distributed-job-scheduler
```

---

## 2️⃣ Setup Environment Variables

Create `.env` file in the root project directory.

Example:

```env
# =========================
# REDIS CONFIG
# =========================
REDIS_HOST=job_scheduler_redis
REDIS_PORT=6379
REDIS_PASSWORD=adminStrongPass

# =========================
# APP CONFIG
# =========================
APP_PORT=8080

# =========================
# MYSQL LOCAL CONFIG
# =========================
DB_USER=root
DB_PASSWORD=Password22Rahasia
DB_HOST=localhost
DB_PORT=3306
DB_NAME=job_scheduler_db

# =========================
# MYSQL DOCKER CONFIG
# =========================
DB_HOST=job_scheduler_db
DB_PORT=3306
DB_HOST_PORT=3307
DB_USER=job_scheduler_user
DB_PASSWORD=Password22Rahasia
DB_NAME=job_scheduler_db
MYSQL_ROOT_PASSWORD=root_password
```

---

Create another `.env` file inside:

```text
worker-service/internal/config/.env
```

Example:

```env
# =========================
# SMTP CONFIG
# =========================
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_EMAIL=example@gmail.com
SMTP_PASSWORD="your_app_password"
```

# 🗄️ Database Migration

This project uses `golang-migrate`.

## Run Migration

```bash
docker compose up migration
```

Or manually:

```bash
migrate -path ./migrations \
-database "mysql://root:password@tcp(localhost:3306)/distributed_job_scheduler" up
```

---

# ▶️ Usage

# 🧪 Development Mode

## Run API Service

```bash
cd api-service
go run cmd/main.go
```

## Run Scheduler Service

```bash
cd scheduler-service
go run cmd/main.go
```

## Run Worker Service

```bash
cd worker-service
go run cmd/main.go
```

---

# 🐳 Production Mode (Docker)

## Build & Run

```bash
docker compose up --build
```

## Stop Containers

```bash
docker compose down
```

---

# 📡 API Endpoints

## Job Endpoints

| Method | Endpoint | Description |
|---|---|---|
| POST | `/api/v1/jobs` | Create new job |
| PUT | `/api/v1/jobs/:id` | Update pending job |
| DELETE | `/api/v1/jobs/:id` | Delete pending job |
| GET | `/api/v1/jobs/:id` | Get job by ID |
| GET | `/api/v1/jobs` | Get all jobs |
| GET | `/api/v1/jobs/:id/logs` | Get job logs |

---

# 📬 Example Email Job Payload

```json
{
  "name": "Send Welcome Email",
  "type": "EMAIL",
  "run_at": "2026-05-12T20:00:00+08:00",
  "payload": {
    "to": "user@gmail.com",
    "subject": "Welcome",
    "body": "Welcome to our platform"
  }
}
```

---

# 🌐 Example Webhook Job Payload

```json
{
  "name": "Webhook Notification",
  "type": "WEBHOOK",
  "run_at": "2026-05-12T20:00:00+08:00",
  "payload": {
    "url": "https://webhook.site/your-id",
    "method": "POST",
    "headers": {
      "Content-Type": "application/json"
    },
    "body": {
      "message": "hello"
    }
  }
}
```

---

# 🧠 Retry Mechanism

When a job execution fails:

- Retry count will increase
- Job status changes back to `pending`
- Scheduler will requeue the job
- Retry delay increases based on retry count
- If retry exceeds `max_retry`, job becomes `failed`

---

# 📁 Project Structure

```text
distributed-job-scheduler/
│
├── api-service/
│   ├── cmd/
│   ├── internal/
│   │   ├── app/
│   │   ├── controller/
│   │   ├── exception/
│   │   ├── helper/
│   │   ├── model/
│   │   │   ├── domain/
│   │   │   └── web/
│   │   ├── repository/
│   │   ├── service/
|   |   └── util
│   └── Dockerfile
│
├── scheduler-service/
│   ├── cmd/
│   ├── internal/
│   │   ├── app/
│   │   ├── model/
│   │   ├── repository/
│   │   └── service/
│   └── Dockerfile
│
├── worker-service/
│   ├── cmd/
│   ├── internal/
│   │   ├── app/
│   │   ├── config/
│   │   ├── helper/
│   │   ├── model/
│   │   ├── repository/
│   │   ├── service/
│   │   └── template/
│   └── Dockerfile
│
├── migrations/
├── docker-compose.yml
├── .env
└── README.md
```
---

# 🕒 Timezone Handling

This project uses:

```text
UTC-based scheduling
```

All jobs are normalized to UTC internally to ensure consistent scheduling across different user timezones and server environments.

---

# 🧪 Testing Redis Queue

Check Redis queue manually:

```bash
docker exec -it job_scheduler_redis redis-cli
```

View queue:

```bash
LRANGE job_queue 0 -1
```

---

# 📜 Contributing

Contributions are welcome.

## Steps

1. Fork repository
2. Create new feature branch

```bash
git checkout -b feature/new-feature
```

3. Commit changes

```bash
git commit -m "Add new feature"
```

4. Push branch

```bash
git push origin feature/new-feature
```

5. Open Pull Request

# 👨‍💻 Author

Developed by Rayhan F.
