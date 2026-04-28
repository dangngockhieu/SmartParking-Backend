package iot_gateway

import (
	"sync"
	"time"
)

const plateCacheTTL = 5 * time.Minute

type plateCacheEntry struct {
	plateNumber string
	storedAt    time.Time
}

// PlateCache lưu biển số tạm thời trong RAM theo gateID.
// Mỗi entry tồn tại tối đa 5 phút, dùng 1 lần rồi xóa (consume).
type PlateCache struct {
	mu    sync.Mutex
	cache map[uint64]plateCacheEntry
}

func NewPlateCache() *PlateCache {
	pc := &PlateCache{
		cache: make(map[uint64]plateCacheEntry),
	}
	go pc.cleanupLoop()
	return pc
}

// Set lưu biển số cho gateID, ghi đè nếu đã tồn tại
func (pc *PlateCache) Set(gateID uint64, plateNumber string) {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	pc.cache[gateID] = plateCacheEntry{
		plateNumber: plateNumber,
		storedAt:    time.Now(),
	}
}

// Consume lấy biển số và xóa khỏi cache (dùng 1 lần).
// Trả về ("", false) nếu không tìm thấy hoặc đã hết TTL.
func (pc *PlateCache) Consume(gateID uint64) (string, bool) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	entry, ok := pc.cache[gateID]
	if !ok {
		return "", false
	}

	delete(pc.cache, gateID)

	if time.Since(entry.storedAt) > plateCacheTTL {
		return "", false
	}

	return entry.plateNumber, true
}

// Peek reads the plate without deleting (for validation before consume)
func (pc *PlateCache) Peek(gateID uint64) (string, bool) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	entry, ok := pc.cache[gateID]
	if !ok {
		return "", false
	}

	if time.Since(entry.storedAt) > plateCacheTTL {
		return "", false
	}

	return entry.plateNumber, true
}

// cleanupLoop dọn dẹp các entry hết TTL mỗi phút
func (pc *PlateCache) cleanupLoop() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		pc.mu.Lock()
		for gateID, entry := range pc.cache {
			if time.Since(entry.storedAt) > plateCacheTTL {
				delete(pc.cache, gateID)
			}
		}
		pc.mu.Unlock()
	}
}
