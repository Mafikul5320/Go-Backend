# 🚗 SpotSync API

SpotSync is a centralized platform for managing parking zones, specifically handling the high-demand reservation of limited EV charging spots. It uses strict Clean Architecture and ensures concurrency safety using database row-level locks.

## 🚀 Live URL
*To be added post-deployment*

## 🛠️ Tech Stack
- **Language**: Go 1.22+
- **Framework**: Echo (`github.com/labstack/echo/v4`)
- **Database**: PostgreSQL (NeonDB / Supabase)
- **ORM**: GORM (`gorm.io/gorm`)
- **Validation**: `go-playground/validator/v10`
- **Authentication**: JWT (`golang-jwt/jwt/v5`) & bcrypt

## 🏛️ Architecture
The project strictly follows **Clean Architecture**:
- **Handler Layer (`/handler`)**: HTTP layer, handles requests, parses JSON, binds & validates DTOs, extracts JWT claims. Handlers do NOT talk to the DB.
- **Service Layer (`/service`)**: Contains all business logic (password hashing, checking zone capacity). Interacts only with Repositories.
- **Repository Layer (`/repository`)**: Data access layer. Handles all GORM queries, including the critical Transaction and `FOR UPDATE` lock to solve the "EV Spot Bottleneck".
- **Models (`/models`)**: Database schemas.
- **DTO (`/dto`)**: Request/Response payloads used by Handlers.

*Dependency Injection happens manually in `main.go`.*

## ⚙️ Setup Instructions

### 1. Prerequisites
- Go 1.22 or higher
- PostgreSQL running locally or remotely

### 2. Installation
```bash
# Clone the repository
git clone <your-repo-url>
cd Assignment-6

# Install dependencies
go mod tidy
```

### 3. Environment Variables
Create a `.env` file in the root directory:
```env
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=spotsync
JWT_SECRET=supersecretkey
```

### 4. Run the Application
```bash
go run main.go
# Or using Air for hot-reloading:
air
```

## 🌐 API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register a new user (driver/admin)
- `POST /api/v1/auth/login` - Login and receive JWT

### Parking Zones
- `POST /api/v1/zones` - Create a parking zone (Admin only)
- `GET /api/v1/zones` - Get all parking zones
- `GET /api/v1/zones/:id` - Get a specific parking zone

### Reservations
- `POST /api/v1/reservations` - Reserve a parking spot
- `GET /api/v1/reservations/my-reservations` - Get my reservations
- `DELETE /api/v1/reservations/:id` - Cancel my reservation
- `GET /api/v1/reservations` - Get all reservations (Admin only)
