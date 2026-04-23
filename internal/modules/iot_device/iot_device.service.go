package iot_device

import (
	"strings"
	"time"

	appErrors "backend/internal/common/errors"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateDevice(req CreateIoTDeviceRequest) (*IoTDeviceResponse, error) {
	req.MacAddress = strings.TrimSpace(req.MacAddress)
	req.DeviceName = strings.TrimSpace(req.DeviceName)

	if req.MacAddress == "" || req.DeviceName == "" {
		return nil, appErrors.NewBadRequest("mac_address va device_name la bat buoc")
	}

	exist, err := s.repo.FindByMacAddress(req.MacAddress)
	if err != nil {
		return nil, appErrors.NewInternal("Khong the kiem tra thiet bi IoT")
	}
	if exist != nil {
		return nil, appErrors.NewBadRequest("Device already exists")
	}

	now := time.Now()
	device := &IoTDevice{
		MacAddress: req.MacAddress,
		DeviceName: &req.DeviceName,
		LotID:      req.LotID,
		Status:     DeviceStatusActive,
		LastSeen:   &now,
	}

	if err := s.repo.CreateDevice(device); err != nil {
		return nil, appErrors.NewInternal("Tao thiet bi IoT that bai")
	}

	createdDevice, err := s.repo.FindByMacAddressWithLot(req.MacAddress)
	if err != nil {
		return nil, appErrors.NewInternal("Khong the lay thiet bi IoT vua tao")
	}
	if createdDevice == nil {
		return nil, appErrors.NewInternal("Khong tim thay thiet bi IoT vua tao")
	}

	resp := ToIoTDeviceResponse(createdDevice)
	return &resp, nil
}

func (s *Service) FindAllDevices(query GetIoTDevicesQuery) ([]IoTDeviceResponse, error) {
	query.Keyword = strings.TrimSpace(query.Keyword)

	devices, err := s.repo.FindAll(query)
	if err != nil {
		return nil, appErrors.NewInternal("Khong the lay danh sach thiet bi IoT")
	}

	return ToIoTDeviceResponses(devices), nil
}

func (s *Service) UpdateDevice(macAddress string, req UpdateIoTDeviceRequest) (*IoTDeviceResponse, error) {
	macAddress = strings.TrimSpace(macAddress)
	if macAddress == "" {
		return nil, appErrors.NewBadRequest("mac_address khong hop le")
	}

	device, err := s.repo.FindByMacAddress(macAddress)
	if err != nil {
		return nil, appErrors.NewInternal("Khong the kiem tra thiet bi IoT")
	}
	if device == nil {
		return nil, appErrors.NewNotFound("Khong tim thay thiet bi IoT")
	}

	updates := map[string]interface{}{}

	if req.DeviceName != nil {
		name := strings.TrimSpace(*req.DeviceName)
		if name == "" {
			return nil, appErrors.NewBadRequest("device_name khong duoc de trong")
		}
		updates["device_name"] = name
	}

	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if req.HasLotID {
		updates["lot_id"] = req.LotID
	}

	if len(updates) == 0 {
		return nil, appErrors.NewBadRequest("Khong co du lieu cap nhat")
	}

	if err := s.repo.UpdateByMacAddress(macAddress, updates); err != nil {
		return nil, appErrors.NewInternal("Cap nhat thiet bi IoT that bai")
	}

	updatedDevice, err := s.repo.FindByMacAddressWithLot(macAddress)
	if err != nil {
		return nil, appErrors.NewInternal("Khong the lay thiet bi IoT sau cap nhat")
	}
	if updatedDevice == nil {
		return nil, appErrors.NewNotFound("Khong tim thay thiet bi IoT")
	}

	resp := ToIoTDeviceResponse(updatedDevice)
	return &resp, nil
}
