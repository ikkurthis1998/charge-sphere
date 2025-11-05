# OCPI Service - ChargeSphere

OCPI 2.3.0 compliant roaming hub for EV charging networks.

## Current Status

### ✅ Completed Modules

#### Credentials Module (Partner Registration)
- Partner registration (CPO and eMSP)
- Token-based authentication
- Credentials management (GET, PUT, DELETE)
- Full unit test coverage

## Architecture

```
services/ocpi-service/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── api/
│   │   ├── handlers/               # HTTP request handlers
│   │   │   ├── credentials_handler.go
│   │   │   └── credentials_handler_test.go
│   │   ├── middleware/             # HTTP middleware
│   │   │   └── auth.go
│   │   └── routes.go               # Route definitions
│   ├── config/
│   │   └── config.go               # Configuration management
│   ├── database/
│   │   └── mongodb.go              # Database connection
│   ├── domain/
│   │   ├── models/                 # Data models
│   │   │   └── credentials.go
│   │   └── services/               # Business logic
│   │       ├── credentials_service.go
│   │       └── credentials_service_test.go
│   └── repository/
│       └── mongodb/                # Database operations
│           └── partner_repository.go
├── config.yaml                     # Configuration file
├── go.mod                          # Go module definition
├── Makefile                        # Build commands
└── README.md                       # This file
```

## Tech Stack

- **Language:** Go 1.21+
- **Framework:** Gin
- **Database:** MongoDB 7
- **Cache:** Redis 7
- **Workflows:** Temporal
- **Testing:** testify

## Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- MongoDB 7
- Redis 7

## Getting Started

### 1. Start Infrastructure

```bash
# Start MongoDB, Redis, and Temporal
make docker-up

# View logs
make docker-logs

# Stop containers
make docker-down
```

### 2. Install Dependencies

```bash
make deps
```

### 3. Run the Service

```bash
make run
```

The service will start on `http://localhost:8080`

### 4. Health Check

```bash
curl http://localhost:8080/health
```

## Testing

### Run Unit Tests (Fast)

Unit tests use mocks and don't require infrastructure:

```bash
make test
# or
make test-unit
```

### Run Integration Tests

Integration tests use real MongoDB and test the full stack:

```bash
make test-integration
```

This will:
1. Start Docker containers (MongoDB, Redis, Temporal)
2. Wait for services to be ready
3. Run integration tests with real database

**Note:** Integration tests require Docker to be running.

### Run All Tests (Unit + Integration)

```bash
make test-all
```

### Run Tests with Coverage

```bash
make test-coverage
```

This generates `coverage.html` that you can open in a browser.

### Test Structure

**Unit Tests:**
- `internal/domain/services/*_test.go` - Business logic tests
- `internal/api/handlers/*_test.go` - HTTP handler tests
- Use mocks for dependencies
- Fast execution (~1 second)

**Integration Tests:**
- `internal/api/integration_test.go` - Full stack tests
- Uses real MongoDB database
- Tests complete HTTP request/response flows
- Tests data persistence
- Slower execution (~5-10 seconds)

## API Documentation

### Base URL

```
http://localhost:8080/ocpi/2.3
```

### Authentication

All endpoints (except POST /credentials for registration) require authentication:

```
Authorization: Token <your_token>
```

---

## Credentials Module

### 1. Register as a Partner

**Endpoint:** `POST /ocpi/2.3/credentials?type=CPO`

**Request:**
```json
{
  "token": "your_token_for_hub_to_call_you",
  "url": "https://your-domain.com/ocpi/2.3",
  "roles": [
    {
      "role": "CPO",
      "party_id": "ABC",
      "country_code": "DE",
      "business_details": {
        "name": "My Charging Company"
      }
    }
  ]
}
```

