# Salary Management Tool

A professional salary management system for HR Managers to manage employee data and gain salary insights. Features secure JWT authentication, role-based access control, and a modern responsive UI built with shadcn/ui.

## Overview

This application consists of:
- **Backend**: Go REST API with JWT authentication and SQLite database
- **Frontend**: React application with shadcn/ui components and Tailwind CSS
- **Database**: SQLite with GORM ORM
- **Authentication**: JWT-based authentication with role-based access control
- **UI**: Modern, responsive design with dark theme support

## Key Features

### Authentication & Security
- **JWT-based authentication** with secure token management
- **Role-based access control** (HR/Admin roles)
- **Password security** with bcrypt hashing and strength validation
- **Protected API endpoints** for all sensitive operations

### Modern UI/UX
- **Professional design** with shadcn/ui components
- **Responsive layout** optimized for all screen sizes
- **Dark theme support** with CSS variables
- **Mobile-friendly navigation** with hamburger menu
- **Interactive components** with hover states and transitions

### Employee Management
- **Complete CRUD operations** for employee records
- **Advanced search** and filtering capabilities
- **Pagination** for large datasets
- **Form validation** with user-friendly error messages

### Analytics & Insights
- **Salary statistics** by country and department
- **Job title analysis** with drill-down capabilities
- **Interactive dashboards** with visual indicators
- **Comprehensive metrics** and KPIs

## Prerequisites

- Node.js (v18+) with npm
- Go (v1.21+)
- Git

## Quick Start

### 1. Backend Setup

```bash
# Navigate to project root
cd /Users/tofiquem/tosif-practice/assingment

# Install Go dependencies
go mod tidy

# Run the backend server
go run cmd/server/main.go
```

The backend will start on `http://localhost:8080`

### 2. Frontend Setup

```bash
# Navigate to web directory
cd web

# Load Node.js (using nvm) and install dependencies
export NVM_DIR="$HOME/.nvm"
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"
nvm use node
npm install

# Run the frontend development server
npm run dev
```

The frontend will start on `http://localhost:5173`

### 3. Running Both Processes Simultaneously

To run both backend and frontend at the same time, open two separate terminal windows:

**Terminal 1 (Backend):**
```bash
cd /Users/tofiquem/tosif-practice/assingment
go run cmd/server/main.go
```

**Terminal 2 (Frontend):**
```bash
cd /Users/tofiquem/tosif-practice/assingment/web
export NVM_DIR="$HOME/.nvm"
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"
nvm use node
npm run dev
```

## Development Workflow

### Backend Development

```bash
# Run tests
make test

# Run specific test packages
make test-unit      # Unit tests only
make test-integration # Integration tests only

# Run with coverage
make test-coverage

# Build the application
make build

# Run the application
make run
```

### Frontend Development

```bash
cd web

# Development server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

### Database Seeding

```bash
# Seed the database with test data
go run cmd/seed/main.go
```

## API Endpoints

### Authentication
- `POST /api/auth/login` - User login with email and password
- `GET /api/auth/me` - Get current user information
- `POST /api/auth/logout` - User logout
- `POST /api/auth/change-password` - Change user password

### Employee Management (Protected)
- `GET /api/employees` - List all employees with pagination
- `POST /api/employees` - Create new employee
- `GET /api/employees/:id` - Get specific employee
- `PUT /api/employees/:id` - Update employee
- `DELETE /api/employees/:id` - Delete employee

### Salary Analytics (Protected)
- `GET /api/analytics/salary/by-country` - Salary statistics by country
- `GET /api/analytics/salary/by-job-title/:country` - Average salary by job title in country
- `GET /api/analytics/salary/department-insights` - Department-wise salary analysis

### Health Check
- `GET /api/health` - Application health status

> **Note**: All employee and analytics endpoints require JWT authentication in the `Authorization: Bearer <token>` header.

## Testing

### Run All Tests
```bash
make test
```

### Test Coverage
```bash
make test-coverage
go tool cover -html=coverage.out
```

### Run Specific Tests
```bash
# Model tests
go test ./internal/models/...

# Handler tests
go test ./internal/handlers/...

