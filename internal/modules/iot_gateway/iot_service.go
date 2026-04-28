package iot_gateway

import (
	"fmt"
	"strings"
	"time"

	appErrors "backend/internal/common/errors"
	"backend/internal/modules/gate"
	"backend/internal/modules/parking_session"
	"backend/internal/modules/parking_slot"
	"backend/internal/modules/rfid_card"
	"backend/internal/modules/user"
)

type Service struct {
	plateCache         *PlateCache
	gateService        *gate.Service
	rfidService        *rfid_card.Service
	sessionService     *parking_session.Service
	parkingSlotService *parking_slot.Service
	userService        *user.Service
}

func NewService(
	plateCache *PlateCache,
	gateService *gate.Service,
	rfidService *rfid_card.Service,
	sessionService *parking_session.Service,
	parkingSlotService *parking_slot.Service,
	userService *user.Service,
) *Service {
	return &Service{
		plateCache:         plateCache,
		gateService:        gateService,
		rfidService:        rfidService,
		sessionService:     sessionService,
		parkingSlotService: parkingSlotService,
		userService:        userService,
	}
}

// HandleCameraPlate lưu biển số tạm trong PlateCache theo gateID
func (s *Service) HandleCameraPlate(req *CameraPlateRequest) (*CameraPlateResponse, error) {
	g, err := s.gateService.FindByID(req.GateID)
	if err != nil {
		if appErr, ok := err.(*appErrors.AppError); ok && appErr.StatusCode == 404 {
			return &CameraPlateResponse{Success: false, Message: "Gate not found"}, nil
		}
		return nil, appErrors.NewInternal("Gate check error")
	}
	if !g.IsActive {
		return &CameraPlateResponse{Success: false, Message: "Gate is inactive"}, nil
	}

	plateNumber := strings.TrimSpace(req.PlateNumber)
	if plateNumber == "" {
		return nil, appErrors.NewBadRequest("Plate number is required")
	}

	s.plateCache.Set(req.GateID, plateNumber)

	return &CameraPlateResponse{
		Success: true,
		Message: "Plate saved from camera",
	}, nil
}

// HandleRfidScan xử lý toàn bộ luồng khi ESP32 quẹt thẻ RFID
func (s *Service) HandleRfidScan(req *RfidScanRequest) (*RfidScanResponse, error) {
	// 1. Validate gate tồn tại và isActive
	g, err := s.gateService.FindByID(req.GateID)
	if err != nil {
		if appErr, ok := err.(*appErrors.AppError); ok && appErr.StatusCode == 404 {
			return rejectResponse("No gate"), nil
		}
		return nil, appErrors.NewInternal("Gate check error")
	}
	if !g.IsActive {
		return rejectResponse("Gate off"), nil
	}

	// 2. Kiểm tra macAddress khớp với gate
	if g.MacAddress != req.MacAddress {
		return rejectResponse("MAC mismatch"), nil
	}

	// 3. Tìm thẻ RFID theo UID
	card, err := s.rfidService.FindByUID(req.RfidUID)
	if err != nil {
		if appErr, ok := err.(*appErrors.AppError); ok && appErr.StatusCode == 404 {
			return rejectResponse("Unknown card"), nil
		}
		return nil, appErrors.NewInternal("RFID check error")
	}
	if !card.IsActive {
		return rejectResponse("Card disabled"), nil
	}

	// 4. Peek biển số từ PlateCache (CHỈ ĐỌC, chưa xóa)
	plateNumber, ok := s.plateCache.Peek(req.GateID)
	if !ok {
		return rejectResponse("No plate"), nil
	}

	// 5. Xử lý theo loại cổng (consume xảy ra bên trong khi mọi thứ hợp lệ)
	switch g.Type {
	case gate.GateTypeEntry:
		// Chỉ kiểm tra trùng biển số ở cổng VÀO
		existingPlate, err := s.sessionService.FindActiveByPlateNumber(plateNumber)
		if err == nil && existingPlate != nil {
			return rejectResponse("Plate in use"), nil
		}
		return s.handleEntry(g, card, plateNumber)
	case gate.GateTypeExit:
		return s.handleExit(g, card, plateNumber)
	default:
		return rejectResponse("Bad gate type"), nil
	}
}