**Response:**
```json
{
  "data": {
    "token": "hub_token_for_you_to_call_hub",
    "url": "http://localhost:8080/ocpi/2.3",
    "roles": [
      {
        "role": "HUB",
        "party_id": "HUB",
        "country_code": "US",
        "business_details": {
          "name": "ChargeSphere Hub"
        }
      }
    ]
  },
  "status_code": 1000,
  "status_message": "Success",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

**Partner Types:**
- `type=CPO` - Charge Point Operator (owns charging stations)
- `type=EMSP` - E-Mobility Service Provider (serves EV drivers)

---

### 2. Get Credentials

**Endpoint:** `GET /ocpi/2.3/credentials`

**Headers:**
```
Authorization: Token your_token
```

**Response:**
```json
{
  "data": {
    "token": "your_token",
    "url": "http://localhost:8080/ocpi/2.3",
    "roles": [
      {
        "role": "HUB",
        "party_id": "HUB",
        "country_code": "US",
        "business_details": {
          "name": "ChargeSphere Hub"
        }
      }
    ]
  },
  "status_code": 1000,
  "status_message": "Success",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

---

### 3. Update Credentials

**Endpoint:** `PUT /ocpi/2.3/credentials`

**Headers:**
```
Authorization: Token your_current_token
```

**Request:**
```json
{
  "token": "your_new_token",
  "url": "https://your-new-domain.com/ocpi/2.3",
  "roles": [
    {
      "role": "CPO",
      "party_id": "ABC",
      "country_code": "DE",
      "business_details": {
        "name": "Updated Company Name"
      }
    }
  ]
}
```

**Response:**
```json
{
  "data": {
    "token": "new_hub_token",
    "url": "http://localhost:8080/ocpi/2.3",
    "roles": [...]
  },
  "status_code": 1000,
  "status_message": "Success",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

---

### 4. Delete Registration

**Endpoint:** `DELETE /ocpi/2.3/credentials`

**Headers:**
```
Authorization: Token your_token
```

**Response:**
```json
{
  "data": null,
  "status_code": 1000,
  "status_message": "Partner registration deleted successfully",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

---

## OCPI Response Format

All responses follow the OCPI standard format:

```json
{
  "data": <response_data>,
  "status_code": <ocpi_status_code>,
  "status_message": "<human_readable_message>",
  "timestamp": "<iso8601_timestamp>"
}
```

**OCPI Status Codes:**
- `1000` - Success
- `1001` - Invalid request
- `2001` - Authentication/Authorization error
- `2002` - Not found
- `3000` - Server error

---

## Example: Complete Registration Flow

### Step 1: CPO Registers with Hub

```bash
curl -X POST http://localhost:8080/ocpi/2.3/credentials?type=CPO \
  -H "Content-Type: application/json" \
  -d '{
    "token": "cpo_secret_token_123",
    "url": "https://cpo-company.com/ocpi/2.3",
    "roles": [{
      "role": "CPO",
      "party_id": "ABC",
      "country_code": "DE",
      "business_details": {
        "name": "ABC Charging Network"
      }
    }]
  }'
```

**Hub returns:**
```json
{
  "data": {
    "token": "hub_token_abc123xyz",
    "url": "http://localhost:8080/ocpi/2.3",
    "roles": [...]
  },
  "status_code": 1000
}
```

### Step 2: CPO Calls Hub (Authenticated)

Now the CPO uses the hub token to call authenticated endpoints:

```bash
curl http://localhost:8080/ocpi/2.3/credentials \
  -H "Authorization: Token hub_token_abc123xyz"
```

---

## Development

### Build

```bash
make build
```

Binary will be created at `bin/server`

### Format Code

```bash
make fmt
```

### Run Linter

```bash
make lint
```

### Clean

```bash
make clean
```

---

## Configuration

Edit `config.yaml` to configure:

```yaml
server:
  host: "0.0.0.0"
  port: 8080
  mode: "debug"  # or "release"

mongodb:
  uri: "mongodb://admin:admin123@localhost:27017"
  database: "chargesphere"
  timeout: 10

redis:
  host: "localhost:6379"
  password: ""
  db: 0

temporal:
  host: "localhost:7233"
  namespace: "default"

ocpi:
  version: "2.3"
  country_code: "US"
  party_id: "CSP"
```

---

## Database Schema

### Partners Collection

```javascript
{
  "_id": ObjectId,
  "partner_id": "DE-ABC",
  "name": "ABC Charging Network",
  "type": "CPO",  // or "EMSP"
  "credentials": {
    "token": "encrypted_token",
    "url": "https://partner.com/ocpi/2.3",
    "roles": [...],
    "version": "2.3"
  },
  "status": "ACTIVE",  // ACTIVE, INACTIVE, SUSPENDED
  "created_at": ISODate,
  "updated_at": ISODate
}
```

**Indexes:**
- `partner_id` (unique)
- `credentials.token` (unique)
- `status`

---

## Testing Strategy

### Unit Tests ✅

- ✅ Service layer tests (business logic)
- ✅ Handler tests (HTTP endpoints)
- ✅ Mock dependencies (repository, services)
- ✅ Edge cases and error handling
- ✅ Fast execution (~1 second)

### Integration Tests ✅

- ✅ Complete partner registration flow
- ✅ Register → Get → Update → Delete lifecycle
- ✅ Multiple partners (CPO and eMSP) registration
- ✅ Duplicate partner registration validation
- ✅ Invalid token authentication
- ✅ Missing authorization header handling
- ✅ Invalid request body validation
- ✅ Invalid role for partner type
- ✅ Invalid country code (must be 2 chars)
- ✅ Invalid party ID (must be 3 chars)
- ✅ Health check endpoint
- ✅ Data persistence verification (MongoDB)
- ✅ Concurrent registrations (10 partners)
- ✅ Token prefix handling (with/without "Token ")
- ✅ Real MongoDB database operations
- ✅ Full HTTP request/response flows

**Coverage:**
- 14 integration test scenarios
- Complete CRUD operations
- Authentication and authorization
- Data validation
- Concurrency handling
- Database persistence

### Load Tests (TODO)

- High-volume concurrent registrations
- Token validation performance under load
- Database query optimization
- Response time benchmarks

---

## Next Modules

### Phase 2: Core Roaming (Weeks 3-6)

- **Locations Module** - Charging station data
- **Tokens Module** - Driver authentication
- **Sessions Module** - Active charging sessions

### Phase 3: Billing (Weeks 7-8)

- **CDRs Module** - Charge Detail Records
- **Tariffs Module** - Pricing information
- **Temporal Workflows** - Reliable billing

### Phase 4: Advanced (Weeks 9+)

- **Commands Module** - Remote control
- **Reservations Module** - Booking
- **Admin Dashboard** - Web interface

---

## Contributing

1. Write tests for all new code
2. Run `make test` before committing
3. Follow Go conventions
4. Update README for new features

---

## Troubleshooting

### MongoDB Connection Error

```bash
# Check if MongoDB is running
docker ps | grep mongodb

# Check MongoDB logs
docker logs charge-sphere-mongodb
```

### Port Already in Use

```bash
# Change port in config.yaml
server:
  port: 8081
```

### Tests Failing

```bash
# Ensure Docker containers are running
make docker-up

# Run tests with verbose output
go test -v ./...
```

---

## License

MIT

---

## Support

For issues and questions, please open an issue on GitHub.
