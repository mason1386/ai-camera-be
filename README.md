# AI Camera Backend

Dự án Backend cho hệ thống AI Camera, được xây dựng bằng Go (Golang) theo kiến trúc Hexagonal / Clean Architecture.

## 🛠 Công nghệ sử dụng

- **Ngôn ngữ**: Go 1.24+
- **Database**: PostgreSQL
- **Cache**: Redis
- **Message Queue**: Kafka
- **Framework**: Gin Gonic
- **Documentation**: Swagger (Swag)
- **Infrastructure**: Docker & Docker Compose

## 🚀 Yêu cầu hệ thống

Trước khi bắt đầu, hãy đảm bảo bạn đã cài đặt:

- [Go](https://go.dev/dl/) (phiên bản 1.24 trở lên)
- [Docker](https://www.docker.com/) & Docker Compose
- [Make](https://www.gnu.org/software/make/) (Thường có sẵn trên Linux/macOS, Windows có thể dùng qua Git Bash hoặc cài đặt riêng)

## 📦 Cài đặt & Chạy dự án

### 1. Cấu hình môi trường

Copy file cấu hình mẫu `.env.example` thành `.env`:

```bash
cp .env.example .env
```

Kiểm tra và chỉnh sửa file `config/config.yaml` nếu bạn muốn thay đổi cấu hình mặc định (Database, Redis, Kafka, Port).

### 2. Khởi động hạ tầng (Database, Redis, Kafka)

Sử dụng Docker Compose để khởi chạy các dịch vụ phụ trợ:

```bash
make docker-up
```

Lệnh này sẽ khởi động Postgres, Redis, Zookeeper và Kafka.

### 3. Khởi tạo Database

Chạy lệnh sau để seed dữ liệu mẫu (nếu đã cấu hình script):

```bash
make seed-db
```

### 4. Chạy ứng dụng

#### API Server
Để chạy API server:

```bash
make run-api
```
Server sẽ lắng nghe tại cổng `8080` (mặc định).

#### Worker (Background Jobs)
Để chạy worker xử lý các tác vụ nền:

```bash
make run-worker
```

## 📚 API Documentation

Sau khi khởi động server, bạn có thể truy cập tài liệu API (Swagger UI) tại đường dẫn:

👉 [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

## 📁 Cấu trúc dự án

Dự án tuân theo cấu trúc Clean Architecture / Hexagonal Architecture:

```
.
├── cmd/                # Entry points của ứng dụng (api, worker)
├── config/             # Các file cấu hình (yaml, go)
├── docs/               # Tài liệu API (Swagger generated)
├── internal/           # Mã nguồn chính (Private application code)
│   ├── adapters/       # Các adapters giao tiếp hạ tầng (Handler HTTP, Postgres, Redis, Kafka...)
│   ├── core/           # Business logic trung tâm
│   │   ├── domain/     # Domain entities (Models)
│   │   ├── ports/      # Interfaces định nghĩa inputs/outputs (Repository interfaces, Service interfaces)
│   │   └── services/   # Implementation của business logic
├── migrations/         # Database migrations
├── pkg/                # Các thư viện dùng chung (Logging, Utils...)
├── scripts/            # Scripts tiện ích (build, seed db...)
└── sql/                # SQL queries (dùng cho SQLC để generate code)
```

## 🛠 Các lệnh Makefile hữu ích

| Lệnh | Mô tả |
|------|-------|
| `make build` | Build ứng dụng ra file binary (vào thư mục `bin/`) |
| `make run-api` | Chạy API Server trực tiếp (go run) |
| `make run-worker` | Chạy Worker trực tiếp (go run) |
| `make docker-up` | Khởi động các containers (Postgres, Redis, Kafka) |
| `make docker-down` | Dừng và xóa các containers |
| `make lint` | Kiểm tra lỗi code (GolangCI-Lint) |
| `make test` | Chạy toàn bộ Unit Tests |
| `make gen-proto` | Generate Go code từ file Protobuf (nếu có sử dụng gRPC) |
