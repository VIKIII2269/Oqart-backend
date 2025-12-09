# Quick Start Guide - OQart Backend

## 🚀 Get Started in 5 Minutes

### Prerequisites
- Docker & Docker Compose installed
- Or: Go 1.21+, PostgreSQL 15+, Redis 7+

### Option 1: Using Docker (Recommended)

```bash
# 1. Start all services
docker-compose up -d

# 2. Run migrations
docker-compose run --rm migrate

# 3. Check health
curl http://localhost:8080/health

# Expected response:
# {
#   "status": "healthy",
#   "timestamp": "2025-12-09T...",
#   "services": {
#     "database": "healthy",
#     "redis": "healthy"
#   }
# }
```

### Option 2: Local Development

```bash
# 1. Copy environment file
cp .env.example .env

# 2. Update JWT_SECRET in .env (or use default for dev)

# 3. Start PostgreSQL and Redis
docker-compose up -d postgres redis

# 4. Run migrations
make migrate-up

# 5. Start the API
make run

# Or directly:
go run cmd/api/main.go
```

## 📝 Test Authentication APIs

### 1. Register a New User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Test@123456",
    "first_name": "John",
    "last_name": "Doe",
    "phone": "9876543210"
  }'
```

**Expected Response:**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "uuid...",
      "email": "test@example.com",
      "phone": "9876543210",
      "role": "customer",
      "status": "active",
      "first_name": "John",
      "last_name": "Doe",
      ...
    },
    "access_token": "eyJhbGciOi...",
    "refresh_token": "eyJhbGciOi...",
    "expires_at": "2025-12-09T..."
  },
  "message": "User registered successfully"
}
```

### 2. Login with Email & Password

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "Test@123456"
  }'
```

### 3. Send Phone OTP

```bash
curl -X POST http://localhost:8080/api/v1/auth/login/phone \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "9876543210"
  }'
```

**Expected Response:**
```json
{
  "success": true,
  "data": {
    "message": "OTP sent successfully to your phone",
    "expires_at": "2025-12-09T..."
  }
}
```

**Note:** In development mode, the OTP is logged to console. Check Docker logs:
```bash
docker-compose logs -f api | grep "OTP generated"
```

### 4. Verify OTP

```bash
# Use the OTP from logs (6-digit code)
curl -X POST http://localhost:8080/api/v1/auth/verify-otp \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "9876543210",
    "code": "123456"
  }'
```

### 5. Access Protected Route

```bash
# Save the access_token from login/register response
TOKEN="your-access-token-here"

curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Authorization: Bearer $TOKEN"
```

### 6. Refresh Token

```bash
REFRESH_TOKEN="your-refresh-token-here"

curl -X POST http://localhost:8080/api/v1/auth/refresh-token \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "'$REFRESH_TOKEN'"
  }'
```

### 7. Forgot Password

```bash
curl -X POST http://localhost:8080/api/v1/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com"
  }'
```

### 8. Reset Password

```bash
# Use the token from forgot password (check logs in dev mode)
curl -X POST http://localhost:8080/api/v1/auth/reset-password \
  -H "Content-Type: application/json" \
  -d '{
    "token": "reset-token-here",
    "password": "NewPassword@123"
  }'
```

## 🔍 Explore the API

### Health Check
```bash
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/health
```

### Ping
```bash
curl http://localhost:8080/api/v1/ping
```

## 🗄️ Database Access

### Using Adminer (Database UI)

```bash
# Start Adminer
docker-compose --profile tools up -d adminer

# Access at: http://localhost:8081
# Server: postgres
# Username: oqart
# Password: oqart
# Database: oqart
```

### Using psql

```bash
# Connect to database
docker-compose exec postgres psql -U oqart

# List tables
\dt

# Query users
SELECT id, email, phone, role, status FROM users;

# Query sessions
SELECT id, user_id, expires_at FROM sessions;

# Query OTPs (dev only)
SELECT phone, code, expires_at, verified FROM otps ORDER BY created_at DESC LIMIT 5;
```

## 📊 Check Logs

```bash
# All services
docker-compose logs -f

# Just API
docker-compose logs -f api

# Just database
docker-compose logs -f postgres

# Filter for specific info
docker-compose logs -f api | grep "OTP"
docker-compose logs -f api | grep "logged in"
```

## 🧪 Test Validation

### Test Invalid Email
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "invalid-email",
    "password": "Test@123456",
    "first_name": "John"
  }'

# Expected: 400 Bad Request with validation error
```

