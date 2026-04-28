package iot_gateway

// CameraPlateRequest là payload camera AI gửi lên
type CameraPlateRequest struct {
	GateID      uint64 `json:"gate_id" binding:"required"`
	PlateNumber string `json:"plate_number" binding:"required"`
}

// CameraPlateResponse là response trả về cho camera
type CameraPlateResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// RfidScanRequest là payload ESP32 + RC522 gửi lên
type RfidScanRequest struct {
	GateID     uint64 `json:"gate_id"     binding:"required"`
	MacAddress string `json:"mac_address" binding:"required"`
	RfidUID    string `json:"rfid_uid"    binding:"required"`
}

// RfidScanResponse là response trả về cho ESP32
// ESP32 dùng action để quyết định mở/đóng barie
// lcd_line1, lcd_line2 hiển thị lên màn hình LCD
type RfidScanResponse struct {
	Success  bool   `json:"success"`
	Action   string `json:"action"`    // "open_barrier" | "reject"
	LCDLine1 string `json:"lcd_line1"` // VD: "BS:51A-123.45"
	LCDLine2 string `json:"lcd_line2"` // VD: "Moi vao!" hoặc "Tam biet!"
	Message  string `json:"message"`   // log nội bộ
}
