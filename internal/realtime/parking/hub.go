package parking

import (
	"encoding/json"
	"log"
	"sync"
)

type Session interface {
	Send([]byte) error
	Close() error
}

type Client struct {
	LotID   uint64
	Session Session
}

type Hub struct {
	mu      sync.RWMutex
	clients map[Session]*Client
}

type EventEnvelope struct {
	Event string `json:"event"`
	Data  any    `json:"data"`
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[Session]*Client),
	}
}

func (h *Hub) Add(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[client.Session] = client
}

func (h *Hub) Remove(s Session) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.clients, s)
}

func (h *Hub) BroadcastToLot(lotID uint64, event string, data any) {
	h.broadcast(event, data, func(client *Client) bool {
		return client.LotID == lotID
	})
}

func (h *Hub) BroadcastAll(event string, data any) {
	h.broadcast(event, data, func(client *Client) bool {
		return true
	})
}

func (h *Hub) broadcast(event string, data any, filter func(*Client) bool) {
	payload, err := json.Marshal(EventEnvelope{
		Event: event,
		Data:  data,
	})
	if err != nil {
		log.Printf("[WT-HUB] marshal failed event=%s err=%v", event, err)
		return
	}

	h.mu.RLock()

	clients := make([]*Client, 0)

	for _, client := range h.clients {
		if filter(client) {
			clients = append(clients, client)
		}
	}

	h.mu.RUnlock()

	for _, client := range clients {
		if err := client.Session.Send(payload); err != nil {
			log.Printf(
				"[WT-HUB] send failed lotID=%d event=%s err=%v",
				client.LotID,
				event,
				err,
			)

			h.Remove(client.Session)
			_ = client.Session.Close()
			continue
		}
	}
}
