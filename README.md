# Salary Management Tool

A professional salary management system for HR Managers to manage employee data and gain salary insights. Features secure JWT authentication, role-based access control, and a modern responsive UI built with shadcn/ui.

---

## Table of Contents

1. [Overview](#overview)
2. [Architecture](#architecture)
3. [Tech Stack](#tech-stack)
4. [Project Structure](#project-structure)
5. [Setup & Installation](#setup--installation)
6. [Usage](#usage)
7. [API Documentation](#api-documentation)
8. [Authentication Flow](#authentication-flow)
9. [Testing](#testing)
10. [Deployment](#deployment)
11. [Environment Variables](#environment-variables)
12. [Troubleshooting](#troubleshooting)

---

## Overview

The Salary Management Tool is a full-stack web application designed for HR departments to efficiently manage employee records and analyze salary data across different dimensions.

### Key Features

| Feature | Description |
|---------|-------------|
| **Authentication** | JWT-based authentication with role-based access control (HR/Admin roles) |
| **Employee Management** | Complete CRUD operations with pagination, search, and filtering |
| **Analytics Dashboard** | Salary statistics by country, job title, and department |
| **Modern UI** | Responsive design using shadcn/ui components with Tailwind CSS |
| **Security** | Password hashing with bcrypt, protected endpoints, input validation |

---

## Architecture

### System Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        Client (Browser)                          │
│                   React + TypeScript + Vite                    │
│                     Tailwind CSS + shadcn/ui                   │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │ HTTP/REST
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Go Backend Server                          │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────────┐  │
│  │   Router    │  │  Middleware │  │    Route Handlers     │  │
│  │  (gorilla)  │──│  (CORS,    │──│  - Auth Handler       │  │
│  │             │  │  JWT Auth) │  │  - Employee Handler   │  │
│  └─────────────┘  └─────────────┘  │  - Analytics Handler  │  │
│                                    └─────────────────────────┘  │
│                              │                                  │
│                              ▼                                  │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────────┐  │
│  │   Models    │  │  Services   │  │   Static File Server    │  │
│  │  (GORM)     │  │  (JWT,      │  │   (SPA Support)         │  │
│  │  - User     │  │  Password)  │  │                         │  │
│  │  - Employee │  │             │  │                         │  │
│  └─────────────┘  └─────────────┘  └─────────────────────────┘  │
│                              │                                  │
│                              ▼                                  │
│                    ┌─────────────────┐                          │
│                    │   SQLite DB     │                          │
│                    │ (GORM ORM)      │                          │
│                    └─────────────────┘                          │
└─────────────────────────────────────────────────────────────────┘
```

### Data Flow

1. **Authentication Flow**
   ```
   Client → POST /api/auth/login → AuthService → JWT Token → localStorage
   ```

2. **Protected Request Flow**
   ```
   Client (with Bearer token) → authMiddleware → Handler → Service → Database
        ↑                              │
        └──────── 401/403 ────────────┘ (if invalid)
   ```

3. **Static Assets Flow**
   ```
   Client → GET / → spaHandler → index.html (React app)
   Client → GET /api/* → API Router → JSON Response
   ```

---

## Tech Stack

### Backend
| Component | Technology |
|-----------|------------|
| Language | Go 1.25 |
| Framework | Gorilla Mux |
| ORM | GORM |
| Database | SQLite |
| Authentication | JWT (golang-jwt) |
| Password Hashing | bcrypt (golang.org/x/crypto) |
| CORS | rs/cors |

### Frontend
| Component | Technology |
|-----------|------------|
| Framework | React 18 |
| Language | TypeScript |
| Build Tool | Vite 6 |
| Styling | Tailwind CSS 4 |
| UI Components | shadcn/ui + Radix UI |
| Routing | React Router 7 |
| Form Handling | React Hook Form + Zod |
| HTTP Client | Axios |
| Icons | Lucide React + Heroicons |

---

## Project Structure

```
/Users/tofiquem/tosif-practice/assingment/
│
├── cmd/                          # Application entry points
│   ├── server/
│   │   └── main.go              # HTTP server initialization
│   └── seed/
│       └── main.go              # Database seeding script
│
├── pkg/                         # Application packages
│   ├── database/
│   │   └── database.go          # SQLite connection & GORM setup
│   ├── handlers/
│   │   ├── auth.go             # Authentication & user management
│   │   ├── employee.go         # Employee CRUD handlers
│   │   └── analytics.go        # Salary analytics handlers
│   ├── models/
│   │   ├── user.go             # User model & auth types
│   │   └── employee.go         # Employee model & request types
│   ├── services/
│   │   ├── auth_service.go     # Business logic for auth
│   │   ├── jwt_service.go      # JWT token generation/validation
│   │   └── password_service.go # Password hashing & validation
│   └── testutils/
│       └── testutils.go        # Test helpers & fixtures
│
├── api/
│   └── index.go                # Vercel serverless function entry
│
├── web/                        # React frontend application
│   ├── src/
│   │   ├── components/         # React components
│   │   │   ├── ui/            # shadcn/ui components
│   │   │   ├── Login.tsx      # Authentication component
│   │   │   ├── Dashboard.tsx  # Main dashboard
│   │   │   ├── EmployeeForm.tsx
│   │   │   ├── Analytics.tsx
│   │   │   └── employees/     # Employee view components
│   │   ├── services/
│   │   │   └── api.ts         # API client with JWT interceptors
│   │   ├── types/
│   │   │   └── index.ts       # TypeScript type definitions
│   │   ├── hooks/
│   │   │   └── useAuth.ts     # Authentication hook
│   │   ├── lib/
│   │   │   └── utils.ts       # Utility functions
│   │   ├── App.tsx            # Root application component
│   │   └── main.tsx           # Application entry point
│   ├── public/                # Static assets
│   ├── package.json           # Frontend dependencies
│   ├── tsconfig.json          # TypeScript configuration
│   └── vite.config.ts         # Vite configuration
│
├── seed/                       # Seed data files
├── public/                     # Production static files
├── go.mod                      # Go module definition
├── go.sum                      # Go dependency checksums
├── Makefile                    # Build & test automation
├── Dockerfile                  # Container configuration
├── entrypoint.sh              # Docker entrypoint script
├── vercel.json                # Vercel deployment config
├── requirements.md            # Project requirements
└── README.md                  # This file
```

---

## Setup & Installation

### Prerequisites

- **Go** (v1.25+)
- **Node.js** (v18+) with npm
- **Git**
- **Make** (optional, for Makefile commands)

### Quick Start

#### Option 1: Run Backend and Frontend Separately

**Terminal 1 - Backend:**
```bash
cd /Users/tofiquem/tosif-practice/assingment

# Install Go dependencies
go mod tidy

# Run the server
go run cmd/server/main.go
```
Backend runs on `http://localhost:8080`

**Terminal 2 - Frontend:**
```bash
cd /Users/tofiquem/tosif-practice/assingment/web

# Install dependencies
npm install

# Run development server
npm run dev
```
Frontend runs on `http://localhost:5173`

#### Option 2: Using Make Commands

```bash
# Run backend only
make run

# Build and run everything (frontend + backend)
make start

# Run tests
make test

# Seed database with test data
make seed
```

### Database Seeding

```bash
# Seed the database with sample employee data
go run cmd/seed/main.go
```

---

## Usage

### Default Login Credentials

| Field | Value |
|-------|-------|
| Email | `admin@company.com` |
| Password | `admin123` |
| Role | Admin (HR Manager) |

### Application Pages

| Route | Description | Access |
|-------|-------------|--------|
| `/login` | Authentication page | Public |
| `/` | Employee dashboard | Protected |
| `/analytics` | Salary analytics | Protected |
| `/add-employee` | Create new employee | Protected (HR/Admin) |
| `/edit-employee/:id` | Edit employee | Protected (HR/Admin) |
| `/employee/:id` | View employee details | Protected |

### Role-Based Access

- **HR Role**: Can manage employees and view analytics
- **Admin Role**: Full access including user management

---

## API Documentation

### Authentication Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/auth/login` | User login | No |
| GET | `/api/auth/me` | Get current user | Yes |
| POST | `/api/auth/logout` | Logout | Yes |
| POST | `/api/auth/change-password` | Change password | Yes |

### Employee Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/employees` | List all employees (paginated) | Yes (HR/Admin) |
| POST | `/api/employees` | Create new employee | Yes (HR/Admin) |
| GET | `/api/employees/:id` | Get specific employee | Yes (HR/Admin) |
| PUT | `/api/employees/:id` | Update employee | Yes (HR/Admin) |
| DELETE | `/api/employees/:id` | Delete employee | Yes (HR/Admin) |

**Query Parameters for List:**
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 20, max: 100)
- `search` - Search by name/email
- `country` - Filter by country
- `department` - Filter by department

### Analytics Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/analytics/salary/by-country` | Salary stats by country | Yes (HR/Admin) |
| GET | `/api/analytics/salary/by-job-title/:country` | Salary by job title | Yes (HR/Admin) |
| GET | `/api/analytics/salary/department-insights` | Department analysis | Yes (HR/Admin) |
| GET | `/api/analytics/salary/department-insights/:country` | Dept analysis by country | Yes (HR/Admin) |

### Admin Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/admin/users` | List all users | Yes (Admin only) |
| POST | `/api/admin/users` | Create new user | Yes (Admin only) |
| PUT | `/api/admin/users/:id` | Update user | Yes (Admin only) |
| DELETE | `/api/admin/users/:id` | Delete user | Yes (Admin only) |

### Health Check

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/health` | Application health status |

---

## Authentication Flow

### JWT Token Flow

```
┌─────────┐                    ┌─────────────┐                    ┌──────────┐
│  Client │────────────────────│   Server    │────────────────────│  Database│
└────┬────┘  POST /api/auth/login└─────┬──────┘                    └──────────┘
     │    {email, password}            │
     │───────────────────────────────>│
     │                                 │ Validate credentials
     │                                 │──────> bcrypt compare
     │                                 │<──────│
     │                                 │ Generate JWT
     │                                 │──────> Sign token
     │                                 │<──────│
     │  {token, expiresAt, user}     │
     │<───────────────────────────────│
     │                                 │
     │ Store token in localStorage     │
     │                                 │
     │─────── Subsequent Requests ─────│
     │                                 │
     │ GET /api/employees              │
     │ Authorization: Bearer <token>   │
     │───────────────────────────────>│
     │                                 │ Validate JWT
     │                                 │──────> Parse & verify
     │                                 │<──────│
     │         {employees data}        │
     │<───────────────────────────────│
```

### Middleware Chain

```go
// Route registration with middleware
protected := r.PathPrefix("/api/employees").Subrouter()
protected.Use(h.authMiddleware)    // Validates JWT token
protected.Use(h.hrMiddleware)      // Checks HR/Admin role
```

---

## Testing

### Test Commands

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run unit tests only
make test-unit

# Run integration tests only
make test-integration

# Run specific test
make test-specific RUN=TestCreateEmployee

# Run with race detection
make test-verbose

# Generate coverage HTML report
make coverage-html
```

### Test Structure

| Package | Test Files | Description |
|---------|------------|-------------|
| `pkg/models` | `*_test.go` | Model validation tests |
| `pkg/handlers` | `*_test.go` | HTTP handler tests |
| `cmd/server` | `main_test.go` | Integration tests |

---

## Deployment

### Docker Deployment

```bash
# Build Docker image
docker build -t salary-management .

# Run container
docker run -p 8080:8080 salary-management
```

### Vercel Deployment

The application is configured for Vercel deployment with:
- **Frontend**: Static build from `web/build`
- **Backend API**: Serverless functions in `api/`

See `vercel.json` for routing configuration.

### Production Build

```bash
# Build backend
make build

# Build frontend
cd web && npm run build

# Run production binary
./salary-management
```

---

## Environment Variables

### Backend Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DATABASE_PATH` | `./salary_management.db` | SQLite database file path |
| `PORT` | `8080` | Server port |
| `PUBLIC_DIR` | `public` | Static files directory |
| `JWT_SECRET` | Auto-generated | JWT signing secret |
| `JWT_EXPIRATION` | `24h` | JWT token lifetime |
| `SEED_DATA_DIR` | `seed` | Seed data directory |

### Frontend Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `VITE_API_URL` | `/api` | Backend API URL (production) |

**Development**: Uses Vite proxy to avoid CORS (configured in `vite.config.ts`)

---

## Troubleshooting

### Common Issues

#### 1. Node.js Not Found
```bash
export NVM_DIR="$HOME/.nvm"
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"
nvm use node
```

#### 2. Port Already in Use
```bash
# Kill process on port 8080
lsof -ti:8080 | xargs kill -9

# Kill process on port 5173
lsof -ti:5173 | xargs kill -9
```

#### 3. Database Connection Issues
```bash
# Remove existing database
rm salary_management.db

# Restart backend to recreate
go run cmd/server/main.go
```

#### 4. Authentication Issues
```bash
# Set JWT secret for production
export JWT_SECRET="your-secret-key-here"

# Clear browser localStorage
# DevTools → Application → Local Storage → Clear
```

#### 5. Frontend Build Issues
```bash
cd web
rm -rf node_modules package-lock.json
npm install
```

#### 6. JWT Token Expired
- Logout and login again
- Default token lifetime: 24 hours

---

## License

This project is for assessment purposes only.