// handleEntry xử lý xe vào
func (s *Service) handleEntry(
	g *gate.Gate,
	card *rfid_card.RfidCard,
	plateNumber string,
) (*RfidScanResponse, error) {

	isAvailable, err := s.parkingSlotService.IsAvailable(g.LotID)
	if err != nil || !isAvailable {
		return rejectResponse("Lot full"), nil
	}
	// Kiểm tra không có session đang active với thẻ này
	existing, err := s.sessionService.FindActiveByCardUID(card.UID)
	if err == nil && existing != nil {
		return rejectResponse("Card in use"), nil
	}

	// Mọi thứ hợp lệ → consume plate khỏi cache
	s.plateCache.Consume(g.ID)

	// Tạo session mới
	_, err = s.sessionService.Create(parking_session.CreateParkingSessionInput{
		LotID:       g.LotID,
		CardUID:     card.UID,
		CardType:    string(card.CardType),
		PlateNumber: plateNumber,
	})
	if err != nil {
		return nil, err
	}

	return &RfidScanResponse{
		Success:  true,
		Action:   "open_barrier",
		LCDLine1: fmt.Sprintf("LP:%s", plateNumber),
		LCDLine2: "Welcome!",
		Message:  "Session created",
	}, nil
}

// handleExit xử lý xe ra
func (s *Service) handleExit(
	g *gate.Gate,
	card *rfid_card.RfidCard,
	plateNumber string,
) (*RfidScanResponse, error) {
	// Tìm session đang active
	session, err := s.sessionService.FindActiveByCardUID(card.UID)
	if err != nil {
		return rejectResponse("No session"), nil
	}

	// Kiểm tra biển số camera quét ở cổng ra khớp với biển số lúc vào
	if session.PlateNumber != plateNumber {
		return rejectResponse("Plate mismatch"), nil
	}

	// Tính phí
	fee := calculateFee(session, card)

	// Mọi thứ hợp lệ → consume plate khỏi cache
	s.plateCache.Consume(g.ID)

	// Kết thúc session
	_, err = s.sessionService.FinishSession(parking_session.FinishParkingSessionInput{
		SessionID: session.ID,
		Fee:       fee,
	})
	if err != nil {
		return nil, err
	}

	if card.UserID != nil {
		if fee > 0 {
			money, err := s.userService.GetWalletBalance(*card.UserID)
			if err != nil {
				return rejectResponse("Error Get Balance"), nil
			}
			if money < fee {
				return rejectResponse("Insufficient balance"), nil
			}
			err = s.userService.WithdrawFromWallet(*card.UserID, fee)
			if err != nil {
				return rejectResponse("Error Withdraw Balance"), nil
			}
		}

		// Nếu trừ tiền thành công hoặc fee = 0, mở barie
		return &RfidScanResponse{
			Success:  true,
			Action:   "open_barrier",
			LCDLine1: fmt.Sprintf("LP:%s", plateNumber),
			LCDLine2: fmt.Sprintf("Fee:%dVND", fee),
			Message:  "Session finished",
		}, nil
	}

	// Trường hợp thẻ vãng lai (Guest, không có UserID) -> Trả tiền mặt, không tự động mở barie
	return &RfidScanResponse{
		Success:  true,
		Action:   "reject", // Không tự mở barie
		LCDLine1: "Please pay cash",
		LCDLine2: fmt.Sprintf("Fee:%dVND", fee),
		Message:  "Guest must pay cash",
	}, nil
}

// ─── helpers ─────────────────────────────────────────────────────────────────

func rejectResponse(msg string) *RfidScanResponse {
	return &RfidScanResponse{
		Success:  false,
		Action:   "reject",
		LCDLine1: "REJECTED",
		LCDLine2: msg,
		Message:  msg,
	}
}

// calculateFee tính phí đơn giản: Nếu là GUEST thì 7k/h, còn Registered là 5k/h, tối thiểu 5000đ
func calculateFee(session *parking_session.ParkingSession, card *rfid_card.RfidCard) int64 {
	const rateGuestPerHour = 7000
	const rateRegisteredPerHour = 5000
	minutes := time.Since(session.EntryTime).Minutes()
	hours := time.Since(session.EntryTime).Hours()
	if minutes != 0 {
		hours++
	}
	if hours == 0 {
		hours = 1
	}
	if card.CardType == rfid_card.CardTypeRegistered && card.IsActive == true {
		return int64(hours * rateRegisteredPerHour)
	}
	return int64(hours * rateGuestPerHour)
}
