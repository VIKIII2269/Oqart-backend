# OQart Backend - Organic Products Marketplace

![Version](https://img.shields.io/badge/version-1.0.0-blue.svg)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-336791?logo=postgresql)
![Redis](https://img.shields.io/badge/Redis-7+-DC382D?logo=redis)
![License](https://img.shields.io/badge/license-MIT-green.svg)

**Production-grade e-commerce backend for India's premier organic products marketplace.**

## 🌟 Features

### Core Features
- **Multi-vendor Platform**: Complete vendor onboarding with document verification
- **Product Management**: Comprehensive catalog with variants, certifications, and reviews
- **Smart Cart & Wishlist**: Real-time inventory checks and price updates
- **Order Processing**: End-to-end order lifecycle management
- **Payment Integration**: Razorpay (primary) and Stripe (backup) with COD support
- **Returns & Refunds**: Automated return processing and refund handling
- **Reviews & Ratings**: Product and vendor rating system

### Technical Features
- **Clean Architecture**: Domain-driven design with clear separation of concerns
- **RESTful API**: 120+ well-documented endpoints
- **JWT Authentication**: Secure token-based auth with refresh tokens
- **Multi-auth Support**: Email, Phone OTP, Google OAuth
- **Role-Based Access Control**: Customer, Vendor, Admin, Super Admin
- **Real-time Caching**: Redis for optimal performance
- **File Upload**: Local and S3-compatible storage
- **Email & SMS**: Notification system with templates
- **Rate Limiting**: API protection against abuse
- **Comprehensive Logging**: Structured logging with zerolog
- **Database Migrations**: Version-controlled schema management

## 📊 Tech Stack

| Technology | Version | Purpose |
|------------|---------|---------|
| **Go** | 1.21+ | Backend language |
| **Gin** | Latest | HTTP framework |
| **PostgreSQL** | 15+ | Primary database |
| **Redis** | 7+ | Caching & sessions |
| **GORM** | Latest | ORM |
| **JWT** | v5 | Authentication |
| **Docker** | Latest | Containerization |
| **Razorpay** | Latest | Payment gateway |

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 15+
- Redis 7+
- Docker & Docker Compose (optional)

### Option 1: Using Docker (Recommended)

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd Oqart-backend
   ```

2. **Copy environment file**
   ```bash
   cp .env.example .env
   ```

3. **Update environment variables** in `.env`
   ```bash
   # Minimum required:
   JWT_SECRET=your-super-secret-key-minimum-32-characters-long
   SMTP_USER=your-email@gmail.com
   SMTP_PASSWORD=your-app-password
   ```

4. **Start services**
   ```bash
   docker-compose up -d
   ```

5. **Run migrations**
   ```bash
   docker-compose run --rm migrate
   ```

6. **Access the application**
   - API: http://localhost:8080
   - Health Check: http://localhost:8080/health
   - Database UI (Adminer): http://localhost:8081 (with `--profile tools`)

### Option 2: Local Development

1. **Install dependencies**
   ```bash
   go mod download
   ```

2. **Setup database**
   ```bash
   # Create database
   createdb oqart

   # Run migrations
   make migrate-up
   ```

3. **Start Redis**
   ```bash
   redis-server
   ```

4. **Run the application**
   ```bash
   make run
   # Or directly:
   go run cmd/api/main.go
   ```

## 📁 Project Structure

```
Oqart-backend/
├── cmd/api/                    # Application entry point
│   └── main.go                 # Main application
├── internal/                   # Private application code
│   ├── domain/                 # Business entities
│   │   ├── user.go
│   │   ├── product.go
│   │   ├── order.go
│   │   └── ...
│   ├── usecase/                # Business logic
│   ├── delivery/http/          # HTTP handlers
│   ├── repository/postgres/    # Database layer
│   └── service/                # External services
├── pkg/                        # Shared packages
│   ├── logger/                 # Logging utilities
│   ├── validator/              # Input validation
│   ├── jwt/                    # JWT utilities
│   ├── cache/                  # Redis cache
│   ├── database/               # Database connection
│   ├── utils/                  # Helper functions
│   └── errors/                 # Error handling
├── config/                     # Configuration
│   └── config.go
├── migrations/                 # Database migrations
│   ├── 000001_*.up.sql
│   ├── 000001_*.down.sql
│   └── ...
├── docs/                       # Documentation
├── tests/                      # Tests
│   ├── unit/
│   └── integration/
├── scripts/                    # Utility scripts
├── .env.example                # Environment template
├── .gitignore
├── Dockerfile                  # Docker image
├── docker-compose.yml          # Docker services
├── Makefile                    # Common commands
├── go.mod                      # Go modules
└── README.md                   # This file
```

## 🔧 Configuration

### Environment Variables

Key configuration options (see `.env.example` for complete list):

#### Server
- `PORT`: Server port (default: 8080)
- `ENV`: Environment (development/production)

#### Database
- `DATABASE_URL`: PostgreSQL connection string
- `DB_MAX_CONNECTIONS`: Max DB connections (default: 25)

#### Redis
- `REDIS_HOST`: Redis host
- `REDIS_PORT`: Redis port (default: 6379)

#### JWT
- `JWT_SECRET`: Secret key for JWT (min 32 chars) **[REQUIRED]**
- `JWT_EXPIRY`: Access token expiry (default: 1h)
- `REFRESH_TOKEN_EXPIRY`: Refresh token expiry (default: 720h)

#### Payments
- `RAZORPAY_KEY_ID`: Razorpay key
- `RAZORPAY_KEY_SECRET`: Razorpay secret

#### Email (SMTP)
- `SMTP_HOST`: SMTP server
- `SMTP_PORT`: SMTP port
- `SMTP_USER`: SMTP username
- `SMTP_PASSWORD`: SMTP password

#### SMS (Twilio)
- `TWILIO_ACCOUNT_SID`: Twilio account SID
- `TWILIO_AUTH_TOKEN`: Twilio auth token

## 🛠️ Makefile Commands

```bash
# Development
make run              # Run the application
make dev              # Run with hot reload (requires air)

# Build
make build            # Build binary
make docker-build     # Build Docker image

# Database
make migrate-up       # Run migrations
make migrate-down     # Rollback last migration
make migrate-create   # Create new migration (NAME=<name>)
make seed             # Seed database

# Testing
make test             # Run tests
make test-coverage    # Run tests with coverage

# Docker
make docker-up        # Start all services
make docker-down      # Stop all services
make docker-logs      # View logs

# Code Quality
make lint             # Run linters
make fmt              # Format code
make tidy             # Tidy go modules

# Utilities
make clean            # Clean build artifacts
make help             # Show help
```

## 📝 API Documentation

### Core API Endpoints (120+ total)

#### Authentication (10 endpoints)
```
POST   /api/v1/auth/register           # Register with email
POST   /api/v1/auth/login              # Login with email/password
POST   /api/v1/auth/login/phone        # Send OTP to phone
POST   /api/v1/auth/verify-otp         # Verify OTP & login
POST   /api/v1/auth/login/google       # Google OAuth login
POST   /api/v1/auth/refresh-token      # Refresh JWT
POST   /api/v1/auth/forgot-password    # Send reset email
POST   /api/v1/auth/reset-password     # Reset with token
POST   /api/v1/auth/verify-email       # Verify email token
POST   /api/v1/auth/logout             # Logout
```

#### User Management (8 endpoints)
```
GET    /api/v1/users/me                # Get current user
PUT    /api/v1/users/me                # Update profile
PUT    /api/v1/users/me/password       # Change password
PUT    /api/v1/users/me/email          # Change email
PUT    /api/v1/users/me/phone          # Change phone
DELETE /api/v1/users/me                # Delete account
GET    /api/v1/users/me/preferences    # Get preferences
PUT    /api/v1/users/me/preferences    # Update preferences
```

#### Products (20 endpoints)
```
GET    /api/v1/products                # List products (filters, search, pagination)
GET    /api/v1/products/:slug          # Get product details
GET    /api/v1/products/featured       # Featured products
GET    /api/v1/products/search         # Advanced search
```

#### Vendor Onboarding (12 endpoints)
```
POST   /api/v1/vendors/onboard/step1   # Basic details
POST   /api/v1/vendors/verify-gstin    # Verify GSTIN
POST   /api/v1/vendors/onboard/step2   # Company & address
POST   /api/v1/vendors/onboard/step3   # Upload documents
POST   /api/v1/vendors/bank-details    # Bank account
POST   /api/v1/vendors/submit          # Final submission
GET    /api/v1/vendors/me              # Vendor profile
```

#### Orders (15 endpoints)
```
POST   /api/v1/orders                  # Place order
GET    /api/v1/orders                  # My orders
GET    /api/v1/orders/:id              # Order details
GET    /api/v1/orders/:id/invoice      # Download invoice PDF
POST   /api/v1/orders/:id/cancel       # Cancel order
GET    /api/v1/orders/:id/track        # Track order
```

#### Cart (8 endpoints)
```
GET    /api/v1/cart                    # Get cart
POST   /api/v1/cart/items              # Add to cart
PUT    /api/v1/cart/items/:id          # Update quantity
DELETE /api/v1/cart/items/:id          # Remove item
DELETE /api/v1/cart                    # Clear cart
```

#### Payments (8 endpoints)
```
POST   /api/v1/payments/initiate       # Create payment intent
POST   /api/v1/payments/verify         # Verify payment
POST   /api/v1/payments/webhook        # Razorpay webhook
GET    /api/v1/payments/:order_id      # Payment status
```

_...and 60+ more endpoints for wishlist, coupons, returns, reviews, admin panel, etc._

### API Response Format

**Success Response:**
```json
{
  "success": true,
  "data": { ... },
  "message": "Success message"
}
```

**Error Response:**
```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable error message"
  }
}
```

**Paginated Response:**
```json
{
  "success": true,
  "data": [...],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total_items": 100,
    "total_pages": 5
  }
}
```

## 🗄️ Database

### Schema

The database consists of 30+ tables organized into:

- **Users & Auth**: users, sessions, password_resets, otps
- **Vendors**: vendors, vendor_documents, vendor_addresses, vendor_bank_details
- **Products**: products, categories, product_images, product_variants
- **Orders**: orders, order_items, order_status_history
- **Payments**: payments, payment_transactions, invoices
- **Returns**: returns, return_items, refunds
- **Reviews**: reviews, review_votes, vendor_ratings
- **Utilities**: notifications, audit_logs, search_logs, serviceable_pincodes

### Migrations

```bash
# Run all migrations
make migrate-up

