# 1337b04rd

**Anonymous Imageboard for Hackers**  
A Go-powered web application that allows users to create threads, post comments, and share images anonymously. Built with **Hexagonal Architecture**, PostgreSQL, and S3-compatible storage.

---

## Table of Contents
- [Introduction](#introduction)
- [Features](#features)
- [Tech Stack](#tech-stack)
- [Architecture](#architecture)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [Frontend Templates](#frontend-templates)
- [Testing](#testing)
- [Logging](#logging)
- [Troubleshooting](#troubleshooting)
- [Contributors](#contributors)
- [License](#license)

---

## Introduction
**1337b04rd** is an anonymous imageboard inspired by early internet forums but with a modern twist. It allows users to create threads, upload images, comment on posts, and interact without registration. Each user is assigned a unique **avatar and name** from the Rick and Morty API, managed through browser sessions using secure cookies.

---

## Features
- ✅ Anonymous posting (no registration required)
- ✅ Create threads with images
- ✅ Comment on posts and reply to other comments
- ✅ Image upload using **S3-compatible storage**
- ✅ PostgreSQL-based persistent storage for posts, comments, and sessions
- ✅ Unique user avatars & names from **Rick and Morty API**
- ✅ **Hexagonal Architecture** for clean separation of concerns
- ✅ **RESTful API** backend in Go
- ✅ **Session management** using secure HTTP cookies
- ✅ Auto-delete posts:
  - Posts without comments → delete after **10 min**
  - Posts with comments → delete **15 min after last comment**
- ✅ Logging with Go's `log/slog`
- ✅ Minimum **20% test coverage**

---

## Tech Stack
- **Language:** Go 1.21+
- **Database:** PostgreSQL
- **Storage:** S3-compatible (e.g., MinIO)
- **Frontend:** HTML templates (6 views provided)
- **External API:** [Rick and Morty API](https://rickandmortyapi.com)
- **Architecture:** Hexagonal (Ports & Adapters)

---

## Architecture
The project follows **Hexagonal Architecture** to separate core business logic from external systems like databases, APIs, and the web layer.

### Layers:
- **Domain Layer:** Core logic (posts, comments, sessions)
- **Application Layer:** Services and use cases
- **Adapters:**
  - **Database Adapter:** PostgreSQL implementation
  - **Storage Adapter:** S3 image storage
  - **External API Adapter:** Rick and Morty avatar provider
  - **HTTP Adapter:** REST API and session handling

Benefits:
- Testable
- Maintainable
- Flexible

---

## Installation

### Prerequisites:
- Go 1.21+
- PostgreSQL
- MinIO or any S3-compatible storage
- Git

### Steps:
```bash
# Clone the repository
git clone https://github.com/yourusername/1337b04rd.git
cd 1337b04rd

# Configure environment variables
cp .env.example .env
# (Fill in the required values)

# Build the application
go build -o 1337b04rd ./cmd/1337b04rd

# Run migrations (create tables in PostgreSQL)
go run ./cmd/migrate

# Start the server
./1337b04rd --port 8080
