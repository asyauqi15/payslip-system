# Payslip System

A scalable Go-based payslip management system that handles employee attendance, overtime tracking, reimbursements, and automated payroll generation.

## Tech Stack

- **Backend**: Go 1.23.8
- **Framework**: Chi Router
- **Database**: PostgreSQL with GORM
- **Authentication**: JWT
- **Migration**: Goose
- **API Documentation**: OpenAPI 3.0 with Swagger UI
- **Configuration**: Viper

## Project Structure

```
payslip-system/
├── api/                        # API specifications
│   └── openapi.yml             # OpenAPI 3.0 specification
├── cmd/                        # CLI commands
│   ├── cmd.go                  # Root command setup
│   ├── http_server.go          # HTTP server command
│   ├── migrate.go              # Database migration command
│   ├── registry.go             # Dependency injection registry
│   └── seed.go                 # Database seeding command
├── config.yml                  # Application configuration
├── db/                         # Database files
│   ├── migrations/             # Database migration files
│   └── seeds/                  # Database seed files
├── internal/                   # Private application code
│   ├── config.go               # Configuration management
│   ├── constant/               # Application constants
│   ├── entity/                 # Database entities/models
│   │   ├── attendance.go
│   │   ├── attendance_period.go
│   │   ├── employee.go
│   │   ├── overtime.go
│   │   ├── payroll.go
│   │   ├── payslip.go
│   │   ├── reimbursement.go
│   │   └── user.go
│   ├── handler/             # HTTP handlers
│   │   ├── admin/           # Admin endpoints
│   │   ├── auth/            # Authentication endpoints
│   │   └── employee/        # Employee endpoints
│   ├── repository/          # Data access layer
│   ├── transport/           # HTTP transport layer
│   │   ├── rest.go          # REST API setup
│   │   ├── swagger.go       # Swagger UI setup
│   │   └── swagger/         # Swagger UI assets
│   └── usecase/             # Business logic layer
│       ├── attendance/
│       ├── attendance_period/
│       ├── auth/
│       ├── overtime/
│       ├── payroll/
│       ├── payslip/
│       └── reimbursement/
├── pkg/                    # Public packages
│   ├── http/               # HTTP utilities
│   ├── jwt-auth/           # JWT authentication
│   └── openapi/            # Generated OpenAPI code
├── Makefile                # Build and development commands
├── go.mod                  # Go module definition
├── go.sum                  # Go module checksums
└── main.go                 # Application entry point
```

## Getting Started

### Prerequisites

- Go 1.23.8 or later
- PostgreSQL 12 or later
- Make (optional, for using Makefile commands)

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd payslip-system
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Install required tools**
   ```bash
   # Install Goose for database migrations
   make install-goose
   
   # Install OpenAPI code generator
   make install-openapi
   ```

### Database Setup

1. **Create a PostgreSQL database**
   ```sql
   CREATE DATABASE payslip_system;
   ```

2. **Update configuration**
   
   Edit `config.yml` with your database connection details:
   ```yaml
   database:
     source: "postgresql://username:password@localhost/payslip_system?sslmode=disable"
   ```

3. **Run database migrations**
   ```bash
   make migrate
   ```

4. **Seed the database with initial data**
   ```bash
   make seed
   ```

### Running the Application

1. **Start the HTTP server**
   ```bash
   make run-http
   ```
   
   The server will start on `http://localhost:8000`

2. **Alternative: Run directly with Go**
   ```bash
   go run main.go http_server
   ```

## API Documentation

### Swagger UI

Once the server is running, you can access the interactive API documentation at:
- **Swagger UI**: http://localhost:8000/static/swagger/

### Authentication

The API uses JWT (JSON Web Tokens) for authentication. To access protected endpoints:

1. **Login to get access token**
   ```bash
   curl -X POST http://localhost:8000/auth/login \
     -H "Content-Type: application/json" \
     -d '{"username": "admin", "password": "password"}'
   ```

2. **Use the access token in subsequent requests**
   ```bash
   curl -X GET http://localhost:8000/admin/employees \
     -H "Authorization: Bearer <your-access-token>"
   ```

### Default Admin Credentials

- **Username**: `admin`
- **Password**: `password`

### Employee Test Credentials

The system comes with 100 pre-seeded employee accounts:
- **Usernames**: `employee001` to `employee100`
- **Password**: `password123` (for all employees)

## API Endpoints

### Authentication
- `POST /auth/login` - User login

### Admin Endpoints
- `POST /admin/attendance-periods` - Create attendance period
- `POST /admin/payrolls` - Run payroll for a period
- `GET /admin/payrolls/{id}` - Get payroll summary

### Employee Endpoints
- `POST /employee/attendance` - Submit daily attendance
- `POST /employee/overtime` - Submit overtime request
- `POST /employee/reimbursement` - Submit reimbursement request
- `GET /employee/payslips/{id}` - Get payslip history

## Business Logic

### Salary Calculation

The system calculates employee compensation using the following components:

1. **Base Salary**: Fixed monthly salary per employee
2. **Prorated Salary**: `(Base Salary × Attendance Days) / Total Working Days`
3. **Overtime Pay**: `Hourly Rate × Overtime Hours × 2`
   - Hourly Rate = `Base Salary / (22 working days × 8 hours)`
   - Overtime is paid at 2x the regular hourly rate
4. **Reimbursements**: Sum of approved expense reimbursements
5. **Total Take Home**: `Prorated Salary + Overtime Pay + Reimbursements`

### Attendance Rules

- Employees can only submit attendance for the current date
- Multiple attendance submissions for the same date are not allowed
- Working days calculation excludes weekends (Saturday and Sunday)

### Overtime Rules

- Overtime can only be submitted for weekdays (Monday-Friday)
- Overtime periods cannot overlap with existing overtime records
- Minimum overtime duration validation can be implemented

## Development

### Code Generation

**Generate OpenAPI client/server code**
```bash
make generate-openapi
```

### Database Operations

**Create a new migration**
```bash
make create-migration NAME=your_migration_name
```

**Create a new seed file**
```bash
make create-seed NAME=your_seed_name
```

**Run migrations**
```bash
make migrate
```

**Rollback latest migration**
```bash
go run main.go migrate --rollback
```

**Run database seeds**
```bash
make seed
```

### Testing

Run the test suite:
```bash
go test ./...
```

### Available Make Commands

```bash
make run-http              # Start HTTP server
make install-goose         # Install Goose migration tool
make install-openapi       # Install OpenAPI code generator
make create-migration      # Create new migration (requires NAME=...)
make create-seed           # Create new seed file (requires NAME=...)
make migrate               # Run database migrations
make seed                  # Run database seeds
make generate-openapi      # Generate OpenAPI code
```

## Configuration

The application uses a YAML configuration file (`config.yml`):

```yaml
database:
  source: "postgresql://postgres:postgres@localhost/payslip_system?sslmode=disable"

http_server:
  port: 8000
  access_token_secret_encoded: "your-access-token-secret"
  refresh_token_secret_encoded: "your-refresh-token-secret"
  access_token_duration: 24h
  refresh_token_duration: 48h
  read_header_timeout: 2s
  read_timeout: 2m
  idle_timeout: 5m
  write_timeout: 0s
```
