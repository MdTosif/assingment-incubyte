# Salary Management Tool

A minimal yet usable salary management tool for an organization with 10,000 employees, designed for HR Managers to manage employee data and gain salary insights.

## Overview

This application consists of:
- **Backend**: Go REST API with SQLite database
- **Frontend**: React application with Vite
- **Database**: SQLite with GORM ORM

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

### Employee Management
- `GET /api/employees` - List all employees with pagination
- `POST /api/employees` - Create new employee
- `GET /api/employees/:id` - Get specific employee
- `PUT /api/employees/:id` - Update employee
- `DELETE /api/employees/:id` - Delete employee

### Salary Analytics
- `GET /api/analytics/salary/by-country` - Salary statistics by country
- `GET /api/analytics/salary/by-job-title/:country` - Average salary by job title in country
- `GET /api/analytics/salary/department-insights` - Department-wise salary analysis

### Health Check
- `GET /api/health` - Application health status

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
|   |-- handlers/         # HTTP handlers
|   |-- models/           # Data models
|   |-- testutils/        # Test utilities
|-- web/                  # React frontend
|   |-- src/              # React components
|   |-- public/           # Static assets
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

### Frontend
- `VITE_API_URL`: Backend API URL (default: `http://localhost:8080`)

## Docker Support

```bash
# Build Docker image
docker build -t salary-management .

# Run container
docker run -p 8080:8080 salary-management
```

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

4. **Frontend build issues**
   ```bash
   cd web
   rm -rf node_modules package-lock.json
   npm install
   ```

## Development Tips

### Backend
- Use `make test` for fast feedback during development
- Check `TESTING.md` for comprehensive testing guidelines
- Use `go run cmd/server/main.go` for hot reload during development

### Frontend
- The app uses Vite for fast development
- Tailwind CSS is configured for styling
- React Router for navigation

### Testing
- All tests should run in under 2 seconds
- Use in-memory databases for test isolation
- Check test coverage before committing changes

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