# Rollback last migration
make migrate-down

# Create new migration
make migrate-create NAME=add_feature_x

# Check current version
make migrate-version

# Force specific version (use with caution)
make migrate-force VERSION=1
```

## 🔒 Security

### Implemented Security Features

- **Password Hashing**: Bcrypt with cost 12
- **JWT Authentication**: HS256 with 1-hour expiry
- **Refresh Tokens**: 30-day expiry, stored in database
- **Rate Limiting**: Configurable per endpoint
- **CORS**: Configurable allowed origins
- **Input Validation**: Comprehensive validation using go-playground/validator
- **SQL Injection Protection**: Parameterized queries via GORM
- **XSS Protection**: Input sanitization
- **HTTPS**: TLS support (configure reverse proxy)

### Security Best Practices

- Never commit `.env` file
- Rotate JWT secrets regularly
- Use strong passwords (min 8 chars, uppercase, lowercase, number, special char)
- Enable 2FA for admin accounts (future feature)
- Monitor audit logs regularly

## 📊 Performance

### Caching Strategy

- **User Sessions**: 1 hour
- **OTP**: 5 minutes
- **Product Details**: 1 hour
- **Product Listings**: 30 minutes
- **Cart**: 24 hours
- **Category Tree**: 6 hours

### Database Optimization

- Comprehensive indexes on frequently queried columns
- Full-text search indexes for products
- Connection pooling (max 25 connections)
- Query optimization with EXPLAIN ANALYZE
- Pagination for large result sets

## 🧪 Testing

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run specific package
go test ./internal/usecase/...

# Run with race detector
go test -race ./...

# Benchmark tests
go test -bench=. ./...
```

