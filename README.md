# Smart Parking Backend

Backend API cho hệ thống bãi xe thông minh, viết bằng Go, sử dụng Gin + GORM + MySQL + Redis.

Link: [Frontend](https://github.com/dangngockhieu/SmartParking-Frontend)

## 1) Tính năng chính

- Xác thực người dùng bằng JWT + refresh token.
- Quản lý bãi xe, cổng (gate), thiết bị IoT, thẻ RFID, vị trí đỗ xe.
- Nhận dữ liệu realtime từ thiết bị:
  - Camera nhận diện biển số.
  - RFID scan tại cổng vào/ra.
  - Sensor cập nhật trạng thái ô đỗ.
- Phát sự kiện realtime trạng thái slot qua hub (WebTransport).
- Dashboard thống kê lưu lượng xe, tỉ lệ lấp đầy theo bãi.
- Ví điện tử + nạp tiền qua PayOS (deposit, webhook, cancel-return).
- Tài liệu API qua Swagger.

## 2) Tech Stack

- Go 1.25 (theo `go.mod`)
- Gin
- GORM + MySQL
- Redis
- Swagger (`swaggo/gin-swagger`)

## 3) Cấu trúc thư mục

```text
backend
├── cmd
│   ├── api/main.go             # Entry point API server
│   └── main/seed.go            # Seed dữ liệu mẫu
│
├── configs
│   └── config.go               # Đọc config từ env
│
├── docs/                       # Docs swagger
│
├── internal
│   ├── auth
│   │   ├── mail/               # Mail Service
│   │   ├── token/              # Manage AccessToken and RefreshToken
│   │   └── ...                 # Register, login, logout, refreshToken, resetPassword
│   │
│   ├── common/errors/          # Format response error
│   │
│   ├── modules
│   │   ├── dashboard/          # Thống kê lưu lượng xe
│   │   ├── gate/               # Manage Gate CRUD
│   │   ├── iot_device/         # Manage Iot_Device CRUD
│   │   ├── iot_gateway/        # API camera and rfid_card post to server
│   │   ├── parking_lot/        # CRUD Parking_lot, get slot by lotId, get gate by lotID
│   │   ├── parking_session/    # Manage Parking_Session, caculator fee
│   │   ├── parking_slot/       # CRUD Parking_slot, update status from sensor realtime
│   │   ├── rfid_card/          # Manage Rfid_Card CRUD
│   │   ├── user/               # Create User, ChangePassword, ChangeRole, GetAllUser
│   │   └── wallet/             # Wallet + PayOS deposit
│   │
│   └── realtime/parking/       # Realtime hub/server
│
├── migrations/                 # Migration Database
│
├── pkg
│   ├── database/               # MySQL, Redis
│   ├── middleware/             # Auth, CORS, error handler
│   └── response/               # Format Response API
│
├── templates/                  # Templates send email
│
└── ...                         # Env, README, ...
```

## 4) Chuẩn bị môi trường

### 4.1 Yêu cầu

- Go (khuyến nghị khớp với `go.mod`)
- MySQL 8+
- Redis 6+

### 4.2 Dùng Docker cho MySQL

Repo đã có `docker-compose.yaml` cho MySQL:

```bash
docker compose up -d
```

Mặc định compose map MySQL ra cổng `3307`.

Mở terminal trong docker và chạy câu lệnh sau để tạo db redis:

```bash
docker run --name parking -p 6379:6379 -d redis
```

Sẽ tạo ra db redis chạy trên cổng `6379`.

## 5) Cấu hình môi trường

Project đọc `.env.local` khi `APP_ENV=local`.

1. Tạo file `.env.local` từ `.env.example`.
2. Điền giá trị phù hợp.

Ví dụ tối thiểu để chạy local:

```env
APP_ENV=local
APP_PORT=8080
FRONTEND_URL=http://localhost:3000

DB_HOST=127.0.0.1
DB_PORT=
DB_USER=
DB_PASS=
DB_NAME=

REDIS_ADDR=127.0.0.1:6379
REDIS_PASSWORD=

JWT_ACCESS_SECRET=
JWT_REFRESH_SECRET=

MAIL_HOST=smtp.example.com
MAIL_PORT=587
MAIL_USER=
MAIL_PASS=
VERIFY_URL=http://localhost:3000/verify
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173

# PayOS
PAYOS_CLIENT_ID=
PAYOS_API_KEY=
PAYOS_CHECKSUM_KEY=
PAYOS_RETURN_URL=http://localhost:3000/wallet/return
PAYOS_CANCEL_URL=http://localhost:3000/wallet/cancel
PAYMENT_UPDATE_STATUS_CANCEL_URL=http://localhost:8080/api/v1/wallets/deposit/cancel-return
```

## 6) Khởi tạo database schema

Hiện repo dùng SQL migration files trong `migrations/`.

Có thể áp dụng thủ công file up bằng MySQL client:

```bash
mysql -h 127.0.0.1 -P 3307 -u admin -p parking < migrations/000001_init.up.sql
```

Rollback thủ công:

```bash
mysql -h 127.0.0.1 -P 3307 -u admin -p parking < migrations/000001_init.down.sql
```

## 7) Seed dữ liệu mẫu

Chạy seed:

```bash
go run ./cmd/main/seed.go
```

Dữ liệu mẫu gồm:

- Admin account: `admin@gmail.com` / `123456`
- Lot: `A`
- Devices:
  - `SENSOR_A_001` (sensor cho parking slots)
  - `GATE_IN_A_001` (controller cổng vào)
  - `GATE_OUT_A_001` (controller cổng ra)
- Gates:
  - `Gate In A` (`ENTRY`, MAC `GATE_IN_A_001`)
  - `Gate Out A` (`EXIT`, MAC `GATE_OUT_A_001`)
- RFID cards:
  - `GUEST001` (`GUEST`)
  - `USER001` (`REGISTERED`)
- 8 slots `A1..A8` gắn với `SENSOR_A_001`

## 8) Chạy server

```bash
go run ./cmd/api
```

Health check:

```http
GET /health
```

Swagger UI:

```text
http://localhost:8080/swagger/index.html
```

## 9) Luồng IoT đang dùng

### 9.1 Camera -> cache biển số

```http
POST /api/v1/iot/camera
Content-Type: application/json

{
  "gate_id": 1,
  "plate_number": "51A-123.45"
}
```

### 9.2 RFID scan -> tạo/kết thúc parking session

```http
POST /api/v1/iot/rfid
Content-Type: application/json

{
  "gate_id": 1,
  "mac_address": "GATE_IN_A_001",
  "rfid_uid": "GUEST001"
}
```

Ghi chú:

- Với gate `ENTRY`: tạo `parking_sessions` mới, dùng `plate_number` lấy từ cache camera.
- Với gate `EXIT`: tìm session active và kết thúc session (tính phí).

### 9.3 Sensor slot -> cập nhật trạng thái ô đỗ

```http
POST /api/v1/parking-slots/sensor
Content-Type: application/json

{
  "mac": "SENSOR_A_001",
  "port": 1,
  "is_occupied": true
}
```

[Luồng hoạt động chi tiết](./operational_flow.md)

## 10) Route groups chính

- `/api/v1/auth`
- `/api/v1/users`
- `/api/v1/parking-lots`
- `/api/v1/parking-slots`
- `/api/v1/parking-sessions`
- `/api/v1/gates`
- `/api/v1/rfid-cards`
- `/api/v1/iot-devices`
- `/api/v1/iot` (camera, rfid)
- `/api/v1/dashboard`
- `/api/v1/wallets` (deposit, transactions, webhook)

## 11) Kiểm tra nhanh chất lượng code

```bash
go test ./...
```

## 12) Notes triển khai

- Nếu muốn bật WebTransport server (`:4433`), cần cung cấp `cert.pem` và `key.pem` ở root project.
- Có thể override đường dẫn TLS bằng `TLS_CERT` và `TLS_KEY`.
- CORS được kiểm soát qua `CORS_ALLOWED_ORIGINS`.
- `VERIFY_URL` là bắt buộc theo config loader.
- `PAYMENT_UPDATE_STATUS_CANCEL_URL` là bắt buộc để cập nhật trạng thái khi hủy thanh toán PayOS.

## 13) Troubleshooting

- Lỗi `APP_PORT is required`: thiếu biến `APP_PORT` trong env.
- Lỗi `VERIFY_URL is required`: thiếu biến `VERIFY_URL` trong env.
- Lỗi kết nối MySQL: kiểm tra đúng `DB_HOST/DB_PORT/DB_USER/DB_PASS/DB_NAME` và container DB đã chạy.
- RFID bị reject do chưa có biển số: cần gọi endpoint camera trước rồi mới quét RFID tại cùng `gate_id`.
