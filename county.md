# 🧱 Clean Architecture Design Document for Job Portal Backend (Go + PostgreSQL + React)

## 📌 Overview

This document outlines the **clean architecture design** and key system responsibilities for building a scalable, secure, and responsive **Online Job Application Portal**. The backend will be built using **Golang**, the database will be **PostgreSQL (Supabase-hosted)**, and the frontend will be powered by **React**.

---

## 🧭 High-Level Architecture (Clean Architecture Layers)

```
┌───────────────────────────┐
│        Handlers/API       │ ◀── REST Controllers (HTTP)
├───────────────────────────┤
│        Services/UseCases  │ ◀── Business Logic
├───────────────────────────┤
│          Repositories     │ ◀── Data Persistence Interfaces
├───────────────────────────┤
│         Models/Entities   │ ◀── Core Domain Entities
└───────────────────────────┘
```

---

## 🔩 Components Breakdown

### 1. **Entities/Models Layer (Domain)**

* `User`, `Job`, `Application`, `Resume`, `Document`, `Dashboard`, `Notification`, `JobStage`, etc.
* Encapsulates domain rules and relationships.
* Contains validation and transformation logic.

### 2. **Repository Interfaces Layer**

* Abstracts database operations: `UserRepository`, `JobRepository`, `ApplicationRepository`, etc.
* Allows easy mocking and swapping databases (e.g., in-memory vs. PostgreSQL).

### 3. **Service Layer (Use Cases)**

* Encodes application workflows:

  * Register user
  * Apply for a job
  * Track application status
  * Post/edit job
  * Manage dashboards
* Handles validation, context timeouts, and coordination across repositories.

### 4. **Handlers (API Layer)**

* Exposes RESTful API endpoints for frontend (React) and mobile compatibility.
* Authenticates requests (JWT or session-based)
* Calls appropriate services with request data.
* Marshals/unmarshals JSON.

---

## 📂 Directory Structure (Example)

```
/job-portal
│
├── /cmd/api              # Entry point
│
├── /internal
│   ├── /handler          # HTTP handlers
│   ├── /service          # Use case logic
│   ├── /repository       # PostgreSQL implementations
│   ├── /model            # Domain models/entities
│   ├── /validator        # Validation logic
│   └── /email            # Email sending logic
│
├── /migrations           # SQL schema and seeds
├── /docs                 # Documentation
└── go.mod
```

---

## 🛠️ Backend Responsibilities

* ✅ **Account creation**: OTP email verification, secure password storage.
* ✅ **Authentication/Authorization**: Session or token-based auth with middleware.
* ✅ **Job listing + filtering**: Support filters like type, location, department, etc.
* ✅ **Profile management**: Upload CVs (PDF/doc), edit profile, mobile number.
* ✅ **Application submission**: Store resume, track application per stage.
* ✅ **Dashboard support**: Role-based dashboards (HR, Admin, Exec).
* ✅ **Notification System**: Email, in-app messages on events.
* ✅ **Admin features**: Job templates, user management, analytics.
* ✅ **Security & Rate Limiting**: HTTPS, JWT, XSS/CSRF protections.

---

## ✅ Functional Modules

| Module          | Responsibilities                                              |
| --------------- | ------------------------------------------------------------- |
| Admin Dashboard | Manage users, jobs, configurations, access control            |
| HR Dashboard    | Post jobs, track applicants, manage interview stages          |
| User Module     | Register, edit profile, apply, upload resume, receive updates |
| Records Module  | Log audit trails, system changes, logs                        |
| Exec Dashboard  | High-level metrics, decisions, recruitment KPIs               |

---

## 🔒 Security Practices

* SSL (always-on Supabase enforced)
* Passwords hashed with bcrypt
* Input validation (backend and frontend)
* JWT expiration and renewal strategies
* Email verification
* File upload constraints (type, size, storage)

---

## 🧪 Performance and Load Testing

* Use tools like `k6`, `Artillery`, or `Locust`.
* Simulate:

  * 1,000 concurrent applications
  * 500 concurrent dashboard users
  * 10K+ search queries per minute
* DB benchmarking: Concurrent reads/writes (e.g., pgbench)

---

## 📊 Key Features Checklist

* [x] Account Registration + Email Verification
* [x] Search and Filter Jobs (Advanced filters)
* [x] Role-based Dashboards
* [x] Application Tracking with Status Timeline
* [x] Resume Upload & Document Parsing
* [x] Audio/Video Interview Scheduling
* [x] ATS-like features for HR
* [x] Analytics: Views, CTR, funnel conversion
* [x] CMS-like Marketing Page Builder (optional)
* [x] Employer Branding Options

---

## 🧱 Technology Stack

| Layer    | Tech                    |
| -------- | ----------------------- |
| Frontend | React, Tailwind, Vite   |
| Backend  | Golang (net/http)       |
| Database | PostgreSQL (Supabase)   |
| Auth     | JWT or Supabase Auth    |
| Cache    | Redis (optional)        |
| Hosting  | Railway, Fly.io, Render |

---

## 🚀 Deployment and CI/CD (Optional Plan)

* GitHub Actions + Docker for testing and deployment.
* Supabase for managed PostgreSQL.
* Railway/Render/Fly.io for app deployment.

---

