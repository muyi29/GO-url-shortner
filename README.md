# URL Shortener 🔗

A production-quality URL shortener built with Go, featuring clean architecture and best practices.

## 🚀 Quick Start

### Prerequisites
- Go 1.21 or higher

### Running the Application

```bash
# Run the server
go run main.go

# The server will start on http://localhost:8080
```

### Testing the Server

```bash
# Health check
curl http://localhost:8080/health

# Home endpoint
curl http://localhost:8080/
```

## 📋 Features (In Progress)

- [x] Basic HTTP server
- [x] Health check endpoint
- [ ] URL shortening
- [ ] URL redirection
- [ ] URL validation
- [ ] Persistence layer
- [ ] Rate limiting
- [ ] Logging
- [ ] Unit tests

## 🏗️ Project Structure

```
url-shortener/
├── main.go           # Application entry point
├── go.mod            # Go module definition
└── README.md         # This file
```

## 🛠️ Tech Stack

- **Language**: Go 1.21
- **HTTP Framework**: net/http (standard library)
- **Storage**: In-memory (will upgrade to PostgreSQL)

## 📝 API Endpoints

### Health Check
```
GET /health
```

Response:
```json
{
  "status": "healthy"
}
```

---

**Status**: 🚧 Work in Progress