## 📦 Deployment

### Docker Deployment

```bash
# Build image
docker build -t oqart-backend:1.0.0 .

# Run container
docker run -p 8080:8080 --env-file .env oqart-backend:1.0.0
```

### Production Checklist

- [ ] Update `JWT_SECRET` with strong random value
- [ ] Set `ENV=production`
- [ ] Configure production database
- [ ] Set up SSL/TLS certificates
- [ ] Configure reverse proxy (Nginx/Traefik)
- [ ] Set up monitoring (Prometheus/Grafana)
- [ ] Configure backup strategy
- [ ] Set up log aggregation
- [ ] Configure CDN for static assets
- [ ] Set up error tracking (Sentry)
- [ ] Enable rate limiting
- [ ] Configure email service (SendGrid/AWS SES)
- [ ] Set up SMS service (Twilio)
- [ ] Configure payment gateways
- [ ] Set up S3 for file uploads

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License.

## 📧 Support

- **Email**: support@oqart.com
- **Documentation**: https://docs.oqart.com
- **Issues**: GitHub Issues

## 🙏 Acknowledgments

- Go community for excellent tooling
- Gin framework for fast HTTP routing
- PostgreSQL for robust database
- Redis for high-performance caching

---

**Built with ❤️ for India's organic products marketplace**
