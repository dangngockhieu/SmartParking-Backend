# Mô tả luồng hoạt động: Embedded ↔ Server ↔ Database (Go backend hiện tại)

Tài liệu này mô tả đúng theo code Go/Gin/GORM đang có trong repository ở thời điểm hiện tại.

## 1. Thiết bị nhúng giao tiếp với server như thế nào?

Thiết bị nhúng vẫn dùng HTTP REST API gửi JSON lên server.

Ví dụ firmware ESP32 gửi trạng thái cảm biến chỗ đỗ:

```cpp
#include <HTTPClient.h>

HTTPClient http;
http.begin("http://192.168.1.100:8080/api/v1/parking-slots/sensor");
http.addHeader("Content-Type", "application/json");

String body = "{\"mac\":\"DEVICE_A_001\",\"port\":3,\"is_occupied\":true}";

int httpCode = http.POST(body);
String response = http.getString();
```

### Endpoint IoT đang hoạt động trong code hiện tại

| Thiết bị     | Endpoint                          | Payload                    |
| ------------ | --------------------------------- | -------------------------- |
| ESP32 + SR04 | POST /api/v1/parking-slots/sensor | { mac, port, is_occupied } |

### Endpoint đã có module nhưng chưa được mount vào main

| Luồng                   | Trạng thái                                                                 |
| ----------------------- | -------------------------------------------------------------------------- |
| POST /api/v1/iot/rfid   | Có code module iot_gateway, nhưng chưa đăng ký route trong cmd/api/main.go |
| POST /api/v1/iot/camera | Chưa thấy route/handler trong code hiện tại                                |

---

## 2. Dữ liệu đi vào server theo các tầng nào?

Luồng request trong Go backend:

1. Gin Router nhận request tại base path /api/v1
2. Global middleware chạy lần lượt:
   - gin.Logger()
   - gin.Recovery()
   - CORS middleware
   - ErrorHandler middleware
3. Route-level middleware (nếu có):
   - Auth middleware (Bearer JWT)
   - Role middleware (RequireRoles)
4. Handler bind JSON vào DTO bằng ShouldBindJSON
5. Service xử lý nghiệp vụ
6. Repository truy vấn MySQL qua GORM
7. Trả response JSON chuẩn qua pkg/response

---

## 3. Server ghi xuống Database như thế nào?

Server không dùng Prisma. Backend hiện tại dùng GORM.

### Kết nối DB

Config đọc từ biến môi trường:

- DB_HOST
- DB_PORT
- DB_USER
- DB_PASS
- DB_NAME

Sau đó dựng DSN và mở kết nối:

```go
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
```

### Cách đọc/ghi

Service gọi Repository, Repository dùng GORM:

```go
// insert
r.db.Create(&entity)

// select
r.db.Where("uid = ?", uid).First(&card)

// update
r.db.Model(&ParkingSlot{}).Where("id = ?", id).Update("status", status)
```

GORM tự sinh SQL tương ứng và gửi đến MySQL.

---

## 4. Luồng Sensor Update đang chạy thật

Thiết bị SR04 gửi:

POST /api/v1/parking-slots/sensor

Body:

```json
{
  "mac": "DEVICE_A_001",
  "port": 3,
  "is_occupied": true
}
```

Luồng xử lý:

1. Handler nhận JSON, bind vào SensorUpdateParkingSlotRequest
2. Service tìm slot theo cặp device_mac + port_number
3. Nếu slot đang MAINTAIN thì trả changed=false
4. Nếu có thay đổi trạng thái:
   - Update status trong bảng parking_slots
   - Broadcast realtime event slotStatusUpdated qua hub
5. Trả về UpdateParkingSlotResponse

---

## 5. Luồng RFID + PlateCache (trạng thái hiện tại)

Module iot_gateway đã có đầy đủ logic RFID:

- Validate gate
- Check MAC với gate
- Check RFID card
- Consume biển số từ PlateCache
- Tạo hoặc đóng parking session theo loại cổng

Tuy nhiên hiện tại có 2 điểm quan trọng:

1. Route iot_gateway chưa được đăng ký trong cmd/api/main.go nên endpoint /api/v1/iot/rfid chưa active.
2. Chưa thấy endpoint camera set dữ liệu vào PlateCache, nên nếu chỉ bật riêng luồng RFID thì dễ gặp reject Chưa có biển số từ camera.

Nói ngắn gọn: logic RFID có sẵn trong code, nhưng wiring runtime chưa hoàn tất.

---

## 6. Dữ liệu lưu ở đâu?

| Dữ liệu          | Nơi lưu                  | Ghi chú                                          |
| ---------------- | ------------------------ | ------------------------------------------------ |
| users            | MySQL                    | Tài khoản hệ thống                               |
| parking_lots     | MySQL                    | Khu bãi xe                                       |
| iot_devices      | MySQL                    | Thiết bị IoT                                     |
| parking_slots    | MySQL                    | Vị trí đỗ + trạng thái                           |
| gates            | MySQL                    | Cổng vào/ra                                      |
| rfid_cards       | MySQL                    | Thẻ RFID                                         |
| parking_sessions | MySQL                    | Phiên gửi xe                                     |
| slot_histories   | MySQL                    | Lịch sử đổi thiết bị/trạng thái slot             |
| vehicle_logs     | MySQL                    | Log xe theo slot                                 |
| PlateCache       | RAM (module iot_gateway) | Cache biển số tạm thời TTL 5 phút, consume 1 lần |

---

## 7. Tổng quan kiến trúc (đúng theo code hiện tại)

```text
ESP32 / Camera (LAN)
        |
        | HTTP JSON
        v
Gin Router (/api/v1)
        |
        +--> Global middleware: Logger, Recovery, CORS, ErrorHandler
        |
        +--> Route middleware (tùy route): Auth JWT, RequireRoles
        |
        v
Handler (bind DTO)
        v
Service (business logic)
        v
Repository (GORM)
        v
MySQL

Realtime:
Service -> parking hub -> WebSocket/WebTransport clients
```

---

## 8. Gợi ý để hoàn chỉnh luồng Embedded đầy đủ như thiết kế ban đầu

1. Mount module iot_gateway trong cmd/api/main.go.
2. Bổ sung endpoint camera để ghi PlateCache.Set(gateID, plateNumber).
3. Thống nhất contract payload giữa firmware và backend (snake_case hay camelCase).
4. Cập nhật Swagger cho tất cả route IoT runtime.
