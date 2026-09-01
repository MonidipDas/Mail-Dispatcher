# 📬 Mail Dispatcher

A lightweight **asynchronous email dispatching service built with Go**.

Mail Dispatcher accepts email requests through an API and processes them in the background using a **Job Queue + Worker Pool** architecture.

Instead of making the API wait for an email to be sent, the email request is converted into a job and placed into a queue. Available workers then pick up jobs from the queue and send the emails concurrently.

## 🏗️ Architecture

```text
                  HTTP Request
                       │
                       ▼
                ┌──────────────┐
                │   API Server │
                └──────┬───────┘
                       │
                       │ Create Job
                       ▼
                ┌──────────────┐
                │   Job Queue  │
                │   Channel    │
                └──────┬───────┘
                       │
             ┌─────────┼─────────┐
             │         │         │
             ▼         ▼         ▼
          Worker 1  Worker 2  Worker 3
             │         │         │
             ▼         ▼         ▼
          Send Mail Send Mail Send Mail
```

## 📋 Job Queue

Every email request becomes a **Job**.

For example:

```text
Job {
    To:      "user@gmail.com"
    Subject: "Welcome"
    Body:    "Hello!"
}
```

The API pushes this job into a Go channel:

```go
jobs <- job
```

The channel acts as the **job queue**.

The API can return immediately without waiting for the email to be sent.

```text
Request
   ↓
Create Job
   ↓
Queue Job
   ↓
Return Response
```

## 👷 Worker Pool

The application starts multiple worker goroutines.

For example, with 3 workers:

```text
Worker 1 ─┐
Worker 2 ─┼──→ Job Queue
Worker 3 ─┘
```

Each worker continuously waits for a job:

```go
for job := range jobs {
    sendEmail(job)
}
```

When a job enters the queue, an available worker picks it up.

```text
             Job Queue
          ┌───┬───┬───┐
          │ J1│ J2│ J3│
          └─┬─┴─┬─┴─┬─┘
            │   │   │
            ▼   ▼   ▼
           W1  W2  W3
```

This allows multiple emails to be sent **concurrently**.

## 🔄 Complete Flow

```text
Client
  │
  │ POST /send
  ▼
API Server
  │
  │ Create Email Job
  ▼
Job Queue
  │
  ├──────────────┐
  │              │
  ▼              ▼
Worker 1       Worker 2
  │              │
  ▼              ▼
SMTP Server    SMTP Server
  │              │
  ▼              ▼
Email Sent     Email Sent
```

### Why use this architecture?

Without a worker pool:

```text
Request → Send Email → Wait → Response
```

With a job queue and worker pool:

```text
Request → Queue Job → Response
                    ↓
              Worker Pool
                    ↓
                Send Email
```

This makes email sending **asynchronous**, allows multiple jobs to be processed concurrently, and prevents slow email operations from blocking the API.

## 🛠️ Tech Stack

* **Go**
* **Goroutines**
* **Channels**
* **Worker Pool**
* **Job Queue**
* **SMTP**
* **REST API**