### Test Weak Password
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test2@example.com",
    "password": "weak",
    "first_name": "John"
  }'

# Expected: 400 Bad Request
# Error: password must be at least 8 characters and contain uppercase, lowercase, number, and special character
```

### Test Invalid Phone
```bash
curl -X POST http://localhost:8080/api/v1/auth/login/phone \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "123"
  }'

# Expected: 400 Bad Request with phone validation error
```

### Test Duplicate Email
```bash
# Register once
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "duplicate@example.com",
    "password": "Test@123456",
    "first_name": "John"
  }'

# Try to register again with same email
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "duplicate@example.com",
    "password": "Test@123456",
    "first_name": "Jane"
  }'

# Expected: 409 Conflict - Email already exists
```

## 🔒 Test Authorization

### Test Accessing Protected Route Without Token
```bash
curl -X POST http://localhost:8080/api/v1/auth/logout

# Expected: 401 Unauthorized
```

### Test With Invalid Token
```bash
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Authorization: Bearer invalid-token-here"

# Expected: 401 Unauthorized
```

## 🛠️ Development Commands

```bash
# Run API
make run

# Run migrations
make migrate-up

# Rollback migration
make migrate-down

# Create new migration
make migrate-create NAME=add_something

# Run tests (when available)
make test

# Format code
make fmt

# Clean build artifacts
make clean

# View all commands
make help
```

## 📱 Test Complete Flow

### User Registration → Login → Protected Action

```bash
# 1. Register
RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "flow-test@example.com",
    "password": "Test@123456",
    "first_name": "Flow",
    "last_name": "Test"
  }')

echo $RESPONSE | jq .

# 2. Extract token
TOKEN=$(echo $RESPONSE | jq -r '.data.access_token')
echo "Token: $TOKEN"

# 3. Use protected endpoint
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Authorization: Bearer $TOKEN"

# Expected: Logout successful
```

### Phone OTP Flow

```bash
# 1. Send OTP
curl -X POST http://localhost:8080/api/v1/auth/login/phone \
  -H "Content-Type: application/json" \
  -d '{"phone": "9999999999"}'

# 2. Check logs for OTP
docker-compose logs api | grep "code" | tail -1

# 3. Verify OTP (replace 123456 with actual code from logs)
RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/auth/verify-otp \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "9999999999",
    "code": "123456"
  }')

echo $RESPONSE | jq .

# 4. Extract and use token
TOKEN=$(echo $RESPONSE | jq -r '.data.access_token')
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Authorization: Bearer $TOKEN"
```

## 🎯 What's Working

✅ **Complete Authentication System (10 APIs)**
- User registration with validation
- Email/password login
- Phone OTP authentication
- JWT token generation
- Token refresh
- Password reset flow
- Email verification
- Secure logout

✅ **Security Features**
- Password hashing (bcrypt)
- JWT with expiry
- Refresh tokens
- OTP rate limiting
- Input validation
- Role-based access control

✅ **Infrastructure**
- PostgreSQL database with migrations
- Redis caching
- Docker containerization
- Health checks
- Structured logging

## 🚧 Coming Next

- User profile management (8 endpoints)
- Vendor onboarding workflow (12 endpoints)
- Product management (20 endpoints)
- Shopping cart & wishlist (13 endpoints)
- Orders & payments (23 endpoints)
- Reviews & ratings (8 endpoints)
- Admin panel (24 endpoints)

## 🐛 Troubleshooting

### Database Connection Failed
```bash
# Check if PostgreSQL is running
docker-compose ps postgres

# Check logs
docker-compose logs postgres

# Restart
docker-compose restart postgres
```

### Redis Connection Failed
```bash
# Check if Redis is running
docker-compose ps redis

# Restart
docker-compose restart redis
```

### Migrations Not Applied
```bash
# Force run migrations
docker-compose run --rm migrate

# Or manually
make migrate-up
```

### Port Already in Use
```bash
# Check what's using port 8080
lsof -i :8080

# Kill process or change PORT in .env
```

## 📚 Additional Resources

- **Full API Documentation**: See README.md
- **Database Schema**: Check migrations/
- **Project Status**: See PROJECT_STATUS.md
- **Environment Variables**: See .env.example

## 🤝 Need Help?

- Check logs: `docker-compose logs -f`
- Verify health: `curl http://localhost:8080/health`
- Database state: `docker-compose exec postgres psql -U oqart`
- Reset everything: `docker-compose down -v && docker-compose up -d`

---

**Ready to build amazing organic products marketplace! 🌿**
