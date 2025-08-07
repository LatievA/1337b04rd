# 1337b04rd

**Anonymous Imageboard for "Hackers"**  
A Go-powered web application that allows users to create threads, post comments, and share images anonymously. Built with **Hexagonal Architecture**, PostgreSQL, and S3-compatible storage.

---

## Installation
```bash
# Clone the repository
git clone https://github.com/LatievA/1337b04rd.git
cd 1337b04rd

# Run docker then use next command
docker-compose up --build

# Wait till all containers are up
# Run migrations (create tables in PostgreSQL)
go run ./cmd/migrations/migrate.go
```
**After app container is up go to the localhost:8081/catalog**

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
- ✅ Auto-archive posts:
  - Posts without comments → archive after **10 min**
  - Posts with comments → archive **15 min after last comment**
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
