package iot_device

import (
	"encoding/json"
	"time"
)

type CreateIoTDeviceRequest struct {
	MacAddress string  `json:"mac_address" binding:"required"`
	DeviceName string  `json:"device_name" binding:"required"`
	LotID      *uint64 `json:"lot_id"`
}

type GetIoTDevicesQuery struct {
	LotID   *uint64       `form:"lot_id"`
	Status  *DeviceStatus `form:"status" binding:"omitempty,oneof=ACTIVE INACTIVE ERROR"`
	Keyword string        `form:"keyword"`
}

type UpdateIoTDeviceRequest struct {
	DeviceName *string       `json:"device_name" binding:"omitempty"`
	Status     *DeviceStatus `json:"status" binding:"omitempty,oneof=ACTIVE INACTIVE ERROR"`
	LotID      *uint64       `json:"lot_id"`
	HasLotID   bool          `json:"-"`
}

type IoTDeviceResponse struct {
	MacAddress string       `json:"mac_address"`
	DeviceName string       `json:"device_name"`
	Status     DeviceStatus `json:"status"`
	LotID      *uint64      `json:"lot_id"`
	LotName    *string      `json:"lot_name"`
	LastSeen   *time.Time   `json:"last_seen"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
}

func ToIoTDeviceResponse(device *IoTDeviceWithLot) IoTDeviceResponse {
	deviceName := ""
	if device.DeviceName != nil {
		deviceName = *device.DeviceName
	}

	return IoTDeviceResponse{
		MacAddress: device.MacAddress,
		DeviceName: deviceName,
		Status:     device.Status,
		LotID:      device.LotID,
		LotName:    device.LotName,
		LastSeen:   device.LastSeen,
		CreatedAt:  device.CreatedAt,
		UpdatedAt:  device.UpdatedAt,
	}
}

func ToIoTDeviceResponses(devices []IoTDeviceWithLot) []IoTDeviceResponse {
	resp := make([]IoTDeviceResponse, 0, len(devices))
	for i := range devices {
		resp = append(resp, ToIoTDeviceResponse(&devices[i]))
	}
	return resp
}

func (r *UpdateIoTDeviceRequest) UnmarshalJSON(data []byte) error {
	type alias struct {
		DeviceName *string       `json:"device_name"`
		Status     *DeviceStatus `json:"status"`
		LotID      *uint64       `json:"lot_id"`
	}

	var decoded alias
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}

	r.DeviceName = decoded.DeviceName
	r.Status = decoded.Status
	r.LotID = decoded.LotID

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	_, r.HasLotID = raw["lot_id"]

	return nil
}
