# UlumFR Backend API

Backend untuk portofolio/CMS menggunakan **Go (Fiber)**, **Neon (PostgreSQL)**, **Prisma**, dan **Cloudflare R2**.

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- Neon PostgreSQL database

### Setup

1. **Clone repository**
   ```bash
   git clone https://github.com/ulumfr/ulumfr-be.git
   cd ulumfr-be
   ```

2. **Copy environment file**
   ```bash
   cp .env.example .env
   # Edit .env dengan credentials Anda
   ```

3. **Install dependencies**
   ```bash
   go mod download
   ```

4. **Generate Prisma client**
   ```bash
   go run github.com/steebchen/prisma-client-go generate
   ```

5. **Push schema ke database**
   ```bash
   go run github.com/steebchen/prisma-client-go db push
   ```

6. **Run server**
   ```bash
   go run ./cmd/api
   ```

7. **Open Swagger docs**
   ```
   http://localhost:8080/swagger/index.html
   ```

## 📁 Project Structure

```
.
├── api/                # Vercel serverless handler
├── cmd/api/            # Local development entry point
├── docs/               # Swagger documentation
├── pkg/
│   ├── config/         # Configuration loader
│   ├── domain/         # Domain entities & DTOs
│   ├── handler/        # HTTP handlers
│   ├── middleware/     # JWT Auth, CORS, Rate limiting
│   ├── repository/     # Database operations
│   ├── service/        # Business logic
│   └── storage/        # Cloudflare R2 client
├── prisma/             # Prisma schema
└── Makefile            # Build commands
```

## 🔌 API Endpoints

### Auth

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/register` | Register new user |
| POST | `/api/v1/auth/login` | Login → get JWT tokens |
| POST | `/api/v1/auth/refresh` | Refresh access token |
| GET | `/api/v1/auth/me` | Get current user |
| POST | `/api/v1/auth/logout` | Logout (invalidate session) |
| POST | `/api/v1/auth/logout-all` | Logout from all devices |

### Public

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/api/v1/public/projects` | List published projects |
| GET | `/api/v1/public/projects/:slug` | Get project by slug |
| GET | `/api/v1/public/categories` | List categories |
| GET | `/api/v1/public/tags` | List tags |
| GET | `/api/v1/public/careers` | List careers |
| GET | `/api/v1/public/educations` | List educations |
| GET | `/api/v1/public/resume` | Get active resume |
| POST | `/api/v1/public/contact` | Submit contact form |

### Admin (Protected - JWT + Admin Role)

| Method | Endpoint | Description |
|--------|----------|-------------|
| CRUD | `/api/v1/admin/projects` | Project management |
| CRUD | `/api/v1/admin/categories` | Category management |
| CRUD | `/api/v1/admin/tags` | Tag management |
| CRUD | `/api/v1/admin/careers` | Career management |
| CRUD | `/api/v1/admin/educations` | Education management |
| CRUD | `/api/v1/admin/resumes` | Resume management |
| CRUD | `/api/v1/admin/contacts` | Contact management |
| POST | `/api/v1/admin/upload-url` | Get presigned upload URL |

📖 **Full API docs:** [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

## 🔐 Authentication

Menggunakan **JWT + Session-based refresh token**:

1. **Login** → Access token (JWT, 15min) + Refresh token (stored in DB, 7 days)
2. **API Request** → Send `Authorization: Bearer <access_token>`
3. **Refresh** → Use refresh token to get new access token
4. **Logout** → Delete session from database

### Make user admin

```sql
UPDATE users SET role = 'ADMIN' WHERE email = 'your@email.com';
```

## 📤 File Upload (R2)

Presigned URL flow:
1. Admin request: `POST /api/v1/admin/upload-url`
2. Backend generate presigned PUT URL
3. Frontend upload langsung ke R2
4. Simpan file URL ke database

## 📝 Make Commands

```bash
make help        # Show all commands
make build       # Build binary
make run         # Build & run
make dev         # Run with hot reload (requires air)
make generate    # Generate Prisma client
make migrate     # Push schema to database
make test        # Run tests
make build-prod  # Production build (Linux)
```

## 🛠 Environment Variables

```env
# Database
DATABASE_URL=postgres://user:pass@host/db?sslmode=require

# Server
PORT=8080
APP_ENV=development

# JWT
JWT_SECRET=your-secret-key-min-32-chars
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h

# Cloudflare R2
R2_ACCOUNT_ID=your_account_id
R2_ACCESS_KEY_ID=your_access_key
R2_SECRET_ACCESS_KEY=your_secret_key
R2_BUCKET_NAME=your_bucket
R2_PUBLIC_URL=https://your-r2-domain.com
```

## 📄 License

MIT