# Integration tests
go test ./cmd/server/...
```

## Project Structure

```
assingment/
|-- cmd/
|   |-- server/           # Backend server entry point
|   |-- seed/             # Database seeding script
|-- internal/
|   |-- database/         # Database connection and setup
|   |-- handlers/         # HTTP handlers (including auth)
|   |-- models/           # Data models (including User model)
|   |-- services/         # Business logic (JWT, Password, Auth services)
|   |-- testutils/        # Test utilities
|-- web/                  # React frontend
|   |-- src/
|   |   |-- components/   # React components (including Navigation, Login)
|   |   |   |-- ui/       # shadcn/ui components
|   |   |-- lib/          # Utility functions
|   |   |-- services/     # API service with JWT interceptors
|   |   |-- types/        # TypeScript type definitions
|   |-- public/           # Static assets
|   |-- components.json   # shadcn/ui configuration
|   |-- tailwind.config.js # Tailwind CSS configuration
|-- Makefile              # Build and test automation
|-- Dockerfile            # Container configuration
|-- requirements.md       # Project requirements
|-- SUMMARY.md            # Project summary
|-- TESTING.md            # Testing documentation
```

## Environment Variables

### Backend
- `DATABASE_PATH`: Path to SQLite database (default: `./salary_management.db`)
- `PORT`: Server port (default: `8080`)
- `PUBLIC_DIR`: Static files directory (default: `public`)
- `JWT_SECRET`: Secret key for JWT token signing (default: auto-generated warning)
- `JWT_EXPIRATION`: JWT token expiration time (default: `24h`)

### Frontend
- `VITE_API_URL`: Backend API URL for production (default: `http://localhost:8080/api`)
  
**Development**: Uses Vite proxy (configured in vite.config.js) to avoid CORS issues
**Production**: Uses `VITE_API_URL` environment variable

To configure the frontend API connection for production:
```bash
cd web
cp .env.example .env
# Edit .env file if needed to change the backend URL
```

## Docker Support

```bash
# Build Docker image
docker build -t salary-management .

# Run container
docker run -p 8080:8080 salary-management
```

## Default Login Credentials

The application comes with a default admin user for testing:

- **Email**: `admin@company.com`
- **Password**: `admin123`
- **Role**: HR Manager

> **Note**: The default user is automatically created when the database is initialized. You can create additional users through the authentication system.

## Troubleshooting

### Common Issues

1. **Node.js not found**
   ```bash
   export NVM_DIR="$HOME/.nvm"
   [ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"
   nvm use node
   ```

2. **Port already in use**
   ```bash
   # Kill process on port 8080
   lsof -ti:8080 | xargs kill -9
   
   # Kill process on port 5173
   lsof -ti:5173 | xargs kill -9
   ```

3. **Database connection issues**
   ```bash
   # Remove existing database
   rm salary_management.db
   
   # Restart backend to recreate database
   go run cmd/server/main.go
   ```

4. **Authentication issues**
   ```bash
   # Check JWT_SECRET is set (recommended for production)
   export JWT_SECRET="your-secret-key-here"
   
   # Clear browser localStorage if stuck
   # Open browser dev tools -> Application -> Local Storage -> Clear
   ```

5. **Frontend build issues**
   ```bash
   cd web
   rm -rf node_modules package-lock.json
   npm install
   ```

6. **JWT Token Expired**
   - If you see authentication errors, simply logout and login again
   - Tokens expire after 24 hours by default

## Development Tips

### Backend
- Use `make test` for fast feedback during development
- Check `TESTING.md` for comprehensive testing guidelines
- Use `go run cmd/server/main.go` for hot reload during development
- Authentication tests use in-memory SQLite databases
- JWT secret can be set via `JWT_SECRET` environment variable

### Frontend
- The app uses Vite for fast development with hot reload
- shadcn/ui components provide professional UI elements
- Tailwind CSS is configured with dark theme support
- React Router handles navigation and protected routes
- TypeScript path aliases (`@/`) for cleaner imports
- JWT tokens are automatically managed in localStorage

### Testing
- All tests should run in under 2 seconds
- Use in-memory databases for test isolation
- Authentication tests cover JWT, password, and auth services
- Check test coverage before committing changes
- Frontend components can be tested with React Testing Library

### Authentication Development
- Default admin user: `admin@company.com` / `admin123`
- Tokens expire after 24 hours by default
- Protected routes automatically redirect to login if not authenticated
- API interceptors handle JWT token injection automatically

## Production Deployment

### Backend
```bash
# Build production binary
make build

# Run production binary
./salary-management
```

### Frontend
```bash
cd web
npm run build
# Deploy the build/ directory to your web server
```

## Contributing

1. Follow the existing code structure
2. Write tests for new features
3. Use the Makefile for common tasks
4. Update documentation as needed

## License

This project is for assessment purposes only.
